package timer

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/datetime"
	"xeno/zohar/core/sched"
)

type TimerManager struct {
	_timewheel     *TimeWheel
	_timewheelSec  *TimeWheel
	_channel       chan cms.ICMS
	_shuttingDown  bool
	_secondLevelTs int64
	_waitGroup     sync.WaitGroup
}

func (ego *TimerManager) _onRunning() {
	ego._shuttingDown = false
	ego._secondLevelTs = datetime.GetRealTimeMilli()
	defer ego._waitGroup.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		select {
		case m := <-ego._channel:
			if m.Id() == cms.CMSID_FINALIZE {
				runtime.Goexit()
			}
		default:
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
}

func (ego *TimerManager) AddAbsTimerMilli(epochMillis int64, repCount int64, repDura uint32, executor uint8, cb sched.TaskFuncType, obj any) *Timer {
	nowTs := datetime.GetRealTimeMilli()
	diff := epochMillis - nowTs
	if diff < 0 {
		diff = 0
	}
	return ego.AddRelTimerMilli(diff, repCount, repDura, executor, cb, obj)
}

func (ego *TimerManager) AddAbsTimerSecond(epochSeconds int64, repCount int64, repDura uint32, executor uint8, cb sched.TaskFuncType, obj any) *Timer {
	nowTs := datetime.GetRealTimeMilli()
	diff := (epochSeconds*1000 - nowTs) / 1000
	if diff < 0 {
		diff = 0
	}
	return ego.AddRelTimerSecond(uint32(diff), repCount, repDura, executor, cb, obj)
}

func (ego *TimerManager) AddRelTimerMilli(millis int64, repCount int64, repDura uint32, executor uint8, cb sched.TaskFuncType, obj any) *Timer {
	d := uint32(millis / 10)
	return ego._timewheel.AddTimer(d, repCount, repDura, executor, cb, obj)
}

func (ego *TimerManager) AddRelTimerSecond(duration uint32, repCount int64, repDura uint32, executor uint8, cb sched.TaskFuncType, obj any) *Timer {
	return ego._timewheelSec.AddTimer(duration, repCount, repDura, executor, cb, obj)
}

func (ego *TimerManager) Wait() {
	ego._waitGroup.Wait()
}

func (ego *TimerManager) Start() int32 {

	ego._waitGroup.Add(1)
	go ego._onRunning()
	return core.MkSuccess(0)
}

func (ego *TimerManager) Stop() int32 {
	ego._shuttingDown = true
	finCMS := cms.NeoFinalize()
	ego._channel <- finCMS
	return core.MkSuccess(0)
}

func NeoTimerManager() *TimerManager {
	s := TimerManager{
		_timewheel:    NeoTimerWheel(10),
		_timewheelSec: NeoTimerWheel(1000),
		_channel:      make(chan cms.ICMS, 4),
		_shuttingDown: false,
	}

	return &s
}
