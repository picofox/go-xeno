package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

func CheckByteBufferListNode(conn IConnection) (*memory.ByteBufferNode, int32) {
	bufferList := conn.BufferBlockList()
	bufNode := bufferList.Back()
	if bufNode == nil {
		bufNode = conn.AllocByteBufferBlock()
		if bufNode == nil {
			return nil, core.MkErr(core.EC_NULL_VALUE, 1)
		}
		bufferList.PushBack(bufNode)
		return bufNode, core.MkSuccess(0)

	} else {
		if bufNode.Buffer().WriteAvailable() <= 0 {
			bufNode = conn.AllocByteBufferBlock()
			if bufNode == nil {
				return nil, core.MkErr(core.EC_NULL_VALUE, 2)
			}
			bufferList.PushBack(bufNode)
			return bufNode, core.MkSuccess(0)
		}
		return bufNode, core.MkSuccess(0)
	}
}
