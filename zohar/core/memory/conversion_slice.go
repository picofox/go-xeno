package memory

import (
	"reflect"
	"strings"
	"time"
)

// ToSlice converts any type slice or array to the specified type slice.
func SliceOfAnyToSliceOfTypeNoRetWithLength[T any](a any, maxlen uint32) ([]T, uint32) {
	r, rlen, _ := SliceOfAnyToSliceOfTypeWithLength[T](a, maxlen)
	return r, rlen
}

// ToSliceE converts any type slice or array to the specified type slice.
// An error will be returned if an error occurred.
func SliceOfAnyToSliceOfTypeWithLength[T any](a any, maxlen uint32) ([]T, uint32, error) {
	if a == nil {
		return nil, 0, nil
	}
	switch v := a.(type) {
	case []T:
		l := min(maxlen, uint32(len(v)))
		return v, l, nil
	case string:
		return SliceOfAnyToSliceOfTypeWithLength[T](strings.Fields(v), maxlen)
	}

	kind := reflect.TypeOf(a).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		// If input is a slice or array.
		v := reflect.ValueOf(a)
		if kind == reflect.Slice && v.IsNil() {
			return nil, 0, nil
		}
		s := make([]T, v.Len())
		l := min(maxlen, uint32(v.Len()))
		var i uint32 = 0
		for ; i < l; i++ {
			val, err := AnyToType[T](v.Index(int(i)).Interface())
			if err != nil {
				return nil, 0, err
			}
			s[i] = val
		}

		return s, l, nil
	default:
		// If input is a single value.
		v, err := AnyToType[T](a)
		if err != nil {
			return nil, 0, err
		}
		return []T{v}, 1, nil
	}
}

// ToSlice converts any type slice or array to the specified type slice.
func SliceOfAnyToSliceOfTypeNoRet[T any](a any) []T {
	r, _ := SliceOfAnyToSliceOfType[T](a)
	return r
}

// ToSliceE converts any type slice or array to the specified type slice.
// An error will be returned if an error occurred.
func SliceOfAnyToSliceOfType[T any](a any) ([]T, error) {
	if a == nil {
		return nil, nil
	}
	switch v := a.(type) {
	case []T:
		return v, nil
	case string:
		return SliceOfAnyToSliceOfType[T](strings.Fields(v))
	}

	kind := reflect.TypeOf(a).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		// If input is a slice or array.
		v := reflect.ValueOf(a)
		if kind == reflect.Slice && v.IsNil() {
			return nil, nil
		}
		s := make([]T, v.Len())
		for i := 0; i < v.Len(); i++ {
			val, err := AnyToType[T](v.Index(i).Interface())
			if err != nil {
				return nil, err
			}
			s[i] = val
		}
		return s, nil
	default:
		// If input is a single value.
		v, err := AnyToType[T](a)
		if err != nil {
			return nil, err
		}
		return []T{v}, nil
	}
}

// ToBoolSlice converts any type to []bool.
func SliceOfAnyToSliceOfBoolNoRet(a any) []bool {
	return SliceOfAnyToSliceOfTypeNoRet[bool](a)
}

// ToBoolSliceE converts any type slice or array to []bool with returned error.
func SliceOfAnyToSliceOfBool(a any) ([]bool, error) {
	return SliceOfAnyToSliceOfType[bool](a)
}

// ToIntSlice converts any type slice or array to []int.
// E.g. covert []strs{"1", "2", "3"} to []int{1, 2, 3}.
func SliceOfAnyToSliceOfIntNoRet(a any) []int {
	return SliceOfAnyToSliceOfTypeNoRet[int](a)
}

// ToIntSliceE converts any type slice or array to []int with returned error..
func SliceOfAnyToSliceOfInt(a any) ([]int, error) {
	return SliceOfAnyToSliceOfType[int](a)
}

// ToInt8Slice converts any type slice or array to []int8.
// E.g. covert []strs{"1", "2", "3"} to []int8{1, 2, 3}.
func SliceOfAnyToSliceOfInt8NoRet(a any) []int8 {
	return SliceOfAnyToSliceOfTypeNoRet[int8](a)
}

// ToInt8SliceE converts any type slice or array to []int8 with returned error.
func SliceOfAnyToSliceOfInt8(a any) ([]int8, error) {
	return SliceOfAnyToSliceOfType[int8](a)
}

// ToInt16Slice converts any type slice or array to []int16.
// For example, covert []strs{"1", "2", "3"} to []int16{1, 2, 3}.
func SliceOfAnyToSliceOfInt16NoRet(a any) []int16 {
	return SliceOfAnyToSliceOfTypeNoRet[int16](a)
}

// ToInt16SliceE converts any type slice or array to []int16 with returned error.
func SliceOfAnyToSliceOfInt16(a any) ([]int16, error) {
	return SliceOfAnyToSliceOfType[int16](a)
}

// ToInt32Slice converts any type slice or array to []int32.
// For example, covert []strs{"1", "2", "3"} to []int32{1, 2, 3}.
func SliceOfAnyToSliceOfInt32NoRet(a any) []int32 {
	return SliceOfAnyToSliceOfTypeNoRet[int32](a)
}

