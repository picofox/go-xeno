package transcomm

import "xeno/zohar/core/inet/message_buffer"

type IServerHandler interface {
	OnReceive(*TCPServerConnection) (message_buffer.INetMessage, int32)
	//Inbound([]IServerHandler, int, *TCPServerConnection, any, any) int32

	Clear()
}
