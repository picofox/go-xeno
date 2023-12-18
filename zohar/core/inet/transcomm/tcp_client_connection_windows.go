package transcomm

import (
	"fmt"
	"net"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type TCPClientConnection struct {
	_index          int
	_conn           *net.TCPConn
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_recvBuffer     *memory.RingBuffer
	_sendBuffer     *memory.LinearBuffer
	_pipeline       []IClientCodecHandler
	_client         *TCPClient
	_isConnected    bool
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
		_recvBuffer:     memory.NeoRingBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_pipeline:       make([]IClientCodecHandler, 0),
		_client:         client,
		_isConnected:    false,
	}
	return &c
}

var _ IConnection = &TCPClientConnection{}
