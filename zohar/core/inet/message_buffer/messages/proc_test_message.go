package messages

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

var sLongText string = strs.CreateSampleString(1024*32, "@", "$")
var sMD5 [16]byte = md5.Sum([]byte(sLongText))

type ProcTestMessage struct {
	TimeStamp     int64          `json:"Timestamp"`
	Str0          string         `json:"Str0"`
	StrEmpty      string         `json:"StrEmpty"`
	I80           int8           `json:"I80"`
	I81           int8           `json:"I81"`
	I160          int16          `json:"I160"`
	I161          int16          `json:"I161"`
	I320          int32          `json:"I320"`
	I321          int32          `json:"I321"`
	I640          int64          `json:"I640"`
	I641          int64          `json:"I641"`
	IsServer      bool           `json:"IsServer"`
	F32           float32        `json:"F32"`
	F64           float64        `json:"F64"`
	U80           uint8          `json:"U80"`
	U160          uint16         `json:"U160"`
	U320          uint32         `json:"U320"`
	U640          uint64         `json:"U640"`
	StrSlice      []string       `json:"StrSlice"`
	StrSliceNull  []string       `json:"StrSliceNull"`
	StrSliceEmpty []string       `json:"StrSliceEmpty"`
	MD5           []byte         `json:"MD5"`
	TDE0          TestDataEntity `json:"TDE0"`
	TextLong      string         `json:"TextLong"`
}

func (ego *ProcTestMessage) Serialize(header memory.ISerializationHeader, buffer memory.IByteBuffer) (int64, int32) {
	var sPos int64 = buffer.WritePos()
	var rc int32 = core.MkSuccess(0)
	if header != nil {
		sPos, _, rc = header.BeginSerializing(buffer)
		if core.Err(rc) {
			return buffer.WritePos() - sPos, rc
		}
	}

	rc = buffer.WriteInt64(ego.TimeStamp)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	rc = buffer.WriteString(ego.Str0)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	rc = buffer.WriteString(ego.StrEmpty)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	rc = buffer.WriteInt8(ego.I80)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt8(ego.I81)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt16(ego.I160)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt16(ego.I161)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt32(ego.I320)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt32(ego.I321)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt64(ego.I640)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteInt64(ego.I641)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteBool(ego.IsServer)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteFloat32(ego.F32)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteFloat64(ego.F64)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	rc = buffer.WriteUInt8(ego.U80)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt16(ego.U160)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt32(ego.U320)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteUInt64(ego.U640)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	rc = buffer.WriteStrings(ego.StrSlice)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteStrings(ego.StrSliceNull)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteStrings(ego.StrSliceEmpty)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}
	rc = buffer.WriteBytes(ego.MD5)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	_, rc = ego.TDE0.Serialize(nil, buffer)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	rc = buffer.WriteString(ego.TextLong)
	if core.Err(rc) {
		return buffer.WritePos() - sPos, rc
	}

	if header != nil {
		_, rc = header.EndSerializing(buffer, sPos, buffer.WritePos()-sPos-header.HeaderLength())
		if core.Err(rc) {
			return buffer.WritePos() - sPos, rc
		}
	}

	return buffer.WritePos() - sPos, rc

}

func (ego *ProcTestMessage) Deserialize(header memory.ISerializationHeader, buffer memory.IByteBuffer) (int64, int32) {
	var rc int32
	var sLen int64 = buffer.ReadAvailable()
	if header != nil {
		_, rc = header.BeginDeserializing(buffer, true)
		if core.Err(rc) {
			return sLen - buffer.ReadAvailable(), rc
		}
	}

	ego.TimeStamp, rc = buffer.ReadInt64()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}

	ego.Str0, rc = buffer.ReadString()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.StrEmpty, rc = buffer.ReadString()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I80, rc = buffer.ReadInt8()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I81, rc = buffer.ReadInt8()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I160, rc = buffer.ReadInt16()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I161, rc = buffer.ReadInt16()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I320, rc = buffer.ReadInt32()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I321, rc = buffer.ReadInt32()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I640, rc = buffer.ReadInt64()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.I641, rc = buffer.ReadInt64()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.IsServer, rc = buffer.ReadBool()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.F32, rc = buffer.ReadFloat32()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.F64, rc = buffer.ReadFloat64()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.U80, rc = buffer.ReadUInt8()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.U160, rc = buffer.ReadUInt16()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.U320, rc = buffer.ReadUInt32()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.U640, rc = buffer.ReadUInt64()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.StrSlice, rc = buffer.ReadStrings()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.StrSliceNull, rc = buffer.ReadStrings()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.StrSliceEmpty, rc = buffer.ReadStrings()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}
	ego.MD5, rc = buffer.ReadBytes()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}

	_, rc = ego.TDE0.Deserialize(nil, buffer)
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}

	ego.TextLong, rc = buffer.ReadString()
	if core.Err(rc) {
		return sLen - buffer.ReadAvailable(), rc
	}

	if header != nil {
		rc = header.EndDeserializing(buffer)
		if core.Err(rc) {
			return sLen - buffer.ReadAvailable(), rc
		}
	}

	clen := sLen - buffer.ReadAvailable()
	c2len := header.BodyLength()
	if clen != c2len {
		panic("xxxxxxx")
	}

	return sLen - buffer.ReadAvailable(), core.MkSuccess(0)
}

