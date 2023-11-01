package datatype

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"xeno/zohar/core"
)

const (
	DT_SINGLE = 0
	DT_LIST   = 1
	DT_DICT   = 2
	DT_COUNT  = 3
)

const (
	T_NULL  = 0
	T_I8    = 1
	T_I16   = 2
	T_I32   = 3
	T_I64   = 4
	T_U8    = 5
	T_U16   = 6
	T_U32   = 7
	T_U64   = 8
	T_BOOL  = 9
	T_F32   = 10
	T_F64   = 11
	T_BYTES = 12
	T_STR   = 13
	T_TLV   = 14
	T_COUNT = 15
)

var s_TypeSizeArr [T_COUNT]uint32 = [T_COUNT]uint32{0, 1, 2, 4, 8, 1, 2, 4, 8, 1, 4, 8, 0, 0, 0}

const TLV_INVALID_INDEX = 4294967295

var SINGLE_TYPE_NAMES = [T_COUNT]string{
	"T_NULL",
	"T_I8",
	"T_I16",
	"T_I32",
	"T_I64",
	"T_U8",
	"T_U16",
	"T_U32",
	"T_U64",
	"T_BOOL",
	"T_F32",
	"T_F64",
	"T_BYTES",
	"T_STR",
	"T_TLV",
}

type TLV struct {
	_type    uint8
	_keyType uint8
	_value   any
}

const (
	PST_ERR   = -1
	PST_EMPTY = 0
	PST_NAME  = 1
	PST_IDX   = 2
)

func parseStrTypeAndParts(s string) (int32, string, uint32) {
	if len(s) <= 0 {
		return PST_EMPTY, "", 0
	}
	start := strings.Index(s, "[")
	if start > -1 {
		ss := string([]rune(s)[strings.Index(s, "[")+1 : len(s)-1])
		if len(ss) < 1 {
			return PST_IDX, string([]rune(s)[0:start]), TLV_INVALID_INDEX
		} else {
			iret, err := strconv.Atoi(ss)
			if err != nil {
				return PST_ERR, "", 0
			}
			return PST_IDX, string([]rune(s)[0:start]), uint32(iret)
		}

	} else {
		return PST_NAME, s, 0
	}

}

func TLVValueAsType[T any](tlv *TLV) (T, error) {
	return AnyToType[T](tlv._value)
}

func (ego *TLV) Length() uint32 {
	ct, st := extractTlVType(ego._type)
	if ct == DT_SINGLE {
		l := s_TypeSizeArr[st]
		if l > 0 {
			return l
		} else if st == T_NULL {
			return 0
		} else if st == T_BYTES {
			return uint32(len(ego._value.([]byte)))
		} else if st == T_STR {
			return uint32(len(ego._value.(string)))
		} else {
			return TLV_INVALID_INDEX
		}

	} else if ct == DT_DICT {
		lfp := sLenMapHandlers[ego._keyType][st]
		if lfp == nil {
			return TLV_INVALID_INDEX
		}
		return uint32(lfp(ego._value))
	} else if ct == DT_LIST {
		lfp := sLenOfListHandlers[st]
		if lfp == nil {
			return TLV_INVALID_INDEX
		}
		return uint32(lfp(ego._value))
	} else {
		return TLV_INVALID_INDEX
	}
}

func (ego *TLV) Type() uint8 {
	return ego._type
}

func (ego *TLV) DimType() uint8 {
	ct, _ := extractTlVType(ego._type)
	return ct
}

func (ego *TLV) SingleType() uint8 {
	_, st := extractTlVType(ego._type)
	return st
}

func (ego *TLV) IsNumeric() bool {
	ct, st := extractTlVType(ego._type)
	if ct == DT_SINGLE {
		if st >= T_I8 && st <= T_F64 {
			return true
		}
	}
	return false
}

func (ego *TLV) CalcDictLength() (uint32, int32) {
	_, st := extractTlVType(ego._type)
	lfp := sLenMapHandlers[ego._keyType][st]
	if lfp == nil {
		return 0, core.MkErr(core.EC_NULL_VALUE, 1)
	}
	return uint32(lfp(ego._value)), core.MkSuccess(0)
}

func (ego *TLV) SetListValue(idx uint32, val any) int32 {
	dt, st := extractTlVType(ego._type)
	if dt == DT_LIST {
		if ego._value == nil {
			ego._value = sEmptyListCreationHandlers[st]()
		}

		if idx >= ego.Length() {
			return core.MkErr(core.EC_INDEX_OOB, 1)
		}

		sSetListValueHandlers[st](ego._value, idx, val)
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_TYPE_MISMATCH, 1)
}

func (ego *TLV) GetListValue(idx uint32) (any, int32) {
	dt, st := extractTlVType(ego._type)
	if dt == DT_LIST {
		if ego._value == nil {
			return nil, core.MkErr(core.EC_NULL_VALUE, 1)
		}

		if idx >= ego.Length() {
			return nil, core.MkErr(core.EC_INDEX_OOB, 1)
		}

		rv, rc := sGetListValueHandlers[st](ego, idx, false)
		return rv, rc

	}
	return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
}

func (ego *TLV) GetDictValue(key any) (any, int32) {
	dt, st := extractTlVType(ego._type)
	if dt == DT_DICT {
		fp := sGetMapValueHandlers[ego._keyType][st]
		if fp != nil {
			rv, ok := fp(ego._value, key)
			if ok {
				return rv, core.MkSuccess(0)
			}
			return nil, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
		}
	}
	return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
}

func (ego *TLV) KeyExist(key any) bool {
	dt, st := extractTlVType(ego._type)
	if dt == DT_DICT {
		fp := sGetMapValueHandlers[ego._keyType][st]
		if fp != nil {
			_, ok := fp(ego._value, key)
			if ok {
				return true
			}
		}
	}
	return false
}

func (ego *TLV) SetDictValue(key any, val any) int32 {
	dt, st := extractTlVType(ego._type)
	if dt == DT_DICT {
		if val == nil {
			cfp := createEmptyMapHandlers[ego._keyType][st]
			if cfp == nil {
				return core.MkErr(core.EC_NULL_VALUE, 1)
			}
			ego._value = cfp()

		}

		ifp := insertMapHandlers[ego._keyType][st]
		if ifp == nil {
			return core.MkErr(core.EC_NULL_VALUE, 2)
		}
		ifp(ego._value, key, val)

		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_TYPE_MISMATCH, 1)
}

func (ego *TLV) SetSingleValue(val any) int32 {
	dt, st := extractTlVType(ego._type)
	if dt == DT_SINGLE {
		rc, _, rtype, _, _, rval := parseSingleTlvParamsHandles[st](0, ego._type, val)
		if core.Err(rc) {
			return rc
		}
		ego._type = rtype
		ego._value = rval
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_TYPE_MISMATCH, 1)
}

func (ego *TLV) GetOrCreateTLVPath(path string) (*TLV, int32, string, uint32) {
	if len(path) <= 0 {
		return ego, PST_EMPTY, "", 0
	}
	var current *TLV = ego
	path = strings.Trim(path, " \t\n\r")
	ss := strings.Split(path, ".")
	lastIndex := len(ss) - 1
	lastTyp, lastMain, lastIdx := parseStrTypeAndParts(ss[lastIndex])
	for i := 0; i < len(ss)-1; i++ {
		typ, main, idx := parseStrTypeAndParts(ss[i])
		if typ == PST_ERR || typ == PST_EMPTY {
			return nil, lastTyp, lastMain, lastIdx
		}

		if typ == PST_IDX {
			rTLV, _ := current.GetDictValue(main)
			if rTLV != nil {
				if rTLV.(*TLV).DimType() != DT_LIST {
					return nil, lastTyp, lastMain, lastIdx
				}
				if idx != TLV_INVALID_INDEX {
					rv, rc := sGetListValueHandlers[T_TLV](rTLV.(*TLV), idx, false)
					if core.Err(rc) {
						return nil, lastTyp, lastMain, lastIdx
					}
					current = rv.(*TLV)
				} else {
					rv := CreateTLV(DT_DICT, T_TLV, T_STR, nil)
					rc := rTLV.(*TLV).PushBack(rv)
					if core.Err(rc) {
						return nil, lastTyp, lastMain, lastIdx
					}
					current = rv
				}
			} else {
				rTLV = CreateTLV(DT_LIST, T_TLV, T_STR, nil)
				if current.DimType() == DT_LIST {
					panic("can not reach here")
				} else if current.DimType() == DT_DICT {
					rc := current.SetDictValue(main, rTLV)
					if core.Err(rc) {
						return nil, lastTyp, lastMain, lastIdx
					}
				}

				rv := CreateTLV(DT_DICT, T_TLV, T_STR, nil)
				rc := rTLV.(*TLV).PushBack(rv)
				if core.Err(rc) {
					return nil, lastTyp, lastMain, lastIdx
				}
				current = rv
			}
		} else if typ == PST_NAME {
			rTLV, _ := current.GetDictValue(main)
			if rTLV != nil {
				current = rTLV.(*TLV)
			} else {
				rTLV = CreateTLV(DT_DICT, T_TLV, T_STR, nil)
				rc := current.SetDictValue(main, rTLV)
				if core.Err(rc) {
					return nil, lastTyp, lastMain, lastIdx
				}
				current = rTLV.(*TLV)
			}
		}
	} // end of for

	return current, lastTyp, lastMain, lastIdx

}

