package cmdline

import (
	"fmt"
	"strings"
	"xeno/zohar/core/memory"
)

type ArgumentSpec struct {
	_shortName  uint8
	_type       uint8
	_optional   uint8
	_keyType    uint8
	_minSACount int
	_maxSACount int
	_longName   string
}

func (ego *ArgumentSpec) name() string {
	return fmt.Sprintf("%s:(%c)", ego._longName, ego._shortName)
}

func (ego *ArgumentSpec) keyType() uint8 {
	return ego._keyType
}

func (ego *ArgumentSpec) ShortCommandString() string {
	return fmt.Sprintf("%c", ego._shortName)
}

func (ego *ArgumentSpec) SingleType() uint8 {
	return ego._type & 0x0f
}

func (ego *ArgumentSpec) ContainerType() uint8 {
	return (ego._type >> 4) & 0x0f
}

func (ego *ArgumentSpec) Optional() bool {
	if ego._optional == 0 {
		return false
	}
	return true
}

func (ego *ArgumentSpec) KeyTypeStr() string {
	if ego.IsDict() {
		if int(ego._keyType) < len(memory.SINGLE_TYPE_NAMES) {
			return memory.SINGLE_TYPE_NAMES[ego._keyType]
		}
	}

	return ""
}

func (ego *ArgumentSpec) ValueTypeStr() string {
	if !ego.IsValid() {
		return ""
	}
	if ego.IsFlag() {
		return ""
	}

	idx := ego._type & 0xf
	if int(idx) < len(memory.SINGLE_TYPE_NAMES) {
		return memory.SINGLE_TYPE_NAMES[idx]
	}

	return ""
}

func (ego *ArgumentSpec) ContainerTypeStr() string {
	if !ego.IsValid() {
		return "INV"
	}

	if ego.IsFlag() {
		return "F"
	} else if ego.IsSingle() {
		return "S"
	} else if ego.IsList() {
		return "L"
	} else if ego.IsDict() {
		return "D"
	}

	return "SNH"
}

func (ego *ArgumentSpec) IsValid() bool {
	if ego.HasShort() || ego.HasLong() {
		u := ego.SingleType()
		if ego.IsFlag() {
			if u == memory.T_NULL {
				return true
			}
		} else {
			if u > memory.T_NULL && u <= memory.T_STR {
				return true
			}
		}
	}

	return false
}

func (ego *ArgumentSpec) String() string {
	return fmt.Sprintf("%d:%d:%s:%d-%d:%s%s%s(%d)", ego._optional, ego._shortName, ego._longName, ego._minSACount, ego._maxSACount, ego.ContainerTypeStr(), ego.KeyTypeStr(), ego.ValueTypeStr(), ego._type)
}

func parseSingleStringType(typeStr string) uint8 {
	if len(typeStr) < 2 {
		return memory.T_NULL
	}

	if strings.Index(typeStr, "i8") >= 0 {
		return memory.T_I8
	} else if strings.Index(typeStr, "i16") >= 0 {
		return memory.T_I16
	} else if strings.Index(typeStr, "i32") >= 0 {
		return memory.T_I32
	} else if strings.Index(typeStr, "i64") >= 0 {
		return memory.T_I64
	} else if strings.Index(typeStr, "u8") >= 0 {
		return memory.T_U8
	} else if strings.Index(typeStr, "u16") >= 0 {
		return memory.T_U16
	} else if strings.Index(typeStr, "u32") >= 0 {
		return memory.T_U32
	} else if strings.Index(typeStr, "u64") >= 0 {
		return memory.T_U64
	} else if strings.Index(typeStr, "f32") >= 0 {
		return memory.T_F32
	} else if strings.Index(typeStr, "f64") >= 0 {
		return memory.T_F64
	} else if strings.Index(typeStr, "str") >= 0 {
		return memory.T_STR
	}

	return memory.T_NULL
}

func (ego *ArgumentSpec) HasLong() bool {
	if len(ego._longName) < 2 {
		return false
	}
	return true
}

func (ego *ArgumentSpec) HasShort() bool {
	if (ego._shortName >= 'A' && ego._shortName <= 'Z') || (ego._shortName >= 'a' && ego._shortName <= 'z') {
		return true
	}
	return false
}

func (ego *ArgumentSpec) IsFlag() bool {
	if ego._maxSACount == 0 {
		return true
	}
	return false
}

func (ego *ArgumentSpec) IsSingle() bool {
	hi := ego.ContainerType()
	if hi == memory.DT_SINGLE {
		return true
	}
	return false
}

func (ego *ArgumentSpec) IsList() bool {
	hi := ego.ContainerType()
	if hi == memory.DT_LIST {
		return true
	}
	return false
}

func (ego *ArgumentSpec) IsDict() bool {
	hi := ego.ContainerType()
	if hi == memory.DT_DICT {
		return true
	}
	return false
}

func parseTypeString(typeStr string) (uint8, uint8) {
	if len(typeStr) < 1 {
		return memory.T_NULL, memory.T_NULL
	}

	idx := strings.Index(typeStr, "D")
	if idx >= 0 {
		dSubStr := typeStr[idx+1:]
		kvp := strings.SplitN(dSubStr, "-", 2)
		if len(kvp) != 2 {
			return memory.T_NULL, memory.T_NULL
		}

		kt := parseSingleStringType(kvp[0])
		vt := parseSingleStringType(kvp[1])
		return (vt & 0xf) | (memory.DT_DICT & 0xf << 4), kt
	} else {
		idx := strings.Index(typeStr, "L")
		if idx >= 0 {
			vt := parseSingleStringType(typeStr)
			return (vt & 0xf) | (memory.DT_LIST << 4), memory.T_NULL
		} else {
			vt := parseSingleStringType(typeStr)
			return (vt & 0xf) | (memory.DT_SINGLE << 4), memory.T_NULL
		}
	}

}

func NeoArgumentSpec(s uint8, typeStr string, opt bool, minSAC int, maxSAC int, longName string) *ArgumentSpec {

	t, kt := parseTypeString(typeStr)
	optional := uint8(0)
	if opt {
		optional = uint8(1)
	}

	return &ArgumentSpec{
		_shortName:  s,
		_type:       t,
		_optional:   optional,
		_keyType:    kt,
		_minSACount: minSAC,
		_maxSACount: maxSAC,
		_longName:   longName,
	}
}
