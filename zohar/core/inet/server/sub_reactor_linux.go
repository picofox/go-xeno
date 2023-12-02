package server

import (
	"runtime"
	"syscall"
	"unsafe"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
)

type EpoolEventDataSubReactor struct {
	Connection *TcpServerConnection
}

func BindSubReactorEventData(ptr unsafe.Pointer, data *EpoolEventDataSubReactor) {
	*(**EpoolEventDataSubReactor)(ptr) = data

}

func ExtractSubReactorEventData(ptr unsafe.Pointer) *EpoolEventDataSubReactor {
	return *(**EpoolEventDataSubReactor)(ptr)
}

type SubReactor struct {
	pollArgs
	_server          *TcpServer
	_epollDescriptor int
}

func (ego *SubReactor) ResetEvent(size int, caps int) {
	ego._size, ego._caps = size, caps
	ego._events = make([]inet.EPollEvent, size)
}

func (ego *SubReactor) onPullIn(evt *inet.EPollEvent) {
	p := ExtractSubReactorEventData(unsafe.Pointer(&evt.Data))
	p.Connection.Read()

}
func (ego *SubReactor) onPullHup(evt *inet.EPollEvent) {
	ego._server.Log(core.LL_INFO, "Sub PullHup:")
}

func (ego *SubReactor) onPullErr(evt *inet.EPollEvent) {
	ego._server.Log(core.LL_ERR, "Sub PullErr:")
}

func (ego *SubReactor) HandlerEvent(evt *inet.EPollEvent) {
	if (evt.Events & syscall.EPOLLIN) != 0 {
		ego.onPullIn(evt)
	} else if (evt.Events & syscall.EPOLLRDHUP) != 0 {
		ego.onPullHup(evt)
	} else if (evt.Events & syscall.EPOLLERR) != 0 {
		ego.onPullErr(evt)
	}
}

func (ego *SubReactor) OnStart() {
	go ego.Loop()
}

func (ego *SubReactor) Loop() int32 {
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

func (ego *SubReactor) AddConnection(conn *TcpServerConnection) {
	ev := inet.EPollEvent{}
	ev.Events = syscall.EPOLLIN | syscall.EPOLLRDHUP | syscall.EPOLLERR | inet.EPOLLET
	info := &EpoolEventDataSubReactor{
		Connection: conn,
	}
	BindSubReactorEventData(unsafe.Pointer(&ev.Data), info)
	inet.EpollCtl(ego._epollDescriptor, syscall.EPOLL_CTL_ADD, conn._fd, &ev)
}

func NeoSubReactor(server *TcpServer) *SubReactor {
	sr := SubReactor{
		pollArgs: pollArgs{
			_caps:   128,
			_size:   128,
			_events: make([]inet.EPollEvent, 128),
		},
		_server:          server,
		_epollDescriptor: -1,
	}
	var err error
	sr._epollDescriptor, err = inet.EpollCreate(0)
	if err != nil {
		server.Log(core.LL_ERR, "Sub Reactor EpollCreate failed. err:()", err.Error())
		return nil
	}
	return &sr
}
