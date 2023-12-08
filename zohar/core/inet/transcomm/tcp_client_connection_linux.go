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
	_pipeline       []IClientHandler
	_client         *TCPClient
	_isConnected    bool
	_reactorIndex   uint32
}

func (ego *TCPClientConnection) ReactorIndex() uint32 {
	return ego._reactorIndex
}

func (ego *TCPClientConnection) SetReactorIndex(u uint32) {
	ego._reactorIndex = u
}

func (ego *TCPClientConnection) reconnect() int32 {
	syscall.Close(ego._fd)
	ego._client.Log(core.LL_SYS, "Reconnect to <%s>", ego._remoteEndPoint.EndPointString())
	ego._isConnected = false
	ego._fd = -1
	ego._recvBuffer.Clear()
	ego._sendBuffer.Clear()
	return ego.Connect()
}

func (ego *TCPClientConnection) OnDisconnected() int32 {
	ego._client._poller.OnConnectionRemove(ego)
	ego._client.Log(core.LL_WARN, "TCPClientConnection <%s> Disconnected.", ego.String())
	return ego.reconnect()
}

func (ego *TCPClientConnection) OnConnectingFailed() int32 {
	ego._client._poller.OnConnectionRemove(ego)
	ego._client.Log(core.LL_WARN, "TCPClientConnection <%s> Connect Failed.", ego.String())
	ego.reconnect()
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnPeerClosed() int32 {
	ego._client._poller.OnConnectionRemove(ego)
	ego._client.Log(core.LL_WARN, "TCPClientConnection <%s> Peer Closed.", ego.String())
	ego.reconnect()
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) OnWritable() int32 {
	ego._isConnected = true
	lsa, _ := syscall.Getsockname(ego._fd)
	ego._localEndPoint = inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, lsa)
	rsa, _ := syscall.Getpeername(ego._fd)
	ego._remoteEndPoint = inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, rsa)
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
func (ego *TCPClientConnection) checkRecvBufferCapacity() int32 {
	if ego._recvBuffer.WriteAvailable() > 0 {
		return core.MkSuccess(0)
	}

	if ego._recvBuffer.Capacity() < MAX_BUFFER_MAX_CAPACITY {
		neoSz := ego._recvBuffer.Capacity() * 2
		if neoSz > MAX_BUFFER_MAX_CAPACITY {
			neoSz = MAX_BUFFER_MAX_CAPACITY
		}
		if ego._recvBuffer.ResizeTo(neoSz) > 0 {
			return core.MkSuccess(0)
		}
	}

	return core.MkErr(core.EC_REACH_LIMIT, 1)
}
func (ego *TCPClientConnection) OnIncomingData() int32 {
	for {
		rc := ego.checkRecvBufferCapacity()
		if core.IsErrType(rc, core.EC_REACH_LIMIT) {
			ego._client.Log(core.LL_ERR, "[SNH] Buffer reach max")
			return core.MkErr(core.EC_REACH_LIMIT, 1) //TODO close connection
		}
		baPtr := ego._recvBuffer.InternalData()
		var nDone int64
		if ego._recvBuffer.WritePos() >= ego._recvBuffer.ReadPos() {
			nDone, rc = inet.SysRead(ego._fd, (*baPtr)[ego._recvBuffer.WritePos():ego._recvBuffer.Capacity()])
		} else {
			nDone, rc = inet.SysRead(ego._fd, (*baPtr)[ego._recvBuffer.WritePos():ego._recvBuffer.ReadPos()])
		}

		if nDone < 0 {
			if core.Err(rc) {
				ego._client.Log(core.LL_SYS, "Connection <%s> SysRead Failed: %d", ego.String(), rc)
			}
			ego.OnDisconnected()
			return core.MkErr(core.EC_TCO_RECV_ERROR, 1)
		} else if nDone == 0 {
			ego._client.Log(core.LL_SYS, "Connection <%s> Closed", ego.String())
			ego.OnPeerClosed()
			return core.MkErr(core.EC_EOF, 1)
		} else {
			var bufParam any = ego._recvBuffer
			var p2 any = nil
			var l int64 = 0
			for _, handler := range ego._pipeline {
				rc, bufParam, l, p2 = handler.OnReceive(ego, bufParam, l, p2)
				if core.Err(rc) {
					return core.MkErr(core.EC_MESSAGE_HANDLING_ERROR, 1)
				}
			}
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
		_index:          index,
		_fd:             -1,
		_localEndPoint:  inet.NeoIPV4EndPointByIdentifier(-1),
		_remoteEndPoint: rAddr,
		_recvBuffer:     memory.NeoRingBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_pipeline:       make([]IClientHandler, 0),
		_client:         client,
		_isConnected:    false,
	}
	return &c
}

var _ IConnection = &TCPClientConnection{}
