package transcomm

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"syscall"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/memory"
)

type TCPClientConnection struct {
	_index                    int
	_conn                     *net.TCPConn
	_localEndPoint            inet.IPV4EndPoint
	_remoteEndPoint           inet.IPV4EndPoint
	_client                   *TCPClient
	_profiler                 *prof.ConnectionProfiler
	_lastPulseTs              int64
	_isConnected              bool
	_sendBuffer               *memory.LinkedListByteBuffer
	_recvBuffer               *memory.LinkedListByteBuffer
	_sendLock                 sync.Mutex
	_outgoingHeader           *memory.O1L31C16Header
	_incomingHeaderBuffer     [6]byte
	_incomingHeaderBufferRIdx int
	_incomingHeader           *memory.O1L31C16Header
	_keepalive                *KeepAlive
}

func (ego *TCPClientConnection) reset() {
	ego._index = -1
	ego._conn = nil
	ego._localEndPoint.SetInvalid()
	ego._remoteEndPoint.SetInvalid()
	ego._client = nil
	ego._isConnected = false
	ego._lastPulseTs = -1
	ego._sendBuffer.Clear()
	ego._recvBuffer.Clear()
	ego._outgoingHeader = nil
	ego._incomingHeaderBufferRIdx = 0
	ego._incomingHeader = nil
}

func (ego *TCPClientConnection) FlushSendingBuffer() int32 {
	ego._sendLock.Lock()
	defer ego._sendLock.Unlock()
	return ego._flushSendingBuffer()
}

func (ego *TCPClientConnection) _flushSendingBuffer() int32 {
	for {
		buf := ego._sendBuffer.InternalDataForReading()
		if buf != nil {
			nDone, err := ego._conn.Write(buf)
			ego._profiler.OnBytesSent(int64(nDone))
			if err != nil {
				if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
					return core.MkErr(core.EC_TRY_AGAIN, 0)
				}
				return core.MkErr(core.EC_TCP_SEND_FAILED, 0)
			} else {
				if !ego._sendBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(nDone)) {
					return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}

				//ego._server.Log()
			}
		} else {
			return core.MkSuccess(0)
		}
	}
}

func (ego *TCPClientConnection) Pulse(ts int64) {
	ego._client.Log(core.LL_DEBUG, "pulse ")
	if ego._keepalive != nil {
		rc := ego._keepalive.Pulse(ego, ts)
		if core.IsErrType(rc, core.EC_TCP_CONNECT_ERROR) {
			ego.OnConnectingFailed()
		}
	}
	strProf := ego._profiler.String()
	ego._client.Log(core.LL_INFO, strProf)
}

func (ego *TCPClientConnection) KeepAliveConfig() *intrinsic.KeepAliveConfig {
	return &ego._client._config.KeepAlive
}

func (ego *TCPClientConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	ego._sendLock.Lock()
	defer ego._sendLock.Unlock()
	ego._outgoingHeader.SetGroupType(msg.GroupType())
	ego._outgoingHeader.SetCommand(msg.Command())
	_, rc := msg.Serialize(ego._outgoingHeader, ego._sendBuffer)
	if core.Err(rc) {
		return core.MkErr(core.EC_SERIALIZE_FIELD_FAIELD, 1)
	}
	if bFlush {
		rc = ego._flushSendingBuffer()
		if core.Err(rc) {
			return core.MkErr(core.EC_TCP_SEND_FAILED, 1)
		}
	}
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) Close() {
	ego._conn.Close()
	ego.reset()
}

func (ego *TCPClientConnection) OnPeerClosed() int32 {
	rc := ego._client.OnPeerClosed(ego)
	ego.Close()
	return rc
}

func (ego *TCPClientConnection) OnDisconnected() int32 {
	rc := ego._client.OnDisconnected(ego)
	ego.Close()
	return rc
}

func (ego *TCPClientConnection) OnIOError() int32 {
	rc := ego._client.OnIOError(ego)
	ego.Close()
	return rc
}

func (ego *TCPClientConnection) OnConnectingFailed() int32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPClientConnection) ReactorIndex() uint32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPClientConnection) SetReactorIndex(u uint32) {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPClientConnection) OnWritable() int32 {
	ego._isConnected = true
	laddr := ego._conn.LocalAddr()
	ego._localEndPoint = inet.NeoIPV4EndPointByAddr(laddr)
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) Type() int8 {
	return CONNTYPE_TCP_CLIENT
}

