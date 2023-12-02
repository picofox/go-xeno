package server

import (
	"sync"
	"sync/atomic"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/logging"
)

type TcpServer struct {
	_mainReactor     *MainReactor
	_subReactors     []*SubReactor
	_config          *config.NetworkServerTCPConfig
	_logger          logging.ILogger
	_connectionMap   sync.Map
	_subReactorIndex atomic.Uint32
}

func (ego *TcpServer) AddConnection(conn *TcpServerConnection) {
	ego.Log(core.LL_INFO, "Incoming Connection [%d] @ <%s -> %s> Added", conn._fd, conn._remoteEndPoint.EndPointString(), conn._localEndPoint.EndPointString())
	ego._connectionMap.Store(conn._fd, conn)
}

func (ego *TcpServer) DispatchConnection(conn *TcpServerConnection) {
	idx := ego._subReactorIndex.Add(1)
	if idx > uint32(len(ego._subReactors)) {
		idx = 0
	}
	ego._subReactors[idx].AddConnection(conn)
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

	return 0
}

func (ego *TcpServer) Wait() {

}

func NeoTcpServer(tcpConfig *config.NetworkServerTCPConfig, logger logging.ILogger) *TcpServer {
	tcpServer := TcpServer{
		_mainReactor: nil,
		_subReactors: make([]*SubReactor, 0),
		_config:      tcpConfig,
		_logger:      logger,
	}
	mr := NeoMainReactor(&tcpServer)

	for i := 0; i < 2; i++ {
		sr := NeoSubReactor(&tcpServer)
		if sr == nil {
			return nil
		}
		tcpServer._subReactors = append(tcpServer._subReactors, sr)
	}

	tcpServer._mainReactor = mr

	return &tcpServer
}
