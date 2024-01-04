package messages

import (
	"encoding/json"
	"strconv"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type KeepAliveMessage struct {
	TT uint64 `json:"TT"`
}

func (ego *KeepAliveMessage) IdentifierString() string {
	return strconv.FormatInt(int64(ego.TT), 10)
}

func (ego *KeepAliveMessage) PiecewiseDeserialize(bufferList *memory.ByteBufferList, bodyLen int64) (int64, int32) {
	var rc int32 = 0
	var logicPacketLength int64 = message_buffer.MAX_PACKET_BODY_SIZE

	rc = SkipHeader(bufferList)
	if core.Err(rc) {
		return bodyLen, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}

	ego.TT, logicPacketLength, bodyLen, rc = DeserializeU64Type(bufferList, logicPacketLength, bodyLen)
	if core.Err(rc) {
		return bodyLen, core.MkErr(core.EC_DESERIALIZE_FIELD_FAIELD, 0)
	}

	return bodyLen, core.MkSuccess(0)
}

func (ego *KeepAliveMessage) BodyLength() int64 {
	return 8
}

func (ego *KeepAliveMessage) PiecewiseSerialize(bufferList *memory.ByteBufferList) (int64, int64, int32) {
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

	rc, curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck = SerializeU64Type(ego.TT, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) { //lm:32756 bl:8
		FreeHeaders(headers)
		return totalIndex, bodyLenCheck, rc
	}

	FreeHeaders(headers)

	return totalIndex, bodyLenCheck, core.MkSuccess(0)
}

func (ego *KeepAliveMessage) SetTimeStamp(ts int64) {
	ego.TT = ego.TT & (1 << 63)
	ego.TT = ego.TT | uint64(ts)
}

func (ego *KeepAliveMessage) TimeStamp() int64 {
	return int64(ego.TT & 0x7FFFFFFFFFFFFFFF)
}

func (ego *KeepAliveMessage) SetIsServerToClient(b bool) {
	if b {
		ego.TT = ego.TT | (1 << 63)
	} else {
		ego.TT = ego.TT & 0x7FFFFFFFFFFFFFFF
	}
}

func (ego *KeepAliveMessage) IsServer() bool {
	if ego.TT&(1<<63) != 0 {
		return true
	}
	return false
}

func (ego *KeepAliveMessage) String() string {
	data, err := json.Marshal(ego)
	if err != nil {
		return "[Marshal_Failed_Msg]"
	}
	return string(data)
}

func (ego *KeepAliveMessage) Serialize(byteBuf memory.IByteBuffer) int64 {
	hdrPos := byteBuf.WritePos()
	byteBuf.WriteInt16(-1)
	byteBuf.WriteInt16(ego.Command())
	byteBuf.WriteUInt64(ego.TT)

	curPos := byteBuf.WritePos()
	var len64 int64 = curPos - hdrPos - message_buffer.O1L15O1T15_HEADER_SIZE
	if len64 <= message_buffer.MAX_PACKET_BODY_SIZE {
		byteBuf.WriterSeek(memory.BUFFER_SEEK_SET, hdrPos)
		byteBuf.WriteInt16(int16(len64))
		byteBuf.WriterSeek(memory.BUFFER_SEEK_SET, curPos)
	}
	return len64 + message_buffer.O1L15O1T15_HEADER_SIZE
}

func (ego *KeepAliveMessage) Deserialize(buffer memory.IByteBuffer) int32 {
	tt, _ := buffer.ReadUInt64()
	ego.TT = tt
	return core.MkSuccess(0)
}

func KeepAliveMessageDeserialize(buffer memory.IByteBuffer) message_buffer.INetMessage {
	m := KeepAliveMessage{}
	rc := m.Deserialize(buffer)
	if core.Err(rc) {
		return nil
	}
	return &m
}

func KeepAliveMessagePiecewiseDeserialize(bufferList *memory.ByteBufferList, bodyLength int64) message_buffer.INetMessage {
	m := KeepAliveMessage{}
	_, rc := m.PiecewiseDeserialize(bufferList, bodyLength)
	if core.Err(rc) {
		return nil
	}
	return &m
}

func _neoKeepAliveData(isServer bool) uint64 {
	if isServer {
		return uint64(chrono.GetRealTimeMilli()) | (1 << 63)
	} else {
		return uint64(chrono.GetRealTimeMilli())
	}
}

func NeoKeepAliveMessage(isServer bool) *KeepAliveMessage {
	m := KeepAliveMessage{
		TT: _neoKeepAliveData(isServer),
	}

	return &m
}

func (ego *KeepAliveMessage) Command() int16 {
	return KEEP_ALIVE_MESSAGE_ID
}

var _ message_buffer.INetMessage = &KeepAliveMessage{}
