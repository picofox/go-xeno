package message_buffer

import "xeno/zohar/core/memory"

type INetMessage interface {
	O1L15O1T15Serialize(memory.IByteBuffer) (int64, int32)
	O1L15O1T15Deserialize(memory.IByteBuffer, int16, int64) (int64, int32)
	Command() int16
	String() string
	MsgGrpType() int8
	IdentifierString() string
}
