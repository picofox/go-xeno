package message_buffer

import "xeno/zohar/core/memory"

type INetMessage interface {
	Serialize([]byte, int64) int32
	Deserialize(memory.IByteBuffer) int32
	Command() int16
}
