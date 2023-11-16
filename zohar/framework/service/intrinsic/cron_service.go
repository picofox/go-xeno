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
	rc := ego.BeginInitializing()
	if core.Err(rc) {
		return rc
	}
	ego.EndInitializing()
	ego.BeginInitialized()
	ego.EndInitialized()
	return core.MkSuccess(0)
}

func (ego *CronService) Finalize() int32 {
	rc := ego.BeginFinalizing()
	if core.Err(rc) {
		return rc
	}
	ego.EndFinalizing()
	ego.BeginUninitialized()
	ego.EndUninitialized()
	return core.MkSuccess(0)
}

func (ego *CronService) Start() int32 {
	rc := ego.BeginStarting()
	if core.Err(rc) {
		return rc
	}
	ego.EndStarting()
	ego._serviceManager.addRef()
	ego._cron.Start()
	ego.BeginStarted()
	ego.EndStarted()
	return core.MkSuccess(0)
}

func (ego *CronService) Stop() int32 {
	rc := ego.BeginStopping()
	if core.Err(rc) {
		return rc
	}
	ego.EndStopping()
	ego._cron.Stop()
	ego.BeginStopped()
	ego.EndStopped()
	return core.MkSuccess(0)
}

func NeoCronService(sm *ServiceManager, location *time.Location) *CronService {
	s := CronService{
		ServiceCommon: ServiceCommon{
			_state: datatype.Uninitialized,
		},
		_cron:           cron.NewWithLocation(location, &sm._waitGroup),
		_serviceManager: sm,
		_shuttingDown:   false,
	}
	return &s
}
