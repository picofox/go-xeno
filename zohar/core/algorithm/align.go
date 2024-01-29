package algorithm

import "xeno/zohar/core"

func AlignSize[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](usize T, cnt T) T {
	return max((usize+cnt-1)&(^(cnt - 1)), cnt)
}

func PowerOf2[T int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64](v T) (T, int32) {
	var ret T = 0
	if (v & (v - 1)) == 0 {
		for {
			v = v / 2
			if v != 0 {
				ret++
			} else {
				break
			}
		}
	} else {
		return 0, core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return ret, core.MkSuccess(0)
}
