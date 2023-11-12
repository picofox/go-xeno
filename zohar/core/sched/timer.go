package sched

import (
	"sync/atomic"
	"xeno/zohar/core/unique"
)

type Timer struct {
	_id          int64
	_duration    uint32
	_expire      uint32
	_repeatCount int64
	_eventFunc   func(*Timer) int32
	_eventObject any
	_cancel      atomic.Bool
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

func (ego *Timer) Execute() {
	if !ego.IsCancelled() {
		if ego._eventFunc != nil {
			go ego._eventFunc(ego)
		}
	}
}

func (ego *Timer) reSchedCheck(timeSlotPointer uint32) bool {
	if ego._repeatCount < 0 {
		ego._expire = ego._duration + (timeSlotPointer)
		return true
	} else if ego._repeatCount == 0 {
		return false
	}
	ego._repeatCount--
	if ego._repeatCount > 0 {
		ego._expire = ego._duration + (timeSlotPointer)
		return true
	}
	return false
}

var s_timerUidGenerator unique.SequentialGenerator

func NeoTimer(dura uint32, repCount int64, cb func(*Timer) int32, obj any) *Timer {
	tm := Timer{
		_id:          s_timerUidGenerator.Next(),
		_duration:    dura,
		_repeatCount: repCount,
		_expire:      0,
		_eventFunc:   cb,
		_eventObject: obj,
	}
	tm._cancel.Store(false)
	return &tm
}
