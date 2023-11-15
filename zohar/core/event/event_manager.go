package event

import (
	"container/list"
	"sync"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/unique"
)

type EventManager struct {
	_events map[string]*list.List
	_seq    unique.SequentialGenerator
	_lock   sync.RWMutex
}

func (ego *EventManager) Register(eventName string, e uint8, f datatype.TaskFuncType, a any) int64 {
	uid := ego._seq.Next()
	t := NeoTask(uid, e, f, a)
	ego.registerTask(eventName, t)
	return uid
}

func (ego *EventManager) registerTask(eventName string, task *Task) {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	tq, ok := ego._events[eventName]
	if !ok {
		tq = list.New()
		tq.PushBack(task)
		ego._events[eventName] = tq
	} else {
		tq.PushBack(task)
	}
}

func (ego *EventManager) Unregister(eventName string, uid int64) bool {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	tq, ok := ego._events[eventName]
	if !ok {
		return false
	} else {
		for i := tq.Front(); i != nil; i = i.Next() {
			tq.Remove(i)
			return true
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
		return true
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
