package intrinsic

import (
	"fmt"
	"runtime"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/sched"
)

type TimerService struct {
	_timewheel      *sched.TimeWheel
	_serviceManager *ServiceManager
	_channel        chan cms.ICMS
	_shuttingDown   bool
}

func (ego *TimerService) _onRunning() {
	ego._shuttingDown = false
	defer ego._serviceManager.delRef()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	for {
		select {
		case m := <-ego._channel:
			if m.Id() == cms.CMSID_GOWORKER_TASK {
				m.(*cms.GoWorkerTask).Exec()
			} else if m.Id() == cms.CMSID_FINALIZE {
				runtime.Goexit()
			}
		default:

		}

		ego._timewheel.UpdateTime()
		time.Sleep(2500 * time.Microsecond)
	}
}

func (ego *TimerService) AddTimer(duration uint32, repCount int32, cb func(any) int32, obj any) int32 {
	return ego._timewheel.AddTimer(duration, repCount, cb, obj)
}

func (ego *TimerService) Initialize() int32 {
	return core.MkSuccess(0)
}

func (ego *TimerService) Finalize() int32 {
	return core.MkSuccess(0)
}

func (ego *TimerService) Start() int32 {
	ego._timewheel.Initialize()
	ego._serviceManager.addRef()
	go ego._onRunning()
	return core.MkSuccess(0)
}

func (ego *TimerService) Stop() int32 {
	ego._shuttingDown = true
	finCMS := cms.NeoFinalize()
	ego._channel <- finCMS
	return core.MkSuccess(0)
}

func NeoTimerService(sm *ServiceManager) *TimerService {
	s := TimerService{
		_timewheel:      sched.NeoTimerWheel(),
		_serviceManager: sm,
		_channel:        make(chan cms.ICMS, 4096),
		_shuttingDown:   false,
	}
	return &s
}
