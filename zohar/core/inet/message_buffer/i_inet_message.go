package message_buffer

import "xeno/zohar/core/memory"

type INetMessage interface {
	Command() uint16
	String() string
	GroupType() int8
	IdentifierString() string
	Validate() int32

	Serialize(memory.ISerializationHeader, memory.IByteBuffer) (int64, int32)
	Deserialize(memory.ISerializationHeader, memory.IByteBuffer) (int64, int32)
}
