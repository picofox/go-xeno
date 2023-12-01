package server

import (
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

type TcpServerConnection struct {
	_fd              int
	_localEndPoint   inet.IPV4EndPoint
	_remoteEndPoint  inet.IPV4EndPoint
	_recvBuffer      []byte
	_recvBufferIndex int64

	_sendBuffer *memory.LinearBuffer
	_pipeline   []IServerHandler
	_lock       sync.Mutex
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
	if ego._sendBuffer.WritePos()+length >= 32768 {
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
	if ego._sendBuffer.WritePos()+length <= 32768 {
		ego._sendBuffer.WriteRawBytes(ba, offset, length)
		return length, core.MkSuccess(0)
	} else if length <= 32768 {
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

func (ego *TcpServerConnection) RecvBufferLength() int64 {
	return ego._recvBufferIndex
}

func (ego *TcpServerConnection) RecvBufferCapacity() int64 {
	return int64(cap(ego._recvBuffer))
}

func NeoTcpServerConnection(fd int, rAddr syscall.Sockaddr, lAddr inet.IPV4EndPoint) *TcpServerConnection {
	ipv4 := strs.IPV4BytesToString(rAddr.(*syscall.SockaddrInet4).Addr[0:4])
	port := uint16(rAddr.(*syscall.SockaddrInet4).Port)
	ra := inet.NeoIPV4EndPointByStrIP(inet.EP_PROTO_TCP, 0, 0, ipv4, port)
	tsc := TcpServerConnection{
		_fd:              fd,
		_localEndPoint:   lAddr,
		_remoteEndPoint:  ra,
		_recvBuffer:      make([]byte, 1024),
		_recvBufferIndex: 0,
		_sendBuffer:      memory.NeoLinearBuffer(1024),
		_pipeline:        make([]IServerHandler, 0),
	}
	return &tsc
}
