package server

import (
	"xeno/zohar/core"
)

type MessageBufferServerHandlers struct {
}

func (ego *MessageBufferServerHandlers) OnReceive(connection *TcpServerConnection, obj any, param1 any) (int32, any, any) {
	paramBA := obj.([]byte)
	paramCMD := param1.(int16)

	if paramBA == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1), nil, nil
	}
	if paramCMD < 0 {
		return core.MkErr(core.EC_INDEX_OOB, 1), nil, nil
	}

	return core.MkSuccess(0), nil, nil
}

func (ego *HandlerRegistration) NeoMessageBufferServerHandlers() *MessageBufferServerHandlers {
	dec := MessageBufferServerHandlers{}
	return &dec
}
