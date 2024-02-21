package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
)

type ServerMessageHandler func(*TCPServerConnection, message_buffer.INetMessage) int32

type DefaultServerMessageRouter struct {
	_server   *TCPServer
	_handlers [2][datatype.UINT16_CAPACITY]ServerMessageHandler
}

func (ego *DefaultServerMessageRouter) RegisterHandler(group int8, cmd uint16, handler ServerMessageHandler) {
	ego._handlers[group][cmd] = handler
}

func (ego *DefaultServerMessageRouter) UnregisterHandler(group int8, cmd uint16) {
	ego._handlers[group][cmd] = nil
}

func (ego *DefaultServerMessageRouter) OnIncomingMessage(conn *TCPServerConnection, message message_buffer.INetMessage) int32 {
	if ego._handlers[message.GroupType()][message.Command()] != nil {
		return ego._handlers[message.GroupType()][message.Command()](conn, message)
	}
	return core.MkErr(core.EC_HANDLER_NOT_FOUND, 1)
}

func (ego *HandlerRegistration) ReflectDefaultServerMessageRouter(s *TCPServer) *DefaultServerMessageRouter {
	dec := DefaultServerMessageRouter{
		_server: s,
	}
	return &dec
}

var _ IServerMessageRouter = &DefaultServerMessageRouter{}
