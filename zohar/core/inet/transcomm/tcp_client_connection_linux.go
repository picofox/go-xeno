package transcomm

import (
	"fmt"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/memory"
)

type TCPClientConnection struct {
	_index          int
	_fd             int
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_recvBuffer     *memory.RingBuffer
	_sendBuffer     *memory.LinearBuffer
	_pipeline       []IServerHandler
	_client         *TCPClient
	_isConnected    bool
}

func (ego *TCPClientConnection) OnWritable() int32 {
	ego._isConnected = true
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) Type() int8 {
	return CONNTYPE_TCP_CLIENT
}

func (ego *TCPClientConnection) Connect() (rc int32) {
	ego._fd, rc = inet.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if core.Err(rc) {
		return rc
	}
	sa := ego._remoteEndPoint.ToSockAddr()
	err := syscall.Connect(ego._fd, sa)
	if err != nil {
		if err != syscall.EINPROGRESS && err != syscall.EALREADY && err != syscall.EINVAL && err != syscall.EISCONN {
			ego._client.Log(core.LL_ERR, "TCP Connection to <%s> Error: %s", ego._remoteEndPoint.EndPointString(), err.Error())
			return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
		}
	}
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnIncomingData() int32 {
	ego._client.Log(core.LL_DEBUG, "OnIncomingData")
	return core.MkSuccess(0)
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
		_fd:             -1,
		_localEndPoint:  inet.NeoIPV4EndPointByIdentifier(-1),
		_remoteEndPoint: rAddr,
		_recvBuffer:     memory.NeoRingBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_pipeline:       make([]IServerHandler, 0),
		_client:         client,
		_isConnected:    false,
	}
	return &c
}

var _ IConnection = &TCPClientConnection{}
