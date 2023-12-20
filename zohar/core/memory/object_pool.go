package memory

import (
	"sync"
)

type ObjectPool[T any] struct {
	ObjectPoolBared[T]
	_lock sync.Mutex
}

func (ego *ObjectPool[T]) String() string {
	return ego.ObjectPoolBared.String()
}

func (ego *ObjectPool[T]) Length() int {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	return ego.ObjectPoolBared.Length()
}

func (ego *ObjectPool[T]) Alloc() *T {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	return ego.ObjectPoolBared.Alloc()
}

func (ego *ObjectPool[T]) Free(element **T) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	return ego.ObjectPoolBared.Free(element)
}

func NeoObjectPool[T any](initialElements int, funInit func(*T)) *ObjectPool[T] {
	op := &ObjectPool[T]{
		ObjectPoolBared: *NeoObjectPoolBared(initialElements, funInit),
	}

	return op
}
