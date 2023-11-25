package message_buffer

type INetMessage interface {
	Serialize([]byte, int64) int32
	Deserialize([]byte, int64) int32
	Neo() INetMessage
	Command() int16
}
