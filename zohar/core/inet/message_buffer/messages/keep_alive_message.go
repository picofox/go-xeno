package messages

import (
	"encoding/json"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type KeepAliveMessage struct {
	TimeStamp int64 `json:"timestamp"`
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
	byteBuf.WriteInt64(ego.TimeStamp)

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
	ts, _ := buffer.ReadInt64()
	ego.TimeStamp = ts

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

func NeoKeepAliveMessage() *KeepAliveMessage {
	m := KeepAliveMessage{
		TimeStamp: chrono.GetRealTimeMilli(),
	}
	return &m
}

func (ego *KeepAliveMessage) Command() int16 {
	return KEEP_ALIVE_MESSAGE_ID
}

//var _ message_buffer.INetMessage = &KeepAliveMessage{}
