package message_buffer

import "xeno/zohar/core/memory"

type INetMessage interface {
	Serialize(memory.IByteBuffer) int64
	Deserialize(memory.IByteBuffer) int32
	Command() int16
}
