package message_buffer

type INetMessage interface {
	Serialize([]byte, int64) int32
	Command() int16
}
