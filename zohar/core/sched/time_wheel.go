package sched

import (
	"fmt"
	"sync"
	"xeno/zohar/core/datetime"
)

const (
	TIME_NEAR_SHIFT  = 8
	TIME_NEAR        = 1 << 8
	TIME_LEVEL_SHIFT = 6
	TIME_LEVEL       = 1 << 6
	TIME_NEAR_MASK   = TIME_NEAR - 1
	TIME_LEVEL_MASK  = TIME_LEVEL - 1
)

type TimeWheel struct {
	_lock         sync.RWMutex
	_near         [TIME_NEAR]*TimerLinkedList
	_t            [4][TIME_LEVEL]*TimerLinkedList
	_time         uint32
	_startTime    uint32
	_current      uint64
	_currentPoint uint64
	_milliInUnit  uint32
}

func (ego *TimeWheel) getTime() uint64 {
	monoNano := datetime.GetMonotonicNano()
	sec := monoNano / 1000000000
	nsec := monoNano - (sec * 1000000000)
	t := sec * (1000 / int64(ego._milliInUnit))
	t = t + nsec/(1000000*int64(ego._milliInUnit))
	return uint64(t)
}

func (ego *TimeWheel) systime(sec *uint32, cs *uint32) {
	tv := datetime.GetRealTimeNano()
	*sec = uint32(tv / 1000000000)
	*cs = uint32((tv - (int64(ego._startTime) * 1000000000)) / 1000000 * int64(ego._milliInUnit))
}

func (ego *TimeWheel) dispatchList(nodeList *TimerLinkedListNode) {
	for {
		if nodeList._data != nil {
			nodeList._data.Execute() //<-fix
			bResched := nodeList._data.reSchedCheck(ego._time)
			if bResched {
				node := NeoTimerLinkedListNode(nodeList._data)
				ego._lock.Lock()
				ego.AddNode(node)
				ego._lock.Unlock()
			}
		}
		nodeList = nodeList._next
		if nodeList == nil {
			break
		}
	}
}

func (ego *TimeWheel) execute() {
	idx := ego._time & TIME_NEAR_MASK
	for ego._near[idx]._head._next != nil {
		current := ego._near[idx].Clear()
		ego._lock.Unlock()
		// dispatch_list don't need lock T
		ego.dispatchList(current)
		ego._lock.Lock()
	}

}

func (ego *TimeWheel) moveList(level uint32, idx uint32) {
	current := ego._t[level][idx].Clear()
	for current != nil {
		tmp := current._next
		ego.AddNode(current)
		current = tmp
	}
}

func (ego *TimeWheel) shift() {
	mask := uint32(TIME_NEAR)
	ego._time++
	ct := ego._time
	if ct == 0 {
		ego.moveList(3, 0)
	} else {
		time := ct >> TIME_NEAR_SHIFT
		var i uint32 = 0
		for (ct & (mask - 1)) == 0 {
			idx := time & TIME_LEVEL_MASK
			if idx != 0 {
				ego.moveList(i, idx)
				break
			}
			mask <<= TIME_LEVEL_SHIFT
			time >>= TIME_LEVEL_SHIFT
			i++
		}
	}
}

func (ego *TimeWheel) update() {
	ego._lock.Lock()
	defer ego._lock.Unlock()

	ego.execute()
	ego.shift()
	ego.execute()
}

func (ego *TimeWheel) UpdateTime() {
	cp := ego.getTime() //how many slot of current time?
	if cp < ego._currentPoint {
		fmt.Sprintf("time diff error: change from %d to %d", cp, ego._currentPoint)
		ego._currentPoint = cp
	} else if cp != ego._currentPoint {
		diff := cp - ego._currentPoint
		ego._currentPoint = cp
		ego._current += diff
		for i := uint64(0); i < diff; i++ {
			ego.update()
		}
	}
}

func (ego *TimeWheel) Initialize() {
	var current uint32 = 0
	ego.systime(&ego._startTime, &current)
	ego._current = uint64(current)
	ego._currentPoint = ego.getTime()
}

func (ego *TimeWheel) AddNode(node *TimerLinkedListNode) {
	var time uint32 = node._data._expire
	var currentTime uint32 = ego._time

	if (time | TIME_NEAR_MASK) == (currentTime | TIME_NEAR_MASK) {
		ego._near[time&TIME_NEAR_MASK].Link(node)
	} else {
		var i int
		var mask uint = TIME_NEAR << TIME_LEVEL_SHIFT
		for i = 0; i < 3; i++ {
			if (time | uint32(mask-1)) == (currentTime | uint32(mask-1)) {
				break
			}
			mask <<= TIME_LEVEL_SHIFT
		}
		idx := (time >> (TIME_NEAR_SHIFT + i*TIME_LEVEL_SHIFT)) & TIME_LEVEL_MASK
		ego._t[i][idx].Link(node)
	}
}

func (ego *TimeWheel) AddTimer(duration uint32, repCount int64, repDura uint32, executor uint8, cb func(any), obj any) *Timer {
	timer := NeoTimer(duration, repCount, repDura, executor, cb, obj)
	node := NeoTimerLinkedListNode(timer)
	ego._lock.Lock()
	defer ego._lock.Unlock()
	node._data._expire = duration + ego._time
	ego.AddNode(node)
	return timer
}

func NeoTimerWheel(millisInUnit uint32) *TimeWheel {
	t := TimeWheel{
		_time:         0,
		_startTime:    0,
		_current:      0,
		_currentPoint: 0,
		_milliInUnit:  millisInUnit,
	}
	for i := 0; i < TIME_NEAR; i++ {
		t._near[i] = NeoTimerLinkedList()
	}

	for i := 0; i < 4; i++ {
		for j := 0; j < TIME_LEVEL; j++ {
			t._t[i][j] = NeoTimerLinkedList()
		}
	}

	return &t
}
