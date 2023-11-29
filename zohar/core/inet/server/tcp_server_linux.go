package server

import (
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/nic"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
)

type TcpServer struct {
	_bindAddress inet.IPV4EndPoint
	_listener    *ListenWrapper
	_config      *config.NetworkServerTCPConfig
	_logger      logging.ILogger
	_eventLoop   *EventLoop
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
	return 0
}

func (ego *TcpServer) Wait() {

}

func NeoTcpServer(tcpConfig *config.NetworkServerTCPConfig, logger logging.ILogger) *TcpServer {
	bindAddr := tcpConfig.BindAddr
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}
	tcpServer := TcpServer{
		_bindAddress: inet.NeoIPV4EndPointByStrIP(inet.EP_PROTO_TCP, 0, 0, bindAddr, tcpConfig.Port),
		_listener:    nil,
		_config:      tcpConfig,
		_logger:      logger,
	}

	if tcpServer._bindAddress.IPV4() != 0 {
		nic.GetNICManager().Update()
		InetAddress := nic.GetNICManager().FindNICByIpV4Address(tcpServer._bindAddress.IPV4())
		if InetAddress == nil {
			tcpServer.Log(core.LL_ERR, "NeoTcpServer FindNICByIpV4Address <%s> Failed", tcpServer._bindAddress.EndPointString())
			return nil
		}
		nm := InetAddress.NetMask()
		m := memory.BytesToUInt32BE(&nm, 0)
		nb := memory.NumberOfOneInInt32(int32(m))
		tcpServer._bindAddress.SetMask(nb)
	}

	tcpServer._listener = NeoListenWrapper(&tcpServer)
	if tcpServer._listener == nil {
		tcpServer.Log(core.LL_ERR, "NeoTcpServer NeoListenWrapper <%s> Failed", tcpServer._bindAddress.EndPointString())
		return nil
	}

	tcpServer.LogFixedWidth(core.LL_SYS, 70, true, "", "NeoTcpServer <%s>", tcpServer._bindAddress.EndPointString())

	return &tcpServer
}
