package server

import (
	"sync"
	"sync/atomic"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/logging"
)

type TcpServerSubReactor struct {
	_connections     sync.Map
	_connectionCount int32
	_server          *TcpServer
}

func (ego *TcpServerSubReactor) Start() {
	go func() {
		for {
			ego._connections.Range(func(key, value interface{}) bool {
				clientId := key.(int64)
				c := value.(*TcpServerConnection)
				rc := c.TryRead()
				if rc < 0 {
					c.Shutdown()
					ego._connections.Delete(clientId)
					logging.Log(core.LL_INFO, "Connection <%s> Removed from Reactor ", c.String())
				}
				c.SetNextReadTimeout(time.Now().Add(ego._server._readTimeout))
				return true
			})
		}
	}()
}

func (ego *TcpServerSubReactor) AddConnection(c *TcpServerConnection) {
	ego._connections.Store(c.Identifier(), c)
	atomic.AddInt32(&(ego._connectionCount), 1)
}

func (ego *TcpServerSubReactor) NumConnections() int32 {
	return ego._connectionCount
}

func NeoTcpServerSubReactor(server *TcpServer) *TcpServerSubReactor {
	r := TcpServerSubReactor{
		_server: server,
	}
	return &r
}