func (ego *TLV) GetTLVPath(path string) *TLV {
	if len(path) <= 0 {
		return ego
	}
	var current *TLV = ego
	path = strings.Trim(path, " \t\n\r")
	ss := strings.Split(path, ".")
	for i := 0; i < len(ss); i++ {
		typ, main, idx := parseStrTypeAndParts(ss[i])
		if typ == PST_ERR || typ == PST_EMPTY {
			return nil
		}

		if typ == PST_IDX {
			rTLV, _ := current.GetDictValue(main)
			if rTLV != nil {
				if rTLV.(*TLV).DimType() != DT_LIST {
					return nil
				}
				if idx >= rTLV.(*TLV).Length() {
					return nil
				}
				rv, rc := sGetListValueHandlers[T_TLV](rTLV.(*TLV), idx, false)
				if core.Err(rc) {
					return nil
				}
				current = rv.(*TLV)
			} else {
				return nil
			}
		} else if typ == PST_NAME {
			rTLV, _ := current.GetDictValue(main)
			if rTLV != nil {
				current = rTLV.(*TLV)
			} else {
				return nil
			}
		}
	} // end of for

	return current

}

func (ego *TLV) PathGet(path string) (any, int32) {
	if len(path) <= 0 {
		return ego._value, core.MkSuccess(0)
	}

	parent := ego.GetTLVPath(path)
	if parent == nil {
		return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
	}

	return parent._value, core.MkSuccess(0)
}

func (ego *TLV) PathSet(path string, val any, dt uint8, st uint8, kt uint8) int32 {
	//var tlv *TLV = nil
	if len(path) <= 0 {
		return core.MkErr(core.EC_TYPE_MISMATCH, 1)
	}

	parent, lastTyp, lastMain, lastIdx := ego.GetOrCreateTLVPath(path)
	if parent == nil {
		return core.MkErr(core.EC_TYPE_MISMATCH, 1)
	}

	tlv := CreateTLV(dt, st, kt, val)
	if tlv == nil {
		return core.MkErr(core.EC_TYPE_MISMATCH, 2)
	}

	if lastTyp == PST_NAME {
		parent.SetDictValue(lastMain, tlv)
	} else if lastTyp == PST_IDX {
		rTLV, _ := parent.GetDictValue(lastMain)
		if rTLV != nil {
			if rTLV.(*TLV).DimType() != DT_LIST {
				return core.MkErr(core.EC_TYPE_MISMATCH, 3)
			}
			if lastIdx != TLV_INVALID_INDEX {
				if lastIdx >= rTLV.(*TLV).Length() {
					return core.MkErr(core.EC_INDEX_OOB, 3)
				}
				rv, rc := sGetListValueHandlers[T_TLV](rTLV.(*TLV), lastIdx, false)
				if core.Err(rc) {
					return core.MkErr(core.EC_INDEX_OOB, 3)
				}
				rv.(*TLV).SetSingleValue(val)
			} else {
				tlv := CreateTLV(dt, st, kt, val)
				if tlv == nil {
					return core.MkErr(core.EC_TYPE_MISMATCH, 20)
				}
				rTLV.(*TLV).PushBack(tlv)
			}
		} else {
			rTLV = CreateTLV(DT_LIST, T_TLV, T_STR, nil)
			if parent.DimType() == DT_LIST {
				panic("can not reach here")
			} else if parent.DimType() == DT_DICT {
				rc := parent.SetDictValue(lastMain, rTLV)
				if core.Err(rc) {
					return core.MkErr(core.EC_TYPE_MISMATCH, 30)
				}
			}

			rv := CreateTLV(dt, st, kt, val)
			rc := rTLV.(*TLV).PushBack(rv)
			if core.Err(rc) {
				return core.MkErr(core.EC_TYPE_MISMATCH, 30)
			}
		}
	}

	return core.MkSuccess(0)
}

func (ego *TLV) stringOfSingleType(indent int) string {
	var sb strings.Builder
	for i := 0; i < indent; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("{\"t\":")
	sb.WriteString(strconv.Itoa(int(ego._type)))

	sb.WriteString(", \"k\":")
	sb.WriteString(strconv.Itoa(int(ego._keyType)))

	sb.WriteString(", \"l\":")
	sb.WriteString(strconv.Itoa(int(ego.Length())))

	if ego.IsNumeric() {
		sb.WriteString(", \"v\":")
		sb.WriteString(ego.AsStringNoRet())
		sb.WriteString("}")
	} else {
		sb.WriteString(", \"v\":\"")
		sb.WriteString(ego.AsStringNoRet())
		sb.WriteString("\"}")
	}

	return sb.String()
}

func (ego *TLV) stringOfDictType(indent int) string {
	_, st := extractTlVType(ego._type)

	var sb strings.Builder
	for i := 0; i < indent; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("{\"t\":")
	sb.WriteString(strconv.Itoa(int(ego._type)))

	sb.WriteString(", \"k\":")
	sb.WriteString(strconv.Itoa(int(ego._keyType)))

	sb.WriteString(", \"l\":")
	sb.WriteString(strconv.Itoa(int(ego.Length())))

	sb.WriteString(", \"v\":")
	sb.WriteString("{")

	if ego._value != nil {
		kind := reflect.TypeOf(ego._value).Kind()
		switch kind {
		case reflect.Map:
			// If input is a slice or array.
			v := reflect.ValueOf(ego._value)
			iter := v.MapRange()
			var idx int32 = 0
			for iter.Next() {
				if idx > 0 {
					sb.WriteString(",")
				}

				sb.WriteString("\"")
				sb.WriteString(fmt.Sprint(iter.Key()))
				sb.WriteString("\":")

				if st == T_NULL {
					sb.WriteString("null")
				} else if st == T_STR {
					sb.WriteString("\"")
					sb.WriteString(fmt.Sprint(iter.Value()))
					sb.WriteString("\"")
				} else if st == T_BYTES {
					sb.WriteString("\"")
					sb.WriteString(fmt.Sprint(iter.Value()))
					sb.WriteString("\"")
				} else if st == T_TLV {
					sb.WriteString(fmt.Sprint(iter.Value()))
				} else {
					sb.WriteString(fmt.Sprint(iter.Value()))
				}
				idx++
			}
		}
	}

	sb.WriteString("}")
	sb.WriteString("}")

	return sb.String()
}

func (ego *TLV) stringOfListType(indent int) string {
	var sb strings.Builder
	for i := 0; i < indent; i++ {
		sb.WriteString("  ")
	}
	sb.WriteString("{\"t\":")
	sb.WriteString(strconv.Itoa(int(ego._type)))

	sb.WriteString(", \"k\":")
	sb.WriteString(strconv.Itoa(int(ego._keyType)))

	sb.WriteString(", \"l\":")
	sb.WriteString(strconv.Itoa(int(ego.Length())))

	sb.WriteString(", \"v\":")
	sb.WriteString("[")
	kind := reflect.TypeOf(ego._value).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		// If input is a slice or array.
		v := reflect.ValueOf(ego._value)
		if kind == reflect.Slice && v.IsNil() {
			sb.WriteString("\"null\"")
		} else {
			for i := 0; i < v.Len(); i++ {
				rr := v.Index(i).Interface()
				val, err := AnyToType[string](rr)
				if err != nil {
					panic("convert failed")
				}

				_, st := extractTlVType(ego._type)
				if st == T_BYTES || st == T_STR {
					sb.WriteString("\"")
					sb.WriteString(val)
					sb.WriteString("\"")
				} else if st == T_NULL {
					sb.WriteString("\"null\"")
				} else if st == T_TLV {
					sb.WriteString(val)
				} else {
					sb.WriteString(val)
				}

				if i < v.Len()-1 {
					sb.WriteString(",")
				}
			}
		}

	default:
		panic("got a signle type")
	}

	sb.WriteString("]")
	sb.WriteString("}")

	return sb.String()
}

func (ego *TLV) AsStringNoRet() string {
	str, _ := ego.AsString()
	return str
}

func (ego *TLV) Value() any {
	return ego._value
}

