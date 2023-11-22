package strs

import (
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

func IPV4UIntToString(ip uint32) string {
	ba := memory.UInt32ToBytesBE(ip)
	str := make([]string, 4)
	for i := 0; i < 4; i++ {
		str[i] = strconv.Itoa(int((*ba)[i]))
	}
	return strings.Join(str, ".")
}

func IPV4Addr2UIntBE(str string) (uint32, int32) {
	is := strings.Split(str, ".")
	if len(is) != 4 {
		return 0, core.MkErr(core.EC_ERROR_COUNT, 1)
	}
	var data []byte = make([]byte, 4)
	for i := 0; i < 4; i++ {
		seg, err := strconv.Atoi(is[i])
		if err != nil || seg < 0 || seg > 255 {
			return 0, core.MkErr(core.EC_ERROR_COUNT, 2)
		}
		data[i] = byte(seg)
	}

	ret := memory.BytesToUInt32BE(&data, 0)
	return ret, core.MkSuccess(0)
}

func IPV4MaskBits2UIntBE(str string) (uint32, int32) {
	oneBits, err := strconv.Atoi(str)
	if err != nil || oneBits > 31 || oneBits < 0 {
		return 0, core.MkErr(core.EC_INDEX_OOB, 1)
	}

	var v uint32 = 0
	for i := 31; i >= 32-oneBits; i-- {
		v |= 1 << i
	}

	return v, core.MkSuccess(0)
}
