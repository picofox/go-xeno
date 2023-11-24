package server

import (
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

const MAX_PACKET_LENGTH = 32 * 1024

type O1L15COT15DecodeServerHandler struct {
	_largeMessageBuffer *memory.LinearBuffer
	_memoryLow          bool
}

func (ego *O1L15COT15DecodeServerHandler) OnReceive(connection *TcpServerConnection, obj any, param1 any) (int32, any, any) {
	if connection.BufferLength() < 4 {
		return core.MkErr(core.EC_TRY_AGAIN, 1), nil, nil
	}
	data := obj.([]byte)
	lba := memory.NeoLinearBufferAdapter(data, 0, connection.BufferLength(), connection.BufferCapacity())
	o1AndLen, _ := lba.ReadUInt16()
	frameLength := int64(o1AndLen & 0x7FFF)
	if lba.ReadAvailable() < int64(frameLength) {
		return core.MkErr(core.EC_TRY_AGAIN, 2), nil, nil
	}
	//one packet is loaded
	o2AndType, _ := lba.ReadUInt16()
	opt1 := (o1AndLen >> 15 & 0x1) == 1
	opt2 := (o2AndType >> 15 & 0x1) == 1
	cmd := int16(o2AndType & 0x7FFF)
	dataRef := lba.SliceOf(frameLength)

	if !opt1 && !opt2 {
		return core.MkSuccess(0), dataRef, cmd
	} else if opt1 && !opt2 { //long message start
		if ego._largeMessageBuffer.Capacity() < frameLength*2 {
			ego._largeMessageBuffer.ResizeTo(frameLength * 2)
		}
		ego._largeMessageBuffer.WriteRawBytes(dataRef, 0, frameLength)
	} else if opt1 && opt2 { //long message trunks
		ego._largeMessageBuffer.WriteRawBytes(dataRef, 0, frameLength)
	} else if !opt1 && opt2 { //long message finished
		ego._largeMessageBuffer.WriteRawBytes(dataRef, 0, frameLength)
		ba, _ := ego._largeMessageBuffer.BytesRef()
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
