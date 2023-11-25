package server

import (
	"github.com/cloudwego/netpoll"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/nic"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
)

type TcpServer struct {
	_bindAddress inet.IPV4EndPoint
	_listener    netpoll.Listener
	_config      *config.NetworkServerTCPConfig
}

func (ego *TcpServer) Start() int32 {
	lis, err := netpoll.CreateListener(ego._bindAddress.ProtoName(), ego._bindAddress.EndPointString())
	if err != nil {
		logging.Log(core.LL_ERR, "TcpServer: Listen Failed of <%s>", ego._bindAddress.String())
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	ego._listener = lis
	return core.MkSuccess(0)
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
