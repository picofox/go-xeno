package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type O1L15COT15CodecClientHandler struct {
	_largeMessageBuffer *memory.LinearBuffer
	_memoryLow          bool
	_packetHeader       message_buffer.MessageHeader
}

func (ego *O1L15COT15CodecClientHandler) OnSend(connection *TCPClientConnection, a any, bFlush bool) int32 {
	var message = a.(message_buffer.INetMessage)
	tLen := message.Serialize(connection._sendBuffer)
	if tLen < 0 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}

	var byteBuf memory.IByteBuffer = connection._sendBuffer
	var cmd int16 = message.Command()

	if tLen <= message_buffer.MAX_PACKET_BODY_SIZE {
		if !bFlush && byteBuf.WriteAvailable() > tLen {
			return core.MkSuccess(1)
		}
		connection.sendImmediately(*(byteBuf.InternalData()), byteBuf.ReadPos(), byteBuf.ReadAvailable())
		byteBuf.Clear()
		return core.MkSuccess(0)

	} else { //large message
		connection.flush()
		rIndex := int64(4)
		ego._packetHeader.Set(true, false, message_buffer.MAX_PACKET_BODY_SIZE, cmd)
		byteBuf.ReaderSeek(memory.BUFFER_SEEK_CUR, message_buffer.O1L15O1T15_HEADER_SIZE)

		for {
			connection.sendImmediately(ego._packetHeader.Data(), 0, message_buffer.O1L15O1T15_HEADER_SIZE)
			connection.sendImmediately(*byteBuf.InternalData(), byteBuf.ReadPos(), message_buffer.MAX_PACKET_BODY_SIZE)

			rIndex += message_buffer.MAX_PACKET_BODY_SIZE
			byteBuf.ReaderSeek(memory.BUFFER_SEEK_SET, rIndex)

			//next loop use non begin version
			ego._packetHeader.Set(true, true, message_buffer.MAX_PACKET_BODY_SIZE, cmd)

			if byteBuf.ReadAvailable() <= message_buffer.MAX_PACKET_BODY_SIZE {
				break
			}
		}
		ego._packetHeader.Set(false, true, message_buffer.MAX_PACKET_BODY_SIZE, cmd)
		connection.sendImmediately(ego._packetHeader.Data(), 0, 4)
		connection.sendImmediately(*byteBuf.InternalData(), byteBuf.ReadPos(), byteBuf.ReadAvailable())
	}
	return core.MkSuccess(0)
}

func (ego *O1L15COT15CodecClientHandler) Reset() {
	ego._largeMessageBuffer.Reset()
}

func (ego *O1L15COT15CodecClientHandler) OnReceive(connection *TCPClientConnection) (any, int32) {
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
			connection._client.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", cmd)
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		endPos := connection._recvBuffer.ReadPos()
		if endPos-beginPos != frameLength {
			connection._client.Log(core.LL_ERR, "Message (CMD:%d) Length Validation Failed, frame length is %d, but got %d read", cmd, frameLength, endPos-beginPos)
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
			connection._client.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", cmd)
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}

		return msg, core.MkSuccess(0)
	}

	return nil, core.MkErr(core.EC_INVALID_STATE, 1)
}

func (ego *O1L15COT15CodecClientHandler) CheckLowMemory() {
	if ego._memoryLow {
		ego._largeMessageBuffer.Reset()
		ego._memoryLow = false
	}
}

func (ego *O1L15COT15CodecClientHandler) OnLowMemory() {
	ego._memoryLow = true
}

func NeoO1L15COT15DecodeClientHandler() *O1L15COT15CodecClientHandler {
	dec := O1L15COT15CodecClientHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
		_packetHeader:       message_buffer.NeoMessageHeader(),
	}
	return &dec
}
func (ego *HandlerRegistration) NeoO1L15COT15DecodeClientHandler() *O1L15COT15CodecClientHandler {
	dec := O1L15COT15CodecClientHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
		_packetHeader:       message_buffer.NeoMessageHeader(),
	}
	return &dec
}

var _ IClientCodecHandler = &O1L15COT15CodecClientHandler{}
