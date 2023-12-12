package transcomm

import (
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/logging"
)

type Poller struct {
	_logger       logging.ILogger
	_mainReactors []*MainReactor
	_subReactors  sync.Map
	_servers      sync.Map
	_stateCode    datatype.StateCode
	_waitGroup    sync.WaitGroup
}

func (ego *Poller) OnIncomingConnection(connection IConnection) {
	sr := ego.NeoSubReactor(connection)
	sr.OnStart()
	ego._subReactors.Store(connection.Identifier(), sr)
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
		mr := ego.neoMainReactor(lis.(*ListenWrapper))
		ego._mainReactors = append(ego._mainReactors, mr)
		mr.OnStart()
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
	ego._stateCode.SetStartStateResult(true)
	return core.MkSuccess(0)
}

func (ego *Poller) Wait() int32 {
	ego._waitGroup.Wait()
	ego._stateCode.SetStopStateResult(true)
	return core.MkSuccess(0)
}

func (ego *Poller) SubReactorCount() int32 {
	var cnt int32 = 0
	ego._subReactors.Range(func(key, value any) bool {
		cnt++
		return true
	})
	return cnt
}

func (ego *Poller) Stop() int32 {
	ego._stateCode.SetStopState()
	for _, subR := range ego._mainReactors {
		subR.OnStop()
	}

	ego._subReactors.Range(func(key, value any) bool {
		value.(*SubReactor).OnStop()
		return true
	})

	ego._stateCode.SetStopStateResult(true)
	return core.MkSuccess(0)
}

func (ego *Poller) NeoSubReactor(connection IConnection) *SubReactor {
	sr := SubReactor{
		_poller:         ego,
		_connection:     connection,
		_commandChannel: make(chan cms.ICMS, 1),
	}
	return &sr
}

func (ego *Poller) SubReactorEnded(sr *SubReactor) {
	ego._subReactors.Delete(sr._connection.Identifier())
}

func (ego *Poller) neoMainReactor(listener *ListenWrapper) *MainReactor {
	mr := MainReactor{
		_poller:         ego,
		_listener:       listener,
		_commandChannel: make(chan cms.ICMS, 1),
	}

	return &mr
}

func (ego *Poller) Initialize() int32 {
	ego._stateCode.SetInitializeState()

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
		_logger:       logging.GetLoggerManager().GetDefaultLogger(),
		_mainReactors: make([]*MainReactor, 0),
		_stateCode:    datatype.StateCode(0),
	}

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
