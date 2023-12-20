package transcomm

import (
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
)

func init() {
	GetDefaultMessageHandlerMapper().Register(messages.KEEP_ALIVE_MESSAGE_ID, KeepAliveMessageHandler)
	GetDefaultMessageHandlerMapper().Register(messages.PROC_TEST_MESSAGE_ID, ProcTestMessageHandler)
}

func ProcTestMessageHandler(connection IConnection, message message_buffer.INetMessage) int32 {
	var m *messages.ProcTestMessage = message.(*messages.ProcTestMessage)
	fmt.Println(m.String())
	if !m.Validate() {
		panic("Message Validation Failed")
	}

	return core.MkErr(core.EC_ALREADY_DONE, 0)
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
			ts := chrono.GetRealTimeMilli()
			delta := ts - pkam.TimeStamp()
			connection.(*TCPClientConnection)._codec.OnKeepAlive(ts, int32(delta))
		}

	} else if connection.Type() == CONNTYPE_TCP_SERVER {
		if pkam.IsServer() {
			ts := chrono.GetRealTimeMilli()
			delta := ts - pkam.TimeStamp()
			connection.(*TCPServerConnection)._codec.OnKeepAlive(ts, int32(delta))
		} else {
			connection.SendMessage(message, true)
		}
	}

	return core.MkErr(core.EC_ALREADY_DONE, 0)

}
