package server

import (
	"xeno/zohar/core/config"
	"xeno/zohar/core/logging"
)

type TcpServer struct {
	_mainReactor *MainReactor
	_config      *config.NetworkServerTCPConfig
	_logger      logging.ILogger
}

func (ego *TcpServer) Log(lv int, fmt string, arg ...any) {
	if ego._logger != nil {
		ego._logger.Log(lv, fmt, arg...)
	}
}

func (ego *TcpServer) LogFixedWidth(lv int, leftLen int, ok bool, failStr string, format string, arg ...any) {
	if ego._logger != nil {
		ego._logger.LogFixedWidth(lv, leftLen, ok, failStr, format, arg...)
	}
}

func (ego *TcpServer) Start() int32 {
	ego._mainReactor.OnStart()
	go ego._mainReactor.Loop()
	return 0
}

func (ego *TcpServer) Wait() {

}

func NeoTcpServer(tcpConfig *config.NetworkServerTCPConfig, logger logging.ILogger) *TcpServer {
	tcpServer := TcpServer{
		_mainReactor: nil,
		_config:      tcpConfig,
		_logger:      logger,
	}
	mr := NeoMainReactor(&tcpServer)
	tcpServer._mainReactor = mr
	return &tcpServer
}
