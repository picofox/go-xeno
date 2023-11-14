package concurrent

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/process"
)

type GoWorker struct {
	_name           string
	_id             int32
	_startedHandler func(*GoWorker)
	_stopHandler    func(*GoWorker)
	_runHandler     func(*GoWorker)
	_data           any
	_waitGroup      *sync.WaitGroup

	_channel       chan cms.ICMS
	_stopped       bool
	_shuttingDown  bool
	_sleepInterval int64
}

func (ego *GoWorker) String() string {
	return fmt.Sprintf("GoWorker_%s_%d: (%t)-(%t)", ego._name, ego._id, ego._shuttingDown, ego._stopped)
}

func (ego *GoWorker) runSafely(t *cms.GoWorkerTask) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logging.Log(core.LL_ERR, "cron: panic running job: %v\n%s", r, buf)
		}
	}()
	t.Exec()
}

func (ego *GoWorker) onRunning() {
	defer ego._waitGroup.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	defer ego.onStopped()
	for {
		if ego._runHandler != nil {
			ego._runHandler(ego)
		}

		select {
		case m := <-ego._channel:
			if m.Id() == cms.CMSID_GOWORKER_TASK {
				ego.runSafely(m.(*cms.GoWorkerTask))
			} else if m.Id() == cms.CMSID_FINALIZE {

				runtime.Goexit()
			}
		default:
		}

		if ego._sleepInterval >= 0 {
			time.Sleep(time.Duration(ego._sleepInterval))
		}
	}

}

func (ego *GoWorker) onStarted() {
	if ego._startedHandler != nil {
		ego._startedHandler(ego)
	}
}

func (ego *GoWorker) onStopped() {
	ego._stopped = true
	ego._shuttingDown = false
	if ego._startedHandler != nil {
		ego._stopHandler(ego)
	}

}

func (ego *GoWorker) Name() string {
	return ego._name
}

func (ego *GoWorker) Id() int32 {
	return ego._id
}

func (ego *GoWorker) IsStopped() bool {
	return ego._stopped
}

func (ego *GoWorker) IsSameGoRoutine(goid int64) bool {
	curGoid := process.GetCurrentGoRoutineId()
	if curGoid == goid {
		return true
	}
	return false
}

func (ego *GoWorker) PostTask(proc func(any), obj any) {
	if ego._shuttingDown || ego._stopped {
		return
	}
	task := cms.NeoCMSGoWorkerTask(proc, obj)
	ego._channel <- task
}

func (ego *GoWorker) Start() {
	ego._stopped = false
	ego._shuttingDown = false
	ego.onStarted()
	ego._waitGroup.Add(1)
	go ego.onRunning()

}

func (ego *GoWorker) Stop() {
	ego._shuttingDown = true
	finCMS := cms.NeoFinalize()
	ego._channel <- finCMS
}

func (ego *GoWorker) SetSleepInterval(milliSecs int64) {
	ego._sleepInterval = milliSecs * int64(time.Millisecond)
}

func (ego *GoWorker) Wait() int32 {
	if ego._waitGroup == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	ego._waitGroup.Wait()
	return core.MkSuccess(0)
}

func NeoGoWorker(name string, id int32, startedHdl func(*GoWorker), runningHdl func(*GoWorker), stoppedHdl func(*GoWorker), data any, sleepInt int64, wg *sync.WaitGroup) *GoWorker {
	sleepInt = sleepInt * int64(time.Millisecond)
	w := &GoWorker{
		_name:           name,
		_id:             id,
		_startedHandler: startedHdl,
		_stopHandler:    stoppedHdl,
		_runHandler:     runningHdl,
		_data:           data,
		_waitGroup:      wg,
		_stopped:        true,
		_shuttingDown:   false,
		_channel:        make(chan cms.ICMS, 1024),
		_sleepInterval:  sleepInt,
	}
	return w
}
