package transcomm

import (
	"fmt"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/memory"
)

type TCPServerConnection struct {
	_fd                       int
	_localEndPoint            inet.IPV4EndPoint
	_remoteEndPoint           inet.IPV4EndPoint
	_recvBuffer               *memory.LinkedListByteBuffer
	_sendBuffer               *memory.LinkedListByteBuffer
	_sendLock                 sync.Mutex
	_server                   *TCPServer
	_reactorIndex             uint32
	_profiler                 *prof.ConnectionProfiler
	_ev                       EPoolEventDataSubReactor
	_outgoingHeader           *memory.O1L31C16Header
	_incomingHeaderBuffer     [6]byte
	_incomingHeaderBufferRIdx int64
	_incomingDataIndex        int64
	_incomingHeader           memory.O1L31C16Header
	_keepalive                *KeepAlive
	_stateCode                uint8
}

func (ego *TCPServerConnection) OnIOError() int32 {
	rc := ego._server.OnIOError(ego)
	ego.Close()
	return rc
}

func (ego *TCPServerConnection) _flushSendingBuffer() int32 {
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

func (ego *TCPServerConnection) FlushSendingBuffer() int32 {
	ego._sendLock.Lock()
	defer ego._sendLock.Unlock()
	return ego._flushSendingBuffer()
}

func (ego *TCPServerConnection) GetEV() *EPoolEventDataSubReactor {
	return &ego._ev
}

func (ego *TCPServerConnection) Pulse(ts int64) {
	if ego._keepalive != nil {
		rc := ego._keepalive.Pulse(ego, ts)
		if core.IsErrType(rc, core.EC_TCP_CONNECT_ERROR) {
			ego.OnDisconnected()
		}
	}
	strProf := ego._profiler.String()
	ego._server.Log(core.LL_INFO, strProf)
}

func (ego *TCPServerConnection) Close() {
	err := syscall.Close(ego._fd)
	if err != nil {
		ego._server.Log(core.LL_ERR, "Close connection <%s> Failed. %s", ego.String(), err.Error())
	}
	ego.reset()
}

func (ego *TCPServerConnection) reset() {
	ego._stateCode = Closed
	ego._fd = -1
	ego._recvBuffer.Clear()
	ego._sendBuffer.Clear()
	ego._profiler.Reset()
	ego._incomingHeaderBufferRIdx = 0
	ego._incomingDataIndex = 0

}

func (ego *TCPServerConnection) KeepAliveConfig() *intrinsic.KeepAliveConfig {
	return &ego._server._config.KeepAlive
}

func (ego *TCPServerConnection) OnKeepAlive(ts int64, delta int32) {
	if ego._keepalive != nil {
		ego._keepalive.OnRoundTripBack(ts)
		if delta >= 0 {
			ego._profiler.GetRTTProf().OnUpdate(delta)
			ego._server.Log(core.LL_DEBUG, "conn %s prof: %s", ego.String(), ego._profiler.String())
		}
	}
}

func (ego *TCPServerConnection) ReactorIndex() uint32 {
	return ego._reactorIndex
}

func (ego *TCPServerConnection) SetReactorIndex(u uint32) {
	ego._reactorIndex = u
}

func (ego *TCPServerConnection) OnDisconnected() int32 {
	ego._server.Log(core.LL_WARN, "TCPClientConnection <%s> Disconnected.", ego.String())
	ego._server.OnDisconnected(ego)
	ego.Close()
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) OnPeerClosed() int32 {
	ego._server.Log(core.LL_WARN, "TCPClientConnection <%s> Peer Closed.", ego.String())
	ego._server.OnDisconnected(ego)
	ego.Close()
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) OnWritable() int32 {
	if ego._stateCode == Connected {
		return core.MkErr(core.EC_NOOP, 0)
	}
	ego._stateCode = Connected
	lsa, _ := syscall.Getsockname(ego._fd)
	ego._localEndPoint = inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, lsa)
	rsa, _ := syscall.Getpeername(ego._fd)
	ego._remoteEndPoint = inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, rsa)

	ego._server.Log(core.LL_DEBUG, "Add conn %s, id %d", ego.String(), ego.Identifier())

	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) Type() int8 {
	return CONNTYPE_TCP_SERVER
}

func (ego *TCPServerConnection) Identifier() int64 {
	return ego.RemoteEndPoint().Identifier()
}

func (ego *TCPServerConnection) PreStop() {

}

func (ego *TCPServerConnection) RemoteEndPoint() *inet.IPV4EndPoint {
	return &ego._remoteEndPoint
}

func (ego *TCPServerConnection) LocalEndPoint() *inet.IPV4EndPoint {
	return &ego._localEndPoint
}

func (ego *TCPServerConnection) String() string {
	return fmt.Sprintf("%s->%s[%d]", ego._remoteEndPoint.EndPointString(), ego._localEndPoint.EndPointString(), ego.Identifier())
}

func (ego *TCPServerConnection) OnIncomingData() int32 {
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
					rc := ego._server.OnIncomingMessage(ego, msg)
					if core.Err(rc) {
						ego._server.Log(core.LL_WARN, "msg %s routing failed err:%s", ego._incomingHeader.String(), core.ErrStr(rc))
					}
				} else {
					if rLen == -1 {
						ego._server.Log(core.LL_WARN, "msg %s not found", ego._incomingHeader.String())
					} else {
						ego._server.Log(core.LL_ERR, "msg %s Deserialize Failed", ego._incomingHeader.String())
					}
				}
				ego._server.Log(core.LL_DEBUG, "Cli-Conn [%x] Got msg [%d-%d] l:%d \n", ego.Identifier(), msg.GroupType(), msg.Command(), ego._incomingHeader.BodyLength())
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
func (ego *TCPServerConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	if ego._stateCode != Connected {
		return core.MkErr(core.EC_NOOP, 0)
	}
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

var _ IConnection = &TCPServerConnection{}
