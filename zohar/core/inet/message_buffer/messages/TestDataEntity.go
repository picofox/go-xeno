package messages

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
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

func (ego *TestDataEntity) O1L15O1T15Serialize(helper *O1L15O1T15SerializationHelper) (int64, int32) {
	var rc int32 = -1
	rc = helper.WriteInt8s(ego.I8s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteUInt8s(ego.U8s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteInt16s(ego.I16s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteUInt16s(ego.U16s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteInt32s(ego.I32s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteUInt32s(ego.U32s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteInt64s(ego.I64s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}
	rc = helper.WriteUInt64s(ego.U64s)
	if core.Err(rc) {
		return helper.DataLength(), rc
	}

	return helper.DataLength(), rc
}

func (ego *TestDataEntity) O1L15O1T15Deserialize(helper *O1L15O1T15DeserializationHelper) (int64, int32) {
	//TODO implement me
	panic("implement me")
}

func (ego *TestDataEntity) String() string {
	//TODO implement me
	panic("implement me")
}

var _ message_buffer.INetDataEntity = &TestDataEntity{}
