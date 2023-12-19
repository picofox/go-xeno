package transcomm

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
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
}

func (ego *TCPServer) Name() string {
	return ego._name
}

func (ego *TCPServer) Listeners() *sync.Map {
	return &ego._listeners
}

func (ego *TCPServer) OnIncomingMessage(conn IConnection, msg any, param any) {

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

func (ego *TCPServer) OnIncomingConnection(listener *ListenWrapper, fd int, rAddr inet.IPV4EndPoint) (IConnection, int32) {
	lsa, _ := syscall.Getsockname(fd)
	lAddr := inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, lsa)
	connection := ego.NeoTCPServerConnection(fd, rAddr, lAddr)
	var output = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+ego._config.Codec)
	if core.Err(rc) {
		panic(fmt.Sprintf("Install Handler Failed %s", ego._config.Codec))

	}
	h := output[0].Interface().(IServerCodecHandler)
	connection._codec = h
	ego.Log(core.LL_INFO, "Incoming Connection [%d] @ <%s -> %s> Added", connection._fd, connection._remoteEndPoint.EndPointString(), connection._localEndPoint.EndPointString())
	ego._connectionMap.Store(connection.Identifier(), connection)
	return connection, core.MkSuccess(0)
}

func (ego *TCPServer) NeoTCPServerConnection(fd int, rAddr inet.IPV4EndPoint, lAddr inet.IPV4EndPoint) *TCPServerConnection {
	connection := TCPServerConnection{
		_fd:             fd,
		_localEndPoint:  lAddr,
		_remoteEndPoint: rAddr,
		_recvBuffer:     memory.NeoRingBuffer(1024),
		_sendBuffer:     memory.NeoLinearBuffer(1024),
		_server:         ego,
		_codec:          nil,
		_profiler:       prof.NeoConnectionProfiler(),
	}

	inet.SysSetTCPNoDelay(fd, true)

	var output []reflect.Value = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Neo"+ego._config.Codec, &connection)
	if core.Err(rc) {
		panic(fmt.Sprintf("Install Handler Failed %s", ego._config.Codec))
	}
	h := output[0].Interface().(IServerCodecHandler)
	connection._codec = h

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
	}

	return &tcpServer
}