func (ego *TLV) AsString() (string, int32) {
	str, err := AnyToString(ego._value)
	if err != nil {
		return "", core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return str, core.MkSuccess(0)
}

func (ego *TLV) AsInt8() (int8, int32) {
	i, err := AnyToInt8(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsInt16() (int16, int32) {
	i, err := AnyToInt16(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsInt32() (int32, int32) {
	i, err := AnyToInt32(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsInt64() (int64, int32) {
	i, err := AnyToInt64(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsUInt8() (uint8, int32) {
	i, err := AnyToUInt8(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsUInt16() (uint16, int32) {
	i, err := AnyToUInt16(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsUInt32() (uint32, int32) {
	i, err := AnyToUInt32(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) AsUInt64() (uint64, int32) {
	i, err := AnyToUInt64(ego._value)
	if err != nil {
		return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return i, core.MkSuccess(0)
}

func (ego *TLV) I8At(idx uint32) (int8, int32) {
	v, rc := ego.At(idx)
	return v.(int8), rc
}

func (ego *TLV) MapInsert(key any, val any) int32 {
	ct, st := extractTlVType(ego._type)
	if ct == DT_DICT {
		h := insertMapHandlers[ego._keyType][st]
		if h == nil {
			return core.MkErr(core.EC_HANDLER_NOT_FOUND, 1)
		}
		h(ego._value, key, val)
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_TYPE_MISMATCH, 1)
}

func (ego *TLV) PushBack(neoVal any) int32 {
	ct, st := extractTlVType(ego._type)
	if ct == DT_LIST {
		if st == T_NULL {
			ego._value = append(ego._value.([]any), nil)
			return core.MkSuccess(0)
		}
		if st == T_I8 {
			ego._value = append(ego._value.([]int8), neoVal.(int8))
		} else if st == T_I16 {
			ego._value = append(ego._value.([]int16), neoVal.(int16))
		} else if st == T_I32 {
			ego._value = append(ego._value.([]int32), neoVal.(int32))
		} else if st == T_I64 {
			ego._value = append(ego._value.([]int64), neoVal.(int64))
		} else if st == T_U8 {
			ego._value = append(ego._value.([]uint8), neoVal.(uint8))
		} else if st == T_U16 {
			ego._value = append(ego._value.([]uint16), neoVal.(uint16))
		} else if st == T_U32 {
			ego._value = append(ego._value.([]uint32), neoVal.(uint32))
		} else if st == T_U64 {
			ego._value = append(ego._value.([]uint64), neoVal.(uint64))
		} else if st == T_BOOL {
			ego._value = append(ego._value.([]bool), neoVal.(bool))
		} else if st == T_F32 {
			ego._value = append(ego._value.([]float32), neoVal.(float32))
		} else if st == T_F64 {
			ego._value = append(ego._value.([]float64), neoVal.(float64))
		} else if st == T_BYTES {
			ego._value = append(ego._value.([][]byte), neoVal.([]byte))
		} else if st == T_STR {
			ego._value = append(ego._value.([]string), neoVal.(string))
		} else if st == T_TLV {
			ego._value = append(ego._value.([]*TLV), neoVal.(*TLV))
		} else {
			return core.MkErr(core.EC_TYPE_CONVERT_FAILED, 2)
		}
	}
	return core.MkSuccess(0)
}

func (ego *TLV) At(idx uint32) (any, int32) {
	ct, st := extractTlVType(ego._type)
	if st >= T_COUNT {
		return nil, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	if ct == DT_LIST {
		if st == T_NULL {
			return nil, core.MkSuccess(0)
		}
		v := reflect.ValueOf(ego._value)
		if idx >= ego.Length() {
			return nil, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
		}
		if st == T_I8 {
			rv, err := AnyToInt8(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_I16 {
			rv, err := AnyToInt16(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 2)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_I32 {
			rv, err := AnyToInt32(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 3)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_I64 {
			rv, err := AnyToInt64(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 4)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_U8 {
			rv, err := AnyToUInt8(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 5)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_U16 {
			rv, err := AnyToUInt16(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 6)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_U32 {
			rv, err := AnyToUInt32(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 7)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_U64 {
			rv, err := AnyToUInt64(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 8)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_BOOL {
			rv, err := AnyToBool(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 9)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_F32 {
			rv, err := AnyToFloat32(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 10)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_F64 {
			rv, err := AnyToFloat64(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 11)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_BYTES {
			rv, err := AnyToByteArray(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 12)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_STR {
			rv, err := AnyToString(v.Index(int(idx)).Interface())
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 13)
			}
			return rv, core.MkSuccess(0)
		} else if st == T_TLV {
			a := v.Index(int(idx)).Interface()
			rv, err := AnyToType[TLV](a)
			if err != nil {
				return 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 13)
			}
			return &rv, core.MkSuccess(0)
		}

	}
	return nil, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 2)
}

func (ego *TLV) String() string {
	ct, st := extractTlVType(ego._type)
	if st >= T_COUNT {
		panic("Invalid TLV Type " + string(ego._type))
	}
	if ct == DT_SINGLE {
		return ego.stringOfSingleType(1)
	} else if ct == DT_LIST {
		return ego.stringOfListType(1)
	} else if ct == DT_DICT {
		return ego.stringOfDictType(1)
	}

	panic("Invalid TLV Type " + string(ego._type))
}

func parseSingleTlvParamsHandleOfNull(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	return core.EC_OK, 0, T_NULL, 0, 0, nil
}

func parseSingleTlvParamsHandleOfI8(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i8, err := AnyToInt8(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 1, T_I8, 0, 0, val
	}
	return core.EC_OK, 1, T_I8, 0, 0, i8
}

func parseSingleTlvParamsHandleOfI16(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToInt16(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 2, T_I16, 0, 0, val
	}
	return core.EC_OK, 2, T_I16, 0, 0, i
}

func parseSingleTlvParamsHandleOfI32(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToInt32(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 4, T_I32, 0, 0, val
	}
	return core.EC_OK, 4, T_I32, 0, 0, i
}

func parseSingleTlvParamsHandleOfI64(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToInt64(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 8, T_I64, 0, 0, val
	}
	return core.EC_OK, 8, T_I64, 0, 0, i
}

func parseSingleTlvParamsHandleOfU8(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i8, err := AnyToUInt8(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 1, T_U8, 0, 0, val
	}
	return core.EC_OK, 1, T_U8, 0, 0, i8
}

func parseSingleTlvParamsHandleOfU16(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToUInt16(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 2, T_U16, 0, 0, val
	}
	return core.EC_OK, 2, T_U16, 0, 0, i
}

func parseSingleTlvParamsHandleOfU32(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToUInt32(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 4, T_U32, 0, 0, val
	}
	return core.EC_OK, 4, T_U32, 0, 0, i
}

func parseSingleTlvParamsHandleOfU64(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToUInt64(val)
	if err != nil {
		return core.MkErr(core.EC_NOOP, 1), 8, T_U64, 0, 0, val
	}
	return core.EC_OK, 8, T_U64, 0, 0, i
}

func parseSingleTlvParamsHandleOfBool(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToBool(val)
	if err != nil {
		return core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1), 1, T_BOOL, 0, 0, val
	}
	return core.EC_OK, 1, T_BOOL, 0, 0, i
}

func parseSingleTlvParamsHandleOfF32(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToFloat32(val)
	if err != nil {
		return core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1), 4, T_F32, 0, 0, val
	}
	return core.EC_OK, 4, T_F32, 0, 0, i
}

func parseSingleTlvParamsHandleOfF64(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToFloat32(val)
	if err != nil {
		return core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1), 8, T_F64, 0, 0, val
	}
	return core.EC_OK, 8, T_F64, 0, 0, i
}

func parseSingleTlvParamsHandleOfBytes(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToByteArray(val)
	if err != nil {
		return core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1), 0, T_BYTES, 0, 0, val
	}
	return core.EC_OK, uint32(len(*i)), T_BYTES, 0, 0, *i
}

func parseSingleTlvParamsHandleOfStr(length uint32, t uint8, val any) (int32, uint32, uint8, uint8, uint16, any) {
	i, err := AnyToString(val)
	if err != nil {
		return core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1), 0, T_STR, 0, 0, ""
	}
	return core.EC_OK, uint32(len(i)), T_STR, 0, 0, i
}

var parseSingleTlvParamsHandles = [T_TLV]func(uint32, uint8, any) (int32, uint32, uint8, uint8, uint16, any){
	parseSingleTlvParamsHandleOfNull,
	parseSingleTlvParamsHandleOfI8,
	parseSingleTlvParamsHandleOfI16,
	parseSingleTlvParamsHandleOfI32,
	parseSingleTlvParamsHandleOfI64,
	parseSingleTlvParamsHandleOfU8,
	parseSingleTlvParamsHandleOfU16,
	parseSingleTlvParamsHandleOfU32,
	parseSingleTlvParamsHandleOfU64,
	parseSingleTlvParamsHandleOfBool,
	parseSingleTlvParamsHandleOfF32,
	parseSingleTlvParamsHandleOfF64,
	parseSingleTlvParamsHandleOfBytes,
	parseSingleTlvParamsHandleOfStr,
}

var sLenOfListHandlers = [T_COUNT]func(v any) uint32{
	func(v any) uint32 {
		return uint32(len(v.([]any)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]int8)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]int16)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]int32)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]int64)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]uint8)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]uint16)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]uint32)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]uint64)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]bool)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]float32)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]float64)))
	},
	func(v any) uint32 {
		return uint32(len(v.([][]byte)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]string)))
	},
	func(v any) uint32 {
		return uint32(len(v.([]*TLV)))
	},
}

