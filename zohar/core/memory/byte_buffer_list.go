package memory

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type ByteBufferList struct {
	_head  *ByteBufferNode
	_tail  *ByteBufferNode
	_count int64
}

func (ego *ByteBufferList) Count() int64 {
	return ego._count
}

func (ego *ByteBufferList) Front() *ByteBufferNode {
	return ego._head
}

func (ego *ByteBufferList) Back() *ByteBufferNode {
	return ego._tail
}

func (ego *ByteBufferList) DeleteNodes(n int64) int64 {
	for n > 0 {
		if ego.PopFront() == nil {
			return n
		}
		n--
		ego._count--
	}
	return n
}

func (ego *ByteBufferList) DeleteUntilReadableNode(node *ByteBufferNode) (int64, *ByteBufferNode) {
	var cnt int64 = 0

	if node != nil {
		for {
			cur := ego._head
			if cur == node {
				if node.ReadAvailable() > 0 {
					return cnt, node

				} else {
					rNode := cur._next
					ego.PopFront()
					return cnt, rNode
				}
			}
			ego.PopFront()
			cnt++
		}
	} else {
		return 0, nil
	}
}

func (ego *ByteBufferList) PopFront() *ByteBufferNode {
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

func (ego *ByteBufferList) PushBack(node *ByteBufferNode) int32 {
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

func (ego *ByteBufferList) PushFront(node *ByteBufferNode) int32 {
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

func NeoByteBufferList() *ByteBufferList {
	return &ByteBufferList{
		_head:  nil,
		_tail:  nil,
		_count: 0,
	}
}
