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
	TimeStamp     int64    `json:"Timestamp"`
	Str0          string   `json:"Str0"`
	StrEmpty      string   `json:"StrEmpty"`
	I80           int8     `json:"I80"`
	I81           int8     `json:"I81"`
	I160          int16    `json:"I160"`
	I161          int16    `json:"I161"`
	I320          int32    `json:"I320"`
	I321          int32    `json:"I321"`
	I640          int64    `json:"I640"`
	I641          int64    `json:"I641"`
	IsServer      bool     `json:"IsServer"`
	F32           float32  `json:"F32"`
	F64           float64  `json:"F64"`
	U80           uint8    `json:"U80"`
	U160          uint16   `json:"U160"`
	U320          uint32   `json:"U320"`
	U640          uint64   `json:"U640"`
	StrSlice      []string `json:"StrSlice"`
	StrSliceNull  []string `json:"StrSliceNull"`
	StrSliceEmpty []string `json:"StrSliceEmpty"`
	MD5           []byte   `json:"MD5"`
	TextLong      string   `json:"TextLong"`
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

func (ego *ProcTestMessage) Validate() bool {
	str0 := "Str0" + strconv.FormatInt(ego.TimeStamp, 10)
	if str0 != ego.Str0 {
		panic("Str0 validation failed.")
		return false
	}
	if ego.StrEmpty != "" {
		panic("StrEmpty failed.")
		return false
	}

	if ego.I80 != -128 {
		panic("I80 failed.")
		return false
	}
	if ego.I81 != 127 {
		panic("I81 failed.")
		return false
	}
	if ego.I160 != -32768 {
		panic("I160 failed.")
		return false
	}
	if ego.I161 != 32767 {
		panic("I161 failed.")
		return false
	}
	if ego.I320 != -(1 << 31) {
		panic("I320 failed.")
		return false
	}
	if ego.I321 != (1<<31)-1 {
		panic("I321 failed.")
		return false
	}
	if ego.I640 != -(1 << 63) {
		panic("I640 failed.")
		return false
	}
	if ego.I641 != (1<<63)-1 {
		panic("I641 failed.")
		return false
	}
	if ego.F32 != 2.71828 {
		panic("F32 failed.")
		return false
	}
	if ego.F64 != 3.141592653 {
		panic("F64 failed.")
		return false
	}
	if ego.U80 != uint8(ego.TimeStamp%256) {
		panic("U80 failed.")
		return false
	}
	if ego.U160 != uint16(ego.TimeStamp%65536) {
		panic("U160 failed.")
		return false
	}
	if ego.U320 != uint32(ego.TimeStamp%0xFFFFFFFF) {
		panic("U320 failed.")
		return false
	}
	if ego.U640 != uint64(ego.TimeStamp) {
		panic("U640 failed.")
		return false
	}

	for i := 0; i < 11; i++ {
		cstr := fmt.Sprintf("StrSlice_%d", ego.TimeStamp)
		if i == 5 {
			if ego.StrSlice[i] != "" {
				panic("StrSlice failed.")
				return false
			}
		} else {
			if ego.StrSlice[i] != cstr {
				panic("StrSlice failed.")
				return false
			}
		}
	}

	if ego.StrSliceNull != nil {
		panic("StrSliceNULL failed.")
		return false
	}

	if len(ego.StrSliceEmpty) != 0 {
		panic("StrSliceEmpty failed.")
		return false
	}

	return true
}

func (ego *ProcTestMessage) O1L15O1T15Serialize(byteBuf memory.IByteBuffer) (int64, int32) {
	sHelper, rc := InitializeSerialization(byteBuf, ego.MsgGrpType(), ego.Command())
	if core.Err(rc) {
		return 0, rc
	}
	defer sHelper.Finalize()

	rc = sHelper.WriteInt64(ego.TimeStamp)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteString(ego.Str0)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}

	rc = sHelper.WriteString(ego.StrEmpty)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}

	rc = sHelper.WriteInt8(ego.I80)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt8(ego.I81)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt16(ego.I160)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt16(ego.I161)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt32(ego.I320)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt32(ego.I321)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt64(ego.I640)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteInt64(ego.I641)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteBool(ego.IsServer)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteFloat32(ego.F32)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteFloat64(ego.F64)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}

	rc = sHelper.WriteUInt8(ego.U80)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteUInt16(ego.U160)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteUInt32(ego.U320)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteUInt64(ego.U640)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}

	rc = sHelper.WriteStrings(ego.StrSlice)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteStrings(ego.StrSliceNull)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteStrings(ego.StrSliceEmpty)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteBytes(ego.MD5)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	rc = sHelper.WriteString(ego.TextLong)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}

	return sHelper.DataLength(), rc
}

func (ego *ProcTestMessage) O1L15O1T15Deserialize(buffer memory.IByteBuffer, length int16, extraLength int64) (int64, int32) {
	dHelper, rc := InitializeDeserialization(buffer, ego.MsgGrpType(), ego.Command(), length, extraLength)
	if core.Err(rc) {
		return 0, rc
	}
	defer dHelper.Finalize()
	origLength := dHelper.DataLength()

	ego.TimeStamp, rc = dHelper.ReadInt64()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.Str0, rc = dHelper.ReadString()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.StrEmpty, rc = dHelper.ReadString()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}

	ego.I80, rc = dHelper.ReadInt8()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I81, rc = dHelper.ReadInt8()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I160, rc = dHelper.ReadInt16()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I161, rc = dHelper.ReadInt16()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I320, rc = dHelper.ReadInt32()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I321, rc = dHelper.ReadInt32()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I640, rc = dHelper.ReadInt64()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.I641, rc = dHelper.ReadInt64()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.IsServer, rc = dHelper.ReadBool()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.F32, rc = dHelper.ReadFloat32()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.F64, rc = dHelper.ReadFloat64()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.U80, rc = dHelper.ReadUInt8()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.U160, rc = dHelper.ReadUInt16()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.U320, rc = dHelper.ReadUInt32()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.U640, rc = dHelper.ReadUInt64()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.StrSlice, rc = dHelper.ReadStrings()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.StrSliceNull, rc = dHelper.ReadStrings()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.StrSliceEmpty, rc = dHelper.ReadStrings()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.MD5, rc = dHelper.ReadBytes()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	ego.TextLong, rc = dHelper.ReadString()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}

	if dHelper.DataLength() != 0 {
		return origLength - dHelper.DataLength(), core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	return origLength, core.MkSuccess(0)
}

func ProcTestMessageDeserialize(buffer memory.IByteBuffer, length int16, extraLength int64) (message_buffer.INetMessage, int64) {
	m := ProcTestMessage{}
	dataLength, rc := m.O1L15O1T15Deserialize(buffer, length, extraLength)
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

	m.TextLong = sLongText
	m.MD5 = sMD5[:]

	for i := 0; i < 11; i++ {
		m.StrSlice[i] = fmt.Sprintf("StrSlice_%d", v)
	}
	m.StrSlice[5] = ""

	m.Str0 = "Str0" + strconv.FormatInt(m.TimeStamp, 10)

	return &m
}

func (ego *ProcTestMessage) Command() int16 {
	return PROC_TEST_MESSAGE_ID
}

func (ego *ProcTestMessage) MsgGrpType() int8 {
	return 0
}

var _ message_buffer.INetMessage = &ProcTestMessage{}
