package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

type O1L15COT15DecodeClientHandler struct {
	_largeMessageBuffer *memory.LinearBuffer
	_memoryLow          bool
}

func (ego *O1L15COT15DecodeClientHandler) Clear() {
	ego._largeMessageBuffer.Clear()
}

func (ego *O1L15COT15DecodeClientHandler) OnReceive(connection *TCPClientConnection, obj any, bufLen int64, param1 any) (int32, any, int64, any) {
	if connection._recvBuffer.ReadAvailable() < 4 {
		return core.MkErr(core.EC_TRY_AGAIN, 1), nil, 0, nil
	}
	o1AndLen, _, _, _ := connection._recvBuffer.PeekUInt16()
	frameLength := int64(o1AndLen & 0x7FFF)
	if connection._recvBuffer.ReadAvailable() < int64(frameLength) {
		return core.MkErr(core.EC_TRY_AGAIN, 2), nil, 0, nil
	}

	connection._recvBuffer.ReadInt16() //skip top half of header

	//one packet is loaded
	o2AndType, _ := connection._recvBuffer.ReadUInt16()
	opt1 := (o1AndLen >> 15 & 0x1) == 1
	opt2 := (o2AndType >> 15 & 0x1) == 1
	cmd := int16(o2AndType & 0x7FFF)

	//connection._recvBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, frameLength)

	if !opt1 && !opt2 {
		return core.MkSuccess(0), connection._recvBuffer, frameLength, cmd

	} else if opt1 && !opt2 { //long message start
		ego._largeMessageBuffer.Clear()
		if ego._largeMessageBuffer.Capacity() < frameLength*2 {
			ego._largeMessageBuffer.ResizeTo(frameLength * 2)
		}
		bs0, bs1 := connection._recvBuffer.BytesRef(frameLength)
		ego._largeMessageBuffer.WriteRawBytes(bs0, 0, int64(len(bs0)))
		if bs1 != nil && len(bs1) > 0 {
			ego._largeMessageBuffer.WriteRawBytes(bs1, 0, int64(len(bs1)))
		}
		connection._recvBuffer.Clear()
		return core.MkErr(core.EC_TRY_AGAIN, 1), nil, 0, nil

	} else if opt1 && opt2 { //long message trunks
		bs0, bs1 := connection._recvBuffer.BytesRef(frameLength)
		ego._largeMessageBuffer.WriteRawBytes(bs0, 0, int64(len(bs0)))
		if bs1 != nil && len(bs1) > 0 {
			ego._largeMessageBuffer.WriteRawBytes(bs1, 0, int64(len(bs1)))
		}
		connection._recvBuffer.Clear()
		return core.MkErr(core.EC_TRY_AGAIN, 1), nil, 0, nil

	} else if !opt1 && opt2 { //long message finished
		bs0, bs1 := connection._recvBuffer.BytesRef(frameLength)
		ego._largeMessageBuffer.WriteRawBytes(bs0, 0, int64(len(bs0)))
		if bs1 != nil && len(bs1) > 0 {
			ego._largeMessageBuffer.WriteRawBytes(bs1, 0, int64(len(bs1)))
		}
		connection._recvBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, frameLength)

		return core.MkSuccess(0), ego._largeMessageBuffer, ego._largeMessageBuffer.ReadAvailable(), cmd
	}

	return core.MkErr(core.EC_INVALID_STATE, 1), nil, 0, nil
}

func (ego *O1L15COT15DecodeClientHandler) CheckLowMemory() {
	if ego._memoryLow {
		ego._largeMessageBuffer.Reset()
		ego._memoryLow = false
	}
}

func (ego *O1L15COT15DecodeClientHandler) OnLowMemory() {
	ego._memoryLow = true
}

func NeoO1L15COT15DecodeClientHandler() *O1L15COT15DecodeClientHandler {
	dec := O1L15COT15DecodeClientHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
	}
	return &dec
}
func (ego *HandlerRegistration) NeoO1L15COT15DecodeClientHandler() *O1L15COT15DecodeClientHandler {
	dec := O1L15COT15DecodeClientHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
	}
	return &dec
}

var _ IClientHandler = &O1L15COT15DecodeClientHandler{}
