package transcomm

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/inet"
)

type EPoolEventDataSubReactor struct {
	Connection IConnection
}

func BindSubReactorEventData(ptr unsafe.Pointer, data *EPoolEventDataSubReactor) {
	*(**EPoolEventDataSubReactor)(ptr) = data

}

func ExtractSubReactorEventData(ptr unsafe.Pointer) *EPoolEventDataSubReactor {
	return *(**EPoolEventDataSubReactor)(ptr)
}

type SubReactor struct {
	pollArgs
	_poller          *Poller
	_epollDescriptor int
	_commandChannel  chan cms.ICMS
	_connections     sync.Map
}

func (ego *SubReactor) ResetEvent(size int) {
	ego._size = size
	ego._events = make([]inet.EPollEvent, size)
}

func (ego *SubReactor) onPullIn(evt *inet.EPollEvent) {
	p := ExtractSubReactorEventData(unsafe.Pointer(&evt.Data))
	p.Connection.OnIncomingData()

}

func (ego *SubReactor) onPullOut(evt *inet.EPollEvent) {
	p := ExtractSubReactorEventData(unsafe.Pointer(&evt.Data))
	p.Connection.OnWritable()
}
func (ego *SubReactor) onPullHup(evt *inet.EPollEvent) {
	ego._poller.Log(core.LL_INFO, "Sub PullHup:")
	p := ExtractSubReactorEventData(unsafe.Pointer(&evt.Data))
	p.Connection.OnPeerClosed()
}

func (ego *SubReactor) onPullErr(evt *inet.EPollEvent) {
	p := ExtractSubReactorEventData(unsafe.Pointer(&evt.Data))
	p.Connection.OnConnectingFailed()
	ego._poller.Log(core.LL_ERR, "Sub PullErr:")
}

func (ego *SubReactor) HandlerEvent(evt *inet.EPollEvent) {
	if (evt.Events & syscall.EPOLLIN) != 0 {
		ego.onPullIn(evt)
	} else if (evt.Events & syscall.EPOLLOUT) != 0 {
		ego.onPullOut(evt)
	} else if (evt.Events & syscall.EPOLLRDHUP) != 0 {
		ego.onPullHup(evt)
	} else if (evt.Events & syscall.EPOLLERR) != 0 {
		ego.onPullErr(evt)
	}
}

func (ego *SubReactor) OnStart() {
	ego._poller.Log(core.LL_SYS, "Sub Reactor Closing")
	ego._poller._waitGroup.Add(1)
	go ego.Loop()
}

func (ego *SubReactor) OnStop() {
	ego._poller.Log(core.LL_SYS, "Sub Reactor Stopping")
	finCMS := cms.NeoFinalize()
	ego._commandChannel <- finCMS
}

func (ego *SubReactor) Loop() int32 {
	defer ego._poller._waitGroup.Done()
	var nReady int = 0
	var err error = nil
	//var msec = 1000

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

		nReady, err = inet.EpollWait(ego._epollDescriptor, ego._events, 1000)
		if err != nil && err != syscall.EINTR {
			return core.MkErr(core.EC_EPOLL_WAIT_ERROR, 1)
		}
		if nReady < 0 {
			//msec = 1000

			runtime.Gosched()
			continue
		} else if nReady == 0 {
			ego._connections.Range(
				func(key, value any) bool {
					ego._poller.Log(core.LL_DEBUG, "pulse....")
					value.(IConnection).Pulse(chrono.GetRealTimeMilli())
					return true
				})

		}
		//msec = 0

		for i := 0; i < nReady; i++ {
			ego.HandlerEvent(&ego._events[i])
		}
	}
}

func (ego *SubReactor) RemoveConnection(conn IConnection) {
	fd := -1
	if conn.Type() == CONNTYPE_TCP_SERVER {
		fd = conn.(*TCPServerConnection)._fd
	} else if conn.Type() == CONNTYPE_TCP_CLIENT {
		fd = conn.(*TCPClientConnection)._fd
	} else {
		return
	}
	ego._connections.Delete(conn.Identifier())
	err := inet.EpollCtl(ego._epollDescriptor, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		ego._poller.Log(core.LL_ERR, "Remove conn %s %d from Poller Failed.", conn.String(), fd)
	}
}

func (ego *SubReactor) AddConnection(conn IConnection) {
	ev := inet.EPollEvent{}
	ev.Events = syscall.EPOLLIN | syscall.EPOLLRDHUP | syscall.EPOLLERR | syscall.EPOLLOUT | inet.EPOLLET
	info := &EPoolEventDataSubReactor{
		Connection: conn,
	}
	BindSubReactorEventData(unsafe.Pointer(&ev.Data), info)
	fd := -1
	if conn.Type() == CONNTYPE_TCP_SERVER {
		fd = conn.(*TCPServerConnection)._fd
	} else if conn.Type() == CONNTYPE_TCP_CLIENT {
		fd = conn.(*TCPClientConnection)._fd
	} else {
		return
	}
	ego._connections.Store(conn.Identifier(), conn)
	inet.EpollCtl(ego._epollDescriptor, syscall.EPOLL_CTL_ADD, fd, &ev)
}
