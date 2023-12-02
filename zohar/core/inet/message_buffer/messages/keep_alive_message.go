package messages

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type KeepAliveMessage struct {
	_timeStamp int64
}

func (ego *KeepAliveMessage) Serialize(data []byte, offset int64) int32 {
	lb := memory.NeoLinearBufferAdapter(data, 0, offset, int64(cap(data))-offset)
	lb.WriteInt64(ego._timeStamp)
	return core.MkSuccess(0)
}

func KeepAliveMessageDeserialize(data []byte, offset int64) message_buffer.INetMessage {
	lb := memory.NeoLinearBufferAdapter(data, offset, int64(len(data))-offset, int64(cap(data))-offset)
	ts, _ := lb.ReadInt64()
	m := KeepAliveMessage{
		_timeStamp: ts,
	}
	return &m
}

func (ego *KeepAliveMessage) Command() int16 {
	return KEEP_ALIVE_MESSAGE_ID
}

//var _ message_buffer.INetMessage = &KeepAliveMessage{}
