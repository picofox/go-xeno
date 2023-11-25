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

func (ego *KeepAliveMessage) Deserialize(data []byte, offset int64) int32 {
	var rc int32
	lb := memory.NeoLinearBufferAdapter(data, offset, int64(len(data))-offset, int64(cap(data))-offset)
	ego._timeStamp, rc = lb.ReadInt64()
	return rc
}

func (ego *KeepAliveMessage) Neo() message_buffer.INetMessage {
	return &KeepAliveMessage{
		_timeStamp: 0,
	}
}

func (ego *KeepAliveMessage) Command() int16 {
	return KEEP_ALIVE_MESSAGE_ID
}

//var _ message_buffer.INetMessage = &KeepAliveMessage{}
