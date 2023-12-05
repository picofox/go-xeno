package transcomm

import (
	"fmt"
	"reflect"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
)

type TcpServerConnection struct {
	_fd             int
	_localEndPoint  inet.IPV4EndPoint
	_remoteEndPoint inet.IPV4EndPoint
	_recvBuffer     *memory.LinearBuffer
	_sendBuffer     *memory.LinearBuffer
	_pipeline       []IServerHandler
	_lock           sync.Mutex
	_server         *TCPServer
}

func (ego *TcpServerConnection) FileDescriptor() int {
	return ego._fd
}

func (ego *TcpServerConnection) checkRecvBufferCapacity() int32 {
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

func (ego *TcpServerConnection) OnIncomingData() {
	for {
		rc := ego.checkRecvBufferCapacity()
		if core.IsErrType(rc, core.EC_REACH_LIMIT) {
			ego._server.Log(core.LL_ERR, "[SNH] Buffer reach max")
			return //TODO close connection
		}
		baPtr := ego._recvBuffer.InternalData()

		nDone, rc := inet.SysRead(ego._fd, (*baPtr)[ego._recvBuffer.WritePos():])
		if nDone < 0 {
			if core.Err(rc) {
				fmt.Println("Sysread return error ")
			}
			return
		} else if nDone == 0 {
			//handle close
			return
		} else {
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
	}

}

func (ego *TcpServerConnection) flush() (int64, int32) {
	ba, _ := ego._sendBuffer.BytesRef()
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

func (ego *TcpServerConnection) sendNImmediately(ba []byte, offset int64, length int64) (int64, int32) {
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

func (ego *TcpServerConnection) sendImmediately(ba []byte, offset int64, length int64) (int64, int32) {
	if ego._sendBuffer.WritePos()+length >= MAX_BUFFER_MAX_CAPACITY {
		ego.flush()
	}
	nLeft, rc := ego.sendNImmediately(ba, offset, length)
	if core.Err(rc) {
		return length - nLeft, rc
	}
	return length - nLeft, core.MkSuccess(0)
}
func (ego *TcpServerConnection) Send(ba []byte, offset int64, length int64) (int64, int32) {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	if ego._sendBuffer.WritePos()+length <= MAX_BUFFER_MAX_CAPACITY {
		ego._sendBuffer.WriteRawBytes(ba, offset, length)
		return length, core.MkSuccess(0)
	} else if length <= MAX_BUFFER_MAX_CAPACITY {
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

func NeoTcpServerConnection(tcpServer *TCPServer, fd int, rAddr syscall.Sockaddr, lAddr inet.IPV4EndPoint) *TcpServerConnection {
	ra := inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, rAddr)
	tsc := TcpServerConnection{
		_fd:             fd,
		_localEndPoint:  lAddr,
		_remoteEndPoint: ra,
		_recvBuffer:     memory.NeoLinearBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_server:         tcpServer,
		_pipeline:       make([]IServerHandler, 0),
	}

	var output []reflect.Value = make([]reflect.Value, 0, 1)
	for _, elem := range tcpServer._config.Handlers {
		rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+elem.Name)
		if core.Err(rc) {
			panic(fmt.Sprintf("Install Handler Failed %s", elem.Name))
		}
		h := output[0].Interface().(IServerHandler)
		tsc._pipeline = append(tsc._pipeline, h)
	}
	return &tsc
}
