package transcomm

import (
	"net"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
)

type ListenWrapper struct {
	_listen      net.Listener
	_bindAddress inet.IPV4EndPoint
	_server      *TCPServer
}

func (ego *ListenWrapper) Server() *TCPServer {
	return ego._server
}

func (ego *ListenWrapper) BindAddr() inet.IPV4EndPoint {
	return ego._bindAddress
}

func (ego *ListenWrapper) Accept() *net.TCPConn {
	conn, err := ego._listen.Accept()
	if err != nil {
		ego._server.Log(core.LL_ERR, "Listener <%s> Accept Failed: (%s)", ego._bindAddress.EndPointString(), err.Error())
		return nil
	}

	return conn.(*net.TCPConn)
}

func (ego *ListenWrapper) PreStrop() {
	if l, ok := ego._listen.(*net.TCPListener); ok {
		// l _is_ a *net.TCPListener inside this block
		l.SetDeadline(time.Now())
	}
}

func NeoListenWrapper(server *TCPServer, ep inet.IPV4EndPoint) *ListenWrapper {
	w := ListenWrapper{
		_listen: nil,

		_bindAddress: ep,
		_server:      server,
	}
	return &w
}
