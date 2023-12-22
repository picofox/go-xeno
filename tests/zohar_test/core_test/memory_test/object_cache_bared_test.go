package memory

import (
	"testing"
	"xeno/zohar/core/memory"
)

func CreateByteBufferNode(param ...any) *memory.ByteBufferNode {
	if len(param) > 0 {
		return memory.NeoByteBufferNode(param[0].(int64))
	}
	return nil
}

func Test_ObjectCacheBared_Functional_Basic(t *testing.T) {
	cache := memory.NeoObjectCacheBared[memory.ByteBufferNode](1024, 1024*1024, CreateByteBufferNode, int64(4096))
	for i := 0; i < 1024; i++ {
		node := cache.Get()
		if node == nil {
			t.Errorf("allocFailed")
		}
	}

	if cache.Count() != 0 {
		t.Errorf("count Failed")
	}

	for i := 0; i < 1024; i++ {
		node := cache.Get()
		if node == nil {
			t.Errorf("allocFailed")
		}
	}

	for i := 0; i < 1024; i++ {
		node := CreateByteBufferNode(int64(4096))
		if node == nil {
			t.Errorf("allocFailed")
		}

		cache.Put(node)
	}

	if cache.Count() != 1024 {
		t.Errorf("count Failed")
	}

	for i := 0; i < 1024*1024; i++ {
		node := CreateByteBufferNode(int64(4096))
		if node == nil {
			t.Errorf("allocFailed")
		}

		cache.Put(node)
	}
	if cache.Count() != 1024*1024 {
		t.Errorf("count Failed")
	}
}
