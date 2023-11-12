package memory

import (
	"reflect"
	"strings"
)

//
// Converts an any element type slice or array to the specified type mapping set.
// Note that the the element type of input don't need to be equal to the map key type.
// For example, []uint64{1, 2, 3} can be converted to map[uint64]struct{}{1:{}, 2:{},3:{}}
// and also can be converted to map[strs]struct{}{"1":{}, "2":{}, "3":{}} if you want.
//

// VectorToSetNoRet converts a slice or array to map[T]struct{}.
// An error will be returned if an error occurred.
func VectorToSetNoRet[T comparable](a any) map[T]struct{} {
	m, _ := VectorToSet[T](a)
	return m
}

// VectorToSet converts a slice or array to map[T]struct{} and returns an error if occurred.
// Note that the the element type of input don't need to be equal to the map key type.
// For example, []uint64{1, 2, 3} can be converted to map[uint64]struct{}{1:{},2:{},3:{}}
// and also can be converted to map[strs]struct{}{"1":{},"2":{},"3":{}} if you want.
// Note that this function is implemented through 1.18 generics, so the element type needs to
// be specified when calling it, e.g. ToSetE[int]([]int{1,2,3}).
func VectorToSet[T comparable](a any) (map[T]struct{}, error) {
	t := reflect.TypeOf(a)
	v := reflect.ValueOf(a)
	if t.Kind() == reflect.Slice && v.IsNil() {
		return nil, nil
	}

	// Execute the conversion.
	mapT := reflect.MapOf(t.Elem(), reflect.TypeOf(struct{}{}))
	mapV := reflect.MakeMapWithSize(mapT, v.Len())
	for i := 0; i < v.Len(); i++ {
		mapV.SetMapIndex(v.Index(i), reflect.ValueOf(struct{}{}))
	}
	if v, ok := mapV.Interface().(map[T]struct{}); ok {
		return v, nil
	}
	// Convert the element to the type T.
	set := make(map[T]struct{}, v.Len())
	for _, k := range mapV.MapKeys() {
		v, err := AnyToType[T](k.Interface())
		if err != nil {
			return nil, err
		}
		set[v] = struct{}{}
	}
	return set, nil
}

// StrToSet convert a strs to map set after split
func StrToSet(s string, sep string) map[string]struct{} {
	if s == "" {
		return nil
	}
	m := make(map[string]struct{})
	for _, v := range strings.Split(s, sep) {
		m[v] = struct{}{}
	}
	return m
}
