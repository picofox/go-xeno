package server

import (
	"xeno/zohar/core/memory"
	"xeno/zohar/core/net"
	"xeno/zohar/core/net/nic"
)

type TcpServer struct {
	_bindAddress net.IPV4EndPoint
}

func NeoTcpServer(ipstr string, port uint16) *TcpServer {
	tcpServer := TcpServer{
		_bindAddress: net.NeoIPV4EndPointByStrIP(net.EP_PROTO_TCP, 0, 0, ipstr, port),
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
