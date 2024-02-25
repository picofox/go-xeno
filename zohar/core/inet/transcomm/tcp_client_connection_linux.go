package transcomm

import (
	"fmt"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/sched/timer"
)

type TCPClientConnection struct {
	_index                    int
	_fd                       int
	_localEndPoint            inet.IPV4EndPoint
	_remoteEndPoint           inet.IPV4EndPoint
	_client                   *TCPClient
	_lastPulseTs              int64
	_stateCode                uint8
	_reactorIndex             uint32
	_profiler                 *prof.ConnectionProfiler
	_ev                       EPoolEventDataSubReactor
	_sendBuffer               *memory.LinkedListByteBuffer
	_recvBuffer               *memory.LinkedListByteBuffer
	_outgoingHeader           *memory.O1L31C16Header
	_incomingHeaderBuffer     [6]byte
	_incomingHeaderBufferRIdx int64
	_incomingDataIndex        int64
	_incomingHeader           memory.O1L31C16Header
	_keepalive                *KeepAlive
	_sendLock                 sync.Mutex
}

func (ego *TCPClientConnection) Logger() logging.ILogger {
	return ego._client._logger
}

func (ego *TCPClientConnection) OnIOError() int32 {
	rc := ego._client.OnIOError(ego)
	ego.Close()
	return rc
}

func NeoTCPClientConnection(index int, client *TCPClient, rAddr inet.IPV4EndPoint) *TCPClientConnection {
	c := TCPClientConnection{
		_index:                    index,
		_fd:                       -1,
		_localEndPoint:            inet.NeoIPV4EndPointByIdentifier(-1),
		_remoteEndPoint:           rAddr,
		_client:                   client,
		_stateCode:                Initialized,
		_reactorIndex:             0xFFFFFFFF,
		_profiler:                 prof.NeoConnectionProfiler(),
		_sendBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_recvBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_outgoingHeader:           memory.NeoO1L31C16Header(0, 0),
		_incomingHeaderBufferRIdx: 0,
		_incomingDataIndex:        0,
	}

	if c.KeepAliveConfig().Enable {
		c._keepalive = NeoKeepAlive(c.KeepAliveConfig(), false)
	}
	return &c
}