var sDefaultListOfAnyHandlers = [T_COUNT]func(uint32, uint8) *TLV{
	func(l uint32, t uint8) *TLV { //null
		rval := make([]any, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //i8
		rval := make([]int8, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //i16
		rval := make([]int16, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //i32
		rval := make([]int32, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //i64
		rval := make([]int64, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //u8
		rval := make([]uint8, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //u16
		rval := make([]uint16, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //u32
		rval := make([]uint32, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //u64
		rval := make([]uint64, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //bool
		rval := make([]bool, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //f32
		rval := make([]float32, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //f64
		rval := make([]float64, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //bytes
		rval := make([][]byte, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //strs
		rval := make([]string, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
	func(l uint32, t uint8) *TLV { //TLV
		rval := make([]*TLV, 0)
		return &TLV{
			_type:    t,
			_keyType: 0,
			_value:   rval,
		}
	},
}

var sGetListValueHandlers = [T_COUNT]func(tlv *TLV, idx uint32, force bool) (any, int32){
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]any)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]any), nil)
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]any)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]int8)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]int8), int8(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]int8)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]int8)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]int16)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]int16), int16(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]int16)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]int16)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]int32)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]int32), int32(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]int32)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]int32)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]int64)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]int64), int64(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]int64)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]int64)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]uint8)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]uint8), uint8(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]uint8)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]uint8)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]uint16)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]uint16), uint16(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]uint16)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]uint16)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]uint32)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]uint32), uint32(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]uint32)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]uint32)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]uint64)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]uint64), uint64(0))
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]uint64)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}

		}
		return tlv._value.([]uint64)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]bool)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]bool), false)
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]bool)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}
		}
		return tlv._value.([]bool)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]float32)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]float32), 0.0)
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]float32)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}
		}
		return tlv._value.([]float32)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]float64)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]float64), 0.0)
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]float64)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}
		}
		return tlv._value.([]float64)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([][]byte)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([][]byte), nil)
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([][]byte)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}
		}
		return tlv._value.([][]byte)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]string)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]string), "")
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]string)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}
		}
		return tlv._value.([]string)[int(idx)], core.MkSuccess(0)
	},
	func(tlv *TLV, idx uint32, force bool) (any, int32) {
		if len(tlv._value.([]*TLV)) <= int(idx) {
			if force {
				tlv._value = append(tlv._value.([]*TLV), nil)
				if tlv._value == nil {
					return nil, core.MkErr(core.EC_NULL_VALUE, 1)
				}
				l := uint32(len(tlv._value.([]*TLV)))
				idx = l - 1
			} else {
				return nil, core.MkErr(core.EC_INDEX_OOB, 1)
			}
		}
		return tlv._value.([]*TLV)[int(idx)], core.MkSuccess(0)
	},
}

var sSetListValueHandlers = [T_COUNT]func(lst any, idx uint32, val any){
	func(lst any, idx uint32, val any) { //null
		lst.([]any)[idx] = nil
	},
	func(lst any, idx uint32, val any) { //i8
		lst.([]int8)[idx] = val.(int8)
	},
	func(lst any, idx uint32, val any) { //i16
		lst.([]int16)[idx] = val.(int16)
	},
	func(lst any, idx uint32, val any) { //i32
		lst.([]int32)[idx] = val.(int32)
	},
	func(lst any, idx uint32, val any) { //i64
		lst.([]int64)[idx] = val.(int64)
	},
	func(lst any, idx uint32, val any) { //i8
		lst.([]uint8)[idx] = val.(uint8)
	},
	func(lst any, idx uint32, val any) { //i16
		lst.([]uint16)[idx] = val.(uint16)
	},
	func(lst any, idx uint32, val any) { //i32
		lst.([]uint32)[idx] = val.(uint32)
	},
	func(lst any, idx uint32, val any) { //i64
		lst.([]uint64)[idx] = val.(uint64)
	},
	func(lst any, idx uint32, val any) { //bool
		lst.([]bool)[idx] = val.(bool)
	},
	func(lst any, idx uint32, val any) { //32
		lst.([]float32)[idx] = val.(float32)
	},
	func(lst any, idx uint32, val any) { //64
		lst.([]float64)[idx] = val.(float64)
	},
	func(lst any, idx uint32, val any) { //[]byte
		lst.([][]byte)[idx] = val.([]byte)
	},
	func(lst any, idx uint32, val any) { //str
		lst.([]string)[idx] = val.(string)
	},
	func(lst any, idx uint32, val any) { //64
		lst.([]float64)[idx] = val.(float64)
	},
}

var sEmptyListCreationHandlers = [T_COUNT]func() any{
	func() any { //null
		return make([]any, 0)
	},
	func() any { //8
		return make([]int8, 0)
	},
	func() any { //16
		return make([]int16, 0)
	},
	func() any { //32
		return make([]int32, 0)
	},
	func() any { //64
		return make([]int64, 0)
	},
	func() any { //8
		return make([]uint8, 0)
	},
	func() any { //16
		return make([]uint16, 0)
	},
	func() any { //32
		return make([]uint32, 0)
	},
	func() any { //64
		return make([]uint64, 0)
	},
	func() any { //bool
		return make([]bool, 0)
	},
	func() any { //32
		return make([]float32, 0)
	},
	func() any { //64
		return make([]float64, 0)
	},
	func() any { //bytes
		return make([][]byte, 0)
	},
	func() any { //bytes
		return make([]string, 0)
	},
	func() any { //bytes
		return make([]*TLV, 0)
	},
}

var ConvertListOfAnyHandlers = [T_COUNT]func(uint8, any, uint32) (any, uint32, int32){
	convertListOfAnyHandlerOfNull,
	convertListOfAnyHandlerOfI8,
	convertListOfAnyHandlerOfI16,
	convertListOfAnyHandlerOfI32,
	convertListOfAnyHandlerOfI64,
	convertListOfAnyHandlerOfU8,
	convertListOfAnyHandlerOfU16,
	convertListOfAnyHandlerOfU32,
	convertListOfAnyHandlerOfU64,
	convertListOfAnyHandlerOfBool,
	convertListOfAnyHandlerOfF32,
	convertListOfAnyHandlerOfF64,
	convertListOfAnyHandlerOfBytes,
	convertListOfAnyHandlerOfStr,
	convertListOfAnyHandlerOfTLV,
}

