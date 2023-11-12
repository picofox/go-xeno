package sched

import "xeno/zohar/core/unique"

type Timer struct {
	_id          int64
	_duration    int64
	_repeatCount int64
	_expire      int64
	_eventFunc   func(any) int32
	_eventObject any
}

func (ego *Timer) ReSchedCheck(timeSlotPointer int32) bool {
	if ego._repeatCount < 0 {
		ego._expire = ego._duration + int64(timeSlotPointer)
		return true
	} else if ego._repeatCount == 0 {
		return false
	}
	ego._repeatCount--
	if ego._repeatCount > 0 {
		ego._expire = ego._duration + int64(timeSlotPointer)
		return true
	}
	return false
}

var s_timerUidGenerator unique.SequentialGenerator

func NeoTimer(dura int64, repCount int64, cb func(any) int32, obj any) *Timer {
	tm := Timer{
		_id:          s_timerUidGenerator.Next(),
		_duration:    dura,
		_repeatCount: repCount,
		_expire:      0,
		_eventFunc:   cb,
		_eventObject: obj,
	}
	return &tm
}
