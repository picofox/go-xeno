package transcomm

import (
	"fmt"
	"net"
	"reflect"
	"time"
	"xeno/zohar/core"
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

func (ego *TCPServerConnection) OnIncomingData() int32 {
	rc := ego.checkRecvBufferCapacity()
	if core.IsErrType(rc, core.EC_REACH_LIMIT) {
		ego._server.Log(core.LL_ERR, "[SNH] Buffer reach max")
		return rc //TODO close connection
	}
	baPtr := ego._recvBuffer.InternalData()
	var nDone int = 0
	var err error

	if ego._recvBuffer.ReadPos() == 1020 {
		fmt.Printf("sss")
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
		return core.MkErr(core.EC_EOF, 0)
	} else {
		src := ego._recvBuffer.WriterSeek(memory.BUFFER_SEEK_CUR, int64(nDone))
		if !src {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}

		msg, rc := ego._codec.OnReceive(ego)
		if core.Err(rc) {
			return rc
		}

		ego._server.OnIncomingMessage(ego, msg, nil)
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
