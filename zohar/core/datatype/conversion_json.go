package datatype

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// JSONToSliceE converts the JSON-encoded data to any type slice with no error returned.
func JSONToSliceNoRet[S ~[]E, E any](data []byte) S {
	s, _ := JSONToSlice[S](data)
	return s
}

// JSONToSliceE converts the JSON-encoded data to any type slice.
// E.g. a JSON value ["foo", "bar", "baz"] can be converted to []strs{"foo", "bar", "baz"}
// when calling JSONToSliceE[[]strs](`["foo", "bar", "baz"]`).
func JSONToSlice[S ~[]E, E any](data []byte) (S, error) {
	var s S
	err := json.Unmarshal(data, &s)
	return s, err
}

//
// Convert map keys and values to slice in indeterminate order.
// E.g. covert map[strs]int{"a":1,"b":2, "c":3} to []strs{"a", "c", "b"} and []int{1, 3, 2}.
//

// MapKeys returns a slice of all the keys in m.
// The keys returned are in indeterminate order.
// You can also use standard library golang.org/x/exp/maps#Keys.
func MapKeys[K comparable, V any, M ~map[K]V](m M) []K {
	s := make([]K, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return s
}

// MapVals returns a slice of all the values in m.
// The values returned are in indeterminate order.
// You can also use standard library golang.org/x/exp/maps#Values.
func MapValue[K comparable, V any, M ~map[K]V](m M) []V {
	s := make([]V, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}

// MapKeyVals returns two slice of all the keys and values in m.
// The keys and values are returned in an indeterminate order.
func MapKeysValues[K comparable, V any, M ~map[K]V](m M) ([]K, []V) {
	ks, vs := make([]K, 0, len(m)), make([]V, 0, len(m))
	for k, v := range m {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	return ks, vs
}

// MapToSlice converts map keys and values to slice in indeterminate order.
func MapToSliceNoRet(a any) (ks any, vs any) {
	ks, vs, _ = MapToSlice(a)
	return
}

// MapToSliceE converts keys and values of map to slice in indeterminate order with error.
func MapToSlice(a any) (ks any, vs any, err error) {
	t := reflect.TypeOf(a)
	if t.Kind() != reflect.Map {
		err = fmt.Errorf("the input %#v of type %T isn't a map", a, a)
		return
	}

	// Convert.
	m := reflect.ValueOf(a)
	keys := m.MapKeys()
	ksT, vsT := reflect.SliceOf(t.Key()), reflect.SliceOf(t.Elem())
	ksV, vsV := reflect.MakeSlice(ksT, 0, m.Len()), reflect.MakeSlice(vsT, 0, m.Len())
	for _, k := range keys {
		ksV = reflect.Append(ksV, k)
		vsV = reflect.Append(vsV, m.MapIndex(k))
	}
	return ksV.Interface(), vsV.Interface(), nil
}

// SplitStrToSlice splits a strs to a slice by the specified separator.
func StrToSliceNoRet[T any](s, sep string) []T {
	v, _ := StrToSlice[T](s, sep)
	return v
}

// SplitStrToSliceE splits a strs to a slice by the specified separator and returns an error if occurred.
// Note that this function is implemented through 1.18 generics, so the element type needs to
// be specified when calling it, e.g. SplitStrToSliceE[int]("1,2,3", ",").
func StrToSlice[T any](s, sep string) ([]T, error) {
	ss := strings.Split(s, sep)
	r := make([]T, len(ss))
	for i := range ss {
		v, err := AnyToType[T](ss[i])
		if err != nil {
			return nil, err
		}
		r[i] = v

	}
	return r, nil
}