func (ego *TCPClientConnection) Connect() int32 {
	var err error
	ego._conn, err = net.DialTCP(ego._remoteEndPoint.ProtoName(), nil, ego._remoteEndPoint.ToTCPAddr())
	if err != nil {
		return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
	}

	ego._conn.SetNoDelay(ego._client._config.NoDelay)

	ego.OnWritable()
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnIncomingData() int32 {
	var nDone int = 0
	var err error
	var nowTs = chrono.GetRealTimeMilli()

	if nowTs-ego._lastPulseTs > int64(intrinsic.GetIntrinsicConfig().Poller.SubReactorPulseInterval) {
		ego.Pulse(nowTs)
		ego._lastPulseTs = nowTs
	}

	readT0 := time.Duration(intrinsic.GetIntrinsicConfig().Poller.SubReactorPulseInterval)
	d := time.Duration(readT0 * time.Millisecond) // 30 seconds
	w := time.Now()                               // from now
	w = w.Add(d)
	ego._conn.SetReadDeadline(w)
	for {
		if ego._incomingHeader == nil {
			nDone, err = ego._conn.Read(ego._incomingHeaderBuffer[ego._incomingHeaderBufferRIdx:6])
			ego._profiler.OnBytesReceived(int64(nDone))
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					return core.MkErr(core.EC_TRY_AGAIN, 1)
				}
				return core.MkErr(core.EC_TCO_RECV_ERROR, 1)
			}
			ego._incomingHeaderBufferRIdx += nDone
			if ego._incomingHeaderBufferRIdx > 6 {
				panic("too large the rader index")
			} else if ego._incomingHeaderBufferRIdx == 6 {
				ego._incomingHeader = memory.NeoO1L31C16HeaderFromBytes(ego._incomingHeaderBuffer[:])
				ego._incomingHeaderBufferRIdx = 0
			}
		} else {
			rAvail := ego._recvBuffer.ReadAvailable()
			bytesToRead := ego._incomingHeader.BodyLength() - rAvail
			if bytesToRead <= 0 {
				msg, rLen := messages.GetDefaultMessageBufferDeserializationMapper().DeserializationDispatch(ego._recvBuffer, ego._incomingHeader)
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
				ego._client.Log(core.LL_DEBUG, "Cli-Conn [%x] Got msg [%d-%d] l:%d \n", ego.Identifier(), msg.GroupType(), msg.Command(), ego._incomingHeader.BodyLength())
				ego._incomingHeader = nil
			} else { //not enough data
				buf := ego._recvBuffer.InternalDataForWriting()
				if buf == nil {
					return core.MkErr(core.EC_NULL_VALUE, 1)
				}
				realLenToRead := min(int(bytesToRead), len(buf))
				nDone, err = ego._conn.Read(buf[:realLenToRead])
				ego._profiler.OnBytesReceived(int64(nDone))
				if err != nil {
					if errors.Is(err, os.ErrDeadlineExceeded) {
						return core.MkErr(core.EC_TRY_AGAIN, 1)
					}
					return core.MkErr(core.EC_TCO_RECV_ERROR, 1)
				}
				ego._recvBuffer.WriterSeek(memory.BUFFER_SEEK_CUR, int64(nDone))
			}
		}
	}

	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnKeepAlive(ts int64, delta int32) {
	if ego._keepalive != nil {
		ego._keepalive.OnRoundTripBack(ts)
		if delta >= 0 {
			ego._profiler.GetRTTProf().OnUpdate(delta)
			ego._client.Log(core.LL_DEBUG, "conn %s prof: %s", ego.String(), ego._profiler.String())
		}
	}
}

func (ego *TCPClientConnection) Identifier() int64 {
	return ego._remoteEndPoint.Identifier()
}

func (ego *TCPClientConnection) String() string {
	return fmt.Sprintf("<%s> -> <%s>", ego._localEndPoint.EndPointString(), ego._remoteEndPoint.EndPointString())
}

func (ego *TCPClientConnection) PreStop() {

}

func (ego *TCPClientConnection) RemoteEndPoint() *inet.IPV4EndPoint {
	return &ego._remoteEndPoint
}

func (ego *TCPClientConnection) LocalEndPoint() *inet.IPV4EndPoint {
	return &ego._localEndPoint
}

func NeoTCPClientConnection(index int, client *TCPClient, rAddr inet.IPV4EndPoint) *TCPClientConnection {
	c := TCPClientConnection{
		_index:                    index,
		_conn:                     nil,
		_localEndPoint:            inet.NeoIPV4EndPointByIdentifier(-1),
		_remoteEndPoint:           rAddr,
		_client:                   client,
		_isConnected:              false,
		_profiler:                 prof.NeoConnectionProfiler(),
		_lastPulseTs:              chrono.GetRealTimeMilli(),
		_sendBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_recvBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_outgoingHeader:           memory.NeoO1L31C16Header(0, 0),
		_incomingHeader:           nil,
		_incomingHeaderBufferRIdx: 0,
	}

	if c.KeepAliveConfig().Enable {
		c._keepalive = NeoKeepAlive(c.KeepAliveConfig(), false)
	}

	return &c
}

var _ IConnection = &TCPClientConnection{}
