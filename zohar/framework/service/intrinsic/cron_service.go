package intrinsic

import (
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/sched/cron"
)

type CronService struct {
	ServiceCommon
	_cron           *cron.Cron
	_serviceManager *ServiceManager
	_shuttingDown   bool
}

func (ego *CronService) AddCron(spec string, cmd datatype.TaskFuncType, a any, executor uint8) int32 {
	err := ego._cron.AddFunc(spec, cmd, a, executor)
	if err != nil {
		return core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	return core.MkSuccess(0)
}

func (ego *CronService) Initialize() int32 {
	rc := ego.SetInitializeState()
	if core.Err(rc) {
		return rc
	}
	ego.SetInitializeStateResult(true)
	return core.MkSuccess(0)
}

func (ego *CronService) Finalize() int32 {
	rc := ego.SetFinalizeState()
	if core.Err(rc) {
		return rc
	}
	ego.SetFinalizeStateResult(true)
	return core.MkSuccess(0)
}

func (ego *CronService) Start() int32 {
	rc := ego.SetStartState()
	if core.Err(rc) {
		return rc
	}
	ego._serviceManager.addRef()
	ego._cron.Start()
	ego.SetStartStateResult(true)
	return core.MkSuccess(0)
}

func (ego *CronService) Stop() int32 {
	rc := ego.SetStopState()
	if core.Err(rc) {
		return rc
	}
	ego._cron.Stop()
	ego.SetStopStateResult(true)
	return core.MkSuccess(0)
}

func NeoCronService(sm *ServiceManager, location *time.Location) *CronService {
	s := CronService{
		ServiceCommon: ServiceCommon{
			_stateCode: datatype.StateCode(0),
		},
		_cron:           cron.NewWithLocation(location, &sm._waitGroup),
		_serviceManager: sm,
		_shuttingDown:   false,
	}
	return &s
}
