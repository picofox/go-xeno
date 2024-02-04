package messages

import (
	"sync"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

const (
	INTERNAL_MSG_GRP_TYPE = int8(0)
	EXTERNAL_MSG_GRP_TYPE = int8(1)
)

type MessageDeserializationHandler func(memory.IByteBuffer, int16, int64) (message_buffer.INetMessage, int64)

type MessageDeserializationMapper struct {
	_mappers [2][32768]MessageDeserializationHandler
}

func (ego *MessageDeserializationMapper) DeserializationDispatch(buffer memory.IByteBuffer, mGrpID int8, cmd int16, logicLength int16, extLength int64) (message_buffer.INetMessage, int64) {
	if cmd < 0 {
		return nil, 0
	}
	if ego._mappers[mGrpID][cmd] != nil {
		return ego._mappers[mGrpID][cmd](buffer, logicLength, extLength)
	}
	return nil, 0
}

func (ego *MessageDeserializationMapper) Register(mGrpId int8, cmd int16, handler MessageDeserializationHandler) {
	ego._mappers[mGrpId][cmd] = handler
}

func NeoMessageDeserializationMapper() *MessageDeserializationMapper {
	m := MessageDeserializationMapper{}
	for i := 0; i < len(m._mappers); i++ {
		for j := 0; j < len(m._mappers[i]); j++ {
			m._mappers[i][j] = nil
		}
	}
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
	GetDefaultMessageBufferDeserializationMapper().Register(INTERNAL_MSG_GRP_TYPE, KEEP_ALIVE_MESSAGE_ID, KeepAliveMessageDeserialize)
	GetDefaultMessageBufferDeserializationMapper().Register(INTERNAL_MSG_GRP_TYPE, PROC_TEST_MESSAGE_ID, ProcTestMessageDeserialize)

}
