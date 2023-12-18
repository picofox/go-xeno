package transcomm

import (
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"syscall"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
)

type TCPServerConnection struct {
	_conn           *net.TCPConn
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_recvBuffer     *memory.RingBuffer
	_sendBuffer     *memory.LinearBuffer
	_codec          IServerCodecHandler
	_server         *TCPServer
	_lock           sync.Mutex
}

func (ego *TCPServerConnection) KeepAliveConfig() *intrinsic.KeepAliveConfig {
	return &ego._server._config.KeepAlive
}

func (ego *TCPServerConnection) flush() (int64, int32) {
	ba, _ := ego._sendBuffer.BytesRef(-1)
	if ba == nil {
		return int64(0), core.MkSuccess(0)
	}
	n, err := ego._conn.Write(ba[ego._sendBuffer.ReadPos():ego._sendBuffer.ReadAvailable()])
	if err != nil {
		if err == syscall.EAGAIN {
			if n > 0 {
				ego._sendBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(n))
			}
			return int64(n), core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		if n > 0 {
			ego._sendBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(n))
		}
		return int64(n), core.MkErr(core.EC_TCP_SEND_FAILED, 0)
	}
	ego._sendBuffer.Clear()
	return int64(n), core.MkSuccess(0)
}

func (ego *TCPServerConnection) sendNImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	var totalRemain int64 = length
	for totalRemain > 0 {
		n, err := ego._conn.Write(ba[offset:totalRemain])
		if err != nil {
			if err == syscall.EAGAIN {
				return totalRemain, core.MkErr(core.EC_TRY_AGAIN, 0)
			}
			return totalRemain, core.MkErr(core.EC_TCP_SEND_FAILED, 1)
		}
		totalRemain -= int64(n)
		offset += int64(n)
	}
	return totalRemain, core.MkSuccess(0)
}

func (ego *TCPServerConnection) sendImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	if ego._sendBuffer.WritePos()+length >= message_buffer.MAX_BUFFER_MAX_CAPACITY {
		ego.flush()
	}
	nLeft, rc := ego.sendNImmediately(ba, offset, length)
	if core.Err(rc) {
		return length - nLeft, rc
	}
	return length - nLeft, core.MkSuccess(0)
}

func (ego *TCPServerConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._codec.OnSend(ego, msg, bFlush)
	if core.Err(rc) {
		return core.MkErr(core.EC_MESSAGE_HANDLING_ERROR, 1)
	}
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) Send(ba []byte, offset int64, length int64) (int64, int32) {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	if ego._sendBuffer.WritePos()+length <= message_buffer.MAX_BUFFER_MAX_CAPACITY {
		ego._sendBuffer.WriteRawBytes(ba, offset, length)
		return length, core.MkSuccess(0)
	} else if length <= message_buffer.MAX_BUFFER_MAX_CAPACITY {
		nDone, rc := ego.flush()
		if core.Err(rc) {
			return nDone, rc
		}
		ego._sendBuffer.WriteRawBytes(ba, offset, length)
		return length, core.MkSuccess(0)
	} else {
		nDone, rc := ego.flush()
		if core.Err(rc) {
			return nDone, rc
		}
		nDone, rc = ego.sendImmediately(ba, offset, length)
		return int64(nDone), rc
	}
}

func (ego *TCPServerConnection) Close() int32 {
	ego._conn.Close()
	ego._recvBuffer.Clear()
	ego._sendBuffer.Clear()
	ego._codec.Reset()
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) OnPeerClosed() int32 {
	return ego._server.OnPeerClosed(ego)
}

func (ego *TCPServerConnection) OnDisconnected() int32 {
	return ego._server.OnDisconnected(ego)
}

