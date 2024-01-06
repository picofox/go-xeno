package transcomm

import (
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type O1L15COT15CodecClientHandler struct {
	_memoryLow           bool
	_keepalive           *KeepAlive
	_connection          *TCPClientConnection
	_hdrDeserializeCache []byte
}

func (ego *O1L15COT15CodecClientHandler) CheckCompletion(byteBuf *memory.ByteBufferNode) (int64, int16, int32) {
	var rc int32 = core.MkSuccess(0)
	var partBodyLen int64 = 0
	var currentFrameLength int64 = 0
	var currentSplitType int8 = 0
	var bodyIndex int64 = 0
	var command int16 = -1

	if byteBuf == nil {
		return partBodyLen, command, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	byteBuf, bodyIndex, currentFrameLength, command, currentSplitType, rc = messages.PeekHeaderContent(ego._hdrDeserializeCache, byteBuf, byteBuf.ReadPos())
	if core.Err(rc) {
		return partBodyLen, command, rc
	}
	if currentFrameLength <= 0 {
		return 0, command, core.MkSuccess(0)
	}

	if currentSplitType == message_buffer.PACKET_SPLITION_TYPE_NONE {
		leftInCurBuffer := byteBuf.ReadAvailable() - bodyIndex
		if leftInCurBuffer >= currentFrameLength {
			return currentFrameLength, command, core.MkSuccess(0)
		}
		partBodyLen = leftInCurBuffer
		byteBuf = byteBuf.Next()
		for byteBuf != nil {
			if partBodyLen+byteBuf.Capacity() >= currentFrameLength {
				partBodyLen += (currentFrameLength - partBodyLen) //todo use abs value currentFrameLength
				return partBodyLen, command, core.MkSuccess(0)
			} else {
				partBodyLen += byteBuf.Capacity()
			}
			byteBuf = byteBuf.Next()
		}
		return 0, command, core.MkErr(core.EC_TRY_AGAIN, 0)

	} else {
		var fakeReaderPos int64 = byteBuf.ReadPos()
		var rBodyLenLarge int64 = 0
		var curBodyLen int64 = 0

		for byteBuf != nil {
			curBodyLen = byteBuf.Capacity() - (fakeReaderPos + bodyIndex)

			if partBodyLen+curBodyLen >= currentFrameLength {
				partBodyLen += currentFrameLength - partBodyLen
				rBodyLenLarge += partBodyLen
				if currentSplitType == message_buffer.PACKET_SPLITION_TYPE_END {
					return rBodyLenLarge, command, core.MkSuccess(0)
				}
				fakeReaderPos = fakeReaderPos + bodyIndex + curBodyLen
				byteBuf, bodyIndex, currentFrameLength, _, currentSplitType, rc = messages.PeekHeaderContent(ego._hdrDeserializeCache, byteBuf, fakeReaderPos)
				fmt.Printf("-> st=%d\n", currentSplitType)
				if core.Err(rc) {
					return rBodyLenLarge, command, rc
				}
				partBodyLen = 0
				fakeReaderPos = 0
				continue

			} else {
				partBodyLen += curBodyLen
				bodyIndex = 0
			}

			byteBuf = byteBuf.Next()
			fakeReaderPos = 0
		}
		return rBodyLenLarge, command, core.MkErr(core.EC_TRY_AGAIN, 0)
	}
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

func (ego *O1L15COT15CodecClientHandler) OnSend(connection *TCPClientConnection, a any, bFlush bool) int32 {
	var message = a.(message_buffer.INetMessage)
	_, bLen, rc := message.PiecewiseSerialize(ego._connection._sendBufferList)
	if bLen != message.BodyLength() {
		//todo remove this check to boost perfermance
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	if core.Err(rc) {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}

	if bFlush {
		var bs int64 = 0
		bs, rc = ego._connection.FlushSendingBuffer()
		ego._connection._client.Log(core.LL_DEBUG, "conn <%s> sent %d bytes", ego._connection.String(), bs)
		if core.Err(rc) {
			return rc
		}
	}

	return core.MkSuccess(0)
}

func (ego *O1L15COT15CodecClientHandler) Reset() {
	if ego._keepalive != nil {
		ego._keepalive.Reset()
	}
}

func (ego *O1L15COT15CodecClientHandler) OnReceive(connection *TCPClientConnection) (any, int32) {
	bodyLen, cmd, rc := ego.CheckCompletion(ego._connection._recvBufferList.Front())
	if core.Err(rc) {
		return nil, rc
	}

	msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(cmd, ego._connection._recvBufferList, bodyLen)
	if msg == nil {
		connection._client.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", cmd)
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}

	rc = GetDefaultMessageHandlerMapper().Handle(connection, msg)
	if core.IsErrType(rc, core.EC_ALREADY_DONE) {
		return nil, core.MkSuccess(0)
	}
	return msg, core.MkSuccess(0)
}

func (ego *O1L15COT15CodecClientHandler) CheckLowMemory() {
	if ego._memoryLow {
		ego._memoryLow = false
	}
}

func (ego *O1L15COT15CodecClientHandler) OnLowMemory() {
	ego._memoryLow = true
}

func (ego *HandlerRegistration) NeoO1L15COT15DecodeClientHandler(c *TCPClientConnection) *O1L15COT15CodecClientHandler {
	dec := O1L15COT15CodecClientHandler{
		_memoryLow:           false,
		_connection:          c,
		_hdrDeserializeCache: make([]byte, message_buffer.O1L15O1T15_HEADER_SIZE),
	}

	if c != nil {
		if c.KeepAliveConfig().Enable {
			dec._keepalive = NeoKeepAlive(c.KeepAliveConfig(), false)
		}
	}

	return &dec
}

var _ IClientCodecHandler = &O1L15COT15CodecClientHandler{}
