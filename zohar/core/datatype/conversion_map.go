package datatype

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Struct2Map converts struct to map[strs]any.
// Such as struct{I int, S strs}{I: 1, S: "a"} to map[I:1 S:a].
// Note that unexported fields of struct can't be converted.
func StructToMap(a any) map[string]any {
	// Check param.
	v := reflect.ValueOf(a)
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := reflect.TypeOf(a)
	var m = make(map[string]any)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			m[t.Field(i).Name] = v.Field(i).Interface()
		}
	}
	return m
}

// Struct2MapStr converts struct to map[strs]strs.
// Such as struct{I int, S strs}{I: 1, S: "a"} to map[I:1 S:a].
// Note that unexported fields of struct can't be converted.
func StructToMapStr(obj any) map[string]string {
	// Check param.
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Struct {
		return nil
	}

	t := reflect.TypeOf(obj)
	var m = make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).IsExported() {
			m[t.Field(i).Name] = AnyToTypeNoRet[string](v.Field(i).Interface())
		}
	}
	return m
}

// ToMapStr converts any type to a map[strs]strs type.
func AnyToMapStrNoRet(a any) map[string]string {
	v, _ := AnyToMapStr(a)
	return v
}

// ToMapStrE converts any type to a map[strs]strs type.
func AnyToMapStr(a any) (map[string]string, error) {
	var m = map[string]string{}

	switch v := a.(type) {
	case map[string]string:
		return v, nil
	case map[string]any:
		for k, val := range v {
			val, err := AnyToType[string](val)
			if err != nil {
				return nil, err
			}
			m[k] = val
		}
	case map[any]string:
		for k, val := range v {
			k, err := AnyToType[string](k)
			if err != nil {
				return nil, err
			}
			m[k] = val
		}
	case map[any]any:
		for k, val := range v {
			k, err := AnyToType[string](k)
			if err != nil {
				return nil, err
			}
			val, err := AnyToType[string](val)
			if err != nil {
				return nil, err
			}
			m[k] = val
		}
	case string:
		if err := jsonStringToObject(v, &m); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unable to convert %#v of type %T to map[strs]strs", a, a)
	}
	return m, nil
}

// jsonStringToObject attempts to unmarshall a strs as JSON into
// the object passed as pointer.
func jsonStringToObject(s string, v any) error {
	data := []byte(s)
	return json.Unmarshal(data, v)
}
