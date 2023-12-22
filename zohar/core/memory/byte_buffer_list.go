package container

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type SinglyLinkedListBared[T ISinglyLinkedListNode] struct {
	_head  *T
	_tail  *T
	_count int64
}

func (ego *SinglyLinkedListBared[T]) Count() int64 {
	return ego._count
}

func (ego *SinglyLinkedListBared[T]) Front() *T {
	return ego._head
}

func (ego *SinglyLinkedListBared[T]) Back() *T {
	return ego._tail
}

func (ego *SinglyLinkedListBared[T]) PopFront() *T {
	if ego._head == nil {
		return nil
	}

	n := ego._head
	ego._head = n.Next()
	if ego._head == nil {
		ego._tail = nil
	}
	n.SetNext(nil)
	ego._count--
	return n
}

func (ego *SinglyLinkedListBared) PushBack(node ISinglyLinkedListNode) int32 {
	if ego._count < datatype.INT64_MAX {
		node.SetNext(nil)
		if ego._head != nil {
			ego._tail.SetNext(node)
		} else {
			ego._head = node
		}
		ego._tail = node
		ego._count++
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_REACH_LIMIT, 1)
}

func (ego *SinglyLinkedListBared) PushFront(node ISinglyLinkedListNode) int32 {
	if ego._count < datatype.INT64_MAX {
		if ego._head == nil {
			node.SetNext(nil)
			ego._tail = node
		} else {
			node.SetNext(ego._head)
		}
		ego._head = node
		ego._count++
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_REACH_LIMIT, 1)
}

func NeoSinglyLinkedList() *SinglyLinkedListBared {
	return &SinglyLinkedListBared{
		_head:  nil,
		_tail:  nil,
		_count: 0,
	}
}
