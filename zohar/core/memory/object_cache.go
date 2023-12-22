package memory

import (
	"sync"
	"sync/atomic"
)

type ObjectCache[T any] struct {
	_pool         *sync.Pool
	_creationFunc func() any
	_count        atomic.Int64
	_initialCount int64
}

func (ego *ObjectCache[T]) Get() *T {
	ego._count.Add(-1)
	return ego._pool.Get().(*T)
}

func (ego *ObjectCache[T]) Count() int64 {
	return ego._count.Load()
}

func (ego *ObjectCache[T]) Put(elem *T) {
	ego._count.Add(1)
	ego._pool.Put(elem)
}

func NeoObjectCache[T any](initialCount int64, cf func() any) *ObjectCache[T] {
	c := ObjectCache[T]{
		_pool: &sync.Pool{
			New: cf,
		},
		_creationFunc: cf,
		_initialCount: initialCount,
	}
	for i := int64(0); i < initialCount; i++ {
		var elem *T = cf().(*T)
		c._pool.Put(elem)
	}
	return &c
}
