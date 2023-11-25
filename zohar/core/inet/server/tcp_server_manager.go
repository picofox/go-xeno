package server

import (
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
)

type TcpServerManager struct {
	_servers map[string]*TcpServer
}

func (ego *TcpServerManager) Initialize(cfg *config.NetworkServerConfig) int32 {
	for cName, perCFG := range cfg.TCP {
		svr := NeoTcpServer(&perCFG)
		ego._servers[cName] = svr
	}
	return core.MkSuccess(0)
}

func (ego *TcpServerManager) Start() int32 {
	for _, v := range ego._servers {
		v.Start()
	}
	return core.MkSuccess(0)
}

var sTcpServerManagerInstance *TcpServerManager
var sTcpServerManagerOnce sync.Once

func GetDefaultTcpServerManager() *TcpServerManager {
	sTcpServerManagerOnce.Do(func() {
		sTcpServerManagerInstance = &TcpServerManager{
			_servers: make(map[string]*TcpServer),
		}
	})

	return sTcpServerManagerInstance
}
