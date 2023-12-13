package transcomm

import (
	"fmt"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type TCPServerConnection struct {
	_fd             int
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_recvBuffer     *memory.RingBuffer
	_sendBuffer     *memory.LinearBuffer
	_codec          IServerHandler
	_lock           sync.Mutex
	_server         *TCPServer
	_reactorIndex   uint32
}

func (ego *TCPServerConnection) ReactorIndex() uint32 {
	return ego._reactorIndex
}

func (ego *TCPServerConnection) SetReactorIndex(u uint32) {
	ego._reactorIndex = u
}

func (ego *TCPServerConnection) OnDisconnected() int32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPServerConnection) OnConnectingFailed() int32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPServerConnection) OnPeerClosed() int32 {
	//TODO implement me
	panic("implement me")
}

func (ego *TCPServerConnection) OnWritable() int32 {
	return core.MkSuccess(0)
}

func (ego *TCPServerConnection) Type() int8 {
	return CONNTYPE_TCP_SERVER
}

func (ego *TCPServerConnection) Identifier() int64 {
	return ego.RemoteEndPoint().Identifier()
}

func (ego *TCPServerConnection) PreStop() {

}

func (ego *TCPServerConnection) RemoteEndPoint() *inet.IPV4EndPoint {
	return &ego._remoteEndPoint
}

func (ego *TCPServerConnection) LocalEndPoint() *inet.IPV4EndPoint {
	return &ego._localEndPoint
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

func (ego *TCPServerConnection) String() string {
	return fmt.Sprintf("%s->%s[%d]", ego._remoteEndPoint.EndPointString(), ego._localEndPoint.EndPointString(), ego.Identifier())
}

func (ego *TCPServerConnection) OnIncomingData() int32 {
	for {
		rc := ego.checkRecvBufferCapacity()
		if core.IsErrType(rc, core.EC_REACH_LIMIT) {
			ego._server.Log(core.LL_ERR, "[SNH] Buffer reach max")
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
				ego._server.Log(core.LL_SYS, "Connection <%s> SysRead Failed: %d", ego.String(), rc)
			}
			return core.MkErr(core.EC_TCO_RECV_ERROR, 1)
		} else if nDone == 0 {
			//handle close
			return core.MkErr(core.EC_EOF, 1)
		} else {
			msg, rc := ego._codec.OnReceive(ego)
			if core.Err(rc) {
				return rc
			}

			ego._server.OnIncomingMessage(ego, msg, nil)
		}
	}

}

//func (ego *TCPServerConnection) handlerExecuteInbound() {
//	ll := len(ego._pipeline)
//	if ll > 0 {
//		rc := ego._pipeline[0].Inbound(ego._pipeline, 0, ego, ego._recvBuffer, nil)
//		if core.Err(rc) {
//			et, em := core.ExErr(rc)
//			ego._server.Log(core.LL_ERR, "connnection_%s : handler pipeline[%d] failed: (%s) ", ego._remoteEndPoint.EndPointString(), em, core.ErrStr(et))
//		}
//	}
//}

func (ego *TCPServerConnection) flush() (int64, int32) {
	ba, _ := ego._sendBuffer.BytesRef(-1)
	n, err := syscall.Write(ego._fd, ba)
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
		n, err := syscall.Write(ego._fd, ba[offset:totalRemain])
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

//func NeoTcpServerConnection(tcpServer *TCPServer, fd int, rAddr syscall.Sockaddr, lAddr inet.IPV4EndPoint) *TCPServerConnection {
//	ra := inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, rAddr)
//	tsc := TCPServerConnection{
//		_fd:             fd,
//		_localEndPoint:  lAddr,
//		_remoteEndPoint: ra,
//		_recvBuffer:     memory.NeoRingBuffer(1024),
//		_sendBuffer:     memory.NeoLinearBuffer(1024),
//		_server:         tcpServer,
//		_pipeline:       make([]IServerHandler, 0),
//	}
//
//	var output []reflect.Value = make([]reflect.Value, 0, 1)
//	for _, elem := range tcpServer._config.Handlers {
//		rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+elem.Name)
//		if core.Err(rc) {
//			panic(fmt.Sprintf("Install Handler Failed %s", elem.Name))
//		}
//		h := output[0].Interface().(IServerHandler)
//		tsc._pipeline = append(tsc._pipeline, h)
//	}
//	return &tsc
//}

var _ IConnection = &TCPServerConnection{}
