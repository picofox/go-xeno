package transcomm

import (
	"runtime"
	"syscall"
	"unsafe"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/inet"
)

type pollArgs struct {
	_size   int
	_events []inet.EPollEvent
}

func BindMainReactorEventData(ptr unsafe.Pointer, data *EPoolEventDataMainReactor) {
	*(**EPoolEventDataMainReactor)(ptr) = data

}

func ExtractMainReactorEventData(ptr unsafe.Pointer) *EPoolEventDataMainReactor {
	return *(**EPoolEventDataMainReactor)(ptr)
}

type EPoolEventDataMainReactor struct {
	FD       int
	Listener *ListenWrapper
}

type MainReactor struct {
	pollArgs
	_epollDescriptor int
	_poller          *Poller
	_commandChannel  chan cms.ICMS
}

func (ego *MainReactor) Accept(listenFd int) (int, syscall.Sockaddr, int32) {
	var fd, sa, err = syscall.Accept(listenFd)
	if err != nil {
		if err == syscall.EAGAIN {
			return -1, nil, core.MkErr(core.EC_TRY_AGAIN, 1)
		}
		return -1, nil, core.MkErr(core.EC_ACCEPT_ERROR, 1)
		ego._poller.Log(core.LL_ERR, "Accept listenFD (%d) error", listenFd)
	}

	syscall.SetNonblock(fd, true)

	return fd, sa, core.MkSuccess(0)
}

func (ego *MainReactor) ResetEvent(size int) {
	ego._size = size
	ego._events = make([]inet.EPollEvent, size)
}

func (ego *MainReactor) onPullIn(evt *inet.EPollEvent) {
	p := ExtractMainReactorEventData(unsafe.Pointer(&evt.Data))
	fd, sa, rc := ego.Accept(p.FD)
	if core.Err(rc) {
		return
	}
	raddr := inet.NeoIPV4EndPointBySockAddr(inet.EP_PROTO_TCP, 0, 0, sa)
	svr := p.Listener.Server()
	var connection IConnection

	connection, rc = svr.OnIncomingConnection(p.Listener, fd, raddr)
	et, _ := core.ExErr(rc)
	if et != core.EC_OK {
		if et == core.EC_NOOP {
			ego._poller.Log(core.LL_SYS, "Connection <%s> is not welcome", raddr.String())
		} else {
			ego._poller.Log(core.LL_SYS, "Make Neo Connection <%s> Failed. %d", raddr.String(), et)
		}
	}
	ego._poller.OnIncomingConnection(connection)

}
func (ego *MainReactor) onPullHup(evt *inet.EPollEvent) {
	p := ExtractMainReactorEventData(unsafe.Pointer(&evt.Data))
	ba := p.Listener.BindAddr()
	ego._poller.Log(core.LL_INFO, "Main React PullHup: <%s>", ba.EndPointString())
}

func (ego *MainReactor) onPullErr(evt *inet.EPollEvent) {
	p := ExtractMainReactorEventData(unsafe.Pointer(&evt.Data))
	ba := p.Listener.BindAddr()
	ego._poller.Log(core.LL_INFO, "Main React Error: <%s>", ba.EndPointString())
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
	defer ego._poller._waitGroup.Done()
	var nReady int = 0
	var err error = nil
	var msec = 1000
	for {
		select {
		case m := <-ego._commandChannel:
			if m.Id() == cms.CMSID_FINALIZE {
				runtime.Goexit()
			}
		default:
		}

		if nReady == ego._size && ego._size < 128*1024 {
			ego.ResetEvent(ego._size << 1)
		}

		nReady, err = inet.EpollWait(ego._epollDescriptor, ego._events, msec)
		if err != nil && err != syscall.EINTR {
			return core.MkErr(core.EC_EPOLL_WAIT_ERROR, 1)
		}
		if nReady < 0 {
			msec = 1000
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
	ego._poller.Log(core.LL_SYS, "Main Reactor Starting")
	ego._poller._waitGroup.Add(1)
	go ego.Loop()
}

func (ego *MainReactor) OnStop() {
	ego._poller.Log(core.LL_SYS, "Main Reactor Stopping")
	finCMS := cms.NeoFinalize()
	ego._commandChannel <- finCMS
}
