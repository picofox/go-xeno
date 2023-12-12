package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type MessageBufferClientHandler struct {
}

func (ego *MessageBufferClientHandler) OnSend(connection *TCPClientConnection, p1 any, tLen int64, bFlush bool) (int32, any, int64, bool) {
	if connection._sendBuffer.Capacity() >= message_buffer.MAX_BUFFER_MAX_CAPACITY && connection._sendBuffer.WriteAvailable() < message_buffer.O1L15O1T15_HEADER_SIZE {
		connection.flush()
	}

	var message message_buffer.INetMessage = p1.(message_buffer.INetMessage)
	tLen = message.Serialize(connection._sendBuffer)
	if tLen < 0 {
		return core.MkErr(core.EC_TRY_AGAIN, 1), message.Command(), tLen, bFlush
	}
	return core.MkSuccess(0), message.Command(), tLen, bFlush
}

func (ego *MessageBufferClientHandler) Clear() {

}

func (ego *MessageBufferClientHandler) OnReceive(connection *TCPClientConnection, obj any, frameLength int64, param1 any) (int32, any, int64, any) {
	paramBA := obj.(memory.IByteBuffer)
	paramCMD := param1.(int16)

	if paramBA == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1), nil, 0, nil
	}
	if paramCMD < 0 {
		return core.MkErr(core.EC_INDEX_OOB, 1), nil, 0, nil
	}

	msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(paramCMD, paramBA)

	connection._client.OnIncomingMessage(connection, msg)

	return core.MkSuccess(0), nil, 0, nil
}

func (ego *HandlerRegistration) NeoMessageBufferClientHandlers() *MessageBufferClientHandler {
	dec := MessageBufferClientHandler{}
	return &dec
}

var _ IClientHandler = &MessageBufferClientHandler{}
