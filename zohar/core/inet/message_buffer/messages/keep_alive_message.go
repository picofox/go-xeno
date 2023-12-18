package messages

import (
	"encoding/json"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/transcomm"
	"xeno/zohar/core/memory"
)

type KeepAliveMessage struct {
	TT uint64 `json:"TT"`
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

//var _ message_buffer.INetMessage = &KeepAliveMessage{}

func KeepAliveMessageHandler(connection transcomm.IConnection, message message_buffer.INetMessage) int32 {
	var pkam *KeepAliveMessage = message.(*KeepAliveMessage)
	if pkam == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	if connection.Type() == transcomm.CONNTYPE_TCP_CLIENT {
		if pkam.IsServer() {
			connection.SendMessage(message, true)
		} else {

		}

	} else if connection.Type() == transcomm.CONNTYPE_TCP_SERVER {
		if pkam.IsServer() {

		} else {
			connection.SendMessage(message, true)
		}
	}

	return core.MkSuccess(0)

}
