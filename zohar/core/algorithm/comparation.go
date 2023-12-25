package algorithm

import "xeno/zohar/core"

func MinValue[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64](v ...T) (int32, int, T) {
	var l int = len(v)
	if l <= 0 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 0), -1, 0
	} else if l == 1 {
		return core.MkSuccess(0), 0, v[0]
	}
	var beg T = v[0]
	var rIdx int = 0
	for i := 1; i < l; i++ {
		if beg > v[i] {
			beg = v[i]
			rIdx = i
		}
	}
	return core.MkSuccess(0), rIdx, beg
}
