package server

import (
	"net"
	"os"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
)

type ListenWrapper struct {
	_listen      net.Listener
	_bindAddress inet.IPV4EndPoint
	_fd          int
	_file        *os.File
}

func NeoListenWrapper(tcpServer *TcpServer, ep inet.IPV4EndPoint) *ListenWrapper {
	l, err := net.Listen(tcpServer._bindAddress.ProtoName(), ep.EndPointString())
	file, err := l.(*net.TCPListener).File()
	if err != nil {
		tcpServer.Log(core.LL_ERR, "File From Listen Failed: %s", err.Error())
		return nil
	}

	fd := int(file.Fd())
	err = syscall.SetNonblock(fd, true)
	if err != nil {
		tcpServer.Log(core.LL_ERR, "SetNonblock of fd %d Failed: %s", fd, err.Error())
		return nil
	}

	w := ListenWrapper{
		_listen:      l,
		_fd:          fd,
		_file:        file,
		_bindAddress: ep,
	}
	return &w
}