func (ego *TCPClientConnection) _flushSendingBuffer() int32 {
	for {
		buf := ego._sendBuffer.InternalDataForReading()
		if buf != nil {
			nDone, rc := inet.SysWriteN(ego._fd, buf)
			ego._profiler.OnBytesSent(int64(nDone))
			if nDone > 0 {
				if !ego._sendBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(nDone)) {
					return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			if core.Err(rc) {
				return rc
			}
		} else {
			return core.MkSuccess(0)
		}
	}
}

func (ego *TCPClientConnection) FlushSendingBuffer() int32 {
	ego._sendLock.Lock()
	defer ego._sendLock.Unlock()
	return ego._flushSendingBuffer()
}

func (ego *TCPClientConnection) Pulse(ts int64) {
	if ts-ego._lastPulseTs > ego._client._config.PulseInterval {
		if ego._keepalive != nil {
			rc := ego._keepalive.Pulse(ego, ts)
			if core.IsErrType(rc, core.EC_TCP_CONNECT_ERROR) {
				ego.OnDisconnected()
			}
		}
		strProf := ego._profiler.String()
		ego._client.Log(core.LL_INFO, strProf)
		ego._lastPulseTs = ts
	}
}

func (ego *TCPClientConnection) KeepAliveConfig() *intrinsic.KeepAliveConfig {
	return &ego._client._config.KeepAlive
}

func (ego *TCPClientConnection) ReactorIndex() uint32 {
	return ego._reactorIndex
}

func (ego *TCPClientConnection) SetReactorIndex(u uint32) {
	ego._reactorIndex = u
}

func doReconnect(s any) int32 {
	conn := s.(*timer.Timer).Object().(*TCPClientConnection)
	conn._client.Log(core.LL_SYS, "Reconnect %s", conn.String())
	conn.Connect()
	return 0
}

func (ego *TCPClientConnection) Close() {
	err := syscall.Close(ego._fd)
	if err != nil {
		ego._client.Log(core.LL_ERR, "Close connection <%s> Failed. %s", ego.String(), err.Error())
	}
	ego._stateCode = Closed
	ego.reset()

	if ego._client._config.AutoReconnect {
		timer.GetDefaultTimerManager().AddRelTimerMilli(3000, 1, 3000, datatype.TASK_EXEC_EXECUTOR_POOL, doReconnect, ego)
	}
}

func (ego *TCPClientConnection) reset() {
	ego._fd = -1
	ego._recvBuffer.Clear()
	ego._sendBuffer.Clear()
	ego._profiler.Reset()
	ego._incomingHeaderBufferRIdx = 0
	ego._incomingDataIndex = 0
	if ego._keepalive != nil {
		ego._keepalive.Reset()
	}

}

func (ego *TCPClientConnection) OnDisconnected() int32 {
	ego._client.Log(core.LL_WARN, "TCPClientConnection <%s> Disconnected.", ego.String())
	ego._client.OnDisconnected(ego)
	ego.Close()
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnConnectingFailed() int32 {
	ego._client.Log(core.LL_WARN, "TCPClientConnection <%s> Connecting Failed.", ego.String())
	ego._client.OnDisconnected(ego)
	ego.Close()
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnPeerClosed() int32 {
	ego._client.Log(core.LL_WARN, "TCPClientConnection <%s> Peer Closed.", ego.String())
	ego._client.OnDisconnected(ego)
	ego.Close()
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnWritable() int32 {
	if ego._stateCode == Connected {
		return core.MkErr(core.EC_NOOP, 0)
	}
	ego._stateCode = Connected
	lsa, _ := syscall.Getsockname(ego._fd)
	ego._localEndPoint = inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, lsa)
	rsa, _ := syscall.Getpeername(ego._fd)
	ego._remoteEndPoint = inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, rsa)

	ego._client.Log(core.LL_DEBUG, "Add conn %s, id %d", ego.String(), ego.Identifier())

	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) Type() int8 {
	return CONNTYPE_TCP_CLIENT
}

func (ego *TCPClientConnection) Connect() (rc int32) {
	ego._client.Log(core.LL_SYS, "Connect %s", ego.String())
	ego._stateCode = Connecting
	ego._fd, rc = inet.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if core.Err(rc) {
		return rc
	}
	//var tv syscall.Timeval
	//tv.Sec = 10
	//er := syscall.SetsockoptTimeval(ego._fd, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &tv)
	//if er != nil {
	//	fmt.Println(er.Error())
	//}
	//er = syscall.SetsockoptTimeval(ego._fd, syscall.SOL_SOCKET, syscall.SO_SNDTIMEO, &tv)
	//if er != nil {
	//	fmt.Println(er.Error())
	//}

	inet.SysSetTCPNoDelay(ego._fd, ego._client._config.NoDelay)

	sa := ego._remoteEndPoint.ToSockAddr()
	err := syscall.Connect(ego._fd, sa)
	if err != nil {
		if err != syscall.EINPROGRESS && err != syscall.EALREADY && err != syscall.EINVAL && err != syscall.EISCONN {
			ego._client.Log(core.LL_ERR, "TCP Connection to <%s> Error: %s", ego._remoteEndPoint.EndPointString(), err.Error())
			return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
		}
	}
	ego._client._poller.OnIncomingConnection(ego)
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnKeepAlive(ts int64, delta int32) {
	if ego._keepalive != nil {
		ego._keepalive.OnRoundTripBack(ts)
		if delta >= 0 {
			ego._profiler.GetRTTProf().OnUpdate(delta)
		}
	}
}

func (ego *TCPClientConnection) OnIncomingData() int32 {
	var nDone int64 = 0
	var rc int32 = 0

	for {
		if ego._incomingHeaderBufferRIdx < ego._incomingHeader.HeaderLength() {
			nDone, rc = inet.SysRead(ego._fd, ego._incomingHeaderBuffer[ego._incomingHeaderBufferRIdx:ego._incomingHeader.HeaderLength()])
			ego._profiler.OnBytesReceived(nDone)
			if core.Err(rc) {
				return rc
			}
			ego._incomingHeaderBufferRIdx += nDone
			if ego._incomingHeaderBufferRIdx == 6 {
				ego._incomingHeader.SetByBytes(ego._incomingHeaderBuffer[:])
			}

		} else if ego._incomingHeaderBufferRIdx == ego._incomingHeader.HeaderLength() {
			bytesToRead := ego._incomingHeader.BodyLength() - ego._incomingDataIndex
			if bytesToRead == 0 {
				msg, rLen := messages.GetDefaultMessageBufferDeserializationMapper().DeserializationDispatch(ego._recvBuffer, &ego._incomingHeader)
				if msg != nil {
					rc := ego._client.OnIncomingMessage(ego, msg)
					if core.Err(rc) {
						ego._client.Log(core.LL_WARN, "msg %s routing failed err:%s", ego._incomingHeader.String(), core.ErrStr(rc))
					}
				} else {
					if rLen == -1 {
						ego._client.Log(core.LL_WARN, "msg %s not found", ego._incomingHeader.String())
					} else {
						ego._client.Log(core.LL_ERR, "msg %s Deserialize Failed", ego._incomingHeader.String())
					}
				}
				//ego._client.Log(core.LL_DEBUG, "Cli-Conn [%x] Got msg [%d-%d] l:%d \n", ego.Identifier(), msg.GroupType(), msg.Command(), ego._incomingHeader.BodyLength())
				ego._incomingDataIndex = 0
				ego._incomingHeaderBufferRIdx = 0
			} else if bytesToRead > 0 { //not enough data
				buf := ego._recvBuffer.InternalDataForWriting()
				if buf == nil {
					return core.MkErr(core.EC_NULL_VALUE, 1)
				}
				realLenToRead := min(int(bytesToRead), len(buf))
				nDone, rc = inet.SysRead(ego._fd, buf[:realLenToRead])
				ego._profiler.OnBytesReceived(nDone)
				ego._incomingDataIndex += nDone
				if core.Err(rc) {
					return rc
				}
				ego._recvBuffer.WriterSeek(memory.BUFFER_SEEK_CUR, nDone)
			} else {
				panic("bytes to read < 0")
			}

		} else {
			panic("[SNH] too long header index found!")
		}
	}

	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) GetEV() *EPoolEventDataSubReactor {
	return &ego._ev
}

func (ego *TCPClientConnection) Identifier() int64 {
	return ego._localEndPoint.Identifier()
}

func (ego *TCPClientConnection) String() string {
	return fmt.Sprintf("[%x]: <%s> --(%s)--> <%s>", ego.Identifier(), ego._localEndPoint.EndPointString(), ConnStateCodeToString(ego._stateCode), ego._remoteEndPoint.EndPointString())
}

func (ego *TCPClientConnection) PreStop() {

}

func (ego *TCPClientConnection) RemoteEndPoint() *inet.IPV4EndPoint {
	return &ego._remoteEndPoint
}

func (ego *TCPClientConnection) LocalEndPoint() *inet.IPV4EndPoint {
	return &ego._localEndPoint
}

func (ego *TCPClientConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	ego._sendLock.Lock()

	if ego._stateCode != Connected {
		return core.MkErr(core.EC_NOOP, 0)
	}
	defer ego._sendLock.Unlock()

	ego._outgoingHeader.SetGroupType(msg.GroupType())
	ego._outgoingHeader.SetCommand(msg.Command())
	_, rc := msg.Serialize(ego._outgoingHeader, ego._sendBuffer)
	if core.Err(rc) {
		return core.MkErr(core.EC_SERIALIZE_FIELD_FAIELD, 1)
	}

	if bFlush {
		rc = ego._flushSendingBuffer()
		return rc
	}
	return core.MkSuccess(0)
}

var _ IConnection = &TCPClientConnection{}
