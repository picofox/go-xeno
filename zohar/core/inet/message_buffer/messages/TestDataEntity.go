package messages

import (
	"encoding/json"
	"math/rand"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/memory"
)

type TestDataEntity struct {
	I8s  []int8   `json:"I8s"`
	U8s  []uint8  `json:"U8s"`
	I16s []int16  `json:"I16s"`
	U16s []uint16 `json:"U16s"`
	I32s []int32  `json:"I32s"`
	U32s []uint32 `json:"U32s"`
	I64s []int64  `json:"I64s"`
	U64s []uint64 `json:"U64s"`
}

func (ego *TestDataEntity) Serialize(header memory.ISerializationHeader, buffer memory.IByteBuffer) (int64, int32) {
	var sPos int64 = buffer.WritePos()
	var rc int32 = core.MkSuccess(0)
	if header != nil {
		sPos, _, rc = header.BeginSerializing(buffer)
		if core.Err(rc) {
			return buffer.WritePos() - sPos, rc
		}
	}
	rc = buffer.WriteInt8s(ego.I8s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt8s(ego.U8s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt16s(ego.I16s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt16s(ego.U16s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt32s(ego.I32s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt32s(ego.U32s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt64s(ego.I64s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt64s(ego.U64s)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	if header != nil {
		_, rc = header.EndSerializing(buffer, sPos, buffer.WritePos()-sPos-datatype.INT64_SIZE)
		if core.Err(rc) {
			return buffer.WritePos() - sPos, rc
		}
	}
	return buffer.WritePos() - sPos, rc
}

func (ego *TestDataEntity) Deserialize(header memory.ISerializationHeader, buffer memory.IByteBuffer) (int64, int32) {
	var sPos int64 = buffer.ReadPos()
	var rc int32 = core.MkSuccess(0)
	if header != nil {
		_, rc = header.BeginDeserializing(buffer, true)
		if core.Err(rc) {
			return buffer.ReadPos() - sPos, rc
		}
	}
	ego.I8s, rc = buffer.ReadInt8s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.U8s, rc = buffer.ReadUInt8s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.I16s, rc = buffer.ReadInt16s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.U16s, rc = buffer.ReadUInt16s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.I32s, rc = buffer.ReadInt32s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.U32s, rc = buffer.ReadUInt32s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.I64s, rc = buffer.ReadInt64s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	ego.U64s, rc = buffer.ReadUInt64s()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	if header != nil {
		rc = header.EndDeserializing(buffer)
		if core.Err(rc) {
			return buffer.ReadPos() - sPos, rc
		}
	}
	return buffer.ReadPos() - sPos, rc
}

func proRandomArr[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](m T) []T {
	it := rand.Intn(18)
	it--
	if it < 0 {
		return nil
	} else if it == 0 {
		return make([]T, 0)
	} else {
		retA := make([]T, it)
		for i := 0; i < it; i++ {
			retA[i] = T(i) % m
		}
		return retA
	}
}

func (ego *TestDataEntity) FillTestData() {
	ego.I8s = proRandomArr[int8](127)
	ego.U8s = proRandomArr[uint8](255)
	ego.I16s = proRandomArr[int16](32767)
	ego.U16s = proRandomArr[uint16](65535)
	ego.I32s = proRandomArr[int32](0x7FFFFFFF)
	ego.U32s = proRandomArr[uint32](0xFFFFFFFF)
	ego.I64s = proRandomArr[int64](0x7FFFFFFFFFFFFFFF)
	ego.U64s = proRandomArr[uint64](0xFFFFFFFFFFFFFFFF)
}

func (ego *TestDataEntity) O1L15O1T15Serialize(helper *O1L15O1T15SerializationHelper) int32 {
	var rc int32 = -1
	rc = helper.WriteInt8s(ego.I8s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteUInt8s(ego.U8s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteInt16s(ego.I16s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteUInt16s(ego.U16s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteInt32s(ego.I32s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteUInt32s(ego.U32s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteInt64s(ego.I64s)
	if core.Err(rc) {
		return rc
	}
	rc = helper.WriteUInt64s(ego.U64s)
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *TestDataEntity) O1L15O1T15Deserialize(helper *O1L15O1T15DeserializationHelper) int32 {
	var rc int32 = -1
	ego.I8s, rc = helper.ReadInt8s()
	if core.Err(rc) {
		return rc
	}
	ego.U8s, rc = helper.ReadUInt8s()
	if core.Err(rc) {
		return rc
	}
	ego.I16s, rc = helper.ReadInt16s()
	if core.Err(rc) {
		return rc
	}
	ego.U16s, rc = helper.ReadUInt16s()
	if core.Err(rc) {
		return rc
	}
	ego.I32s, rc = helper.ReadInt32s()
	if core.Err(rc) {
		return rc
	}
	ego.U32s, rc = helper.ReadUInt32s()
	if core.Err(rc) {
		return rc
	}
	ego.I64s, rc = helper.ReadInt64s()
	if core.Err(rc) {
		return rc
	}
	ego.U64s, rc = helper.ReadUInt64s()
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *TestDataEntity) String() string {
	data, err := json.Marshal(ego)
	if err != nil {
		return "[Marshal_Failed_Msg]"
	}
	return string(data)
}

var _ INetDataEntity = &TestDataEntity{}
var _ memory.ISerializable = &TestDataEntity{}
