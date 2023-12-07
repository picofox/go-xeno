package transcomm

import (
	"fmt"
	"net"
	"reflect"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
	"xeno/zohar/core/xplatform"
)

type TCPServerConnection struct {
	_conn           *net.TCPConn
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_recvBuffer     *memory.RingBuffer
	_sendBuffer     *memory.LinearBuffer
	_pipeline       []IServerHandler
	_server         *TCPServer
	_fd             uintptr
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

func (ego *TCPServerConnection) FileDescriptor() xplatform.FileDescriptor {
	return xplatform.FileDescriptor(ego._fd)
}

func (ego *TCPServerConnection) Shutdown() {
	ego._conn.Close()
}
func (ego *TCPServerConnection) checkRecvBufferCapacity() int32 {
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
	if ego._recvBuffer.WritePos() >= ego._recvBuffer.ReadPos() {
		nDone, err = ego._conn.Read((*baPtr)[ego._recvBuffer.WritePos():ego._recvBuffer.Capacity()])
	} else {
		nDone, err = ego._conn.Read((*baPtr)[ego._recvBuffer.WritePos():ego._recvBuffer.ReadPos()])
	}

	if nDone < 0 {
		if err != nil {
			ego._server.Log(core.LL_SYS, "Connection <%s> SysRead Failed: %s", ego.String(), err.Error())
		}
		return core.MkErr(core.EC_EOF, 0)
	} else if nDone == 0 {
		return core.MkErr(core.EC_EOF, 0)
	} else {
		var bufParam any = ego._recvBuffer
		var p2 any = nil
		var l int64 = 0
		for _, handler := range ego._pipeline {
			rc, bufParam, l, p2 = handler.OnReceive(ego, bufParam, l, p2)
			if core.Err(rc) {
				return rc
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

func NeoTCPServerConnection(conn *net.TCPConn, config *config.NetworkServerTCPConfig) *TCPServerConnection {
	c := TCPServerConnection{
		_conn:           conn,
		_localEndPoint:  inet.NeoIPV4EndPointByAddr(conn.LocalAddr()),
		_remoteEndPoint: inet.NeoIPV4EndPointByAddr(conn.RemoteAddr()),
		_recvBuffer:     memory.NeoRingBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_server:         nil,
		_pipeline:       make([]IServerHandler, 0),
	}

	var output []reflect.Value = make([]reflect.Value, 0, 1)
	for _, elem := range config.Handlers {
		rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+elem.Name)
		if core.Err(rc) {
			panic(fmt.Sprintf("Install Handler Failed %s", elem.Name))
		}
		h := output[0].Interface().(IServerHandler)
		c._pipeline = append(c._pipeline, h)
	}

	file, err := c._conn.File()
	if err != nil {
		c._server.Log(core.LL_ERR, "Get File From connection <%s> Failed.", c._conn.RemoteAddr().String())
		return nil
	}
	c._fd = file.Fd()

	return &c
}

var _ IConnection = &TCPServerConnection{}
