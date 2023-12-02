package server

import (
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

const O1L15O1T15_HEADER_SIZE = 4

type O1L15COT15DecodeServerHandler struct {
	_largeMessageBuffer *memory.LinearBuffer
	_memoryLow          bool
}

func (ego *O1L15COT15DecodeServerHandler) OnReceive(connection *TcpServerConnection, obj any, param1 any) (int32, any, any) {
	if connection._recvBuffer.ReadAvailable() < 4 {
		return core.MkErr(core.EC_TRY_AGAIN, 1), nil, nil
	}
	o1AndLen, _, _, _ := connection._recvBuffer.PeekUInt16()
	frameLength := int64(o1AndLen & 0x7FFF)
	if connection._recvBuffer.ReadAvailable() < int64(frameLength) {
		return core.MkErr(core.EC_TRY_AGAIN, 2), nil, nil
	}

	connection._recvBuffer.PeekUInt16() //skip top half of header

	//one packet is loaded
	o2AndType, _ := connection._recvBuffer.ReadUInt16()
	opt1 := (o1AndLen >> 15 & 0x1) == 1
	opt2 := (o2AndType >> 15 & 0x1) == 1
	cmd := int16(o2AndType & 0x7FFF)
	dataRef := connection._recvBuffer.SliceOf(frameLength)
	if !opt1 && !opt2 {
		return core.MkSuccess(0), dataRef, cmd
	} else if opt1 && !opt2 { //long message start
		if ego._largeMessageBuffer.Capacity() < frameLength*2 {
			ego._largeMessageBuffer.ResizeTo(frameLength * 2)
		}
		ego._largeMessageBuffer.WriteRawBytes(dataRef, 0, frameLength)
		connection._recvBuffer.Clear()
	} else if opt1 && opt2 { //long message trunks
		ego._largeMessageBuffer.WriteRawBytes(dataRef, 0, frameLength)
		connection._recvBuffer.Clear()
	} else if !opt1 && opt2 { //long message finished
		ego._largeMessageBuffer.WriteRawBytes(dataRef, 0, frameLength)
		ba, _ := ego._largeMessageBuffer.BytesRef()
		connection._recvBuffer.Clear()
		return core.MkSuccess(0), ba, cmd
	}

	return core.MkErr(core.EC_INVALID_STATE, 1), nil, nil
}

func (ego *O1L15COT15DecodeServerHandler) CheckLowMemory() {
	if ego._memoryLow {
		ego._largeMessageBuffer.Reset()
		ego._memoryLow = false
	}
}

func (ego *O1L15COT15DecodeServerHandler) OnLowMemory() {
	ego._memoryLow = true
}

func NeoO1L15COT15DecodeServerHandler() *O1L15COT15DecodeServerHandler {
	dec := O1L15COT15DecodeServerHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
	}
	return &dec
}
func (ego *HandlerRegistration) NeoO1L15COT15DecodeServerHandler() *O1L15COT15DecodeServerHandler {
	dec := O1L15COT15DecodeServerHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
	}
	return &dec
}
