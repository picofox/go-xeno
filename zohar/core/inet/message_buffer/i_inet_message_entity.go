package message_buffer

import (
	"xeno/zohar/core/inet/message_buffer/messages"
)

type INetDataEntity interface {
	O1L15O1T15Serialize(*messages.O1L15O1T15SerializationHelper) (int64, int32)
	O1L15O1T15Deserialize(*messages.O1L15O1T15DeserializationHelper) (int64, int32)
	String() string
}
