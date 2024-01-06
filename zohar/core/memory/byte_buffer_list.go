package memory

import (
	"fmt"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type ByteBufferList struct {
	_head  *ByteBufferNode
	_tail  *ByteBufferNode
	_count int64
}

var sSplittingTypeString []string = []string{"Sgl", "END", "BEG", "MID"}

func SplittingTypeString(t int8) string {
	return sSplittingTypeString[t]
}

func (ego *ByteBufferList) String() string {
	var ss strings.Builder
	var idx int64 = 0
	var offset int64 = 0
	var frameLength int64 = 0
	var totalFrameLen int64 = 0
	var cmd int16 = 0
	var st int8
	cur := ego._head
	for cur != nil {
		if offset >= cur.ReadAvailable() {
			return ss.String()
		}
		ss.WriteString(fmt.Sprintf("%8d -> ", idx))
		if frameLength == 0 { //header parse

			i0 := BytesToInt16BE(&cur._data, offset)
			i1 := BytesToInt16BE(&cur._data, offset+2)
			frameLength = int64(i0 & 0x7FFF)
			cmd = i1 & 0x7FFF
			o1 := int8(i0 >> 15 & 0x1)
			o2 := int8(i1 >> 15 & 0x1)
			st = (o1 << 1) | o2
			offset += 4
			if cmd != 32767 {
				panic("cmd error")
			}
		}
		rl := cur.ReadAvailableByOffset(offset)
		if frameLength <= rl {
			totalFrameLen += frameLength
			offset += frameLength
			if st == 1 {
				s := fmt.Sprintf("Fin,%d - %d", totalFrameLen, frameLength)
				ss.WriteString(s)
				totalFrameLen = 0
				if offset < 4096 {
					frameLength = 0
					continue
				}

			} else if st == 3 {
				s := fmt.Sprintf("MID,%d - %d", totalFrameLen, frameLength)
				ss.WriteString(s)
				if offset < 4096 {
					frameLength = 0
					continue
				}

			} else if st == 2 {
				s := fmt.Sprintf("BEG,%d - %d", totalFrameLen, frameLength)
				ss.WriteString(s)
				if offset < 4096 {
					frameLength = 0
					continue
				}
			} else {
				panic("type error")
			}

			frameLength = 0

		} else {
			totalFrameLen += rl
			frameLength -= rl
			offset = 0
			s := fmt.Sprintf("%s,%0d - %d - %d", SplittingTypeString(st), totalFrameLen, frameLength, rl)
			ss.WriteString(s)
		}

		ss.WriteString("\n")

		cur = cur._next
		offset = 0
		idx++
	}
	return ss.String()
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

//func (ego *ByteBufferList) DeleteNodes(n int64) int64 {
//	for n > 0 {
//		if ego.PopFront() == nil {
//			return n
//		}
//		n--
//		ego._count--
//	}
//	return n
//}

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
