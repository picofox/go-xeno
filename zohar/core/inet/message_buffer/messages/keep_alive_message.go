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

func (ego *KeepAliveMessage) O1L15O1T15Serialize(byteBuf memory.IByteBuffer) (int64, int32) {
	sHelper, rc := InitializeSerialization(byteBuf, ego.MsgGrpType(), ego.Command())
	if core.Err(rc) {
		return 0, rc
	}
	defer sHelper.Finalize()
	rc = sHelper.WriteUInt64(ego.TT)
	if core.Err(rc) {
		return sHelper.DataLength(), rc
	}
	return sHelper.DataLength(), rc
}

func (ego *KeepAliveMessage) O1L15O1T15Deserialize(buffer memory.IByteBuffer, length int16, extraLength int64) (int64, int32) {
	dHelper, rc := InitializeDeserialization(buffer, ego.MsgGrpType(), ego.Command(), length, extraLength)
	if core.Err(rc) {
		return 0, rc
	}
	defer dHelper.Finalize()
	origLength := dHelper.DataLength()

	ego.TT, rc = dHelper.ReadUInt64()
	if core.Err(rc) {
		return origLength - dHelper.DataLength(), rc
	}
	if dHelper.DataLength() != 0 {
		return origLength - dHelper.DataLength(), core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	return origLength, core.MkSuccess(0)
}

func KeepAliveMessageDeserialize(buffer memory.IByteBuffer, length int16, extraLength int64) (message_buffer.INetMessage, int64) {
	m := KeepAliveMessage{}
	l, rc := m.O1L15O1T15Deserialize(buffer, length, extraLength)
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

func (ego *KeepAliveMessage) Command() int16 {
	return KEEP_ALIVE_MESSAGE_ID
}

func (ego *KeepAliveMessage) MsgGrpType() int8 {
	return 0
}

var _ message_buffer.INetMessage = &KeepAliveMessage{}
