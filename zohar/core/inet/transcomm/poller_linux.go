package transcomm

import (
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/logging"
)

type Poller struct {
	_logger          logging.ILogger
	_mainReactor     *MainReactor
	_subReactors     []*SubReactor
	_config          *intrinsic.PollerConfig
	_subReactorIndex atomic.Uint32
	_servers         sync.Map
	_stateCode       datatype.StateCode
	_waitGroup       sync.WaitGroup
}

func (ego *Poller) OnIncomingConnection(connection IConnection) {
	idx := ego._subReactorIndex.Add(1) % uint32(len(ego._subReactors))
	connection.SetReactorIndex(idx)
	ego._subReactors[idx].AddConnection(connection)
}

func (ego *Poller) OnConnectionRemove(connection IConnection) {
	ego._subReactors[connection.ReactorIndex()].RemoveConnection(connection)
}

func (ego *Poller) RegisterTCPServer(svr *TCPServer) {
	ego._servers.Store(svr.Name(), svr)
}

func (ego *Poller) OnServerStart(svr *TCPServer) int32 {
	s, ok := ego._servers.Load(svr.Name())
	if !ok {
		return core.EC_ELEMENT_NOT_FOUND
	}

	lisMap := s.(*TCPServer).Listeners()
	lisMap.Range(func(k2, lis any) bool {
		ev := inet.EPollEvent{}
		ev.Events = syscall.EPOLLIN | syscall.EPOLLRDHUP | syscall.EPOLLERR | inet.EPOLLET
		info := &EPoolEventDataMainReactor{
			FD:       lis.(*ListenWrapper).FileDescriptor(),
			Listener: lis.(*ListenWrapper),
		}
		BindMainReactorEventData(unsafe.Pointer(&ev.Data), info)
		err := inet.EpollCtl(ego._mainReactor._epollDescriptor, syscall.EPOLL_CTL_ADD, info.Listener.FileDescriptor(), &ev)
		if err != nil {
			ego._logger.Log(core.LL_ERR, "Add socket fd %d to Main reactor Failed", info.Listener.FileDescriptor())
		}
		ego._logger.Log(core.LL_SYS, "Add socket fd %d to Main reactor", info.Listener.FileDescriptor())
		return true
	},
	)
	return core.MkSuccess(0)
}

func (ego *Poller) Start() int32 {
	rc := ego._stateCode.SetStartState()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_NOOP) {
			return core.MkSuccess(0)
		}
		return rc
	}
	ego._mainReactor.OnStart()
	for _, ra := range ego._subReactors {
		ra.OnStart()
	}

	ego._stateCode.SetStartStateResult(true)
	return core.MkSuccess(0)
}

func (ego *Poller) Wait() int32 {
	ego._waitGroup.Wait()
	ego._stateCode.SetStopStateResult(true)
	return core.MkSuccess(0)
}

func (ego *Poller) Stop() int32 {
	ego._stateCode.SetStopState()
	ego._mainReactor.OnStop()
	for _, subR := range ego._subReactors {
		subR.OnStop()
	}

	return core.MkSuccess(0)
}

func (ego *Poller) NeoSubReactor() *SubReactor {
	sr := SubReactor{
		pollArgs: pollArgs{
			_size:   128,
			_events: make([]inet.EPollEvent, 128),
		},
		_poller:          ego,
		_epollDescriptor: -1,
		_commandChannel:  make(chan cms.ICMS, 1),
	}
	var err error
	sr._epollDescriptor, err = inet.EpollCreate(0)
	if err != nil {
		ego.Log(core.LL_ERR, "Sub Reactor EpollCreate failed. err:()", err.Error())
		return nil
	}
	return &sr
}

func (ego *Poller) neoMainReactor() *MainReactor {
	mr := MainReactor{
		pollArgs: pollArgs{
			_size:   128,
			_events: make([]inet.EPollEvent, 128),
		},
		_epollDescriptor: -1,
		_poller:          ego,
		_commandChannel:  make(chan cms.ICMS, 1),
	}

	var err error
	mr._epollDescriptor, err = inet.EpollCreate(0)
	if err != nil {
		ego._logger.Log(core.LL_ERR, "EpollCreate failed. err:()", err.Error())
		return nil
	}

	return &mr
}

func (ego *Poller) SubReactorCount() int32 {
	return int32(len(ego._subReactors))
}

func (ego *Poller) Initialize() int32 {
	ego._stateCode.SetInitializeState()
	ego._mainReactor = ego.neoMainReactor()

	var subCnt = int32(runtime.GOMAXPROCS(0)/20 + 1)
	if ego._config != nil && ego._config.SubReactorCount > 0 {
		subCnt = ego._config.SubReactorCount
	}

	for i := int32(0); i < subCnt; i++ {
		sr := ego.NeoSubReactor()
		if sr == nil {
			return core.MkErr(core.EC_NULL_VALUE, 1)
		}
		ego._subReactors = append(ego._subReactors, sr)
	}

	ego._servers.Range(
		func(key, svr any) bool {
			lisMap := svr.(*TCPServer).Listeners()
			lisMap.Range(
				func(k2, lis any) bool {
					ev := inet.EPollEvent{}
					ev.Events = syscall.EPOLLIN | syscall.EPOLLRDHUP | syscall.EPOLLERR | inet.EPOLLET
					info := &EPoolEventDataMainReactor{
						FD:       lis.(*ListenWrapper).FileDescriptor(),
						Listener: lis.(*ListenWrapper),
					}
					BindMainReactorEventData(unsafe.Pointer(&ev.Data), info)
					err := inet.EpollCtl(ego._mainReactor._epollDescriptor, syscall.EPOLL_CTL_ADD, info.Listener.FileDescriptor(), &ev)
					if err != nil {
						ego._logger.Log(core.LL_ERR, "Add socket fd %d to Main reactor Failed", info.Listener.FileDescriptor())
					}
					ego._logger.Log(core.LL_SYS, "Add socket fd %d to Main reactor", info.Listener.FileDescriptor())
					return true
				},
			)
			return true
		},
	)

	ego._stateCode.SetInitializeStateResult(true)
	return core.MkSuccess(0)
}

func (ego *Poller) Log(lv int, fmt string, arg ...any) {
	if ego._logger != nil {
		ego._logger.Log(lv, fmt, arg...)
	}
}

func (ego *Poller) LogFixedWidth(lv int, leftLen int, ok bool, failStr string, format string, arg ...any) {
	if ego._logger != nil {
		ego._logger.LogFixedWidth(lv, leftLen, ok, failStr, format, arg...)
	}
}

func NeoPoller() *Poller {
	p := Poller{
		_logger:      logging.GetLoggerManager().GetDefaultLogger(),
		_mainReactor: nil,
		_subReactors: make([]*SubReactor, 0),
		_config:      &intrinsic.GetIntrinsicConfig().Poller,
		_stateCode:   datatype.StateCode(0),
	}
	p._subReactorIndex.Store(0)

	return &p
}

var sDefaultPollerInstance *Poller
var sDefaultPollerInstanceOnce sync.Once

func GetDefaultPoller() *Poller {
	sDefaultPollerInstanceOnce.Do(
		func() {
			sDefaultPollerInstance = NeoPoller()
		})
	return sDefaultPollerInstance
}
