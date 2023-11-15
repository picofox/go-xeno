package event

import (
	"sync"
	"xeno/zohar/core/datatype"
)

type EventManager struct {
	_events map[string][]*Task
	_lock   sync.RWMutex
}

func (ego *EventManager) Register(eventName string, e uint8, f datatype.TaskFuncType, a any) {
	t := NeoTask(e, f, a)
	ego.RegisterTask(eventName, t)
}

func (ego *EventManager) RegisterTask(eventName string, task *Task) {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	tq, ok := ego._events[eventName]
	if !ok {
		tq = make([]*Task, 1)
		tq[0] = task
		ego._events[eventName] = tq
	} else {
		tq = append(tq, task)
	}
}

func (ego *EventManager) Fire(eventName string, overrideExecutor uint8) {
	var task *Task = nil
	{
		ego._lock.RLock()
		defer ego._lock.RUnlock()
		tq, ok := ego._events[eventName]
		if !ok {
			return
		}
		for _, elem := range tq {
			task = elem
		}
	}
	if task != nil {
		if overrideExecutor > datatype.TASK_EXEC_NEO_ROUTINE {
			task.Execute()

		} else {
			task.ExecuteBy(overrideExecutor)
		}
	}
}

var sEventManager *EventManager
var sEventManagerOnce sync.Once

func GetDefaultEventManager() *EventManager {
	sEventManagerOnce.Do(func() {
		sEventManager = NeoEventManager()
	},
	)
	return sEventManager
}

func NeoEventManager() *EventManager {
	evm := EventManager{
		_events: make(map[string][]*Task),
	}
	return &evm
}