var createEmptyMapHandlers = [T_COUNT][T_COUNT]func() any{
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{func() any { return make(map[int8]any) }, func() any { return make(map[int8]int8) }, func() any { return make(map[int8]int16) }, func() any { return make(map[int8]int32) }, func() any { return make(map[int8]int64) }, func() any { return make(map[int8]uint8) }, func() any { return make(map[int8]uint16) }, func() any { return make(map[int8]uint32) }, func() any { return make(map[int8]uint64) }, func() any { return make(map[int8]bool) }, func() any { return make(map[int8]float32) }, func() any { return make(map[int8]float64) }, func() any { return make(map[int8][]byte) }, func() any { return make(map[int8]string) }, func() any { return make(map[int8]*TLV) }},
	{func() any { return make(map[int16]any) }, func() any { return make(map[int16]int8) }, func() any { return make(map[int16]int16) }, func() any { return make(map[int16]int32) }, func() any { return make(map[int16]int64) }, func() any { return make(map[int16]uint8) }, func() any { return make(map[int16]uint16) }, func() any { return make(map[int16]uint32) }, func() any { return make(map[int16]uint64) }, func() any { return make(map[int16]bool) }, func() any { return make(map[int16]float32) }, func() any { return make(map[int16]float64) }, func() any { return make(map[int16][]byte) }, func() any { return make(map[int16]string) }, func() any { return make(map[int16]*TLV) }},
	{func() any { return make(map[int32]any) }, func() any { return make(map[int32]int8) }, func() any { return make(map[int32]int16) }, func() any { return make(map[int32]int32) }, func() any { return make(map[int32]int64) }, func() any { return make(map[int32]uint8) }, func() any { return make(map[int32]uint16) }, func() any { return make(map[int32]uint32) }, func() any { return make(map[int32]uint64) }, func() any { return make(map[int32]bool) }, func() any { return make(map[int32]float32) }, func() any { return make(map[int32]float64) }, func() any { return make(map[int32][]byte) }, func() any { return make(map[int32]string) }, func() any { return make(map[int32]*TLV) }},
	{func() any { return make(map[int64]any) }, func() any { return make(map[int64]int8) }, func() any { return make(map[int64]int16) }, func() any { return make(map[int64]int32) }, func() any { return make(map[int64]int64) }, func() any { return make(map[int64]uint8) }, func() any { return make(map[int64]uint16) }, func() any { return make(map[int64]uint32) }, func() any { return make(map[int64]uint64) }, func() any { return make(map[int64]bool) }, func() any { return make(map[int64]float32) }, func() any { return make(map[int64]float64) }, func() any { return make(map[int64][]byte) }, func() any { return make(map[int64]string) }, func() any { return make(map[int64]*TLV) }},
	{func() any { return make(map[uint8]any) }, func() any { return make(map[uint8]int8) }, func() any { return make(map[uint8]int16) }, func() any { return make(map[uint8]int32) }, func() any { return make(map[uint8]int64) }, func() any { return make(map[uint8]uint8) }, func() any { return make(map[uint8]uint16) }, func() any { return make(map[uint8]uint32) }, func() any { return make(map[uint8]uint64) }, func() any { return make(map[uint8]bool) }, func() any { return make(map[uint8]float32) }, func() any { return make(map[uint8]float64) }, func() any { return make(map[uint8][]byte) }, func() any { return make(map[uint8]string) }, func() any { return make(map[uint8]*TLV) }},
	{func() any { return make(map[uint16]any) }, func() any { return make(map[uint16]int8) }, func() any { return make(map[uint16]int16) }, func() any { return make(map[uint16]int32) }, func() any { return make(map[uint16]int64) }, func() any { return make(map[uint16]uint8) }, func() any { return make(map[uint16]uint16) }, func() any { return make(map[uint16]uint32) }, func() any { return make(map[uint16]uint64) }, func() any { return make(map[uint16]bool) }, func() any { return make(map[uint16]float32) }, func() any { return make(map[uint16]float64) }, func() any { return make(map[uint16][]byte) }, func() any { return make(map[uint16]string) }, func() any { return make(map[uint16]*TLV) }},
	{func() any { return make(map[uint32]any) }, func() any { return make(map[uint32]int8) }, func() any { return make(map[uint32]int16) }, func() any { return make(map[uint32]int32) }, func() any { return make(map[uint32]int64) }, func() any { return make(map[uint32]uint8) }, func() any { return make(map[uint32]uint16) }, func() any { return make(map[uint32]uint32) }, func() any { return make(map[uint32]uint64) }, func() any { return make(map[uint32]bool) }, func() any { return make(map[uint32]float32) }, func() any { return make(map[uint32]float64) }, func() any { return make(map[uint32][]byte) }, func() any { return make(map[uint32]string) }, func() any { return make(map[uint32]*TLV) }},
	{func() any { return make(map[uint64]any) }, func() any { return make(map[uint64]int8) }, func() any { return make(map[uint64]int16) }, func() any { return make(map[uint64]int32) }, func() any { return make(map[uint64]int64) }, func() any { return make(map[uint64]uint8) }, func() any { return make(map[uint64]uint16) }, func() any { return make(map[uint64]uint32) }, func() any { return make(map[uint64]uint64) }, func() any { return make(map[uint64]bool) }, func() any { return make(map[uint64]float32) }, func() any { return make(map[uint64]float64) }, func() any { return make(map[uint64][]byte) }, func() any { return make(map[uint64]string) }, func() any { return make(map[uint64]*TLV) }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}, //bool
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}, //f32
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}, //f64
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}, //bytes
	{func() any { return make(map[string]any) }, func() any { return make(map[string]int8) }, func() any { return make(map[string]int16) }, func() any { return make(map[string]int32) }, func() any { return make(map[string]int64) }, func() any { return make(map[string]uint8) }, func() any { return make(map[string]uint16) }, func() any { return make(map[string]uint32) }, func() any { return make(map[string]uint64) }, func() any { return make(map[string]bool) }, func() any { return make(map[string]float32) }, func() any { return make(map[string]float64) }, func() any { return make(map[string][]byte) }, func() any { return make(map[string]string) }, func() any { return make(map[string]*TLV) }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

var insertMapHandlers = [T_COUNT][T_COUNT]func(m any, k any, v any){
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{func(m any, k any, v any) { m.(map[int8]any)[k.(int8)] = nil }, func(m any, k any, v any) { m.(map[int8]int8)[k.(int8)] = v.(int8) }, func(m any, k any, v any) { m.(map[int8]int16)[k.(int8)] = v.(int16) }, func(m any, k any, v any) { m.(map[int8]int32)[k.(int8)] = v.(int32) }, func(m any, k any, v any) { m.(map[int8]int64)[k.(int8)] = v.(int64) }, func(m any, k any, v any) { m.(map[int8]uint8)[k.(int8)] = v.(uint8) }, func(m any, k any, v any) { m.(map[int8]uint16)[k.(int8)] = v.(uint16) }, func(m any, k any, v any) { m.(map[int8]uint32)[k.(int8)] = v.(uint32) }, func(m any, k any, v any) { m.(map[int8]uint64)[k.(int8)] = v.(uint64) }, func(m any, k any, v any) { m.(map[int8]bool)[k.(int8)] = v.(bool) }, func(m any, k any, v any) { m.(map[int8]float32)[k.(int8)] = v.(float32) }, func(m any, k any, v any) { m.(map[int8]float64)[k.(int8)] = v.(float64) }, func(m any, k any, v any) { m.(map[int8][]byte)[k.(int8)] = v.([]byte) }, func(m any, k any, v any) { m.(map[int8]string)[k.(int8)] = v.(string) }, func(m any, k any, v any) { m.(map[int8]*TLV)[k.(int8)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[int16]any)[k.(int16)] = nil }, func(m any, k any, v any) { m.(map[int16]int8)[k.(int16)] = v.(int8) }, func(m any, k any, v any) { m.(map[int16]int16)[k.(int16)] = v.(int16) }, func(m any, k any, v any) { m.(map[int16]int32)[k.(int16)] = v.(int32) }, func(m any, k any, v any) { m.(map[int16]int64)[k.(int16)] = v.(int64) }, func(m any, k any, v any) { m.(map[int16]uint8)[k.(int16)] = v.(uint8) }, func(m any, k any, v any) { m.(map[int16]uint16)[k.(int16)] = v.(uint16) }, func(m any, k any, v any) { m.(map[int16]uint32)[k.(int16)] = v.(uint32) }, func(m any, k any, v any) { m.(map[int16]uint64)[k.(int16)] = v.(uint64) }, func(m any, k any, v any) { m.(map[int16]bool)[k.(int16)] = v.(bool) }, func(m any, k any, v any) { m.(map[int16]float32)[k.(int16)] = v.(float32) }, func(m any, k any, v any) { m.(map[int16]float64)[k.(int16)] = v.(float64) }, func(m any, k any, v any) { m.(map[int16][]byte)[k.(int16)] = v.([]byte) }, func(m any, k any, v any) { m.(map[int16]string)[k.(int16)] = v.(string) }, func(m any, k any, v any) { m.(map[int16]*TLV)[k.(int16)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[int32]any)[k.(int32)] = nil }, func(m any, k any, v any) { m.(map[int32]int8)[k.(int32)] = v.(int8) }, func(m any, k any, v any) { m.(map[int32]int16)[k.(int32)] = v.(int16) }, func(m any, k any, v any) { m.(map[int32]int32)[k.(int32)] = v.(int32) }, func(m any, k any, v any) { m.(map[int32]int64)[k.(int32)] = v.(int64) }, func(m any, k any, v any) { m.(map[int32]uint8)[k.(int32)] = v.(uint8) }, func(m any, k any, v any) { m.(map[int32]uint16)[k.(int32)] = v.(uint16) }, func(m any, k any, v any) { m.(map[int32]uint32)[k.(int32)] = v.(uint32) }, func(m any, k any, v any) { m.(map[int32]uint64)[k.(int32)] = v.(uint64) }, func(m any, k any, v any) { m.(map[int32]bool)[k.(int32)] = v.(bool) }, func(m any, k any, v any) { m.(map[int32]float32)[k.(int32)] = v.(float32) }, func(m any, k any, v any) { m.(map[int32]float64)[k.(int32)] = v.(float64) }, func(m any, k any, v any) { m.(map[int32][]byte)[k.(int32)] = v.([]byte) }, func(m any, k any, v any) { m.(map[int32]string)[k.(int32)] = v.(string) }, func(m any, k any, v any) { m.(map[int32]*TLV)[k.(int32)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[int64]any)[k.(int64)] = nil }, func(m any, k any, v any) { m.(map[int64]int8)[k.(int64)] = v.(int8) }, func(m any, k any, v any) { m.(map[int64]int16)[k.(int64)] = v.(int16) }, func(m any, k any, v any) { m.(map[int64]int32)[k.(int64)] = v.(int32) }, func(m any, k any, v any) { m.(map[int64]int64)[k.(int64)] = v.(int64) }, func(m any, k any, v any) { m.(map[int64]uint8)[k.(int64)] = v.(uint8) }, func(m any, k any, v any) { m.(map[int64]uint16)[k.(int64)] = v.(uint16) }, func(m any, k any, v any) { m.(map[int64]uint32)[k.(int64)] = v.(uint32) }, func(m any, k any, v any) { m.(map[int64]uint64)[k.(int64)] = v.(uint64) }, func(m any, k any, v any) { m.(map[int64]bool)[k.(int64)] = v.(bool) }, func(m any, k any, v any) { m.(map[int64]float32)[k.(int64)] = v.(float32) }, func(m any, k any, v any) { m.(map[int64]float64)[k.(int64)] = v.(float64) }, func(m any, k any, v any) { m.(map[int64][]byte)[k.(int64)] = v.([]byte) }, func(m any, k any, v any) { m.(map[int64]string)[k.(int64)] = v.(string) }, func(m any, k any, v any) { m.(map[int64]*TLV)[k.(int64)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[uint8]any)[k.(uint8)] = nil }, func(m any, k any, v any) { m.(map[uint8]int8)[k.(uint8)] = v.(int8) }, func(m any, k any, v any) { m.(map[uint8]int16)[k.(uint8)] = v.(int16) }, func(m any, k any, v any) { m.(map[uint8]int32)[k.(uint8)] = v.(int32) }, func(m any, k any, v any) { m.(map[uint8]int64)[k.(uint8)] = v.(int64) }, func(m any, k any, v any) { m.(map[uint8]uint8)[k.(uint8)] = v.(uint8) }, func(m any, k any, v any) { m.(map[uint8]uint16)[k.(uint8)] = v.(uint16) }, func(m any, k any, v any) { m.(map[uint8]uint32)[k.(uint8)] = v.(uint32) }, func(m any, k any, v any) { m.(map[uint8]uint64)[k.(uint8)] = v.(uint64) }, func(m any, k any, v any) { m.(map[uint8]bool)[k.(uint8)] = v.(bool) }, func(m any, k any, v any) { m.(map[uint8]float32)[k.(uint8)] = v.(float32) }, func(m any, k any, v any) { m.(map[uint8]float64)[k.(uint8)] = v.(float64) }, func(m any, k any, v any) { m.(map[uint8][]byte)[k.(uint8)] = v.([]byte) }, func(m any, k any, v any) { m.(map[uint8]string)[k.(uint8)] = v.(string) }, func(m any, k any, v any) { m.(map[uint8]*TLV)[k.(uint8)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[uint16]any)[k.(uint16)] = nil }, func(m any, k any, v any) { m.(map[uint16]int8)[k.(uint16)] = v.(int8) }, func(m any, k any, v any) { m.(map[uint16]int16)[k.(uint16)] = v.(int16) }, func(m any, k any, v any) { m.(map[uint16]int32)[k.(uint16)] = v.(int32) }, func(m any, k any, v any) { m.(map[uint16]int64)[k.(uint16)] = v.(int64) }, func(m any, k any, v any) { m.(map[uint16]uint8)[k.(uint16)] = v.(uint8) }, func(m any, k any, v any) { m.(map[uint16]uint16)[k.(uint16)] = v.(uint16) }, func(m any, k any, v any) { m.(map[uint16]uint32)[k.(uint16)] = v.(uint32) }, func(m any, k any, v any) { m.(map[uint16]uint64)[k.(uint16)] = v.(uint64) }, func(m any, k any, v any) { m.(map[uint16]bool)[k.(uint16)] = v.(bool) }, func(m any, k any, v any) { m.(map[uint16]float32)[k.(uint16)] = v.(float32) }, func(m any, k any, v any) { m.(map[uint16]float64)[k.(uint16)] = v.(float64) }, func(m any, k any, v any) { m.(map[uint16][]byte)[k.(uint16)] = v.([]byte) }, func(m any, k any, v any) { m.(map[uint16]string)[k.(uint16)] = v.(string) }, func(m any, k any, v any) { m.(map[uint16]*TLV)[k.(uint16)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[uint32]any)[k.(uint32)] = nil }, func(m any, k any, v any) { m.(map[uint32]int8)[k.(uint32)] = v.(int8) }, func(m any, k any, v any) { m.(map[uint32]int16)[k.(uint32)] = v.(int16) }, func(m any, k any, v any) { m.(map[uint32]int32)[k.(uint32)] = v.(int32) }, func(m any, k any, v any) { m.(map[uint32]int64)[k.(uint32)] = v.(int64) }, func(m any, k any, v any) { m.(map[uint32]uint8)[k.(uint32)] = v.(uint8) }, func(m any, k any, v any) { m.(map[uint32]uint16)[k.(uint32)] = v.(uint16) }, func(m any, k any, v any) { m.(map[uint32]uint32)[k.(uint32)] = v.(uint32) }, func(m any, k any, v any) { m.(map[uint32]uint64)[k.(uint32)] = v.(uint64) }, func(m any, k any, v any) { m.(map[uint32]bool)[k.(uint32)] = v.(bool) }, func(m any, k any, v any) { m.(map[uint32]float32)[k.(uint32)] = v.(float32) }, func(m any, k any, v any) { m.(map[uint32]float64)[k.(uint32)] = v.(float64) }, func(m any, k any, v any) { m.(map[uint32][]byte)[k.(uint32)] = v.([]byte) }, func(m any, k any, v any) { m.(map[uint32]string)[k.(uint32)] = v.(string) }, func(m any, k any, v any) { m.(map[uint32]*TLV)[k.(uint32)] = v.(*TLV) }},
	{func(m any, k any, v any) { m.(map[uint64]any)[k.(uint64)] = nil }, func(m any, k any, v any) { m.(map[uint64]int8)[k.(uint64)] = v.(int8) }, func(m any, k any, v any) { m.(map[uint64]int16)[k.(uint64)] = v.(int16) }, func(m any, k any, v any) { m.(map[uint64]int32)[k.(uint64)] = v.(int32) }, func(m any, k any, v any) { m.(map[uint64]int64)[k.(uint64)] = v.(int64) }, func(m any, k any, v any) { m.(map[uint64]uint8)[k.(uint64)] = v.(uint8) }, func(m any, k any, v any) { m.(map[uint64]uint16)[k.(uint64)] = v.(uint16) }, func(m any, k any, v any) { m.(map[uint64]uint32)[k.(uint64)] = v.(uint32) }, func(m any, k any, v any) { m.(map[uint64]uint64)[k.(uint64)] = v.(uint64) }, func(m any, k any, v any) { m.(map[uint64]bool)[k.(uint64)] = v.(bool) }, func(m any, k any, v any) { m.(map[uint64]float32)[k.(uint64)] = v.(float32) }, func(m any, k any, v any) { m.(map[uint64]float64)[k.(uint64)] = v.(float64) }, func(m any, k any, v any) { m.(map[uint64][]byte)[k.(uint64)] = v.([]byte) }, func(m any, k any, v any) { m.(map[uint64]string)[k.(uint64)] = v.(string) }, func(m any, k any, v any) { m.(map[uint64]*TLV)[k.(uint64)] = v.(*TLV) }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{func(m any, k any, v any) { m.(map[string]any)[k.(string)] = nil }, func(m any, k any, v any) { m.(map[string]int8)[k.(string)] = v.(int8) }, func(m any, k any, v any) { m.(map[string]int16)[k.(string)] = v.(int16) }, func(m any, k any, v any) { m.(map[string]int32)[k.(string)] = v.(int32) }, func(m any, k any, v any) { m.(map[string]int64)[k.(string)] = v.(int64) }, func(m any, k any, v any) { m.(map[string]uint8)[k.(string)] = v.(uint8) }, func(m any, k any, v any) { m.(map[string]uint16)[k.(string)] = v.(uint16) }, func(m any, k any, v any) { m.(map[string]uint32)[k.(string)] = v.(uint32) }, func(m any, k any, v any) { m.(map[string]uint64)[k.(string)] = v.(uint64) }, func(m any, k any, v any) { m.(map[string]bool)[k.(string)] = v.(bool) }, func(m any, k any, v any) { m.(map[string]float32)[k.(string)] = v.(float32) }, func(m any, k any, v any) { m.(map[string]float64)[k.(string)] = v.(float64) }, func(m any, k any, v any) { m.(map[string][]byte)[k.(string)] = v.([]byte) }, func(m any, k any, v any) { m.(map[string]string)[k.(string)] = v.(string) }, func(m any, k any, v any) { m.(map[string]*TLV)[k.(string)] = v.(*TLV) }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

var sLenMapHandlers = [T_COUNT][T_COUNT]func(m any) int{
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{func(m any) int { return len(m.(map[int8]any)) }, func(m any) int { return len(m.(map[int8]int8)) }, func(m any) int { return len(m.(map[int8]int16)) }, func(m any) int { return len(m.(map[int8]int32)) }, func(m any) int { return len(m.(map[int8]int64)) }, func(m any) int { return len(m.(map[int8]uint8)) }, func(m any) int { return len(m.(map[int8]uint16)) }, func(m any) int { return len(m.(map[int8]uint32)) }, func(m any) int { return len(m.(map[int8]uint64)) }, func(m any) int { return len(m.(map[int8]bool)) }, func(m any) int { return len(m.(map[int8]float32)) }, func(m any) int { return len(m.(map[int8]float64)) }, func(m any) int { return len(m.(map[int8][]byte)) }, func(m any) int { return len(m.(map[int8]string)) }, func(m any) int { return len(m.(map[int8]*TLV)) }},
	{func(m any) int { return len(m.(map[int16]any)) }, func(m any) int { return len(m.(map[int16]int8)) }, func(m any) int { return len(m.(map[int16]int16)) }, func(m any) int { return len(m.(map[int16]int32)) }, func(m any) int { return len(m.(map[int16]int64)) }, func(m any) int { return len(m.(map[int16]uint8)) }, func(m any) int { return len(m.(map[int16]uint16)) }, func(m any) int { return len(m.(map[int16]uint32)) }, func(m any) int { return len(m.(map[int16]uint64)) }, func(m any) int { return len(m.(map[int16]bool)) }, func(m any) int { return len(m.(map[int16]float32)) }, func(m any) int { return len(m.(map[int16]float64)) }, func(m any) int { return len(m.(map[int16][]byte)) }, func(m any) int { return len(m.(map[int16]string)) }, func(m any) int { return len(m.(map[int16]*TLV)) }},
	{func(m any) int { return len(m.(map[int32]any)) }, func(m any) int { return len(m.(map[int32]int8)) }, func(m any) int { return len(m.(map[int32]int16)) }, func(m any) int { return len(m.(map[int32]int32)) }, func(m any) int { return len(m.(map[int32]int64)) }, func(m any) int { return len(m.(map[int32]uint8)) }, func(m any) int { return len(m.(map[int32]uint16)) }, func(m any) int { return len(m.(map[int32]uint32)) }, func(m any) int { return len(m.(map[int32]uint64)) }, func(m any) int { return len(m.(map[int32]bool)) }, func(m any) int { return len(m.(map[int32]float32)) }, func(m any) int { return len(m.(map[int32]float64)) }, func(m any) int { return len(m.(map[int32][]byte)) }, func(m any) int { return len(m.(map[int32]string)) }, func(m any) int { return len(m.(map[int32]*TLV)) }},
	{func(m any) int { return len(m.(map[int64]any)) }, func(m any) int { return len(m.(map[int64]int8)) }, func(m any) int { return len(m.(map[int64]int16)) }, func(m any) int { return len(m.(map[int64]int32)) }, func(m any) int { return len(m.(map[int64]int64)) }, func(m any) int { return len(m.(map[int64]uint8)) }, func(m any) int { return len(m.(map[int64]uint16)) }, func(m any) int { return len(m.(map[int64]uint32)) }, func(m any) int { return len(m.(map[int64]uint64)) }, func(m any) int { return len(m.(map[int64]bool)) }, func(m any) int { return len(m.(map[int64]float32)) }, func(m any) int { return len(m.(map[int64]float64)) }, func(m any) int { return len(m.(map[int64][]byte)) }, func(m any) int { return len(m.(map[int64]string)) }, func(m any) int { return len(m.(map[int64]*TLV)) }},
	{func(m any) int { return len(m.(map[uint8]any)) }, func(m any) int { return len(m.(map[uint8]int8)) }, func(m any) int { return len(m.(map[uint8]int16)) }, func(m any) int { return len(m.(map[uint8]int32)) }, func(m any) int { return len(m.(map[uint8]int64)) }, func(m any) int { return len(m.(map[uint8]uint8)) }, func(m any) int { return len(m.(map[uint8]uint16)) }, func(m any) int { return len(m.(map[uint8]uint32)) }, func(m any) int { return len(m.(map[uint8]uint64)) }, func(m any) int { return len(m.(map[uint8]bool)) }, func(m any) int { return len(m.(map[uint8]float32)) }, func(m any) int { return len(m.(map[uint8]float64)) }, func(m any) int { return len(m.(map[uint8][]byte)) }, func(m any) int { return len(m.(map[uint8]string)) }, func(m any) int { return len(m.(map[uint8]*TLV)) }},
	{func(m any) int { return len(m.(map[uint16]any)) }, func(m any) int { return len(m.(map[uint16]int8)) }, func(m any) int { return len(m.(map[uint16]int16)) }, func(m any) int { return len(m.(map[uint16]int32)) }, func(m any) int { return len(m.(map[uint16]int64)) }, func(m any) int { return len(m.(map[uint16]uint8)) }, func(m any) int { return len(m.(map[uint16]uint16)) }, func(m any) int { return len(m.(map[uint16]uint32)) }, func(m any) int { return len(m.(map[uint16]uint64)) }, func(m any) int { return len(m.(map[uint16]bool)) }, func(m any) int { return len(m.(map[uint16]float32)) }, func(m any) int { return len(m.(map[uint16]float64)) }, func(m any) int { return len(m.(map[uint16][]byte)) }, func(m any) int { return len(m.(map[uint16]string)) }, func(m any) int { return len(m.(map[uint16]*TLV)) }},
	{func(m any) int { return len(m.(map[uint32]any)) }, func(m any) int { return len(m.(map[uint32]int8)) }, func(m any) int { return len(m.(map[uint32]int16)) }, func(m any) int { return len(m.(map[uint32]int32)) }, func(m any) int { return len(m.(map[uint32]int64)) }, func(m any) int { return len(m.(map[uint32]uint8)) }, func(m any) int { return len(m.(map[uint32]uint16)) }, func(m any) int { return len(m.(map[uint32]uint32)) }, func(m any) int { return len(m.(map[uint32]uint64)) }, func(m any) int { return len(m.(map[uint32]bool)) }, func(m any) int { return len(m.(map[uint32]float32)) }, func(m any) int { return len(m.(map[uint32]float64)) }, func(m any) int { return len(m.(map[uint32][]byte)) }, func(m any) int { return len(m.(map[uint32]string)) }, func(m any) int { return len(m.(map[uint32]*TLV)) }},
	{func(m any) int { return len(m.(map[uint64]any)) }, func(m any) int { return len(m.(map[uint64]int8)) }, func(m any) int { return len(m.(map[uint64]int16)) }, func(m any) int { return len(m.(map[uint64]int32)) }, func(m any) int { return len(m.(map[uint64]int64)) }, func(m any) int { return len(m.(map[uint64]uint8)) }, func(m any) int { return len(m.(map[uint64]uint16)) }, func(m any) int { return len(m.(map[uint64]uint32)) }, func(m any) int { return len(m.(map[uint64]uint64)) }, func(m any) int { return len(m.(map[uint64]bool)) }, func(m any) int { return len(m.(map[uint64]float32)) }, func(m any) int { return len(m.(map[uint64]float64)) }, func(m any) int { return len(m.(map[uint64][]byte)) }, func(m any) int { return len(m.(map[uint64]string)) }, func(m any) int { return len(m.(map[uint64]*TLV)) }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{func(m any) int { return len(m.(map[string]any)) }, func(m any) int { return len(m.(map[string]int8)) }, func(m any) int { return len(m.(map[string]int16)) }, func(m any) int { return len(m.(map[string]int32)) }, func(m any) int { return len(m.(map[string]int64)) }, func(m any) int { return len(m.(map[string]uint8)) }, func(m any) int { return len(m.(map[string]uint16)) }, func(m any) int { return len(m.(map[string]uint32)) }, func(m any) int { return len(m.(map[string]uint64)) }, func(m any) int { return len(m.(map[string]bool)) }, func(m any) int { return len(m.(map[string]float32)) }, func(m any) int { return len(m.(map[string]float64)) }, func(m any) int { return len(m.(map[string][]byte)) }, func(m any) int { return len(m.(map[string]string)) }, func(m any) int { return len(m.(map[string]*TLV)) }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

var sGetMapValueHandlers = [T_COUNT][T_COUNT]func(m any, key any) (any, bool){
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil}, {func(m any, key any) (any, bool) { v, ok := m.(map[int8]any)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]int8)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]int16)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]int32)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]int64)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]uint8)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]uint16)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]uint32)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]uint64)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]bool)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]float32)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]float64)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8][]byte)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]string)[key.(int8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int8]*TLV)[key.(int8)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[int16]any)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]int8)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]int16)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]int32)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]int64)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]uint8)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]uint16)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]uint32)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]uint64)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]bool)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]float32)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]float64)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16][]byte)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]string)[key.(int16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int16]*TLV)[key.(int16)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[int32]any)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]int8)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]int16)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]int32)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]int64)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]uint8)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]uint16)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]uint32)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]uint64)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]bool)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]float32)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]float64)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32][]byte)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]string)[key.(int32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int32]*TLV)[key.(int32)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[int64]any)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]int8)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]int16)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]int32)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]int64)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]uint8)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]uint16)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]uint32)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]uint64)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]bool)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]float32)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]float64)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64][]byte)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]string)[key.(int64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[int64]*TLV)[key.(int64)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[uint8]any)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]int8)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]int16)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]int32)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]int64)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]uint8)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]uint16)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]uint32)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]uint64)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]bool)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]float32)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]float64)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8][]byte)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]string)[key.(uint8)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint8]*TLV)[key.(uint8)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[uint16]any)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]int8)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]int16)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]int32)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]int64)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]uint8)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]uint16)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]uint32)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]uint64)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]bool)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]float32)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]float64)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16][]byte)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]string)[key.(uint16)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint16]*TLV)[key.(uint16)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[uint32]any)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]int8)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]int16)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]int32)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]int64)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]uint8)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]uint16)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]uint32)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]uint64)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]bool)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]float32)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]float64)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32][]byte)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]string)[key.(uint32)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint32]*TLV)[key.(uint32)]; return v, ok }},
	{func(m any, key any) (any, bool) { v, ok := m.(map[uint64]any)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]int8)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]int16)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]int32)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]int64)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]uint8)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]uint16)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]uint32)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]uint64)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]bool)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]float32)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]float64)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64][]byte)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]string)[key.(uint64)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[uint64]*TLV)[key.(uint64)]; return v, ok }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
	{func(m any, key any) (any, bool) { v, ok := m.(map[string]any)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]int8)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]int16)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]int32)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]int64)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]uint8)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]uint16)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]uint32)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]uint64)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]bool)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]float32)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]float64)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string][]byte)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]string)[key.(string)]; return v, ok }, func(m any, key any) (any, bool) { v, ok := m.(map[string]*TLV)[key.(string)]; return v, ok }},
	{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
}

