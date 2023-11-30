package server

import (
	"runtime"
	"syscall"
	"unsafe"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/nic"
	"xeno/zohar/core/memory"
)

type pollArgs struct {
	_size   int
	_caps   int
	_events []inet.EPollEvent
}

func BindMainReactorEventData(ptr unsafe.Pointer, data *EpoolEventDataMainReactor) {
	*(**EpoolEventDataMainReactor)(ptr) = data

}

type EpoolEventDataMainReactor struct {
	FD int
}

type MainReactor struct {
	pollArgs
	_listener        []*ListenWrapper
	_server          *TcpServer
	_epollDescriptor int
}

func (ego *MainReactor) ResetEvent(size int, caps int) {
	ego._size, ego._caps = size, caps
	ego._events = make([]inet.EPollEvent, size)
}

func (ego *MainReactor) onPullIn(evt *inet.EPollEvent) {
	ego._server.Log(core.LL_DEBUG, "PullIn: fd:%d")
}
func (ego *MainReactor) onPullHup(evt *inet.EPollEvent) {
	ego._server.Log(core.LL_INFO, "PullHup:")
}

func (ego *MainReactor) onPullErr(evt *inet.EPollEvent) {
	ego._server.Log(core.LL_ERR, "PullErr:")
}

func (ego *MainReactor) HandlerEvent(evt *inet.EPollEvent) {
	if (evt.Events & syscall.EPOLLIN) != 0 {
		ego.onPullIn(evt)
	} else if (evt.Events & syscall.EPOLLRDHUP) != 0 {
		ego.onPullHup(evt)
	} else if (evt.Events & syscall.EPOLLERR) != 0 {
		ego.onPullErr(evt)
	}
}

func (ego *MainReactor) Loop() int32 {
	var nReady int = 0
	var err error = nil
	var msec = -1
	for {
		if nReady == ego._size && ego._size < 128*1024 {
			ego.ResetEvent(ego._size<<1, ego._caps)
		}
		nReady, err = inet.EpollWait(ego._epollDescriptor, ego._events, msec)
		if err != nil && err != syscall.EINTR {
			return core.MkErr(core.EC_EPOLL_WAIT_ERROR, 1)
		}
		if nReady < 0 {
			msec = -1
			runtime.Gosched()
			continue
		}
		msec = 0

		for i := 0; i < nReady; i++ {
			ego.HandlerEvent(&ego._events[i])
		}
	}
}

func (ego *MainReactor) OnStart() {
	for _, eps := range ego._server._config.ListenerEndPoints {
		bAddr := inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, eps)
		if !bAddr.Valid() {
			ego._server.Log(core.LL_ERR, "Convert ipport string %s to endpoint failed.", eps)
		}

		if bAddr.IPV4() != 0 {
			nic.GetNICManager().Update()
			InetAddress := nic.GetNICManager().FindNICByIpV4Address(bAddr.IPV4())
			if InetAddress == nil {
				ego._server.Log(core.LL_ERR, "NeoTcpServer FindNICByIpV4Address <%s> Failed", bAddr.EndPointString())
			}
			nm := InetAddress.NetMask()
			m := memory.BytesToUInt32BE(&nm, 0)
			nb := memory.NumberOfOneInInt32(int32(m))
			bAddr.SetMask(nb)
		}

		lis := NeoListenWrapper(ego._server, bAddr)
		ego._listener = append(ego._listener, lis)

		ev := inet.EPollEvent{}
		ev.Events = syscall.EPOLLIN | syscall.EPOLLRDHUP | syscall.EPOLLERR
		info := &EpoolEventDataMainReactor{
			FD: lis._fd,
		}
		BindMainReactorEventData(unsafe.Pointer(&ev.Data), info)

		inet.EpollCtl(ego._epollDescriptor, syscall.EPOLL_CTL_ADD, lis._fd, &ev)
		ego._server.Log(core.LL_SYS, "Add socket fd %d to main reactor", lis._fd)
	}

}

func NeoMainReactor(server *TcpServer) *MainReactor {
	ll := len(server._config.ListenerEndPoints)
	if ll < 1 {
		server.Log(core.LL_ERR, "No endpoint info is configured")
		return nil
	}

	mr := MainReactor{
		pollArgs: pollArgs{
			_caps:   128,
			_size:   128,
			_events: make([]inet.EPollEvent, 128),
		},
		_listener: make([]*ListenWrapper, 0),
		_server:   server,
	}

	var err error
	mr._epollDescriptor, err = inet.EpollCreate(0)
	if err != nil {
		server.Log(core.LL_ERR, "EpollCreate failed. err:()", err.Error())
		return nil
	}

	return &mr
}
