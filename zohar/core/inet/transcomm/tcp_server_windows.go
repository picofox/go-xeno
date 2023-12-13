package transcomm

import (
	"fmt"
	"net"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/nic"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
)

type TCPServer struct {
	_name          string
	_poller        *Poller
	_config        *config.NetworkServerTCPConfig
	_logger        logging.ILogger
	_listeners     sync.Map
	_connectionMap sync.Map
}

func (ego *TCPServer) OnPeerClosed(connection *TCPServerConnection) int32 {
	ego.Log(core.LL_SYS, "Connection Peer <%s> Closed.", connection.String())
	ego._connectionMap.Delete(connection.Identifier())
	ego._poller.SubReactorEnded(connection.Identifier())
	connection.Close()
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnDisconnected(connection *TCPServerConnection) int32 {
	ego.Log(core.LL_SYS, "Connection Peer <%s> Disconnected.", connection.String())
	ego._connectionMap.Delete(connection.Identifier())
	ego._poller.SubReactorEnded(connection.Identifier())
	connection.Close()
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnIOError(connection *TCPServerConnection) int32 {
	ego.Log(core.LL_SYS, "Connection IO <%s> Error.", connection.String())
	ego._connectionMap.Delete(connection.Identifier())
	ego._poller.SubReactorEnded(connection.Identifier())
	connection.Close()
	return core.MkSuccess(0)
}

func (ego *TCPServer) Name() string {
	return ego._name
}

func (ego *TCPServer) Listeners() *sync.Map {
	return &ego._listeners
}

func (ego *TCPServer) OnIncomingMessage(conn IConnection, msg message_buffer.INetMessage, param any) {
	fmt.Printf("Got msg [%v] \n", msg.String())
	conn.(*TCPServerConnection).SendMessage(msg, true)
}

func (ego *TCPServer) Initialize() int32 {
	for _, eps := range ego._config.ListenerEndPoints {
		bAddr := inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, eps)
		if !bAddr.Valid() {
			ego._poller.Log(core.LL_ERR, "Convert IP&Port string %s to endpoint failed.", eps)
		}

		if bAddr.IPV4() != 0 {
			nic.GetNICManager().Update()
			InetAddress := nic.GetNICManager().FindNICByIpV4Address(bAddr.IPV4())
			if InetAddress == nil {
				ego._poller.Log(core.LL_ERR, "NeoTcpServer FindNICByIpV4Address <%s> Failed", bAddr.EndPointString())
			}
			nm := InetAddress.NetMask()
			m := memory.BytesToUInt32BE(&nm, 0)
			nb := memory.NumberOfOneInInt32(int32(m))
			bAddr.SetMask(nb)
		}

		lis := NeoListenWrapper(ego, bAddr)
		ego._listeners.Store(lis._bindAddress.Identifier(), lis)

	}
	return core.MkSuccess(0)
}

func (ego *TCPServer) Log(lv int, fmt string, arg ...any) {
	if ego._logger != nil {
		ego._logger.Log(lv, fmt, arg...)
	}
}

func (ego *TCPServer) OnIncomingConnection(connection IConnection) (IConnection, int32) {
	ego.Log(core.LL_INFO, "Incoming Connection [%d] @ <%s -> %s> Added", connection.Identifier(), connection.RemoteEndPoint().EndPointString(), connection.LocalEndPoint().EndPointString())
	ego._connectionMap.Store(connection.Identifier(), connection)
	return connection, core.MkSuccess(0)
}

func (ego *TCPServer) LogFixedWidth(lv int, leftLen int, ok bool, failStr string, format string, arg ...any) {
	if ego._logger != nil {
		ego._logger.LogFixedWidth(lv, leftLen, ok, failStr, format, arg...)
	}
}

func (ego *TCPServer) Start() int32 {
	rc := int32(0)
	ego._listeners.Range(func(key, value any) bool {
		lis := value.(*ListenWrapper)
		l, err := net.Listen(lis._bindAddress.ProtoName(), lis._bindAddress.EndPointString())
		if err != nil {
			ego.Log(core.LL_ERR, "Listen on <%s> Failed:(%s)", lis._bindAddress.String(), err.Error())
			rc = core.MkErr(core.EC_LISTEN_ERROR, 1)
			return false
		}
		lis._listen = l
		return true
	})

	if core.Err(rc) {
		return core.MkErr(core.EC_LISTEN_ERROR, 1)
	}

	ego._poller.OnServerStart(ego)
	return core.MkSuccess(0)
}

func (ego *TCPServer) Stop() int32 {
	return 0
}

func (ego *TCPServer) Wait() {

}

func NeoTcpServer(name string, tcpConfig *config.NetworkServerTCPConfig, logger logging.ILogger) *TCPServer {
	tcpServer := TCPServer{
		_name:   name,
		_poller: GetDefaultPoller(),
		_config: tcpConfig,
		_logger: logger,
	}

	return &tcpServer
}
