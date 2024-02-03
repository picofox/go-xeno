package messages

import (
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

func IsMessageComplete(buffer memory.IByteBuffer) (bool, int16, int16, int64, int32) {
	lenAndO1, rc := buffer.PeekInt16(0)
	if core.Err(rc) {
		return false, -1, -1, -1, rc
	}
	cmdAndO2, rc := buffer.PeekInt16(2)
	if core.Err(rc) {
		return false, -1, -1, -1, rc
	}
	o1 := lenAndO1>>15&0x1 == 1
	isInternal := cmdAndO2>>15&0x1 == 1
	cmd := int16(cmdAndO2 & 0x7FFF)
	l := int16(lenAndO1 & 0x7FFF)
	if isInternal {
		panic("internal")
	}
	if cmd != 1 {
		panic("cmd")
	}
	if !o1 {
		if buffer.ReadAvailable() >= message_buffer.O1L15O1T15_HEADER_SIZE+int64(l) {
			return isInternal, cmd, l, 0, core.MkSuccess(0)
		} else {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 1)
		}
	} else {
		off := message_buffer.O1L15O1T15_HEADER_SIZE + int64(l)
		if buffer.ReadAvailable() < off+8 {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 2)
		}
		eLen, rc := buffer.PeekInt64(off)
		if core.Err(rc) {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_INCOMPLETE_DATA, 3)
		}
		if buffer.ReadAvailable() < off+8+eLen {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 4)
		}
		return isInternal, cmd, l, eLen, core.MkSuccess(0)
	}
}

type O1L15O1T15DeserializationHelper struct {
	_command         int16
	_isInternal      bool
	_logicDataLength int16
	_extDataLength   int64
	_buffer          memory.IByteBuffer
	_temp            []byte
}

func O1L15O1T15DeserializationHelperCreator() any {
	return &O1L15O1T15DeserializationHelper{
		_command:         -1,
		_logicDataLength: 0,
		_extDataLength:   0,
		_buffer:          nil,
		_temp:            make([]byte, 8),
	}
}

var sO1L15O1T15DeserializationHelperCache *memory.ObjectCache[O1L15O1T15DeserializationHelper] = memory.NeoObjectCache[O1L15O1T15DeserializationHelper](16, O1L15O1T15DeserializationHelperCreator)

func (ego *O1L15O1T15DeserializationHelper) _init(buffer memory.IByteBuffer, isInternal bool, cmd int16, logicLength int16, extraLength int64) int32 {
	ego._command = cmd
	ego._isInternal = isInternal
	ego._logicDataLength = logicLength
	ego._extDataLength = extraLength
	ego._buffer = buffer
	ego._buffer.ReadInt32()
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15DeserializationHelper) BufferRemain() int64 {
	return ego._buffer.ReadAvailable()
}

