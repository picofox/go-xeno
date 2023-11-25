package server

import (
	"github.com/cloudwego/netpoll"
	"xeno/zohar/core/inet"
)

type TcpServer struct {
	_bindAddress inet.IPV4EndPoint
	_listener    netpoll.Listener
	_config      *config.NetworkServerTCPConfig
}

func (ego *TcpServer) Start() {
	lis, err := netpoll.CreateListener(ego._bindAddress.ProtoName(), ego._bindAddress.EndPointString())
	if err != nil {
		logging.Log(core.LL_ERR, "TcpServer: Listen Failed of <%s>", ego._bindAddress.String())
		return nil
	}
	ego._listener = lis
}

func NeoTcpServer(tcpConfig *config.NetworkServerTCPConfig) *TcpServer {
	bindAddr := tcpConfig.BindAddr
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}
	tcpServer := TcpServer{
		_bindAddress: inet.NeoIPV4EndPointByStrIP(inet.EP_PROTO_TCP, 0, 0, bindAddr, tcpConfig.Port),
		_listener:    nil,
		_config:      tcpConfig,
	}

	if tcpServer._bindAddress.IPV4() != 0 {
		nic.GetNICManager().Update()
		InetAddress := nic.GetNICManager().FindNICByIpV4Address(tcpServer._bindAddress.IPV4())
		if InetAddress == nil {
			return nil
		}
		nm := InetAddress.NetMask()
		m := memory.BytesToUInt32BE(&nm, 0)
		nb := memory.NumberOfOneInInt32(int32(m))
		tcpServer._bindAddress.SetMask(nb)
	}

	return &tcpServer
}
