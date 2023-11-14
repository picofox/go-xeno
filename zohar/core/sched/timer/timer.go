package timer

import (
	"fmt"
	"sync/atomic"
	"xeno/zohar/core/concurrent"
	"xeno/zohar/core/unique"
)

const (
	TASK_EXEC_EXECUTOR_POOL   = uint8(0)
	TASK_EXEC_CURRENT_ROUTINE = uint8(1)
	TASK_EXEC_NEO_ROUTINE     = uint8(2)
)

type Timer struct {
	_id             int64
	_duration       uint32
	_expire         uint32
	_repeatCount    int64
	_repeatDuration uint32
	_executor       uint8
	_eventFunc      func(any)
	_eventObject    any
	_cancel         atomic.Bool
}

func (ego *Timer) String() string {
	return fmt.Sprintf("id:%d dura:%d exp:%d remainCnt:%d", ego._id, ego._duration, ego._expire, ego.RemainCount())
}

func (ego *Timer) RemainCount() int64 {
	return ego._repeatCount
}

func (ego *Timer) Cancel() {
	ego._cancel.Store(true)
}

func (ego *Timer) IsCancelled() bool {
	return ego._cancel.Load()
}

func (ego *Timer) Id() int64 {
	return ego._id
}

func (ego *Timer) Object() any {
	return ego._eventObject
}

var sTimerExecMethodsArr = [3]func(*Timer){
	func(timer *Timer) {
		concurrent.GetDefaultGoExecutorPool().PostTask(timer._eventFunc, timer)
	},
	func(timer *Timer) {
		timer._eventFunc(timer)
	},
	func(timer *Timer) {
		go timer._eventFunc(timer)
	},
}

func (ego *Timer) Execute() {
	if !ego.IsCancelled() {
		if ego._eventFunc != nil {
			sTimerExecMethodsArr[ego._executor](ego)
		}
	}
}

func (ego *Timer) reSchedCheck(timeSlotPointer uint32) bool {
	if ego._repeatCount < 0 {
		ego._expire = ego._repeatDuration + (timeSlotPointer)
		return true
	} else if ego._repeatCount == 0 {
		return false
	}
	ego._repeatCount--
	if ego._repeatCount > 0 {
		ego._expire = ego._repeatDuration + (timeSlotPointer)
		return true
	}
	return false
}

var s_timerUidGenerator unique.SequentialGenerator

func NeoTimer(dura uint32, repCount int64, repDura uint32, executor uint8, cb func(any), obj any) *Timer {
	if executor > TASK_EXEC_NEO_ROUTINE {
		executor = TASK_EXEC_EXECUTOR_POOL
	}
	tm := Timer{
		_id:             s_timerUidGenerator.Next(),
		_duration:       dura,
		_repeatCount:    repCount,
		_expire:         0,
		_repeatDuration: repDura,
		_executor:       executor,
		_eventFunc:      cb,
		_eventObject:    obj,
	}
	tm._cancel.Store(false)
	return &tm
}
