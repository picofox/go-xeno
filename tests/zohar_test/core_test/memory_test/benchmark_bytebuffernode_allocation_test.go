package memory

import (
	"sync"
	"testing"
	"xeno/zohar/core/memory"
)

func Benchmark_ByteBufferNode_Allocation_Base(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go allocation(&wg)
	}
	wg.Wait()
}

func allocationP(wg *sync.WaitGroup) {
	defer wg.Done()
	list := memory.NeoByteBufferList()
	for i := 0; i < 1024*128; i++ {
		n := memory.GetByteBuffer4KCache().Get()
		list.PushBack(n)
	}
}

func allocation(wg *sync.WaitGroup) {
	defer wg.Done()
	list := memory.NeoByteBufferList()
	for i := 0; i < 1024*128; i++ {
		n := memory.NeoByteBufferNode(4096)
		list.PushBack(n)
	}
}

func Benchmark_ByteBufferNode_Allocation(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go allocationP(&wg)
	}
	wg.Wait()
}
