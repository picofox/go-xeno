package memory

import (
	"container/list"
)

type ObjectPool[T any] struct {
	_freeList    *list.List
	_eleCreation func() *T
	_totalCount  int
	_freeCount   int
}

func (ego *ObjectPool[T]) Length() int {
	return ego._freeList.Len()
}

func (ego *ObjectPool[T]) Get() *T {
	f := ego._freeList.Front()
	if f == nil {
		ele := ego._eleCreation()
		ego._totalCount++
		ego._freeCount--
		return ele
	}
	ret := f.Value
	ego._freeList.Remove(f)
	ego._freeCount--
	return ret.(*T)
}

func (ego *ObjectPool[T]) Put(element *T) {
	ego._freeList.PushBack(element)
	ego._freeCount++
}

func NeoObjectPool[T any](initialElements int, elementCreation func() *T) *ObjectPool[T] {
	l := list.New()
	op := &ObjectPool[T]{
		_freeList:    l,
		_eleCreation: elementCreation,
		_totalCount:  0,
	}
	for i := 0; i < initialElements; i++ {
		ele := op._eleCreation()
		op._freeList.PushBack(ele)
	}
	op._totalCount = initialElements
	op._freeCount = initialElements
	return op
}
