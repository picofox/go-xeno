package messages

import (
	"encoding/json"
	"strconv"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type ProcTestMessage struct {
	TimeStamp int64  `json:"Timestamp"`
	Str0      string `json:"Str0"`
	StrEmpty  string `json:"StrEmpty"`
	I80       int8   `json:"I80"`
	I81       int8   `json:"I81"`
	I160      int16  `json:"I160"`
	I161      int16  `json:"I161"`
	I320      int32  `json:"I320"`
	I321      int32  `json:"I321"`
	I640      int64  `json:"I640"`
	I641      int64  `json:"I641"`
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
	if ego.I320 != -(2 ^ 31) {
		panic("I320 failed.")
		return false
	}
	if ego.I321 != (2^31)-1 {
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
	return true
}

func (ego *ProcTestMessage) Serialize(byteBuf memory.IByteBuffer) int64 {
	hdrPos := byteBuf.WritePos()
	byteBuf.WriteInt16(-1)
	byteBuf.WriteInt16(ego.Command())
	byteBuf.WriteInt64(ego.TimeStamp)
	byteBuf.WriteString(ego.Str0)
	byteBuf.WriteString(ego.StrEmpty)
	byteBuf.WriteInt8(ego.I80)
	byteBuf.WriteInt8(ego.I81)
	byteBuf.WriteInt16(ego.I160)
	byteBuf.WriteInt16(ego.I161)
	byteBuf.WriteInt32(ego.I320)
	byteBuf.WriteInt32(ego.I321)
	byteBuf.WriteInt64(ego.I640)
	byteBuf.WriteInt64(ego.I641)

	curPos := byteBuf.WritePos()
	var len64 int64 = curPos - hdrPos - message_buffer.O1L15O1T15_HEADER_SIZE
	if len64 <= message_buffer.MAX_PACKET_BODY_SIZE {
		byteBuf.WriterSeek(memory.BUFFER_SEEK_SET, hdrPos)
		byteBuf.WriteInt16(int16(len64))
		byteBuf.WriterSeek(memory.BUFFER_SEEK_SET, curPos)
	}
	return len64 + message_buffer.O1L15O1T15_HEADER_SIZE
}

func (ego *ProcTestMessage) Deserialize(buffer memory.IByteBuffer) int32 {
	var rc = int32(0)
	ts, _ := buffer.ReadInt64()
	ego.TimeStamp = ts
	if ego.Str0, rc = buffer.ReadString(); core.Err(rc) {
		return rc
	}
	if ego.StrEmpty, rc = buffer.ReadString(); core.Err(rc) {
		return rc
	}
	if ego.I80, rc = buffer.ReadInt8(); core.Err(rc) {
		return rc
	}
	if ego.I81, rc = buffer.ReadInt8(); core.Err(rc) {
		return rc
	}
	if ego.I160, rc = buffer.ReadInt16(); core.Err(rc) {
		return rc
	}
	if ego.I161, rc = buffer.ReadInt16(); core.Err(rc) {
		return rc
	}
	if ego.I320, rc = buffer.ReadInt32(); core.Err(rc) {
		return rc
	}
	if ego.I321, rc = buffer.ReadInt32(); core.Err(rc) {
		return rc
	}
	if ego.I640, rc = buffer.ReadInt64(); core.Err(rc) {
		return rc
	}
	if ego.I641, rc = buffer.ReadInt64(); core.Err(rc) {
		return rc
	}

	return rc
}

func ProcTestMessageDeserialize(buffer memory.IByteBuffer) message_buffer.INetMessage {
	m := ProcTestMessage{}
	rc := m.Deserialize(buffer)
	if core.Err(rc) {
		return nil
	}
	return &m
}

func NeoProcTestMessage() message_buffer.INetMessage {
	m := ProcTestMessage{
		TimeStamp: chrono.GetRealTimeMilli(),
		Str0:      "",
		I80:       -128,
		I81:       127,
		I160:      -32768,
		I161:      32767,
		I320:      -2147483648,
		I321:      2147483647,
		I640:      -(2 ^ 63),
		I641:      (2 ^ 63) - 1,
	}

	m.Str0 = "Str0" + strconv.FormatInt(m.TimeStamp, 10)

	return &m
}

func (ego *ProcTestMessage) Command() int16 {
	return PROC_TEST_MESSAGE_ID
}

var _ message_buffer.INetMessage = &ProcTestMessage{}