func (ego *O1L15O1T15DeserializationHelper) String() string {
	var ss strings.Builder
	ss.WriteString("\nIsInternal=")
	ss.WriteString(strconv.FormatBool(ego._isInternal))
	ss.WriteString("\nLogicDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._logicDataLength)))
	ss.WriteString("\nExtDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._extDataLength)))
	return ss.String()
}

func InitializeDeserialization(buffer memory.IByteBuffer, isInternal bool, cmd int16, logicLength int16, extLength int64) (*O1L15O1T15DeserializationHelper, int32) {
	helper := sO1L15O1T15DeserializationHelperCache.Get()
	return helper, helper._init(buffer, isInternal, cmd, logicLength, extLength)
}

func (ego *O1L15O1T15DeserializationHelper) FinalizeDeserialization() int32 {
	defer sO1L15O1T15DeserializationHelperCache.Put(ego)
	if ego._logicDataLength != 0 || ego._extDataLength != 0 {
		panic("xxxx")
	}
	return 0
}
func (ego *O1L15O1T15DeserializationHelper) ReadUInt8() (uint8, int32) {
	v, r := ego.ReadInt8()
	return uint8(v), r
}

func (ego *O1L15O1T15DeserializationHelper) ReadUInt16() (uint16, int32) {
	v, r := ego.ReadInt16()
	return uint16(v), r
}

func (ego *O1L15O1T15DeserializationHelper) ReadUInt32() (uint32, int32) {
	v, r := ego.ReadInt32()
	return uint32(v), r
}

func (ego *O1L15O1T15DeserializationHelper) ReadUInt64() (uint64, int32) {
	v, r := ego.ReadInt64()
	return uint64(v), r
}

func (ego *O1L15O1T15DeserializationHelper) ReadBool() (bool, int32) {
	v, r := ego.ReadInt8()
	if v == 0 {
		return false, r
	} else if v == 1 {
		return true, r
	} else {
		return false, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadInt8() (int8, int32) {
	curTurnReadBytes := min(int64(ego._logicDataLength), datatype.INT8_SIZE)
	if curTurnReadBytes > 0 {
		iv, rc := ego._buffer.ReadInt8()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength--
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return iv, core.MkSuccess(0)

	} else {
		iv, rc := ego._buffer.ReadInt8()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength--
		return iv, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadInt16() (int16, int32) {
	curTurnReadBytes := min(int64(ego._logicDataLength), datatype.INT16_SIZE)
	if curTurnReadBytes >= datatype.INT16_SIZE {
		iv, rc := ego._buffer.ReadInt16()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(curTurnReadBytes)
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return iv, core.MkSuccess(0)

	} else if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(ego._temp, 0, curTurnReadBytes, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength = 0
		remainRLen := datatype.INT16_SIZE - curTurnReadBytes
		ego._buffer.ReadInt64()
		_, rc = ego._buffer.ReadRawBytes(ego._temp, curTurnReadBytes, remainRLen, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= remainRLen
		return memory.BytesToInt16BE(&ego._temp, 0), core.MkSuccess(0)

	} else {
		iv, rc := ego._buffer.ReadInt16()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= datatype.INT16_SIZE
		return iv, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadFloat32() (float32, int32) {
	curTurnReadBytes := min(int64(ego._logicDataLength), datatype.FLOAT32_SIZE)
	if curTurnReadBytes >= datatype.FLOAT32_SIZE {
		fv, rc := ego._buffer.ReadFloat32()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(curTurnReadBytes)
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return fv, core.MkSuccess(0)

	} else if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(ego._temp, 0, curTurnReadBytes, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength = 0
		remainRLen := datatype.FLOAT32_SIZE - curTurnReadBytes
		ego._buffer.ReadInt64()
		_, rc = ego._buffer.ReadRawBytes(ego._temp, curTurnReadBytes, remainRLen, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= remainRLen
		return memory.BytesToFloat32BE(&ego._temp, 0), core.MkSuccess(0)

	} else {
		fv, rc := ego._buffer.ReadFloat32()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= datatype.FLOAT32_SIZE
		return fv, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadFloat64() (float64, int32) {
	curTurnReadBytes := min(int64(ego._logicDataLength), datatype.FLOAT64_SIZE)
	if curTurnReadBytes >= datatype.FLOAT64_SIZE {
		fv, rc := ego._buffer.ReadFloat64()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(curTurnReadBytes)
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return fv, core.MkSuccess(0)

	} else if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(ego._temp, 0, curTurnReadBytes, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength = 0
		remainRLen := datatype.FLOAT64_SIZE - curTurnReadBytes
		ego._buffer.ReadInt64()
		_, rc = ego._buffer.ReadRawBytes(ego._temp, curTurnReadBytes, remainRLen, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= remainRLen
		return memory.BytesToFloat64BE(&ego._temp, 0), core.MkSuccess(0)

	} else {
		fv, rc := ego._buffer.ReadFloat64()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= datatype.FLOAT64_SIZE
		return fv, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadInt32() (int32, int32) {
	curTurnReadBytes := min(int64(ego._logicDataLength), datatype.INT32_SIZE)
	if curTurnReadBytes >= datatype.INT32_SIZE {
		iv, rc := ego._buffer.ReadInt32()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(curTurnReadBytes)
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return iv, core.MkSuccess(0)

	} else if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(ego._temp, 0, curTurnReadBytes, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength = 0
		remainRLen := datatype.INT32_SIZE - curTurnReadBytes
		ego._buffer.ReadInt64()
		_, rc = ego._buffer.ReadRawBytes(ego._temp, curTurnReadBytes, remainRLen, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= remainRLen
		return memory.BytesToInt32BE(&ego._temp, 0), core.MkSuccess(0)

	} else {
		iv, rc := ego._buffer.ReadInt32()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= datatype.INT32_SIZE
		return iv, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadInt64() (int64, int32) {
	curTurnReadBytes := min(int64(ego._logicDataLength), datatype.INT64_SIZE)
	if curTurnReadBytes >= datatype.INT64_SIZE {
		iv, rc := ego._buffer.ReadInt64()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(curTurnReadBytes)
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return iv, core.MkSuccess(0)

	} else if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(ego._temp, 0, curTurnReadBytes, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength = 0
		remainRLen := datatype.INT64_SIZE - curTurnReadBytes
		ego._buffer.ReadInt64()
		_, rc = ego._buffer.ReadRawBytes(ego._temp, curTurnReadBytes, remainRLen, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= remainRLen
		return memory.BytesToInt64BE(&ego._temp, 0), core.MkSuccess(0)

	} else {
		iv, rc := ego._buffer.ReadInt64()
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= datatype.INT64_SIZE
		return iv, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadString() (string, int32) {
	rBA, rc := ego.ReadBytes()
	if core.Err(rc) {
		return "", rc
	}
	if rBA == nil {
		return "", core.MkErr(core.EC_NULL_VALUE, 1)
	} else if len(rBA) == 0 {
		return "", core.MkSuccess(0)
	}
	return memory.StringRef(rBA), core.MkSuccess(0)
}

func (ego *O1L15O1T15DeserializationHelper) ReadBytes() ([]byte, int32) {
	l, rc := ego.ReadInt32()
	if core.Err(rc) {
		return nil, rc
	}
	if l == -1 {
		return nil, core.MkSuccess(0)
	} else if l == 0 {
		return memory.ConstEmptyBytes(), core.MkSuccess(0)
	} else if l > 0 {
		rBA := make([]byte, l)
		rc = ego.ReadRawBytes(rBA, 0, int64(l))
		if core.Err(rc) {
			return nil, rc
		}
		return rBA, core.MkSuccess(0)
	} else {
		panic("l < -1 data broken")
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadRawBytes(bs []byte, baOff int64, readLength int64) int32 {
	curTurnReadBytes := min(int64(ego._logicDataLength), readLength)
	if curTurnReadBytes >= readLength {
		_, rc := ego._buffer.ReadRawBytes(bs, baOff, readLength, true)
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(readLength)
		if ego._logicDataLength == 0 {
			if ego._extDataLength > 0 {
				ego._buffer.ReadInt64()
			}
		} else if ego._logicDataLength < 0 {
			panic("ego._logicDataLength < 0")
		}
		return core.MkSuccess(0)

	} else if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(bs, baOff, curTurnReadBytes, true)
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength = 0
		remainRLen := readLength - curTurnReadBytes
		ego._buffer.ReadInt64()
		_, rc = ego._buffer.ReadRawBytes(bs, baOff+curTurnReadBytes, remainRLen, true)
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= remainRLen
		return core.MkSuccess(0)

	} else {
		_, rc := ego._buffer.ReadRawBytes(bs, baOff, readLength, true)
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._extDataLength -= readLength
		return core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadStrings() ([]string, int32) {
	ac, rc := ego.ReadInt32()
	if ac == -1 {
		return nil, core.MkSuccess(0)
	} else if ac == 0 {
		return memory.ConstEmptyStringArr(), core.MkSuccess(0)
	} else {
		var ret []string = make([]string, ac)
		for i := int32(0); i < ac; i++ {
			ret[i], rc = ego.ReadString()
			if core.Err(rc) {
				return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
			}
		}
		return ret, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadInt8Array() ([]int8, int32) {
	ac, rc := ego.ReadInt32()
	if ac == -1 {
		return nil, core.MkSuccess(0)
	} else if ac == 0 {
		return memory.ConstEmptyInt8Arr(), core.MkSuccess(0)
	} else {
		var ret []int8 = make([]int8, ac)
		for i := int32(0); i < ac; i++ {
			ret[i], rc = ego.ReadInt8()
			if core.Err(rc) {
				return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
			}
		}
		return ret, core.MkSuccess(0)
	}
}

func (ego *O1L15O1T15DeserializationHelper) ReadUInt8Array() ([]uint8, int32) {
	ac, rc := ego.ReadInt32()
	if ac == -1 {
		return nil, core.MkSuccess(0)
	} else if ac == 0 {
		return memory.ConstEmptyUInt8Arr(), core.MkSuccess(0)
	} else {
		var ret []uint8 = make([]uint8, ac)
		for i := int32(0); i < ac; i++ {
			ret[i], rc = ego.ReadUInt8()
			if core.Err(rc) {
				return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
			}
		}
		return ret, core.MkSuccess(0)
	}
}
