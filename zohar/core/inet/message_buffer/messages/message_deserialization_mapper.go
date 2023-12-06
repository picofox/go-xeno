package messages

import (
	"sync"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type MessageDeserializationHandler func(buffer memory.IByteBuffer) message_buffer.INetMessage

type MessageDeserializationMapper struct {
	_mapper [32768]MessageDeserializationHandler
}

func (ego *MessageDeserializationMapper) Deserialize(cmd int16, buffer memory.IByteBuffer) message_buffer.INetMessage {
	if cmd < 0 {
		return nil
	}
	if ego._mapper[cmd] != nil {
		return ego._mapper[cmd](buffer)
	}
	return nil
}

func (ego *MessageDeserializationMapper) Register(cmd int16, handler MessageDeserializationHandler) {
	ego._mapper[cmd] = handler
}

func NeoMessageDeserializationMapper() *MessageDeserializationMapper {
	m := MessageDeserializationMapper{}

	return &m
}

var sMessageDeserializationMapper *MessageDeserializationMapper
var sMessageDeserializationMapperOnce sync.Once

func GetDefaultMessageBufferDeserializationMapper() *MessageDeserializationMapper {
	sMessageDeserializationMapperOnce.Do(func() {
		sMessageDeserializationMapper = NeoMessageDeserializationMapper()
	})
	return sMessageDeserializationMapper
}

func init() {
	GetDefaultMessageBufferDeserializationMapper().Register(KEEP_ALIVE_MESSAGE_ID, KeepAliveMessageDeserialize)
}
