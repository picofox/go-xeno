package transcomm

import "xeno/zohar/core/inet/message_buffer"

type IServerCodecHandler interface {
	OnReceive(*TCPServerConnection) (message_buffer.INetMessage, int32)
	//Inbound([]IServerCodecHandler, int, *TCPServerConnection, any, any) int32

	Reset()
}
