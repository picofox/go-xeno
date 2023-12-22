package messages

import (
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

func AllocByteBufferBlock() *memory.ByteBufferNode {
	n := memory.GetByteBuffer4KCache().Get()
	n.Clear()
	return n
}

func CheckByteBufferListNode(bufferList *memory.ByteBufferList) (*memory.ByteBufferNode, int32) {
	bufNode := bufferList.Back()
	if bufNode == nil {
		bufNode = AllocByteBufferBlock()
		if bufNode == nil {
			return nil, core.MkErr(core.EC_NULL_VALUE, 1)
		}
		bufferList.PushBack(bufNode)
		return bufNode, core.MkSuccess(0)

	} else {
		if bufNode.Buffer().WriteAvailable() <= 0 {
			bufNode = AllocByteBufferBlock()
			if bufNode == nil {
				return nil, core.MkErr(core.EC_NULL_VALUE, 2)
			}
			bufferList.PushBack(bufNode)
			return bufNode, core.MkSuccess(0)
		}
		return bufNode, core.MkSuccess(0)
	}
}
