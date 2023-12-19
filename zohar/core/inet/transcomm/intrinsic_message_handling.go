package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
)

func init() {
	GetDefaultMessageHandlerMapper().Register(messages.KEEP_ALIVE_MESSAGE_ID, KeepAliveMessageHandler)
}

func KeepAliveMessageHandler(connection IConnection, message message_buffer.INetMessage) int32 {
	var pkam *messages.KeepAliveMessage = message.(*messages.KeepAliveMessage)
	if pkam == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	if connection.Type() == CONNTYPE_TCP_CLIENT {
		if pkam.IsServer() {
			connection.SendMessage(message, true)
		} else {
			connection.(*TCPClientConnection)._codec.OnKeepAlive(chrono.GetRealTimeMilli())
		}

	} else if connection.Type() == CONNTYPE_TCP_SERVER {
		if pkam.IsServer() {
			connection.(*TCPServerConnection)._codec.OnKeepAlive(chrono.GetRealTimeMilli())
		} else {
			connection.SendMessage(message, true)
		}
	}

	return core.MkSuccess(0)

}
