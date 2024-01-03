package transcomm

import (
	"fmt"
	"net"
	"reflect"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
)

type TCPClientConnection struct {
	_index          int
	_conn           *net.TCPConn
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_codec          IClientCodecHandler
	_client         *TCPClient
	_profiler       *prof.ConnectionProfiler
	_sendBufferList *memory.ByteBufferList
	_recvBufferList *memory.ByteBufferList
	_isConnected    bool
}

func (ego *TCPClientConnection) FlushSendingBuffer() (int64, int32) {
	var sentBytes int64 = 0
	byteBuf := ego._sendBufferList.Front()
	for byteBuf != nil {
		ba, _ := byteBuf.BytesRef(-1)
		if ba == nil {
			return sentBytes, core.MkErr(core.EC_NULL_VALUE, 2)
		}

		remainLength := len(ba)
		if remainLength == 0 {
			ego._client.Log(core.LL_ERR, "Found 0 Len buffer")
			ego._sendBufferList.PopFront()
			memory.GetByteBuffer4KCache().Put(byteBuf)
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
					memory.GetByteBuffer4KCache().Put(byteBuf)
				}
			}
		}
		byteBuf = ego._sendBufferList.Front()
	}

	return sentBytes, core.MkSuccess(0)
}

func (ego *TCPClientConnection) BufferBlockList() *memory.ByteBufferList {
	return ego._sendBufferList
}

func (ego *TCPClientConnection) reset() {

}

func (ego *TCPClientConnection) Pulse(ts int64) {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPClientConnection) KeepAliveConfig() *intrinsic.KeepAliveConfig {
	return &ego._client._config.KeepAlive
}

func (ego *TCPClientConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	return 0
}

func (ego *TCPClientConnection) OnPeerClosed() int32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPClientConnection) OnDisconnected() int32 {
	//TODO implement me
	panic("implement me")
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

func (ego *TCPClientConnection) flush() (int64, int32) {
	return 0, core.MkSuccess(0)
}

func (ego *TCPClientConnection) sendImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	return 0, core.MkSuccess(0)
}

func (ego *TCPClientConnection) SendImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	return ego.sendImmediately(ba, offset, length)
}

func (ego *TCPClientConnection) OnIncomingData() int32 {
	//TODO implement me
	panic("implement me")
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
		_index:          index,
		_conn:           nil,
		_localEndPoint:  inet.NeoIPV4EndPointByIdentifier(-1),
		_remoteEndPoint: rAddr,
		_codec:          nil,
		_client:         client,
		_isConnected:    false,
		_sendBufferList: memory.NeoByteBufferList(),
		_recvBufferList: memory.NeoByteBufferList(),
		_profiler:       prof.NeoConnectionProfiler(),
	}

	c._conn.SetNoDelay(c._client._config.NoDelay)
	var output = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+c._client._config.Codec)
	if core.Err(rc) {
		return nil
	}
	h := output[0].Interface().(IClientCodecHandler)
	c._codec = h
	return &c
}

var _ IConnection = &TCPClientConnection{}
