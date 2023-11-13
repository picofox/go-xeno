package intrinsic

import (
	"fmt"
	"runtime"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/datetime"
	"xeno/zohar/core/sched"
)

type TimerService struct {
	_timewheel      *sched.TimeWheel
	_timewheelSec   *sched.TimeWheel
	_serviceManager *ServiceManager
	_channel        chan cms.ICMS
	_shuttingDown   bool
	_secondLevelTs  int64
}

func (ego *TimerService) _onRunning() {
	ego._shuttingDown = false
	ego._secondLevelTs = datetime.GetRealTimeMilli()
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

		curr := datetime.GetRealTimeMilli()
		diff := curr - ego._secondLevelTs
		if diff > 1000 {
			ego._timewheelSec.UpdateTime()
			ego._secondLevelTs = curr
		}
	}
}

func (ego *TimerService) AddTimer(duration uint32, repCount int64, cb func(*sched.Timer) int32, obj any) *sched.Timer {
	return ego._timewheel.AddTimer(duration, repCount, cb, obj)
}

func (ego *TimerService) AddSecondTimer(duration uint32, repCount int64, cb func(*sched.Timer) int32, obj any) *sched.Timer {
	return ego._timewheelSec.AddTimer(duration, repCount, cb, obj)
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
		_timewheel:      sched.NeoTimerWheel(10),
		_timewheelSec:   sched.NeoTimerWheel(1000),
		_serviceManager: sm,
		_channel:        make(chan cms.ICMS, 4096),
		_shuttingDown:   false,
	}

	return &s
}
