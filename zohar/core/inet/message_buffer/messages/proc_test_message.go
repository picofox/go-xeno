package messages

import (
	"strconv"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type ProcTestMessage struct {
	_timeStamp int64
	_str0      string
}

func (ego *ProcTestMessage) Serialize(byteBuf memory.IByteBuffer) int64 {
	hdrPos := byteBuf.WritePos()
	byteBuf.WriteInt16(-1)
	byteBuf.WriteInt16(ego.Command())
	byteBuf.WriteInt64(ego._timeStamp)
	byteBuf.WriteString(ego._str0)

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
	ego._timeStamp = ts
	ego._str0, rc = buffer.ReadString()

	return rc
}

func ProcTestMessageDeserialize(buffer memory.IByteBuffer) message_buffer.INetMessage {
	m := KeepAliveMessage{}
	rc := m.Deserialize(buffer)
	if core.Err(rc) {
		return nil
	}
	return &m
}

func NeoProcTestMessage() message_buffer.INetMessage {
	m := ProcTestMessage{
		_timeStamp: chrono.GetRealTimeMilli(),
		_str0:      "",
	}

	m._str0 = "_str0" + strconv.FormatInt(m._timeStamp, 10)

	return &m
}

func (ego *ProcTestMessage) Command() int16 {
	return PROC_TEST_MESSAGE_ID
}

//var _ message_buffer.INetMessage = &KeepAliveMessage{}
