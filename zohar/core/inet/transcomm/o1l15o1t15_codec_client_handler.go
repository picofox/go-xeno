package transcomm

import (
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type O1L15COT15CodecClientHandler struct {
	_largeMessageBuffer *memory.LinearBuffer
	_memoryLow          bool
	_packetHeader       message_buffer.MessageHeader
	_keepalive          *KeepAlive
	_connection         *TCPClientConnection
	_sendingBuffer      *memory.LinearBuffer
}

func (ego *O1L15COT15CodecClientHandler) OnKeepAlive(ts int64, delta int32) {
	if ego._keepalive != nil {
		ego._keepalive.OnRoundTripBack(ts)
		if delta >= 0 {
			ego._connection._profiler.GetRTTProf().OnUpdate(delta)
			ego._connection._client.Log(core.LL_DEBUG, "conn %s prof: %s", ego._connection.String(), ego._connection._profiler.String())
		}
	}
}

func (ego *O1L15COT15CodecClientHandler) Pulse(conn IConnection, nowTs int64) {
	if ego._keepalive != nil {
		rc := ego._keepalive.Pulse(conn, nowTs)
		if core.IsErrType(rc, core.EC_TCP_CONNECT_ERROR) {
			ego._connection.OnConnectingFailed()
		}
	}
}

func (ego *O1L15COT15CodecClientHandler) OnSend(connection *TCPClientConnection, a any, bflush bool) int32 {

	var message = a.(message_buffer.INetMessage)
	ego._sendingBuffer.Clear()
	tLen := message.Serialize(ego._sendingBuffer)
	if tLen < 0 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}

	var byteBuf memory.IByteBuffer = connection._sendBuffer
	var cmd int16 = message.Command()
	var offset int64 = 0
	var ba *[]byte = nil
	if tLen <= message_buffer.MAX_BUFFER_MAX_CAPACITY {
		for {
			curBB, rc := CheckByteBufferListNode(connection)
			if core.Err(rc) {
				return rc
			}
			writableBytes := curBB.Buffer().WriteAvailable()
			ba = ego._sendingBuffer.InternalData()
			if writableBytes >= tLen {
				rc = curBB.Buffer().WriteRawBytes(*ba, offset, tLen)
				if bflush {
					ego._connection.FlushSendingBuffer()
				}
				return rc
			} else {
				rc = curBB.Buffer().WriteRawBytes(*ba, offset, writableBytes)
				if core.Err(rc) {
					return rc
				}
				offset += writableBytes
			}
		}

		if bflush {
			ego._connection.FlushSendingBuffer()
		}
		return core.MkSuccess(0)
	} else { //large message
		rIndex := int64(message_buffer.O1L15O1T15_HEADER_SIZE)
		ego._packetHeader.Set(true, false, message_buffer.MAX_PACKET_BODY_SIZE, cmd)
		if !byteBuf.ReaderSeek(memory.BUFFER_SEEK_CUR, message_buffer.O1L15O1T15_HEADER_SIZE) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}

		for {
			_, rc := connection.sendImmediately(ego._packetHeader.Data(), 0, message_buffer.O1L15O1T15_HEADER_SIZE)
			fmt.Printf("write header %s \n", ego._packetHeader.String())
			if core.Err(rc) {
				fmt.Printf("write header %s Failed\n", ego._packetHeader.String())
				return rc
			}
			_, rc = connection.sendImmediately(*byteBuf.InternalData(), byteBuf.ReadPos(), message_buffer.MAX_PACKET_BODY_SIZE)
			fmt.Printf("write body from %d, len %d\n", byteBuf.ReadPos(), message_buffer.MAX_PACKET_BODY_SIZE)
			if core.Err(rc) {
				fmt.Printf("write body from %d, len %d Failed\n", byteBuf.ReadPos(), message_buffer.MAX_PACKET_BODY_SIZE)
				return rc
			}
			rIndex += message_buffer.MAX_PACKET_BODY_SIZE
			byteBuf.ReaderSeek(memory.BUFFER_SEEK_SET, rIndex)

			//next loop use non begin version
			ego._packetHeader.Set(true, true, message_buffer.MAX_PACKET_BODY_SIZE, cmd)

			if byteBuf.ReadAvailable() <= message_buffer.MAX_PACKET_BODY_SIZE {
				break
			}
		}
		ego._packetHeader.Set(false, true, message_buffer.MAX_PACKET_BODY_SIZE, cmd)
		_, rc := connection.sendImmediately(ego._packetHeader.Data(), 0, message_buffer.O1L15O1T15_HEADER_SIZE)
		fmt.Printf("write header %s \n", ego._packetHeader.String())
		if core.Err(rc) {
			fmt.Printf("write header %s Failed\n", ego._packetHeader.String())
			return rc
		}
		_, rc = connection.sendImmediately(*byteBuf.InternalData(), byteBuf.ReadPos(), byteBuf.ReadAvailable())
		fmt.Printf("write body from %d, len %d \n", byteBuf.ReadPos(), message_buffer.MAX_PACKET_BODY_SIZE)
		if core.Err(rc) {
			fmt.Printf("write body from %d, len %d Failed\n", byteBuf.ReadPos(), message_buffer.MAX_PACKET_BODY_SIZE)
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *O1L15COT15CodecClientHandler) Reset() {
	ego._largeMessageBuffer.Reset()
	if ego._keepalive != nil {
		ego._keepalive.Reset()
	}

}

func (ego *O1L15COT15CodecClientHandler) OnReceive(connection *TCPClientConnection) (any, int32) {
	if connection._recvBuffer.ReadAvailable() <= message_buffer.O1L15O1T15_HEADER_SIZE {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	o1AndLen, _, _, _ := connection._recvBuffer.PeekUInt16()
	frameLength := int64(o1AndLen & 0x7FFF)
	if connection._recvBuffer.ReadAvailable() < int64(frameLength)+message_buffer.O1L15O1T15_HEADER_SIZE {
		connection._recvBuffer.ResizeTo(message_buffer.MAX_BUFFER_MAX_CAPACITY)
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
		delta := endPos - beginPos
		if endPos < beginPos {
			delta = connection._recvBuffer.Capacity() - beginPos + endPos
		}
		if delta != frameLength {
			connection._client.Log(core.LL_ERR, "Message (CMD:%d) Length Validation Failed, frame length is %d, but got %d read", cmd, frameLength, delta)
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}

		rc := GetDefaultMessageHandlerMapper().Handle(connection, msg)
		if core.IsErrType(rc, core.EC_ALREADY_DONE) {
			return nil, core.MkSuccess(0)
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

func (ego *HandlerRegistration) NeoO1L15COT15DecodeClientHandler(c *TCPClientConnection) *O1L15COT15CodecClientHandler {
	dec := O1L15COT15CodecClientHandler{
		_largeMessageBuffer: memory.NeoLinearBuffer(0),
		_memoryLow:          false,
		_packetHeader:       message_buffer.NeoMessageHeader(),
		_connection:         c,
		_sendingBuffer:      memory.NeoLinearBuffer(32768),
	}

	if c.KeepAliveConfig().Enable {
		dec._keepalive = NeoKeepAlive(c.KeepAliveConfig(), false)
	}

	return &dec
}

var _ IClientCodecHandler = &O1L15COT15CodecClientHandler{}
