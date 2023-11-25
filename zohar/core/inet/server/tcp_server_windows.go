package server

import (
	"net"
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/nic"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
)

const (
	INITIAL_REACTORS   = 2
	NCONNS_PER_REACTOR = 32
	MAX_REACTORS       = 1024
)

type TcpServer struct {
	_bindAddress inet.IPV4EndPoint
	_listener    net.Listener
	_subReactors []*TcpServerSubReactor
	_lock        sync.RWMutex
	_readTimeout time.Duration
	_config      *config.NetworkServerTCPConfig
}

func (ego *TcpServer) AddConnectionFailOver(c *TcpServerConnection) int32 {
	r := NeoTcpServerSubReactor(ego)
	r.Start()
	r.AddConnection(c)
	ego._lock.Lock()
	defer ego._lock.Unlock()
	ego._subReactors = append(ego._subReactors, r)
	return core.MkSuccess(0)
}

func (ego *TcpServer) AddConnection(c *TcpServerConnection) int32 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	ll := len(ego._subReactors)
	for i := 0; i < ll; i++ {
		nConns := ego._subReactors[i].NumConnections()
		if nConns < NCONNS_PER_REACTOR {
			ego._subReactors[i].AddConnection(c)
			return core.MkSuccess(0)
		}
	}

	if len(ego._subReactors) >= MAX_REACTORS {
		return core.MkErr(core.EC_TRY_AGAIN, 1)
	}

	return core.MkErr(core.EC_REACH_LIMIT, 1)
}

func (ego *TcpServer) listen() *TcpServer {
	lis, err := net.Listen(ego._bindAddress.ProtoName(), ego._bindAddress.EndPointString())
	if err != nil {
		logging.Log(core.LL_ERR, "TcpServer: Listen Failed of <%s>", ego._bindAddress.String())
		return nil
	}
	ego._listener = lis
	return ego
}

func (ego *TcpServer) Stop() {
	ego._listener.Close()

}

func (ego *TcpServer) Start() {
	ego.listen()
	for {
		conn, err := ego._listener.Accept()
		if err != nil {
			logging.Log(core.LL_ERR, "Accept Failed: (%s)", err.Error())
			break
		}
		conn.SetReadDeadline(time.Now().Add(ego._readTimeout))
		connWrapper := NeoTcpServerConnection(conn, ego._config)
		rc := ego.AddConnection(connWrapper)
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				rc = ego.AddConnectionFailOver(connWrapper)
				if core.Err(rc) {
					logging.Log(core.LL_WARN, "Can NOT handle neo connection %s", connWrapper.String())
				}
			} else {
				logging.Log(core.LL_ERR, "Tcp Server Clients Reach Max")
			}
		}
	}
}
func NeoTcpServer(tcpConfig *config.NetworkServerTCPConfig) *TcpServer {
	bindAddr := tcpConfig.BindAddr
	if bindAddr == "" {
		bindAddr = "0.0.0.0"
	}
	tcpServer := TcpServer{
		_bindAddress: inet.NeoIPV4EndPointByStrIP(inet.EP_PROTO_TCP, 0, 0, bindAddr, tcpConfig.Port),
		_listener:    nil,
		_subReactors: make([]*TcpServerSubReactor, 0, 1024),
		_readTimeout: 1 * time.Millisecond,
		_config:      tcpConfig,
	}

	for i := 0; i < INITIAL_REACTORS; i++ {
		r := NeoTcpServerSubReactor(&tcpServer)
		r.Start()
		tcpServer._subReactors = append(tcpServer._subReactors, r)
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
