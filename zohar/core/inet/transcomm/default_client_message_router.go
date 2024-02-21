package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
)

type ClientMessageHandler func(*TCPClientConnection, message_buffer.INetMessage) int32

type DefaultClientMessageRouter struct {
	_client   *TCPClient
	_handlers [2][datatype.UINT16_CAPACITY]ClientMessageHandler
}

func (ego *DefaultClientMessageRouter) RegisterHandler(group int8, cmd uint16, handler ClientMessageHandler) {
	ego._handlers[group][cmd] = handler
}

func (ego *DefaultClientMessageRouter) UnregisterHandler(group int8, cmd uint16) {
	ego._handlers[group][cmd] = nil
}

func (ego *DefaultClientMessageRouter) OnIncomingMessage(conn *TCPClientConnection, message message_buffer.INetMessage) int32 {
	if ego._handlers[message.GroupType()][message.Command()] != nil {
		return ego._handlers[message.GroupType()][message.Command()](conn, message)
	}
	return core.MkErr(core.EC_HANDLER_NOT_FOUND, 1)
}

func (ego *HandlerRegistration) ReflectDefaultClientMessageRouter(c *TCPClient) *DefaultClientMessageRouter {
	dec := DefaultClientMessageRouter{
		_client: c,
	}
	return &dec
}

var _ IClientMessageRouter = &DefaultClientMessageRouter{}