func (ego *ProcTestMessage) IdentifierString() string {
	return ego.Str0
}

func (ego *ProcTestMessage) String() string {
	data, err := json.Marshal(ego)
	if err != nil {
		return "[Marshal_Failed_Msg]"
	}
	return string(data)
}

func (ego *ProcTestMessage) Validate() int32 {
	str0 := "Str0" + strconv.FormatInt(ego.TimeStamp, 10)
	if str0 != ego.Str0 {
		panic("Str0 validation failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.StrEmpty != "" {
		panic("StrEmpty failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}

	if ego.I80 != -128 {
		panic("I80 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I81 != 127 {
		panic("I81 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I160 != -32768 {
		panic("I160 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I161 != 32767 {
		panic("I161 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I320 != -(1 << 31) {
		panic("I320 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I321 != (1<<31)-1 {
		panic("I321 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I640 != -(1 << 63) {
		panic("I640 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.I641 != (1<<63)-1 {
		panic("I641 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.F32 != 2.71828 {
		panic("F32 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.F64 != 3.141592653 {
		panic("F64 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.U80 != uint8(ego.TimeStamp%256) {
		panic("U80 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.U160 != uint16(ego.TimeStamp%65536) {
		panic("U160 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.U320 != uint32(ego.TimeStamp%0xFFFFFFFF) {
		panic("U320 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ego.U640 != uint64(ego.TimeStamp) {
		panic("U640 failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}

	for i := 0; i < 11; i++ {
		cstr := fmt.Sprintf("StrSlice_%d", ego.TimeStamp)
		if i == 5 {
			if ego.StrSlice[i] != "" {
				panic("StrSlice failed.")
				return core.MkErr(core.EC_INVALID_STATE, 1)
			}
		} else {
			if ego.StrSlice[i] != cstr {
				panic("StrSlice failed.")
				return core.MkErr(core.EC_INVALID_STATE, 1)
			}
		}
	}

	if ego.StrSliceNull != nil {
		panic("StrSliceNULL failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}

	if len(ego.StrSliceEmpty) != 0 {
		panic("StrSliceEmpty failed.")
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}

	if ego.TextLong[0] != '@' || ego.TextLong[len(ego.TextLong)-1] != '$' {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}

	return core.MkSuccess(0)
}

func ProcTestMessageDeserialize(buffer memory.IByteBuffer, header memory.ISerializationHeader) (message_buffer.INetMessage, int64) {
	m := ProcTestMessage{}
	dataLength, rc := m.Deserialize(header, buffer)
	if core.Err(rc) {
		return nil, dataLength
	}
	return &m, dataLength
}

func NeoProcTestMessage(isClient bool) message_buffer.INetMessage {
	v := chrono.GetRealTimeMilli()
	m := ProcTestMessage{
		TimeStamp:     v,
		Str0:          "",
		I80:           -128,
		I81:           127,
		I160:          -32768,
		I161:          32767,
		I320:          -2147483648,
		I321:          2147483647,
		I640:          -(1 << 63),
		I641:          (1 << 63) - 1,
		F32:           2.71828,
		F64:           3.141592653,
		U80:           uint8(v % 256),
		U160:          uint16(v % 65536),
		U320:          uint32(v % (0xFFFFFFFF)),
		U640:          uint64(v),
		StrSlice:      make([]string, 11),
		StrSliceNull:  nil,
		StrSliceEmpty: make([]string, 0),
		IsServer:      isClient,
	}
	m.TDE0.FillTestData()
	m.TextLong = sLongText
	m.MD5 = sMD5[:]

	for i := 0; i < 11; i++ {
		m.StrSlice[i] = fmt.Sprintf("StrSlice_%d", v)
	}
	m.StrSlice[5] = ""

	m.Str0 = "Str0" + strconv.FormatInt(m.TimeStamp, 10)

	return &m
}

func (ego *ProcTestMessage) Command() uint16 {
	return PROC_TEST_MESSAGE_ID
}

func (ego *ProcTestMessage) GroupType() int8 {
	return 0
}

var _ message_buffer.INetMessage = &ProcTestMessage{}

var _ memory.ISerializable = &ProcTestMessage{}
