package message_buffer

import "xeno/zohar/core/memory"

type INetMessage interface {
	Serialize(memory.IByteBuffer) int64
	Deserialize(memory.IByteBuffer) int32
	Command() int16
	String() string
	IdentifierString() string
	BodyLength() int64
	PiecewiseSerialize(bufferList *memory.ByteBufferList) (int64, int64, int32)
	PiecewiseDeserialize(bufferList *memory.ByteBufferList, bodyLength int64) (int64, int32)
}
