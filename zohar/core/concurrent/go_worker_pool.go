package concurrent

import (
	"math/rand"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/process"
)

type GoWorkerPool struct {
	_name           string
	_initCount      int32
	_startedHandler func(*GoWorker)
	_stopHandler    func(*GoWorker)
	_runHandler     func(*GoWorker)
	_workers        []*GoWorker
	_counter        int32
	_object         any
	_sleepInterval  int64
	_waitGroup      sync.WaitGroup
	_config         *intrinsic.GoWorkerPoolConfig
	_shuttingDown   bool
}

func (ego *GoWorkerPool) IsCurrentGoRoutineInPool() bool {
	if ego._workers == nil {
		return false
	}
	for i := 0; i < len(ego._workers); i++ {
		w := ego._workers[i]
		if w != nil {
			if w.IsSameGoRoutine(process.GetCurrentGoRoutineId()) {
				return true
			}
		}
	}
	return false
}

func (ego *GoWorkerPool) Wait() {
	ego._waitGroup.Wait()
}

func (ego *GoWorkerPool) WorkerCount() int {
	if ego._workers == nil {
		return 0
	}
	return len(ego._workers)
}

func (ego *GoWorkerPool) PostTask(proc datatype.TaskFuncType, obj any) {
	if ego._shuttingDown {
		return
	}
	wc := ego.WorkerCount()
	if wc < 1 {
		return
	}
	idx := rand.Intn(wc)
	ego._workers[idx].PostTask(proc, obj)
}

func (ego *GoWorkerPool) BroadcastTask(proc datatype.TaskFuncType, obj any) {
	if ego._shuttingDown {
		return
	}
	wc := ego.WorkerCount()
	if wc < 1 {
		return
	}
	for i := 0; i < wc; i++ {
		ego._workers[i].PostTask(proc, obj)
	}
}

func (ego *GoWorkerPool) Start() int32 {
	return ego.SetWorkerCount(ego._config.InitialCount)
}

func (ego *GoWorkerPool) SetWorkerCount(cnt int32) int32 {
	c := max(cnt, 0)
	for len(ego._workers) < int(c) {
		w := NeoGoWorker(ego._config.Name, ego._counter, ego._startedHandler, ego._runHandler, ego._stopHandler, ego._object, ego._sleepInterval, &ego._waitGroup)
		if w == nil {
			return core.MkErr(core.EC_ERROR_COUNT, 1)
		}
		w.Start()
		ego._workers = append(ego._workers, w)
		ego._counter++
	}

	for len(ego._workers) > int(c) {
		idx := len(ego._workers) - 1
		w := ego._workers[idx]
		ego._workers = append(ego._workers[:idx], ego._workers[idx+1:]...)
		w.Stop()
	}

	ego._counter = c
	return core.MkSuccess(0)
}

func (ego *GoWorkerPool) Stop() {
	ego._shuttingDown = true
	ego.SetWorkerCount(0)
}

func NeoGoWorkerPool(startedHdl func(*GoWorker), runningHdl func(*GoWorker),
	stoppedHdl func(*GoWorker), data any, cfg *intrinsic.GoWorkerPoolConfig) *GoWorkerPool {
	if cfg == nil {
		cfg = &intrinsic.GoWorkerPoolConfig{
			Name:          "WorkerPool",
			PulseInterval: 1000,
			InitialCount:  2,
		}
	}

	wp := &GoWorkerPool{
		_name:           cfg.Name,
		_initCount:      cfg.InitialCount,
		_startedHandler: startedHdl,
		_runHandler:     runningHdl,
		_stopHandler:    stoppedHdl,
		_workers:        nil,
		_object:         data,
		_sleepInterval:  cfg.PulseInterval,
		_config:         cfg,
	}
	return wp
}
