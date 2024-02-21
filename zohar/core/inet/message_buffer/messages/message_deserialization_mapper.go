package messages

import (
	"sync"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

const (
	INTERNAL_MSG_GRP_TYPE = int8(0)
	EXTERNAL_MSG_GRP_TYPE = int8(1)
)

type MessageDeserializationHandler func(memory.IByteBuffer, memory.ISerializationHeader) (message_buffer.INetMessage, int64)

type MessageDeserializationMapper struct {
	_mappers [2][datatype.UINT16_CAPACITY]MessageDeserializationHandler
}

func (ego *MessageDeserializationMapper) DeserializationDispatch(buffer memory.IByteBuffer, header *memory.O1L31C16Header) (message_buffer.INetMessage, int64) {
	if ego._mappers[header.GroupType()][header.Command()] != nil {
		return ego._mappers[header.GroupType()][header.Command()](buffer, header)
	}
	return nil, -1
}

func (ego *MessageDeserializationMapper) Register(mGrpId int8, cmd uint16, handler MessageDeserializationHandler) {
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
