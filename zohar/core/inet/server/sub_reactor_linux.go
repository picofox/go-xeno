package server

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
)

type SubReactor struct {
	pollArgs
	_server          *TcpServer
	_epollDescriptor int
}

func NeoSubReactor(server *TcpServer) *SubReactor {
	sr := SubReactor{
		pollArgs: pollArgs{
			_caps:   128,
			_size:   128,
			_events: make([]inet.EPollEvent, 128),
		},
		_server:          server,
		_epollDescriptor: -1,
	}
	var err error
	sr._epollDescriptor, err = inet.EpollCreate(0)
	if err != nil {
		server.Log(core.LL_ERR, "Sub Reactor EpollCreate failed. err:()", err.Error())
		return nil
	}
	return &sr
}
