package transcomm

import (
	"net"
	"os"
	"xeno/zohar/core/inet"
)

type ListenWrapper struct {
	_listen      net.Listener
	_bindAddress inet.IPV4EndPoint
	_fd          int
	_file        *os.File
	_server      *TCPServer
}

func (ego *ListenWrapper) FileDescriptor() int {
	return ego._fd
}

func (ego *ListenWrapper) Server() *TCPServer {
	return ego._server
}

func (ego *ListenWrapper) BindAddr() inet.IPV4EndPoint {
	return ego._bindAddress
}

func NeoListenWrapper(server *TCPServer, ep inet.IPV4EndPoint) *ListenWrapper {
	//l, err := net.Listen(ep.ProtoName(), ep.EndPointString())
	//file, err := l.(*net.TCPListener).File()
	//if err != nil {
	//	server.Log(core.LL_ERR, "File From Listen Failed: %s", err.Error())
	//	return nil
	//}
	//
	//fd := int(file.Fd())
	//err = syscall.SetNonblock(fd, true)
	//if err != nil {
	//	server.Log(core.LL_ERR, "SetNonblock of fd %d Failed: %s", fd, err.Error())
	//	return nil
	//}

	w := ListenWrapper{
		_listen:      nil,
		_fd:          -1,
		_file:        nil,
		_bindAddress: ep,
		_server:      server,
	}
	return &w
}
