package transcomm

import (
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
)

type O1L15COT15CodecServerHandler struct {
	_memoryLow           bool
	_keepalive           *KeepAlive
	_connection          *TCPServerConnection
	_hdrDeserializeCache []byte
}

func (ego *O1L15COT15CodecServerHandler) OnKeepAlive(ts int64, delta int32) {
	if ego._keepalive != nil {
		ego._keepalive.OnRoundTripBack(ts)
		if delta >= 0 {
			ego._connection._profiler.GetRTTProf().OnUpdate(delta)
			ego._connection._server.Log(core.LL_DEBUG, "conn %s prof: %s", ego._connection.String(), ego._connection._profiler.String())
		}
	}
}

func (ego *O1L15COT15CodecServerHandler) Pulse(conn IConnection, nowTs int64) {
	if ego._keepalive != nil {
		rc := ego._keepalive.Pulse(conn, nowTs)
		if core.IsErrType(rc, core.EC_TCP_CONNECT_ERROR) {
			ego._connection.OnConnectingFailed()
		}
	}
}

func (ego *O1L15COT15CodecServerHandler) Reset() {
	if ego._keepalive != nil {
		ego._keepalive.Reset()
	}
}

func (ego *O1L15COT15CodecServerHandler) CheckCompletion(byteBuf *memory.ByteBufferNode) (int64, int16, int32) {
	var idx int64 = 0
	var frameLength int64 = 0
	var totalFrameLen int64 = 0
	var cmd int16 = 0
	var st int8
	cur := byteBuf
	var offset int64 = 0

	if cur != nil {
		offset = cur.ReadPos()
	}

	for cur != nil {
		if offset >= cur.ReadAvailable() && cur.Next() == nil {
			return -1, cmd, core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		if frameLength == 0 { //header parse
			i0 := memory.BytesToInt16BE(cur.DataRef(), offset)
			i1 := memory.BytesToInt16BE(cur.DataRef(), offset+2)
			frameLength = int64(i0 & 0x7FFF)
			cmd = i1 & 0x7FFF
			o1 := int8(i0 >> 15 & 0x1)
			o2 := int8(i1 >> 15 & 0x1)
			st = (o1 << 1) | o2
			offset += message_buffer.O1L15O1T15_HEADER_SIZE
			if cmd != 32767 {
				fmt.Printf("xxx\n")
			}
		}
		rl := cur.ReadAvailableByOffset(offset)
		if frameLength <= rl {
			totalFrameLen += frameLength
			offset += frameLength
			if st == 0 {
				return totalFrameLen, cmd, core.MkSuccess(0)

			} else if st == 1 {
				return totalFrameLen, cmd, core.MkSuccess(0)

			} else if st == 3 {
				if offset < cur.Capacity() {
					frameLength = 0
					continue
				}

			} else if st == 2 {

				if offset < cur.Capacity() {
					frameLength = 0
					continue
				}
			} else {
				panic("type error")
			}

			frameLength = 0

		} else {
			totalFrameLen += rl
			frameLength -= rl
			offset = 0
		}

		cur = cur.Next()
		offset = 0
		idx++
	}
	return -1, cmd, core.MkErr(core.EC_TRY_AGAIN, 0)
}

func (ego *O1L15COT15CodecServerHandler) OnReceive(connection *TCPServerConnection) (any, int32) {
	bodyLen, cmd, rc := ego.CheckCompletion(ego._connection._recvBufferList.Front())
	if core.Err(rc) {
		return nil, rc
	}

	msg := messages.GetDefaultMessageBufferDeserializationMapper().Deserialize(cmd, ego._connection._recvBufferList, bodyLen)
	if msg == nil {
		connection._server.Log(core.LL_ERR, "Deserialize Message (CMD:%d) error.", cmd)
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}

	rc = GetDefaultMessageHandlerMapper().Handle(connection, msg)
	if core.IsErrType(rc, core.EC_ALREADY_DONE) {
		return nil, core.MkSuccess(0)
	}
	return msg, core.MkSuccess(0)
}

func (ego *O1L15COT15CodecServerHandler) OnSend(connection *TCPServerConnection, a any, bFlush bool) int32 {
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
		ego._connection._server.Log(core.LL_DEBUG, "conn <%s> sent %d bytes", ego._connection.String(), bs)
		if core.Err(rc) {
			return rc
		}
	}

	return core.MkSuccess(0)
}

func (ego *O1L15COT15CodecServerHandler) CheckLowMemory() {
	if ego._memoryLow {
		ego._memoryLow = false
	}
}

func (ego *O1L15COT15CodecServerHandler) OnLowMemory() {
	ego._memoryLow = true
}

func (ego *HandlerRegistration) NeoO1L15COT15DecodeServerHandler(c *TCPServerConnection) *O1L15COT15CodecServerHandler {
	dec := O1L15COT15CodecServerHandler{
		_memoryLow:           false,
		_connection:          c,
		_hdrDeserializeCache: make([]byte, message_buffer.O1L15O1T15_HEADER_SIZE),
	}

	if c != nil {
		if c.KeepAliveConfig().Enable {
			dec._keepalive = NeoKeepAlive(c.KeepAliveConfig(), true)
		}
	}

	return &dec
}

var _ IServerCodecHandler = &O1L15COT15CodecServerHandler{}
