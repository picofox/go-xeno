package memory

import "container/list"

type ObjectCacheBared[T any] struct {
	_freeList     *list.List
	_creationFunc func(a ...any) *T
	_initCount    int
	_maxCount     int
	_args         []any
}

func (ego *ObjectCacheBared[T]) Count() int {
	return ego._freeList.Len()
}

func (ego *ObjectCacheBared[T]) Get() *T {
	if ego._freeList.Front() == nil {
		elem := ego._creationFunc(ego._args...)
		return elem
	} else {
		fn := ego._freeList.Front()
		elem := fn.Value.(*T)
		ego._freeList.Remove(fn)
		return elem
	}
}

func (ego *ObjectCacheBared[T]) Put(elem *T) {
	c := ego._freeList.Len()
	if c < ego._maxCount {
		ego._freeList.PushBack(elem)
	}
}

func NeoObjectCacheBared[T any](initialCount int, maxCount int, cf func(a ...any) *T, a ...any) *ObjectCacheBared[T] {
	c := ObjectCacheBared[T]{
		_freeList:     list.New(),
		_creationFunc: cf,
		_initCount:    initialCount,
		_maxCount:     maxCount,
		_args:         a,
	}

	for i := 0; i < initialCount; i++ {
		var elem *T = cf(c._args...)
		c._freeList.PushBack(elem)
	}
	return &c
}