func convertMapToTLVMap(st uint8, kt uint8, m any) (any, uint32) {
	if m == nil {
		fp := createEmptyMapHandlers[kt][st]
		if fp == nil {
			return nil, 0
		}
		return fp(), 0
	}
	switch v := m.(type) {
	case map[int8]int8:
		if kt == T_I8 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[int8]int16:
		if kt == T_I8 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[int8]int32:
		if kt == T_I8 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[int8]int64:
		if kt == T_I8 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[int8]uint8:
		if kt == T_I8 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[int8]uint16:
		if kt == T_I8 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[int8]uint32:
		if kt == T_I8 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[int8]uint64:
		if kt == T_I8 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[int8]bool:
		if kt == T_I8 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[int8]float32:
		if kt == T_I8 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[int8]float64:
		if kt == T_I8 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[int8][]byte:
		if kt == T_I8 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[int8]string:
		if kt == T_I8 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[int8]TLV:
		if kt == T_I8 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[int16]int8:
		if kt == T_I16 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[int16]int16:
		if kt == T_I16 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[int16]int32:
		if kt == T_I16 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[int16]int64:
		if kt == T_I16 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[int16]uint8:
		if kt == T_I16 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[int16]uint16:
		if kt == T_I16 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[int16]uint32:
		if kt == T_I16 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[int16]uint64:
		if kt == T_I16 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[int16]bool:
		if kt == T_I16 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[int16]float32:
		if kt == T_I16 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[int16]float64:
		if kt == T_I16 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[int16][]byte:
		if kt == T_I16 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[int16]string:
		if kt == T_I16 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[int16]TLV:
		if kt == T_I16 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[int32]int8:
		if kt == T_I32 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[int32]int16:
		if kt == T_I32 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[int32]int32:
		if kt == T_I32 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[int32]int64:
		if kt == T_I32 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[int32]uint8:
		if kt == T_I32 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[int32]uint16:
		if kt == T_I32 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[int32]uint32:
		if kt == T_I32 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[int32]uint64:
		if kt == T_I32 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[int32]bool:
		if kt == T_I32 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[int32]float32:
		if kt == T_I32 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[int32]float64:
		if kt == T_I32 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[int32][]byte:
		if kt == T_I32 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[int32]string:
		if kt == T_I32 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[int32]TLV:
		if kt == T_I32 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[int64]int8:
		if kt == T_I64 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[int64]int16:
		if kt == T_I64 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[int64]int32:
		if kt == T_I64 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[int64]int64:
		if kt == T_I64 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[int64]uint8:
		if kt == T_I64 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[int64]uint16:
		if kt == T_I64 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[int64]uint32:
		if kt == T_I64 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[int64]uint64:
		if kt == T_I64 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[int64]bool:
		if kt == T_I64 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[int64]float32:
		if kt == T_I64 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[int64]float64:
		if kt == T_I64 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[int64][]byte:
		if kt == T_I64 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[int64]string:
		if kt == T_I64 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[int64]TLV:
		if kt == T_I64 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[uint8]int8:
		if kt == T_U8 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[uint8]int16:
		if kt == T_U8 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[uint8]int32:
		if kt == T_U8 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[uint8]int64:
		if kt == T_U8 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[uint8]uint8:
		if kt == T_U8 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[uint8]uint16:
		if kt == T_U8 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[uint8]uint32:
		if kt == T_U8 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[uint8]uint64:
		if kt == T_U8 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[uint8]bool:
		if kt == T_U8 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[uint8]float32:
		if kt == T_U8 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[uint8]float64:
		if kt == T_U8 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[uint8][]byte:
		if kt == T_U8 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[uint8]string:
		if kt == T_U8 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[uint8]TLV:
		if kt == T_U8 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[uint16]int8:
		if kt == T_U16 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[uint16]int16:
		if kt == T_U16 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[uint16]int32:
		if kt == T_U16 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[uint16]int64:
		if kt == T_U16 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[uint16]uint8:
		if kt == T_U16 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[uint16]uint16:
		if kt == T_U16 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[uint16]uint32:
		if kt == T_U16 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[uint16]uint64:
		if kt == T_U16 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[uint16]bool:
		if kt == T_U16 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[uint16]float32:
		if kt == T_U16 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[uint16]float64:
		if kt == T_U16 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[uint16][]byte:
		if kt == T_U16 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[uint16]string:
		if kt == T_U16 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[uint16]TLV:
		if kt == T_U16 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[uint32]int8:
		if kt == T_U32 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[uint32]int16:
		if kt == T_U32 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[uint32]int32:
		if kt == T_U32 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[uint32]int64:
		if kt == T_U32 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[uint32]uint8:
		if kt == T_U32 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[uint32]uint16:
		if kt == T_U32 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[uint32]uint32:
		if kt == T_U32 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[uint32]uint64:
		if kt == T_U32 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[uint32]bool:
		if kt == T_U32 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[uint32]float32:
		if kt == T_U32 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[uint32]float64:
		if kt == T_U32 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[uint32][]byte:
		if kt == T_U32 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[uint32]string:
		if kt == T_U32 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[uint32]TLV:
		if kt == T_U32 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[uint64]int8:
		if kt == T_U64 && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[uint64]int16:
		if kt == T_U64 && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[uint64]int32:
		if kt == T_U64 && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[uint64]int64:
		if kt == T_U64 && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[uint64]uint8:
		if kt == T_U64 && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[uint64]uint16:
		if kt == T_U64 && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[uint64]uint32:
		if kt == T_U64 && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[uint64]uint64:
		if kt == T_U64 && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[uint64]bool:
		if kt == T_U64 && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[uint64]float32:
		if kt == T_U64 && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[uint64]float64:
		if kt == T_U64 && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[uint64][]byte:
		if kt == T_U64 && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[uint64]string:
		if kt == T_U64 && st == T_STR {
			return v, uint32(len(v))
		}
	case map[uint64]TLV:
		if kt == T_U64 && st == T_TLV {
			return v, uint32(len(v))
		}

	case map[string]int8:
		if kt == T_STR && st == T_I8 {
			return v, uint32(len(v))
		}
	case map[string]int16:
		if kt == T_STR && st == T_I16 {
			return v, uint32(len(v))
		}
	case map[string]int32:
		if kt == T_STR && st == T_I32 {
			return v, uint32(len(v))
		}
	case map[string]int64:
		if kt == T_STR && st == T_I64 {
			return v, uint32(len(v))
		}
	case map[string]uint8:
		if kt == T_STR && st == T_U8 {
			return v, uint32(len(v))
		}
	case map[string]uint16:
		if kt == T_STR && st == T_U16 {
			return v, uint32(len(v))
		}
	case map[string]uint32:
		if kt == T_STR && st == T_U32 {
			return v, uint32(len(v))
		}
	case map[string]uint64:
		if kt == T_STR && st == T_U64 {
			return v, uint32(len(v))
		}
	case map[string]bool:
		if kt == T_STR && st == T_BOOL {
			return v, uint32(len(v))
		}
	case map[string]float32:
		if kt == T_STR && st == T_F32 {
			return v, uint32(len(v))
		}
	case map[string]float64:
		if kt == T_STR && st == T_F64 {
			return v, uint32(len(v))
		}
	case map[string][]byte:
		if kt == T_STR && st == T_BYTES {
			return v, uint32(len(v))
		}
	case map[string]string:
		if kt == T_STR && st == T_STR {
			return v, uint32(len(v))
		}
	case map[string]TLV:
		if kt == T_STR && st == T_TLV {
			return v, uint32(len(v))
		}

	}

	return nil, uint32(0)
}

