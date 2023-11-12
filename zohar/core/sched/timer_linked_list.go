package sched

type TimerLinkedListNode struct {
	_next        *TimerLinkedListNode
	_repeatCount int32
	_duration    uint32
	_expire      uint32
	_eventFunc   func(any) int32
	_eventObject any
}

func NeoTimerLinkedListNode(dura uint32, rcnt int32, f func(any) int32, obj any) *TimerLinkedListNode {
	return &TimerLinkedListNode{
		_next:        nil,
		_repeatCount: rcnt,
		_duration:    dura,
		_expire:      0,
		_eventFunc:   f,
		_eventObject: obj,
	}
}

type TimerLinkedList struct {
	_head TimerLinkedListNode
	_tail *TimerLinkedListNode
}

func NeoTimerLinkedList() *TimerLinkedList {
	sl := TimerLinkedList{
		_tail: nil,
	}
	sl._tail = &sl._head
	return &sl
}

func (ego *TimerLinkedList) Tail() *TimerLinkedListNode {
	return ego._tail
}

func (ego *TimerLinkedList) Link(node *TimerLinkedListNode) {
	ego._tail._next = node
	ego._tail = node
	node._next = nil
}

func (ego *TimerLinkedList) Clear() *TimerLinkedListNode {
	ret := ego._head._next
	ego._head._next = nil
	ego._tail = &(ego._head)
	return ret
}
