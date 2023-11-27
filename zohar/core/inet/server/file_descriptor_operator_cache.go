package server

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

func newOperatorCache() *operatorCache {
	return &operatorCache{
		cache:    make([]*FileDescriptorOperator, 0, 1024),
		freelist: make([]int32, 0, 1024),
	}
}

type operatorCache struct {
	locked int32
	first  *FileDescriptorOperator
	cache  []*FileDescriptorOperator
	// freelist store the freeable operator
	// to reduce GC pressure, we only store op index here
	freelist   []int32
	freelocked int32
}

func (c *operatorCache) alloc() *FileDescriptorOperator {
	lock(&c.locked)
	if c.first == nil {
		const opSize = unsafe.Sizeof(FileDescriptorOperator{})
		n := block4k / opSize
		if n == 0 {
			n = 1
		}
		index := int32(len(c.cache))
		for i := uintptr(0); i < n; i++ {
			pd := &FileDescriptorOperator{index: index}
			c.cache = append(c.cache, pd)
			pd.next = c.first
			c.first = pd
			index++
		}
	}
	op := c.first
	c.first = op.next
	unlock(&c.locked)
	return op
}

// freeable mark the operator that could be freed
// only poller could do the real free action
func (c *operatorCache) freeable(op *FileDescriptorOperator) {
	// reset all state
	op.unused()
	op.reset()
	lock(&c.freelocked)
	c.freelist = append(c.freelist, op.index)
	unlock(&c.freelocked)
}

func (c *operatorCache) free() {
	lock(&c.freelocked)
	defer unlock(&c.freelocked)
	if len(c.freelist) == 0 {
		return
	}

	lock(&c.locked)
	for _, idx := range c.freelist {
		op := c.cache[idx]
		op.next = c.first
		c.first = op
	}
	c.freelist = c.freelist[:0]
	unlock(&c.locked)
}

func lock(locked *int32) {
	for !atomic.CompareAndSwapInt32(locked, 0, 1) {
		runtime.Gosched()
	}
}

func unlock(locked *int32) {
	atomic.StoreInt32(locked, 0)
}
