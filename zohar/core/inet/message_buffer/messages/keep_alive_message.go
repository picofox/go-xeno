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

func (ego *KeepAliveMessage) Command() uint16 {
	return KEEP_ALIVE_MESSAGE_ID
}

func (ego *KeepAliveMessage) GroupType() int8 {
	return INTERNAL_MSG_GRP_TYPE
}

func (ego *KeepAliveMessage) Serialize(header memory.ISerializationHeader, buffer memory.IByteBuffer) (int64, int32) {
	var sPos int64 = buffer.WritePos()
	var rc int32 = core.MkSuccess(0)
	if header != nil {
		sPos, _, rc = header.BeginSerializing(buffer)
		if core.Err(rc) {
			return buffer.WritePos() - sPos, rc
		}
	}
	rc = buffer.WriteUInt64(ego.TT)
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

func (ego *KeepAliveMessage) Deserialize(header memory.ISerializationHeader, buffer memory.IByteBuffer) (int64, int32) {
	var sPos int64 = buffer.ReadPos()
	var rc int32
	if header != nil {
		_, rc = header.BeginDeserializing(buffer, true)
		if core.Err(rc) {
			return buffer.ReadPos() - sPos, rc
		}
	}
	ego.TT, rc = buffer.ReadUInt64()
	if core.Err(rc) {
		return buffer.ReadPos() - sPos, rc
	}
	if header != nil {
		rc = header.EndDeserializing(buffer)
		if core.Err(rc) {
			return buffer.ReadPos() - sPos, rc
		}
	}
	return buffer.ReadPos() - sPos, core.MkSuccess(0)
}

func (ego *KeepAliveMessage) Validate() int32 {
	return core.MkSuccess(0)
}

func (ego *KeepAliveMessage) IdentifierString() string {
	return strconv.FormatInt(int64(ego.TT), 10)
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

func KeepAliveMessageDeserialize(buffer memory.IByteBuffer, header memory.ISerializationHeader) (message_buffer.INetMessage, int64) {
	m := KeepAliveMessage{}
	l, rc := m.Deserialize(header, buffer)
	if core.Err(rc) {
		return nil, l
	}
	return &m, l
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

var _ message_buffer.INetMessage = &KeepAliveMessage{}
