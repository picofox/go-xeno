package transcomm

import (
	"fmt"
	"reflect"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
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
	_packetHeader   message_buffer.MessageHeader
	_lock           sync.Mutex
}

func (ego *TCPClientConnection) ReactorIndex() uint32 {
	return ego._reactorIndex
}

func (ego *TCPClientConnection) SetReactorIndex(u uint32) {
	ego._reactorIndex = u
}

func (ego *TCPClientConnection) reconnect() int32 {
	for _, handler := range ego._pipeline {
		handler.Clear()
	}
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

func (ego *TCPClientConnection) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	var param any = msg
	var rc int32 = 0
	var pLen int64

	for i := len(ego._pipeline) - 1; i >= 0; i-- {
		rc, param, pLen, bFlush = ego._pipeline[i].OnSend(ego, param, pLen, bFlush)
		if core.Err(rc) {
			return core.MkErr(core.EC_MESSAGE_HANDLING_ERROR, 1)
		}
	}
	return core.MkSuccess(0)
}

func (ego *TCPClientConnection) sendNImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	var totalRemain int64 = offset + length
	if length < 0 {
		totalRemain = int64(len(ba))
	}

	if totalRemain <= 0 {
		return totalRemain, core.MkSuccess(0)
	}

	n, err := syscall.Write(ego._fd, ba[offset:totalRemain])
	if err != nil {
		if err == syscall.EAGAIN {
			return totalRemain, core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		return totalRemain, core.MkErr(core.EC_TCP_SEND_FAILED, 1)
	}
	totalRemain -= int64(n)
	offset += int64(n)

	return totalRemain, core.MkSuccess(0)
}

func (ego *TCPClientConnection) sendImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	if ego._sendBuffer.WritePos()+length >= message_buffer.MAX_BUFFER_MAX_CAPACITY {
		ego.flush()
	}
	nLeft, rc := ego.sendNImmediately(ba, offset, length)
	if core.Err(rc) {
		return length - nLeft, rc
	}
	return length - nLeft, core.MkSuccess(0)
}

func (ego *TCPClientConnection) SendImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	return ego.sendImmediately(ba, offset, length)
}

func (ego *TCPClientConnection) Send(ba []byte, offset int64, length int64) (int64, int32) {
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

func (ego *TCPClientConnection) flush() (int64, int32) {
	if ego._sendBuffer.ReadAvailable() <= 0 {
		return 0, core.MkSuccess(0)
	}
	ba, _ := ego._sendBuffer.BytesRef(-1)
	n, err := inet.SysWriteN(ego._fd, ba)
	if core.Err(err) {
		if n > 0 {
			ego._sendBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, int64(n))
		}
		if core.IsErrType(err, core.EC_TRY_AGAIN) {
			return int64(n), core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		return int64(n), core.MkErr(core.EC_TCP_SEND_FAILED, 0)
	}
	ego._sendBuffer.Clear()
	return int64(n), core.MkSuccess(0)
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
		_packetHeader:   message_buffer.NeoMessageHeader(),
	}
	var output = make([]reflect.Value, 0, 1)
	for _, elem := range c._client._config.Handlers {
		rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+elem.Name)
		if core.Err(rc) {
			panic(fmt.Sprintf("Install Handler Failed %s", elem.Name))
		}
		h := output[0].Interface().(IClientHandler)
		c._pipeline = append(c._pipeline, h)
	}
	return &c
}

var _ IConnection = &TCPClientConnection{}
