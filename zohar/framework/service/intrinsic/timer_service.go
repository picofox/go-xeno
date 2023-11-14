package intrinsic

import (
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/sched/cron"
)

type CronService struct {
	_cron           *cron.Cron
	_serviceManager *ServiceManager
	_shuttingDown   bool
}

func (ego *CronService) Initialize() int32 {
	return core.MkSuccess(0)
}

func (ego *CronService) Finalize() int32 {
	return core.MkSuccess(0)
}

func (ego *CronService) Start() int32 {
	ego._cron.Start()
	return core.MkSuccess(0)
}

func (ego *CronService) Stop() int32 {
	ego._cron.Stop()
	return core.MkSuccess(0)
}

func NeoCronService(sm *ServiceManager, location *time.Location) *CronService {
	s := CronService{
		_cron:           cron.NewWithLocation(location),
		_serviceManager: sm,
		_shuttingDown:   false,
	}

	return &s
}
