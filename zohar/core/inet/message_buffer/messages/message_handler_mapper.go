package messages

import (
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/transcomm"
)

type MessageHandler func(connection transcomm.IConnection, message message_buffer.INetMessage) int32

type MessageHandlerMapper struct {
	_mapper [32768]MessageHandler
}

func (ego *MessageHandlerMapper) Handle(connection transcomm.IConnection, message message_buffer.INetMessage) int32 {
	if ego._mapper[message.Command()] != nil {
		return ego._mapper[message.Command()](connection, message)
	} else {
		return core.MkErr(core.EC_NOOP, 1)
	}

}

func (ego *MessageHandlerMapper) Register(cmd int16, handler MessageHandler) {
	ego._mapper[cmd] = handler
}

func NeoMessageHandlerMapper() *MessageHandlerMapper {
	m := MessageHandlerMapper{}

	return &m
}

var sMessageHandlerMapperMapper *MessageHandlerMapper
var sMessageHandlerMapperMapperOnce sync.Once

func GetDefaultMessageHandlerMapper() *MessageHandlerMapper {
	sMessageHandlerMapperMapperOnce.Do(func() {
		sMessageHandlerMapperMapper = NeoMessageHandlerMapper()
	})
	return sMessageHandlerMapperMapper
}

func init() {
	GetDefaultMessageHandlerMapper().Register(KEEP_ALIVE_MESSAGE_ID, KeepAliveMessageHandler)
}
