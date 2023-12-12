package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type MessageBufferServerHandler struct {
}

func (ego *MessageBufferServerHandler) Clear() {

}

func (ego *MessageBufferServerHandler) OnReceive(connection *TCPServerConnection, obj any, frameLength int64, param1 any) (int32, any, int64, any) {
	paramBA := obj.(memory.IByteBuffer)
	paramCMD := param1.(int16)

	if paramBA == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1), nil, 0, nil
	}
	if paramCMD < 0 {
		return core.MkErr(core.EC_INDEX_OOB, 1), nil, 0, nil
	}

	beginPos := paramBA.ReadPos()
	msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(paramCMD, paramBA)
	if msg == nil {
		connection._server.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", paramCMD)
		return core.MkErr(core.EC_NULL_VALUE, 1), nil, 0, nil
	}
	endPos := paramBA.ReadPos()

	if endPos-beginPos != frameLength {
		connection._server.Log(core.LL_ERR, "Message (CMD:%d) Length Validation Failed, frame length is %d, but got %d read", paramCMD, frameLength, endPos-beginPos)
	}

	connection._server.OnIncomingMessage(connection, msg, nil)

	return core.MkSuccess(0), nil, 0, nil
}

func (ego *HandlerRegistration) NeoMessageBufferServerHandler() *MessageBufferServerHandler {
	dec := MessageBufferServerHandler{}
	return &dec
}

var _ IServerHandler = &MessageBufferServerHandler{}
