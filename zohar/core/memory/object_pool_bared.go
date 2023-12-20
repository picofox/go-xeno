package memory

import (
	"container/list"
	"fmt"
	"xeno/zohar/core"
)

type ObjectPoolBared[T any] struct {
	_freeList    *list.List
	_usedMap     map[*T]int8
	_eleInitFunc func(*T)
	_initCount   int
}

func (ego *ObjectPoolBared[T]) String() string {
	return fmt.Sprintf("i:%d,f:%d,u:%d", ego._initCount, ego._freeList.Len(), len(ego._usedMap))
}

func (ego *ObjectPoolBared[T]) Length() int {
	return ego._freeList.Len()
}

func (ego *ObjectPoolBared[T]) Alloc() *T {
	f := ego._freeList.Front()
	var rv *T = nil
	if f != nil {
		rv = f.Value.(*T)
		ego._freeList.Remove(f)
	} else {
		rv = new(T)
	}
	ego._usedMap[rv] = 1
	return rv
}

func (ego *ObjectPoolBared[T]) Free(element **T) int32 {
	_, ok := ego._usedMap[*element]
	if ok {
		delete(ego._usedMap, *element)
		ego._freeList.PushBack(element)
		*element = nil
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)

}

func NeoObjectPoolBared[T any](initialElements int, funInit func(*T)) *ObjectPoolBared[T] {
	l := list.New()
	op := &ObjectPoolBared[T]{
		_freeList:    l,
		_usedMap:     make(map[*T]int8),
		_eleInitFunc: funInit,
		_initCount:   initialElements,
	}
	for i := 0; i < initialElements; i++ {
		ele := new(T)
		if op._eleInitFunc != nil {
			op._eleInitFunc(ele)
			op._freeList.PushBack(ele)
		}
		op._freeList.PushBack(ele)
	}

	return op
}
