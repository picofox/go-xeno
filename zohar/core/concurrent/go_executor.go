package concurrent

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/process"
)

type GoExecutor struct {
	_id           int32
	_waitGroup    *sync.WaitGroup
	_channel      chan cms.ICMS
	_stopped      bool
	_shuttingDown bool
	_context      context.Context
	_cancel       context.CancelFunc
}

func (ego *GoExecutor) String() string {
	return fmt.Sprintf("GoExecutor_%d: (%t)-(%t)", ego._id, ego._shuttingDown, ego._stopped)
}

func (ego *GoExecutor) onRunning() {
	defer ego._waitGroup.Done()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	timeout := time.After(1000 * time.Millisecond)
	for {

		select {
		case <-timeout:

		case m := <-ego._channel:
			if m.Id() == cms.CMSID_GOWORKER_TASK {
				m.(*cms.GoWorkerTask).Exec()
			}

		case <-ego._context.Done():
			ego._shuttingDown = false
			ego._stopped = true
			runtime.Goexit()

		}
	}
}

func (ego *GoExecutor) Id() int32 {
	return ego._id
}

func (ego *GoExecutor) IsStopped() bool {
	return ego._stopped
}

func (ego *GoExecutor) IsSameGoRoutine(goid int64) bool {
	curGoid := process.GetCurrentGoRoutineId()
	if curGoid == goid {
		return true
	}
	return false
}

func (ego *GoExecutor) Start() {
	ego._stopped = false
	ego._shuttingDown = false
	ego._waitGroup.Add(1)
	go ego.onRunning()

}

func (ego *GoExecutor) Stop() {
	ego._shuttingDown = true
	ego._cancel()
}

func NeoGoExecutor(id int32, ch chan cms.ICMS, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) *GoExecutor {
	w := &GoExecutor{
		_id:           id,
		_waitGroup:    wg,
		_stopped:      true,
		_shuttingDown: false,
		_channel:      ch,
		_context:      ctx,
		_cancel:       cancel,
	}
	return w
}