func convertListOfAnyHandlerOfTLV(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[*TLV](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfStr(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[string](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfBytes(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[[]byte](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfF64(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[float64](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfF32(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[float32](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfBool(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[bool](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfU64(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[uint64](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfU32(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[uint32](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfU16(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[uint16](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfU8(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[uint8](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfI64(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[int64](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfI32(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[int32](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfI16(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[int16](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfI8(t uint8, val any, l uint32) (any, uint32, int32) {
	val, rlen, err := SliceOfAnyToSliceOfTypeWithLength[int8](val, l)
	if err != nil {
		return nil, 0, core.MkErr(core.EC_TYPE_CONVERT_FAILED, 1)
	}
	return val, rlen, core.MkSuccess(0)
}

func convertListOfAnyHandlerOfNull(t uint8, val any, l uint32) (any, uint32, int32) {
	return make([]any, l), l, core.MkSuccess(0)
}

func MakeTLVType(ct uint8, st uint8) uint8 {
	return ((ct & 3) << 6) | (st & 0x3f)
}

func extractTlVType(t uint8) (uint8, uint8) {
	ct := t >> 6
	st := t & 0x3f
	return ct, st
}

func parseTlVParams(t uint8, length uint32, val any) (int32, uint32, uint8, uint8, uint16, any) {
	ct, st := extractTlVType(t)
	if st >= T_COUNT {
		panic("Invalid TLV Type " + string(t))
	}
	var rlen uint32 = 0
	var rtype uint8 = 0
	var rudt8 uint8 = 0
	var rudt16 uint16 = 0
	var rval any = nil
	var rc int32 = core.EC_OK
	if ct == DT_SINGLE {
		rc, rlen, rtype, rudt8, rudt16, rval = parseSingleTlvParamsHandles[st](length, t, val)
		return rc, rlen, rtype, rudt8, rudt16, rval
	} else if ct == DT_LIST {

	} else if ct == DT_DICT {
		panic("Invalid TLV Type " + string(t))
	}

	return core.EC_OK, rlen, rtype, rudt8, rudt16, rval
}

func CreateTLV3(t uint8, val any) *TLV {
	if t == T_NULL || val == nil {
		return &TLV{_type: t, _keyType: 0, _value: nil}
	}

	dt, st := extractTlVType(t)
	return CreateTLV(dt, st, 0, val)
}

func CreateTLV(dt uint8, st uint8, kt uint8, val any) *TLV {
	if dt >= DT_COUNT {
		return nil
	}
	if st >= T_COUNT {
		return nil
	}

	t := MakeTLVType(dt, st)
	var rtype uint8 = 0
	var rudt8 uint8 = 0
	var rval any = nil
	var rc int32 = core.EC_OK
	if dt == DT_SINGLE {
		rc, _, rtype, rudt8, _, rval = parseSingleTlvParamsHandles[st](0, t, val)
		if rc != core.EC_OK {
			return nil
		}
	} else if dt == DT_LIST {
		if val == nil {
			return sDefaultListOfAnyHandlers[st](0, t)
		}
		rval, _, rc = ConvertListOfAnyHandlers[st](t, val, 0)
		if core.Err(rc) {
			return nil
		}
		rtype = t
		rudt8 = 0
		rval = val

	} else if dt == DT_DICT {
		rudt8 = kt
		rval, _ = convertMapToTLVMap(st, kt, val)
		if rval == nil {
			return nil
		}
		rtype = t
	} else {
		panic("Invalid TLV Type " + string(t))
	}

	return &TLV{
		_type:    rtype,
		_keyType: rudt8,
		_value:   rval,
	}
}
