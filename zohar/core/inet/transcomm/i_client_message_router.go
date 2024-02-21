package transcomm

import "xeno/zohar/core/inet/message_buffer"

type IClientMessageRouter interface {
	OnIncomingMessage(*TCPClientConnection, message_buffer.INetMessage) int32
	RegisterHandler(int8, uint16, ClientMessageHandler)
	UnregisterHandler(int8, uint16)
}
