package transcomm

import "xeno/zohar/core/inet/message_buffer"

type IServerMessageRouter interface {
	OnIncomingMessage(*TCPServerConnection, message_buffer.INetMessage) int32
	RegisterHandler(int8, uint16, ServerMessageHandler)
	UnregisterHandler(int8, uint16)
}
