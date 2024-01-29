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
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
)

type TCPServerConnection struct {
	_conn           *net.TCPConn
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_sendBuffer     *memory.LinearBuffer
	_codec          IServerCodecHandler
	_server         *TCPServer
	_profiler       *prof.ConnectionProfiler
	_sendBufferList *memory.ByteBufferList
	_recvBufferList *memory.ByteBufferList
	_lastPulseTs    int64
	_lock           sync.Mutex
}

func (ego *TCPServerConnection) GetBufferNodeForReceiving() *memory.ByteBufferNode {
	byteBuf := ego._recvBufferList.Back()
	if byteBuf == nil || byteBuf.ReadAvailable() <= 0 {
		byteBuf = memory.NeoByteBufferNode(4096)
		if byteBuf == nil {
			ego._server.Log(core.LL_ERR, "Get ByteBufferNode for writing Failed.")
			return nil
		}
		ego._recvBufferList.PushBack(byteBuf)
	}
	return byteBuf
}

func (ego *TCPServerConnection) FlushSendingBuffer() (int64, int32) {
	ego._lock.Lock()
	defer ego._lock.Unlock()

	var sentBytes int64 = 0
	byteBuf := ego._sendBufferList.Front()
	for byteBuf != nil {
		ba, _ := byteBuf.BytesRef(-1)
		if ba == nil {
			return sentBytes, core.MkErr(core.EC_NULL_VALUE, 2)
		}

		remainLength := len(ba)
		if remainLength == 0 {
			ego._server.Log(core.LL_ERR, "Found 0 Len buffer")
			ego._sendBufferList.PopFront()
			memory.GetByteBufferCache().Put(byteBuf)
			continue
		}
		for remainLength > 0 {
			nDone, err := ego._conn.Write(ba)
			if err != nil {
				if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
					return sentBytes, core.MkErr(core.EC_TRY_AGAIN, 0)
				}
				return sentBytes, core.MkErr(core.EC_TCP_SEND_FAILED, 0)
			} else {
				sentBytes += int64(nDone)
				byteBuf.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(nDone))
				remainLength -= nDone
				if byteBuf.ReadAvailable() <= 0 {
					ego._sendBufferList.PopFront()
					memory.GetByteBufferCache().Put(byteBuf)
				}
			}
		}

		byteBuf = ego._sendBufferList.Front()
	}

	return sentBytes, core.MkSuccess(0)

}

func (ego *TCPServerConnection) BufferBlockList() *memory.ByteBufferList {
	return ego._sendBufferList
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

func (ego *TCPServerConnection) clearBufferList() {
	for {
		byteBuf := ego._sendBufferList.PopFront()
		if byteBuf != nil {
			memory.GetByteBufferCache().Put(byteBuf)
		} else {
			break
		}
	}
	for {
		byteBuf := ego._recvBufferList.PopFront()
		if byteBuf != nil {
			memory.GetByteBufferCache().Put(byteBuf)
		} else {
			break
		}
	}
}

func (ego *TCPServerConnection) Close() int32 {
	ego._conn.Close()
	ego._sendBuffer.Clear()
	ego.clearBufferList()
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

func (ego *TCPServerConnection) PreStop() {
	ego._conn.SetReadDeadline(time.Now())
}

func (ego *TCPServerConnection) Pulse(ts int64) {
	ego._codec.Pulse(ego, ts)
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

	byteBuffer := ego.GetBufferNodeForReceiving()
	if byteBuffer == nil {
		//todo close conn
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}

	nDone, err = ego._conn.Read((*byteBuffer.InternalData())[byteBuffer.WritePos():byteBuffer.Capacity()])
	if err != nil {
		fmt.Printf("read Failed %d (%s)\n", nDone, err.Error())
	} else {
		fmt.Printf("read OK %d \n", nDone)
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
						if nowTs-ego._lastPulseTs > int64(intrinsic.GetIntrinsicConfig().Poller.SubReactorPulseInterval) {
							ego.Pulse(nowTs)
							ego._lastPulseTs = nowTs
						}
						return core.MkSuccess(0)
					}
				}
			}
			return core.MkErr(core.EC_TCO_RECV_ERROR, 1)
		}
	} else {
		bOk := byteBuffer.WriterSeek(memory.BUFFER_SEEK_CUR, int64(nDone))
		if !bOk {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		for {
			msg, rc := ego._codec.OnReceive(ego)
			if core.Err(rc) || msg == nil {
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
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_server:         listener.Server(),
		_codec:          nil,
		_profiler:       prof.NeoConnectionProfiler(),
		_sendBufferList: memory.NeoByteBufferList(),
		_recvBufferList: memory.NeoByteBufferList(),
		_lastPulseTs:    chrono.GetRealTimeMilli(),
	}

	c._conn.SetNoDelay(c._server._config.NoDelay)

	var output []reflect.Value = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+listener.Server()._config.Codec, &c)
	if core.Err(rc) {
		panic(fmt.Sprintf("Install Handler Failed %s", listener.Server()._config.Codec))
	}
	h := output[0].Interface().(IServerCodecHandler)
	c._codec = h
	return &c
}

var _ IConnection = &TCPServerConnection{}
