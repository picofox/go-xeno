package transcomm

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"reflect"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
)

type TCPServerConnection struct {
	_conn           net.Conn
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndpoint inet.IPV4EndPoint
	_recvBuffer     *memory.LinearBuffer
	_sendBuffer     *memory.LinearBuffer
	_pipeline       []IServerHandler
	_server         *TCPServer
}

func (ego *TCPServerConnection) Shutdown() {
	ego._conn.Close()
}
func (ego *TCPServerConnection) checkRecvBufferCapacity() int32 {
	if ego._recvBuffer.WriteAvailable() < O1L15O1T15_HEADER_SIZE {
		wa := ego._recvBuffer.Compact()
		if wa >= O1L15O1T15_HEADER_SIZE {
			return core.MkSuccess(0)
		}
		if ego._recvBuffer.Capacity() < MAX_BUFFER_MAX_CAPACITY {
			if ego._recvBuffer.ResizeTo(ego._recvBuffer.Capacity()*2) > 0 {
				return core.MkSuccess(0)
			}
			return core.MkErr(core.EC_RESPACE_FAILED, 1)
		} else {
			return core.MkErr(core.EC_REACH_LIMIT, 1)
		}
	}

	return core.MkSuccess(0)
}
func (ego *TCPServerConnection) ioRoutine() {
	rc := ego.checkRecvBufferCapacity()
	if core.IsErrType(rc, core.EC_REACH_LIMIT) {
		ego._server.Log(core.LL_ERR, "[SNH] Buffer reach max")
		return //TODO close connection
	}
	nDone, err := ego._conn.Read((*ego._recvBuffer.InternalData())[ego._recvBuffer.WritePos():])
	if nDone < 0 {
		if err == nil {
			//TODO close remove this one
		} else {
			//TODO close remove this one
		}
		ego._conn.Close()
		return
	} else if nDone == 0 {
		ego._conn.Close()
		return
	}
	var bufParam any = ego._recvBuffer
	var p2 any = nil
	for _, handler := range ego._pipeline {
		rc, bufParam, p2 = handler.OnReceive(ego, bufParam, p2)
		if core.Err(rc) {
			return
		}
	}
	ego._server.OnIncomingMessage(ego, bufParam, p2)
}

func (ego *TCPServerConnection) TryRead() int {
	n, err := ego._conn.Read(ego._buffer[ego._length:])

	if err != nil {
		if err == io.EOF {
			logging.Log(core.LL_ERR, "Read Conn <%s> Closed", ego.String())
			return -1
		} else if errors.Is(err, os.ErrDeadlineExceeded) {
			return 0
		} else {
			logging.Log(core.LL_ERR, "Read Conn <%s> Error: %s", ego.String(), err.Error())
			return -2
		}
	}
	ego._length = ego._length + int64(n)

	return n
}

func (ego *TCPServerConnection) String() string {
	return fmt.Sprintf("<%s> -> <%s>", ego._remoteEndpoint.EndPointString(), ego._localEndPoint.EndPointString())
}

func (ego *TCPServerConnection) Identifier() int64 {
	return ego._remoteEndpoint.Identifier()
}

func NeoTCPServerConnection(conn net.Conn, config *config.NetworkServerTCPConfig) *TCPServerConnection {
	c := TCPServerConnection{
		_conn:           conn,
		_localEndPoint:  inet.NeoIPV4EndPointByAddr(conn.LocalAddr()),
		_remoteEndpoint: inet.NeoIPV4EndPointByAddr(conn.RemoteAddr()),
		_recvBuffer:     memory.NeoLinearBuffer(1024),
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

	return &c
}
