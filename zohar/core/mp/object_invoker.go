package mp

import (
	"reflect"
	"sync"
	"xeno/zohar/core"
)

type ObjectInvoker struct {
	_types sync.Map
}

func (ego *ObjectInvoker) RegisterClass(name string, cls any) {
	ego._types.Store(name, reflect.ValueOf(cls))
}

func (ego *ObjectInvoker) UnregisterClass(name string) {
	ego._types.Delete(name)
}

func (ego *ObjectInvoker) Invoke(output *[]reflect.Value, objName string, mName string, param ...any) int32 {
	var args = make([]reflect.Value, 0)
	for i := 0; i < len(param); i++ {
		args = append(args, reflect.ValueOf(param[i]))
	}
	obj, ok := ego._types.Load(objName)
	if !ok {
		return core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	mm := obj.(reflect.Value).MethodByName(mName)
	if !mm.IsValid() {
		return core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
	}

	*output = mm.Call(args)
	return core.MkSuccess(0)
}

var sObjectInvokerInstance *ObjectInvoker
var sObjectInvokerInstanceOnce sync.Once

func GetDefaultObjectInvoker() *ObjectInvoker {
	sObjectInvokerInstanceOnce.Do(func() {
		sObjectInvokerInstance = &ObjectInvoker{}
	})
	return sObjectInvokerInstance
}