// ToInt32SliceE converts any type slice or array []int32 with returned error.
func SliceOfAnyToSliceOfInt32(a any) ([]int32, error) {
	return SliceOfAnyToSliceOfType[int32](a)
}

// ToInt64Slice converts any type slice or array to []int64 slice.
// For example, covert []strs{"1", "2", "3"} to []int64{1, 2, 3}.
func SliceOfAnyToSliceOfInt64NoRet(a any) []int64 {
	return SliceOfAnyToSliceOfTypeNoRet[int64](a)
}

// ToInt64SliceE converts any type slice or array to []int64 slice with returned error.
func SliceOfAnyToSliceOfInt64(a any) ([]int64, error) {
	return SliceOfAnyToSliceOfType[int64](a)
}

// ToUintSlice converts any type slice or array to []uint.
// For example, covert []strs{"1", "2", "3"} to []uint{1, 2, 3}.
func SliceOfAnyToSliceOfUIntNoRet(a any) []uint {
	return SliceOfAnyToSliceOfTypeNoRet[uint](a)
}

// ToUintSliceE converts any type slice or array to []uint with returned error.
func SliceOfAnyToSliceOfUInt(a any) ([]uint, error) {
	return SliceOfAnyToSliceOfType[uint](a)
}

// ToUint8Slice converts any type slice or array to []uint8.
// E.g. covert []strs{"1", "2", "3"} to []uint8{1, 2, 3}.
func SliceOfAnyToSliceOfUInt8NoRet(a any) []uint8 {
	return SliceOfAnyToSliceOfTypeNoRet[uint8](a)
}

// ToUint8SliceE converts any type slice or array to []uint8 slice with returned error.
func SliceOfAnyToSliceOfUInt8(a any) ([]uint8, error) {
	return SliceOfAnyToSliceOfType[uint8](a)
}

// ToByteSlice converts an type slice or array to []byte.
// E.g. covert []strs{"1", "2", "3"} to []byte{1, 2, 3}.
func SliceOfAnyToSliceOfByteNoRet(a any) []byte {
	return SliceOfAnyToSliceOfUInt8NoRet(a)
}

// ToByteSliceE converts any type slice or array to []byte with returned error.
func SliceOfAnyToSliceOfByte(a any) ([]byte, error) {
	return SliceOfAnyToSliceOfUInt8(a)
}

// ToUint16Slice converts any type slice or array to []uint16.
// For example, covert []strs{"1", "2", "3"} to []uint16{1, 2, 3}.
func SliceOfAnyToSliceOfUInt16NoRet(a any) []uint16 {
	return SliceOfAnyToSliceOfTypeNoRet[uint16](a)
}

// ToUint16SliceE converts any type slice or array to []uint16 slice with returned error.
func SliceOfAnyToSliceOfUInt16(a any) ([]uint16, error) {
	return SliceOfAnyToSliceOfType[uint16](a)
}

// ToUint32Slice converts any type slice or array to []uint32.
// For example, covert []strs{"1", "2", "3"} to []uint32{1, 2, 3}.
func SliceOfAnyToSliceOfUInt32NoRet(a any) []uint32 {
	return SliceOfAnyToSliceOfTypeNoRet[uint32](a)
}

// ToUint32SliceE converts any type slice or array to []uint32 slice with returned error.
func SliceOfAnyToSliceOfUInt32(a any) ([]uint32, error) {
	return SliceOfAnyToSliceOfType[uint32](a)
}

// ToUint64Slice converts any type slice or array to []uint64.
// For example, covert []strs{"1", "2", "3"} to []uint64{1, 2, 3}.
func SliceOfAnyToSliceOfUInt64NoRet(a any) []uint64 {
	return SliceOfAnyToSliceOfTypeNoRet[uint64](a)
}

// ToUint64SliceE converts any type slice or array to []uint64 with returned error.
func SliceOfAnyToSliceOfUInt64(a any) ([]uint64, error) {
	return SliceOfAnyToSliceOfType[uint64](a)
}

// ToDurationSlice converts any type slice or array to []time.Duration.
func SliceOfAnyToSliceOfDurationNoRet(a any) []time.Duration {
	return SliceOfAnyToSliceOfTypeNoRet[time.Duration](a)
}

// ToDurationSliceE converts any type to []time.Duration with returned error.
func SliceOfAnyToSliceOfDuration(a any) ([]time.Duration, error) {
	return SliceOfAnyToSliceOfType[time.Duration](a)
}

// ToStrSlice converts any type slice or array to []strs.
// For example, covert []int{1, 2, 3} to []strs{"1", "2", "3"}.
func SliceOfAnyToSliceOfStrNoRet(a any) []string {
	return SliceOfAnyToSliceOfTypeNoRet[string](a)
}

// ToStrSliceE converts any type slice or array to []strs with returned error.
func SliceOfAnyToSliceOfStr(a any) ([]string, error) {
	return SliceOfAnyToSliceOfType[string](a)
}
