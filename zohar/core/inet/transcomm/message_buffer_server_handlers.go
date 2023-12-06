package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type MessageBufferServerHandlers struct {
}

func (ego *MessageBufferServerHandlers) OnReceive(connection *TCPServerConnection, obj any, frameLength int64, param1 any) (int32, any, int64, any) {
	paramBA := obj.(memory.IByteBuffer)
	paramCMD := param1.(int16)

	if paramBA == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1), nil, 0, nil
	}
	if paramCMD < 0 {
		return core.MkErr(core.EC_INDEX_OOB, 1), nil, 0, nil
	}

	msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(paramCMD, paramBA)

	connection._server.OnIncomingMessage(connection, msg, nil)

	return core.MkSuccess(0), nil, 0, nil
}

func (ego *HandlerRegistration) NeoMessageBufferServerHandlers() *MessageBufferServerHandlers {
	dec := MessageBufferServerHandlers{}
	return &dec
}