func (ego *TCPServerConnection) OnConnectingFailed() int32 {
	return ego._server.OnDisconnected(ego)
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
func (ego *TCPServerConnection) checkRecvBufferCapacity() int32 {
	if ego._recvBuffer.WriteAvailable() > 0 {
		return core.MkSuccess(0)
	}

	if ego._recvBuffer.Capacity() < message_buffer.MAX_BUFFER_MAX_CAPACITY {
		neoSz := ego._recvBuffer.Capacity() * 2
		if neoSz > message_buffer.MAX_BUFFER_MAX_CAPACITY {
			neoSz = message_buffer.MAX_BUFFER_MAX_CAPACITY
		}
		if ego._recvBuffer.ResizeTo(neoSz) > 0 {
			return core.MkSuccess(0)
		}
	}

	return core.MkErr(core.EC_REACH_LIMIT, 1)
}

func (ego *TCPServerConnection) PreStop() {
	ego._conn.SetReadDeadline(time.Now())
}

func (ego *TCPServerConnection) Pulse() int32 {
	ego._codec.Pulse(ego, chrono.GetRealTimeMilli())
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) OnIncomingData() int32 {
	rc := ego.checkRecvBufferCapacity()
	if core.IsErrType(rc, core.EC_REACH_LIMIT) {
		ego._server.Log(core.LL_ERR, "[SNH] Buffer reach max")
		return rc //TODO close connection
	}
	baPtr := ego._recvBuffer.InternalData()
	var nDone int = 0
	var err error

	if ego._server._config.KeepAlive.Enable {
		readT0 := time.Duration(ego._server._config.KeepAlive.IntervalMillis)
		d := time.Duration(readT0 * time.Millisecond) // 30 seconds
		w := time.Now()                               // from now
		w = w.Add(d)
		ego._conn.SetReadDeadline(w)
	}

	if ego._recvBuffer.WritePos() >= ego._recvBuffer.ReadPos() {
		nDone, err = ego._conn.Read((*baPtr)[ego._recvBuffer.WritePos():ego._recvBuffer.Capacity()])
	} else {
		nDone, err = ego._conn.Read((*baPtr)[ego._recvBuffer.WritePos():ego._recvBuffer.ReadPos()])
	}

	if nDone < 0 {
		if err != nil {
			ego._server.Log(core.LL_SYS, "Connection <%s> SysRead Failed: %s", ego.String(), err.Error())
		}
		return core.MkErr(core.EC_TCO_RECV_ERROR, 0)
	} else if nDone == 0 {
		if err != nil {
			if err == io.EOF {
				return core.MkErr(core.EC_EOF, 0)
			} else {
				var e *net.OpError
				ok := errors.As(err, &e)
				if ok {
					if e.Timeout() {
						rc = ego.Pulse()
						if !core.Err(rc) {
							return core.MkSuccess(0)
						} else {
							if core.IsErrType(rc, core.EC_TRY_AGAIN) {
								return core.MkSuccess(0)
							}
						}
					}
				}
			}
			return core.MkErr(core.EC_TCO_RECV_ERROR, 1)
		}

	} else {
		src := ego._recvBuffer.WriterSeek(memory.BUFFER_SEEK_CUR, int64(nDone))
		if !src {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}

		for {
			msg, rc := ego._codec.OnReceive(ego)
			if core.Err(rc) {
				return rc
			}

			ego._server.OnIncomingMessage(ego, msg.(message_buffer.INetMessage), nil)
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

func NeoTCPServerConnection(conn *net.TCPConn, listener *ListenWrapper) *TCPServerConnection {
	c := TCPServerConnection{
		_conn:           conn,
		_localEndPoint:  inet.NeoIPV4EndPointByAddr(conn.LocalAddr()),
		_remoteEndPoint: inet.NeoIPV4EndPointByAddr(conn.RemoteAddr()),
		_recvBuffer:     memory.NeoRingBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_server:         listener.Server(),
		_codec:          nil,
	}

	c._conn.SetNoDelay(c._server._config.NoDelay)

	var output []reflect.Value = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+listener.Server()._config.Codec)
	if core.Err(rc) {
		panic(fmt.Sprintf("Install Handler Failed %s", listener.Server()._config.Codec))
	}
	h := output[0].Interface().(IServerCodecHandler)
	c._codec = h
	return &c
}

var _ IConnection = &TCPServerConnection{}
