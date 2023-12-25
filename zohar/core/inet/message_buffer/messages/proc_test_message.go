package messages

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

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

func (ego *ProcTestMessage) BodyLength() int64 {
	var sz int = 0
	var idx int = 0
	var tmpLen int = 0

	sz += 8 //TimeStamp
	sz += 4
	sz += len(ego.Str0) //Str0
	sz += 4
	sz += len(ego.StrEmpty) //StrEmpty
	sz += 1
	sz += 1
	sz += 2
	sz += 2
	sz += 4
	sz += 4
	sz += 8
	sz += 8
	sz += 1
	sz += 4
	sz += 8
	sz += 1
	sz += 2
	sz += 4
	sz += 8

	sz += 4
	tmpLen = len(ego.StrSlice)
	if ego.StrSlice != nil {
		for idx = 0; idx < tmpLen; idx++ {
			sz += 4
			sz += len(ego.StrSlice[idx])
		}
	}

	sz += 4
	tmpLen = len(ego.StrSliceNull)
	if ego.StrSliceNull != nil {
		for idx = 0; idx < tmpLen; idx++ {
			sz += 4
			sz += len(ego.StrSliceNull[idx])
		}
	}

	sz += 4
	tmpLen = len(ego.StrSliceEmpty)
	if ego.StrSliceEmpty != nil {
		for idx = 0; idx < tmpLen; idx++ {
			sz += 4
			sz += len(ego.StrSliceEmpty[idx])
		}
	}

	sz += 4
	sz += len(ego.MD5)

	sz += 4
	sz += len(ego.TextLong)
	return int64(sz)
}

func (ego *ProcTestMessage) SerializeToList(bufferList *memory.ByteBufferList) (int64, int32) {
	var totalIndex int64 = 0
	var bodyLenCheck int64 = 0
	var rc int32 = 0
	var curNode *memory.ByteBufferNode = nil
	var preCalBodyLen int64 = 0
	var logicPacketCount int64 = 0
	var logicPacketRemain int64 = 0
	var lastPackBytes int64 = 0
	var headers []*message_buffer.MessageHeader = nil
	var headerIdx int = 0

	preCalBodyLen = ego.BodyLength()
	logicPacketCount = (preCalBodyLen / message_buffer.MAX_PACKET_BODY_SIZE) + 1
	lastPackBytes = preCalBodyLen % message_buffer.MAX_PACKET_BODY_SIZE
	headers = AllocHeaders(logicPacketCount, lastPackBytes, ego.Command())
	logicPacketRemain = 0

	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI64Type(ego.TimeStamp, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeStringType(ego.Str0, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeStringType(ego.StrEmpty, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI8Type(ego.I80, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI8Type(ego.I81, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI16Type(ego.I160, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI16Type(ego.I161, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI32Type(ego.I320, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI32Type(ego.I321, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI64Type(ego.I640, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeI64Type(ego.I641, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeBoolType(ego.IsServer, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeF32Type(ego.F32, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeF64Type(ego.F64, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeU8Type(ego.U80, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeU16Type(ego.U160, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeU32Type(ego.U320, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeU64Type(ego.U640, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeStringsType(ego.StrSlice, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeStringsType(ego.StrSliceNull, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeStringsType(ego.StrSliceEmpty, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeBytesType(ego.MD5, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}
	rc, curNode, headerIdx, totalIndex, _, logicPacketRemain, bodyLenCheck = SerializeStringType(ego.TextLong, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return totalIndex, rc
	}

	return totalIndex, core.MkSuccess(0)
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

	md5 := md5.Sum([]byte(ego.TextLong))
	if bytes.Compare(md5[:], ego.MD5) != 0 {
		panic("StrSliceEmpty failed.")
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
	byteBuf.WriteBool(ego.IsServer)
	byteBuf.WriteFloat32(ego.F32)
	byteBuf.WriteFloat64(ego.F64)
	byteBuf.WriteUInt8(ego.U80)
	byteBuf.WriteUInt16(ego.U160)
	byteBuf.WriteUInt32(ego.U320)
	byteBuf.WriteUInt64(ego.U640)
	byteBuf.WriteStrings(ego.StrSlice)
	byteBuf.WriteStrings(ego.StrSliceNull)
	byteBuf.WriteStrings(ego.StrSliceEmpty)
	byteBuf.WriteBytes(ego.MD5)
	byteBuf.WriteString(ego.TextLong)
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
	if ego.IsServer, rc = buffer.ReadBool(); core.Err(rc) {
		return rc
	}
	if ego.F32, rc = buffer.ReadFloat32(); core.Err(rc) {
		return rc
	}
	if ego.F64, rc = buffer.ReadFloat64(); core.Err(rc) {
		return rc
	}
	if ego.U80, rc = buffer.ReadUInt8(); core.Err(rc) {
		return rc
	}
	if ego.U160, rc = buffer.ReadUInt16(); core.Err(rc) {
		return rc
	}
	if ego.U320, rc = buffer.ReadUInt32(); core.Err(rc) {
		return rc
	}
	if ego.U640, rc = buffer.ReadUInt64(); core.Err(rc) {
		return rc
	}
	if ego.StrSlice, rc = buffer.ReadStrings(); core.Err(rc) {
		return rc
	}
	if ego.StrSliceNull, rc = buffer.ReadStrings(); core.Err(rc) {
		return rc
	}
	if ego.StrSliceEmpty, rc = buffer.ReadStrings(); core.Err(rc) {
		return rc
	}
	if ego.MD5, rc = buffer.ReadBytes(); core.Err(rc) {
		return rc
	}
	if ego.TextLong, rc = buffer.ReadString(); core.Err(rc) {
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
	var ss strings.Builder
	for i := 0; i < 1024*16; i++ {
		ss.WriteString("@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789*")
	}
	m.TextLong = ss.String()
	ba := md5.Sum([]byte(m.TextLong))
	m.MD5 = ba[:]

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

var _ message_buffer.INetMessage = &ProcTestMessage{}
