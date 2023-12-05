package transcomm

import (
	"sync"
	"xeno/zohar/core/cms"
)

type EventLoop struct {
	sync.Mutex
	_handlers []IServerHandler
	_ctrlChan chan cms.ICMS
	_server   *TCPServer
}

func NeoEventLoop(server *TCPServer) *EventLoop {
	el := EventLoop{
		_server:   server,
		_ctrlChan: make(chan cms.ICMS, 1),
	}
	return &el
}
