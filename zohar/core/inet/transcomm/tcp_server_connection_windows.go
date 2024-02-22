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

type TCPServerConnection struct {
	_conn                     *net.TCPConn
	_localEndPoint            inet.IPV4EndPoint
	_remoteEndPoint           inet.IPV4EndPoint
	_server                   *TCPServer
	_profiler                 *prof.ConnectionProfiler
	_lastPulseTs              int64
	_sendBuffer               *memory.LinkedListByteBuffer
	_recvBuffer               *memory.LinkedListByteBuffer
	_sendLock                 sync.Mutex
	_outgoingHeader           *memory.O1L31C16Header
	_incomingHeaderBuffer     [6]byte
	_incomingHeaderBufferRIdx int
	_incomingDataIndex        int64
	_incomingHeader           *memory.O1L31C16Header
	_keepalive                *KeepAlive
}

func (ego *TCPServerConnection) FlushSendingBuffer() int32 {
	ego._sendLock.Lock()
	defer ego._sendLock.Unlock()
	return ego._flushSendingBuffer()
}

func (ego *TCPServerConnection) KeepAliveConfig() *intrinsic.KeepAliveConfig {
	return &ego._server._config.KeepAlive
}

func (ego *TCPServerConnection) _flushSendingBuffer() int32 {
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

func (ego *TCPServerConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
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

func (ego *TCPServerConnection) Close() int32 {
	ego._conn.Close()
	ego.reset()
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) OnPeerClosed() int32 {
	rc := ego._server.OnPeerClosed(ego)
	ego.Close()
	return rc
}

func (ego *TCPServerConnection) OnDisconnected() int32 {
	rc := ego._server.OnDisconnected(ego)
	ego.Close()
	return rc
}

func (ego *TCPServerConnection) OnIOError() int32 {
	rc := ego._server.OnIOError(ego)
	ego.Close()
	return rc
}

func (ego *TCPServerConnection) OnConnectingFailed() int32 {
	rc := ego._server.OnDisconnected(ego)
	ego.reset()
	return rc
}

func (ego *TCPServerConnection) ReactorIndex() uint32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPServerConnection) SetReactorIndex(u uint32) {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPServerConnection) OnWritable() int32 {
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) Type() int8 {
	return CONNTYPE_TCP_SERVER
}

func (ego *TCPServerConnection) RemoteEndPoint() *inet.IPV4EndPoint {
	return &ego._remoteEndPoint
}

func (ego *TCPServerConnection) LocalEndPoint() *inet.IPV4EndPoint {
	return &ego._localEndPoint
}

func (ego *TCPServerConnection) Shutdown() {
	ego._conn.Close()
}

func (ego *TCPServerConnection) PreStop() {
	ego._conn.SetReadDeadline(time.Now())
}

func (ego *TCPServerConnection) Pulse(ts int64) {
	if ego._keepalive != nil {
		rc := ego._keepalive.Pulse(ego, ts)
		if core.IsErrType(rc, core.EC_TCP_CONNECT_ERROR) {
			ego.OnConnectingFailed()
		}
	}

	strProf := ego._profiler.String()
	ego._server.Log(core.LL_INFO, strProf)
}

func (ego *TCPServerConnection) OnIncomingData() int32 {
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
			bytesToRead := ego._incomingHeader.BodyLength() - ego._incomingDataIndex
			if bytesToRead <= 0 {
				msg, rLen := messages.GetDefaultMessageBufferDeserializationMapper().DeserializationDispatch(ego._recvBuffer, ego._incomingHeader)
				if msg != nil {
					ego._server.Log(core.LL_DEBUG, "Got msg %s", ego._incomingHeader.String())
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
				ego._server.Log(core.LL_DEBUG, "Svr-Conn [%x] Got msg [%d-%d] l:%d \n", ego.Identifier(), msg.GroupType(), msg.Command(), ego._incomingHeader.BodyLength())
				ego._incomingHeader = nil
				ego._incomingDataIndex = 0
			} else { //not enough data
				buf := ego._recvBuffer.InternalDataForWriting()
				if buf == nil {
					return core.MkErr(core.EC_NULL_VALUE, 1)
				}
				realLenToRead := min(int(bytesToRead), len(buf))
				nDone, err = ego._conn.Read(buf[:realLenToRead])
				ego._profiler.OnBytesReceived(int64(nDone))
				ego._incomingDataIndex += int64(nDone)
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

func (ego *TCPServerConnection) String() string {
	return fmt.Sprintf("<%s> -> <%s>", ego._remoteEndPoint.EndPointString(), ego._localEndPoint.EndPointString())
}

func (ego *TCPServerConnection) Identifier() int64 {
	return ego._remoteEndPoint.Identifier()
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

func (ego *TCPServerConnection) reset() {
	ego._conn = nil
	ego._localEndPoint.SetInvalid()
	ego._remoteEndPoint.SetInvalid()
	ego._server = nil
	ego._profiler = nil
	ego._lastPulseTs = -1
	ego._sendBuffer.Clear()
	ego._recvBuffer.Clear()
	ego._outgoingHeader = nil
	ego._incomingHeaderBufferRIdx = 0
	ego._incomingDataIndex = 0
	ego._incomingHeader = nil
	ego._keepalive = nil
}

func NeoTCPServerConnection(conn *net.TCPConn, listener *ListenWrapper) *TCPServerConnection {
	c := TCPServerConnection{
		_conn:                     conn,
		_localEndPoint:            inet.NeoIPV4EndPointByAddr(conn.LocalAddr()),
		_remoteEndPoint:           inet.NeoIPV4EndPointByAddr(conn.RemoteAddr()),
		_server:                   listener.Server(),
		_profiler:                 prof.NeoConnectionProfiler(),
		_lastPulseTs:              chrono.GetRealTimeMilli(),
		_sendBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_recvBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_outgoingHeader:           memory.NeoO1L31C16Header(0, 0),
		_incomingHeader:           nil,
		_incomingHeaderBufferRIdx: 0,
		_incomingDataIndex:        0,
		_keepalive:                nil,
	}

	if c.KeepAliveConfig().Enable {
		c._keepalive = NeoKeepAlive(c.KeepAliveConfig(), true)
	}

	c._conn.SetNoDelay(c._server._config.NoDelay)

	return &c
}

var _ IConnection = &TCPServerConnection{}
