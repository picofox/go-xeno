package concurrent

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/config"
	"xeno/zohar/core/process"
)

type GoExecutorPool struct {
	_name         string
	_initCount    int32
	_channel      chan cms.ICMS
	_workers      []*GoExecutor
	_counter      int32
	_waitGroup    sync.WaitGroup
	_config       *config.GoExecutorPoolConfig
	_lock         sync.Mutex
	_shuttingDown bool
}

func (ego *GoExecutorPool) String() string {
	var ss strings.Builder
	ss.WriteString("GoExecutorPool:\n")
	ss.WriteString("\tWorker Count: ")
	ss.WriteString(strconv.Itoa(int(ego._counter)))
	ss.WriteString("\n\tShutting Down: ")
	ss.WriteString(strconv.FormatBool(ego._shuttingDown))
	return ss.String()
}

func (ego *GoExecutorPool) IsCurrentGoRoutineInPool() bool {
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

func (ego *GoExecutorPool) Wait() {
	fmt.Printf("call wait gid : %d\n", process.GetCurrentGoRoutineId())
	ego._waitGroup.Wait()
}

func (ego *GoExecutorPool) WorkerCount() int {
	if ego._workers == nil {
		return 0
	}
	return len(ego._workers)
}

func (ego *GoExecutorPool) PostTask(proc func(any), obj any) {
	if ego._shuttingDown {
		return
	}
	wc := int32(ego.WorkerCount())
	if wc < 1 {
		return
	}

	qLen := len(ego._channel)
	if qLen > ego._config.HighWaterMark {
		ego._lock.Lock()
		if wc < ego._config.MaxCount {
			ego.SetWorkerCount(int32(wc) + 1)
		}
		ego._lock.Unlock()
	} else if qLen <= ego._config.LowWaterMark {
		ego._lock.Lock()
		if wc > ego._config.MinCount {
			ego.SetWorkerCount(int32(wc) - 1)
		}
		ego._lock.Unlock()
	}

	task := cms.NeoCMSGoWorkerTask(proc, obj)
	ego._channel <- task
}

func (ego *GoExecutorPool) Start() int32 {
	return ego.SetWorkerCount(ego._config.InitialCount)
}

func (ego *GoExecutorPool) SetWorkerCount(cnt int32) int32 {
	c := max(cnt, 0)
	for len(ego._workers) < int(c) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		w := NeoGoExecutor(ego._counter, ego._channel, &ego._waitGroup, ctx, cancel)
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

func (ego *GoExecutorPool) Stop() {
	ego._shuttingDown = true
	ego.SetWorkerCount(0)
}

func NeoGoExecutorPool() *GoExecutorPool {
	cfg := &config.GetIntrinsicConfig().GoExecutorPool
	wp := &GoExecutorPool{
		_name:         cfg.Name,
		_initCount:    cfg.InitialCount,
		_channel:      make(chan cms.ICMS, cfg.QueueSize),
		_workers:      nil,
		_counter:      0,
		_config:       cfg,
		_waitGroup:    sync.WaitGroup{},
		_shuttingDown: false,
	}

	return wp
}
