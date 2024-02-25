package transcomm

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/config"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/inet/nic"
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/mp"
)

type TCPServer struct {
	_name          string
	_poller        *Poller
	_config        *config.NetworkServerTCPConfig
	_logger        logging.ILogger
	_listeners     sync.Map
	_connectionMap sync.Map
	_router        IServerMessageRouter
}

func (ego *TCPServer) Name() string {
	return ego._name
}

func (ego *TCPServer) Listeners() *sync.Map {
	return &ego._listeners
}

func (ego *TCPServer) OnIncomingMessage(conn *TCPServerConnection, msg message_buffer.INetMessage) int32 {
	return ego._router.OnIncomingMessage(conn, msg)
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

func (ego *TCPServer) LogFixedWidth(lv int, leftLen int, ok bool, failStr string, format string, arg ...any) {
	if ego._logger != nil {
		ego._logger.LogFixedWidth(lv, leftLen, ok, failStr, format, arg...)
	}
}

func (ego *TCPServer) ConnectedConnectionCount() int {
	var count int = 0
	ego._connectionMap.Range(func(key, value any) bool {
		c := value.(*TCPServerConnection)
		if c._stateCode == Connected {
			count++
		}

		return true
	})
	return count
}

func (ego *TCPServer) BroadCastMessage(message message_buffer.INetMessage, bFlush bool) int32 {
	ego._connectionMap.Range(func(key, value any) bool {
		c := value.(*TCPServerConnection)
		rc := c.SendMessage(message, bFlush)
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {

			} else {
				panic(fmt.Sprintf("%d", core.ErrStr(rc)))
			}
		}
		return true
	})
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnPeerClosed(connection *TCPServerConnection) int32 {
	ego.Log(core.LL_SYS, "Connection Peer <%s> Closed.", connection.String())
	ego._connectionMap.Delete(connection.Identifier())
	ego._poller.OnConnectionRemove(connection)
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnDisconnected(connection *TCPServerConnection) int32 {
	ego.Log(core.LL_SYS, "Connection Peer <%s> Disconnected.", connection.String())
	ego._connectionMap.Delete(connection.Identifier())
	ego._poller.OnConnectionRemove(connection)
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnIOError(connection *TCPServerConnection) int32 {
	ego.Log(core.LL_SYS, "Connection IO <%s> Error.", connection.String())
	ego._connectionMap.Delete(connection.Identifier())
	ego._poller.OnConnectionRemove(connection)
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnKeepAliveMessage(conn *TCPServerConnection, message message_buffer.INetMessage) int32 {
	var pkam *messages.KeepAliveMessage = message.(*messages.KeepAliveMessage)
	if pkam.IsServer() {
		ts := chrono.GetRealTimeMilli()
		delta := ts - pkam.TimeStamp()
		conn.OnKeepAlive(ts, int32(delta))
		ego.Log(core.LL_DEBUG, "Got KA back")
	} else {
		conn.SendMessage(message, true)
	}
	return core.MkSuccess(0)
}

var s_sproctestCount int = 0

func (ego *TCPServer) OnProcTestMessage(conn *TCPServerConnection, message message_buffer.INetMessage) int32 {
	m := message.(*messages.ProcTestMessage)
	if m.IsServer {
		if core.Err(m.Validate()) {
			panic("invalid msg")
		}
		s_sproctestCount++
		if s_sproctestCount%1000 == 0 {
			ego.Log(core.LL_DEBUG, "Got Pro Message %d", s_sproctestCount)
		}
	} else {
		//ego.Log(core.LL_DEBUG, "echo proc test mesg")
		rc := conn.SendMessage(m, true)
		if core.Err(rc) {
			if !core.IsErrType(rc, core.EC_TRY_AGAIN) {

				panic(fmt.Sprintf("echo proc test msg failed. %s", core.ErrStr(rc)))

			}
		}
	}
	return core.MkSuccess(0)
}

func (ego *TCPServer) OnIncomingConnection(listener *ListenWrapper, fd int, rAddr inet.IPV4EndPoint) (IConnection, int32) {
	lsa, _ := syscall.Getsockname(fd)
	lAddr := inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, lsa)
	connection := ego.NeoTCPServerConnection(fd, rAddr, lAddr)
	ego._connectionMap.Store(connection.Identifier(), connection)
	return connection, core.MkSuccess(0)
}

func (ego *TCPServer) NeoTCPServerConnection(fd int, rAddr inet.IPV4EndPoint, lAddr inet.IPV4EndPoint) *TCPServerConnection {
	connection := TCPServerConnection{
		_fd:                       fd,
		_localEndPoint:            lAddr,
		_remoteEndPoint:           rAddr,
		_recvBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_sendBuffer:               memory.NeoLinkedListByteBuffer(datatype.SIZE_4K),
		_server:                   ego,
		_profiler:                 prof.NeoConnectionProfiler(),
		_outgoingHeader:           memory.NeoO1L31C16Header(0, 0),
		_incomingHeaderBufferRIdx: 0,
		_incomingDataIndex:        0,
		_keepalive:                nil,
		_stateCode:                Initialized,
	}

	inet.SysSetTCPNoDelay(fd, true)

	if connection.KeepAliveConfig().Enable {
		connection._keepalive = NeoKeepAlive(connection.KeepAliveConfig(), true)
	}

	return &connection
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
		file, err := l.(*net.TCPListener).File()
		if err != nil {
			ego.Log(core.LL_ERR, "File From Listener %v Failed %s", l.Addr(), err.Error())
			rc = core.MkErr(core.EC_GET_LOW_FD_ERROR, 1)
			return false
		}
		fd := int(file.Fd())
		err = syscall.SetNonblock(fd, true)
		if err != nil {
			ego.Log(core.LL_ERR, "Set File Descriptor %d NB mode Failed %s", fd, err.Error())
			rc = core.MkErr(core.EC_SET_NONBLOCK_ERROR, 1)
			return false
		}
		lis._listen = l
		lis._file = file
		lis._fd = fd
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
		_router: nil,
	}

	var output []reflect.Value = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Reflect"+tcpConfig.Codec, &tcpServer)
	if core.Err(rc) {
		panic(fmt.Sprintf("Install Handler Failed %s", tcpConfig.Codec))
	}
	h := output[0].Interface().(IServerMessageRouter)
	tcpServer._router = h
	tcpServer._router.RegisterHandler(messages.INTERNAL_MSG_GRP_TYPE, messages.KEEP_ALIVE_MESSAGE_ID, tcpServer.OnKeepAliveMessage)
	tcpServer._router.RegisterHandler(messages.INTERNAL_MSG_GRP_TYPE, messages.PROC_TEST_MESSAGE_ID, tcpServer.OnProcTestMessage)

	return &tcpServer
}
