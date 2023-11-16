package event

import (
	"container/list"
	"reflect"
	"sync"
	"xeno/zohar/core/datatype"
)

type EventManager struct {
	_events map[string]*list.List
	_lock   sync.RWMutex
}

func (ego *EventManager) Register(eventName string, e uint8, f datatype.TaskFuncType, a any) {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	t := NeoTask(e, f, a)
	ego.registerTask(eventName, t)
}

func (ego *EventManager) registerTask(eventName string, task *Task) {
	tq, ok := ego._events[eventName]
	if !ok {
		tq = list.New()
		tq.PushBack(task)
		ego._events[eventName] = tq
	} else {
		tq.PushBack(task)
	}
}

func (ego *EventManager) Unregister(eventName string, f datatype.TaskFuncType) bool {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	tq, ok := ego._events[eventName]
	if !ok {
		return false
	} else {
		for i := tq.Front(); i != nil; i = i.Next() {
			if reflect.ValueOf(i.Value.(*Task)._function).Pointer() == reflect.ValueOf(f).Pointer() {
				tq.Remove(i)
			}
		}
		if tq.Len() < 1 {
			delete(ego._events, eventName)
		}

	}
	return false
}

func (ego *EventManager) Fire(eventName string, overrideExecutor uint8) bool {
	var tq *list.List = nil
	var ok bool = false
	ego._lock.Lock()
	defer ego._lock.Unlock()
	tq, ok = ego._events[eventName]
	if !ok {
		return false
	}
	for i := tq.Front(); i != nil; i = i.Next() {
		if i.Value != nil {
			if overrideExecutor > datatype.TASK_EXEC_NEO_ROUTINE {
				i.Value.(*Task).Execute()
			} else {
				i.Value.(*Task).ExecuteBy(overrideExecutor)
			}
		}
	}
	return true
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
		_events: make(map[string]*list.List),
	}
	return &evm
}
