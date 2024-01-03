package memory_test

import (
	"testing"
	"xeno/zohar/core/memory"
)

var testCount int64 = 10240

func Test_ByteBufferList_Functional_Basic(t *testing.T) {
	list := memory.NeoByteBufferList()

	var todel *memory.ByteBufferNode = nil
	for i := int64(0); i < testCount; i++ {
		cur := memory.NeoByteBufferNode(4096)
		list.PushBack(cur)

		if i == 192 {
			todel = cur
			todel.WriteInt64(int64(i))
		}
	}

	list.DeleteUntilReadableNode(todel)

	if list.Count() != testCount-192 {
		t.Errorf("del failed 0")
	}
	if list.Front() != todel {
		t.Errorf("del failed 1")
	}

}
