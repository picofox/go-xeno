package transcomm

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

func (ego *TcpServer) AddConnectionFailOver(c *TCPServerConnection) int32 {
	r := NeoTcpServerSubReactor(ego)
	r.Start()
	r.AddConnection(c)
	ego._lock.Lock()
	defer ego._lock.Unlock()
	ego._subReactors = append(ego._subReactors, r)
	return core.MkSuccess(0)
}

func (ego *TcpServer) AddConnection(c *TCPServerConnection) int32 {
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
		logging.Log(core.LL_ERR, "TCPServer: Listen Failed of <%s>", ego._bindAddress.String())
		return nil
	}
	ego._listener = lis
	return ego
}

func (ego *TcpServer) Stop() {
	ego._listener.Close()

}

func (ego *TcpServer) Start() int32 {
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
	return 0
}
func NeoTcpServer(tcpConfig *config.NetworkServerTCPConfig, logger logging.ILogger) *TcpServer {
	bindAddr := inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, tcpConfig.ListenerEndPoints[0])
	tcpServer := TcpServer{
		_bindAddress: bindAddr,
		_listener:    nil,
		_subReactors: make([]*TcpServerSubReactor, 0, 1024),
		_readTimeout: 1 * time.Millisecond,
		_config:      tcpConfig,
		_logger:      logger,
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