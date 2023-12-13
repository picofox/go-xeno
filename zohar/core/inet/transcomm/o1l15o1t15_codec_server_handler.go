package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type O1L15COT15CodecServerHandler struct {
	_largeMessageBuffer *memory.LinearBuffer
	_memoryLow          bool
}

func (ego *O1L15COT15CodecServerHandler) Reset() {
	ego._largeMessageBuffer.Reset()
}

func (ego *O1L15COT15CodecServerHandler) OnReceive(connection *TCPServerConnection) (message_buffer.INetMessage, int32) {
	if connection._recvBuffer.ReadAvailable() < 4 {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	o1AndLen, _, _, _ := connection._recvBuffer.PeekUInt16()
	frameLength := int64(o1AndLen & 0x7FFF)
	if connection._recvBuffer.ReadAvailable() < int64(frameLength) {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 2)
	}

	connection._recvBuffer.ReadInt16() //skip top half of header

	//one packet is loaded
	o2AndType, _ := connection._recvBuffer.ReadUInt16()
	opt1 := (o1AndLen >> 15 & 0x1) == 1
	opt2 := (o2AndType >> 15 & 0x1) == 1
	cmd := int16(o2AndType & 0x7FFF)

	//connection._recvBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, frameLength)

	if !opt1 && !opt2 {
		beginPos := connection._recvBuffer.ReadPos()
		msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(cmd, connection._recvBuffer)
		if msg == nil {
			connection._server.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", cmd)
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		endPos := connection._recvBuffer.ReadPos()
		if endPos-beginPos != frameLength {
			connection._server.Log(core.LL_ERR, "Message (CMD:%d) Length Validation Failed, frame length is %d, but got %d read", cmd, frameLength, endPos-beginPos)
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}

		return msg, core.MkSuccess(0)

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
		return nil, core.MkErr(core.EC_TRY_AGAIN, 2)

	} else if opt1 && opt2 { //long message trunks
		bs0, bs1 := connection._recvBuffer.BytesRef(frameLength)
		ego._largeMessageBuffer.WriteRawBytes(bs0, 0, int64(len(bs0)))
		if bs1 != nil && len(bs1) > 0 {
			ego._largeMessageBuffer.WriteRawBytes(bs1, 0, int64(len(bs1)))
		}
		connection._recvBuffer.Clear()
		return nil, core.MkErr(core.EC_TRY_AGAIN, 2)

	} else if !opt1 && opt2 { //long message finished
		bs0, bs1 := connection._recvBuffer.BytesRef(frameLength)
		ego._largeMessageBuffer.WriteRawBytes(bs0, 0, int64(len(bs0)))
		if bs1 != nil && len(bs1) > 0 {
			ego._largeMessageBuffer.WriteRawBytes(bs1, 0, int64(len(bs1)))
		}
		connection._recvBuffer.ReaderSeek(memory.BUFFER_SEEK_CUR, frameLength)

		msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(cmd, ego._largeMessageBuffer)
		if msg == nil {
			connection._server.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", cmd)
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}

		return msg, core.MkSuccess(0)
	}

	return nil, core.MkErr(core.EC_INVALID_STATE, 1)
}

func (ego *O1L15COT15CodecServerHandler) CheckLowMemory() {
	if ego._memoryLow {
		ego._largeMessageBuffer.Reset()
		ego._memoryLow = false
	}
}

func (ego *O1L15COT15CodecServerHandler) OnLowMemory() {
	ego._memoryLow = true
}

func NeoO1L15COT15DecodeServerHandler() *O1L15COT15CodecServerHandler {
	dec := O1L15COT15CodecServerHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
	}
	return &dec
}
func (ego *HandlerRegistration) NeoO1L15COT15DecodeServerHandler() *O1L15COT15CodecServerHandler {
	dec := O1L15COT15CodecServerHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
	}
	return &dec
}

var _ IServerCodecHandler = &O1L15COT15CodecServerHandler{}
