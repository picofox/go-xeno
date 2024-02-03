package messages

import (
	"fmt"
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

type O1L15O1T15SerializationHelper struct {
	_command         int16
	_logicDataLength int16
	_extDataLength   int64
	_headerPos       int64
	_extraLengthPos  int64
	_headerData      [message_buffer.O1L15O1T15_HEADER_SIZE]byte
	_buffer          memory.IByteBuffer
	_temp            []byte
}

func O1L15O1T15SerializationHelperCreator() any {
	return &O1L15O1T15SerializationHelper{
		_command:         -1,
		_logicDataLength: 0,
		_extDataLength:   0,
		_headerPos:       -1,
		_extraLengthPos:  -1,
		_buffer:          nil,
		_temp:            make([]byte, 8),
	}
}

func (ego *O1L15O1T15SerializationHelper) ReadableBytes() int64 {
	return ego._buffer.ReadAvailable()
}

var sO1L15O1T15SerializationHelperCache *memory.ObjectCache[O1L15O1T15SerializationHelper] = memory.NeoObjectCache[O1L15O1T15SerializationHelper](16, O1L15O1T15SerializationHelperCreator)

func (ego *O1L15O1T15SerializationHelper) _init(buffer memory.IByteBuffer, isInternal bool, cmd int16) int32 {
	ego._command = cmd
	ego._logicDataLength = 0
	ego._extDataLength = 0
	ego._headerPos = buffer.WritePos()
	ego._extraLengthPos = -1
	ego.SetHeader(false, isInternal, 0)
	ego._buffer = buffer
	rc := buffer.WriteRawBytes(ego._headerData[:], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15SerializationHelper) MarkHeaderAsLargeMessage() {
	ego._headerData[0] |= 1 << 7
}

func (ego *O1L15O1T15SerializationHelper) SetHeaderLength(length int16) {
	u0 := int16(ego._headerData[0]&0xff)<<8 | int16(ego._headerData[1]&0xff)
	if uint16(u0)&0x8000 != 0 {
		length = int16(uint16(length&0x7FFF) | uint16(1)<<15)
	}
	ego._headerData[0] = byte((length >> 8) & 0xFF)
	ego._headerData[1] = byte((length & 0xFF) & 0xFF)
}

func (ego *O1L15O1T15SerializationHelper) SetHeaderOptAndLength(o0 bool, length int16) {
	var lenAndO0 int16 = length
	if o0 {
		iv := 1 << 15
		lenAndO0 = length | int16(iv)
	}
	ego._headerData[0] = byte((lenAndO0 >> 8) & 0xFF)
	ego._headerData[1] = byte((lenAndO0 & 0xFF) & 0xFF)
}

func (ego *O1L15O1T15SerializationHelper) SetHeader(o0 bool, o1 bool, length int16) {
	var lenAndO0 int16 = length
	var cmdAndO1 int16 = ego._command
	if o0 {
		iv := 1 << 15
		lenAndO0 = length | int16(iv)
	}
	if o1 {
		iv := 1 << 15
		cmdAndO1 = ego._command | int16(iv)
	}

	ego._headerData[0] = byte((lenAndO0 >> 8) & 0xFF)
	ego._headerData[1] = byte((lenAndO0 & 0xFF) & 0xFF)
	ego._headerData[2] = byte((cmdAndO1 >> 8) & 0xFF)
	ego._headerData[3] = byte((cmdAndO1 & 0xFF) & 0xFF)
}

func (ego *O1L15O1T15SerializationHelper) String() string {
	var ss strings.Builder
	ss.WriteString(fmt.Sprintf("Header= %02x %02x %02x %02x", ego._headerData[0], ego._headerData[1], ego._headerData[2], ego._headerData[3]))
	ss.WriteString("\nAppDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._logicDataLength)))
	ss.WriteString("\nExtDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._extDataLength)))
	ss.WriteString("\n_headerPos=")
	ss.WriteString(strconv.FormatInt(ego._headerPos, 10))
	ss.WriteString("\n_extraLengthPos=")
	ss.WriteString(strconv.FormatInt(ego._extraLengthPos, 10))
	return ss.String()
}

func InitializeSerialization(buffer memory.IByteBuffer, isInternal bool, cmd int16) (*O1L15O1T15SerializationHelper, int32) {
	helper := sO1L15O1T15SerializationHelperCache.Get()
	return helper, helper._init(buffer, isInternal, cmd)
}

func (ego *O1L15O1T15SerializationHelper) FinalizeSerialization() int32 {
	defer sO1L15O1T15SerializationHelperCache.Put(ego)
	if ego._extDataLength > 0 {
		rc := ego._buffer.SetInt64(ego._extraLengthPos, ego._extDataLength)
		if core.Err(rc) {
			return rc
		}
		if ego._logicDataLength != message_buffer.MAX_PACKET_BODY_SIZE {
			panic("[SNH] length invalid")
		}
		ego.SetHeaderOptAndLength(true, ego._logicDataLength)
		rc = ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}

	} else {
		ego.SetHeaderOptAndLength(false, ego._logicDataLength)
		rc := ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}
	}

	return core.MkSuccess(0)
}

func (ego *O1L15O1T15SerializationHelper) WriteString(s string) int32 {
	bLen := len(s)
	if bLen == 0 {
		return ego.WriteInt32(0)
	}
	ba := memory.ByteRef(s, 0, int(bLen))
	rc := ego.WriteBytes(ba)
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15SerializationHelper) WriteBytes(srcBA []byte) int32 {
	var rc int32 = 0
	if srcBA == nil {
		rc = ego.WriteInt32(int32(-1))
	} else {
		blen := len(srcBA)
		rc = ego.WriteInt32(int32(blen))
		if core.Err(rc) {
			return rc
		}
		if blen > 0 {
			rc = ego.WriteRawBytes(srcBA, 0, int64(blen))
		}
	}

	return rc
}

func (ego *O1L15O1T15SerializationHelper) WriteBool(b bool) int32 {
	if b {
		return ego.WriteInt8(1)
	} else {
		return ego.WriteInt8(0)
	}
}
func (ego *O1L15O1T15SerializationHelper) WriteInt8(iv int8) int32 {
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), datatype.INT8_SIZE)
	if curTurnWriteByte > 0 {
		rc := ego._buffer.WriteInt8(iv)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength++
		return core.MkSuccess(0)
	}
	rc := ego._buffer.WriteInt8(iv)
	if core.Err(rc) {
		return rc
	}
	ego._extDataLength++

	return core.MkSuccess(0)
}
func (ego *O1L15O1T15SerializationHelper) WriteInt16(iv int16) int32 {
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), datatype.INT16_SIZE)
	if curTurnWriteByte >= datatype.INT16_SIZE {
		rc := ego._buffer.WriteInt16(iv)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength += datatype.INT16_SIZE
		return core.MkSuccess(0)
	} else if curTurnWriteByte > 0 {
		memory.Int16IntoBytesBE(iv, &ego._temp, 0)
		rc := ego._buffer.WriteRawBytes(ego._temp, 0, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength = message_buffer.MAX_PACKET_BODY_SIZE
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0)
		reaminWLen := datatype.INT16_SIZE - curTurnWriteByte
		rc = ego._buffer.WriteRawBytes(ego._temp, curTurnWriteByte, reaminWLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += reaminWLen
	} else {
		rc := ego._buffer.WriteInt16(iv)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += datatype.INT16_SIZE
	}
	return core.MkSuccess(0)
}
func (ego *O1L15O1T15SerializationHelper) WriteInt32(iv int32) int32 {
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), datatype.INT32_SIZE)
	if curTurnWriteByte >= datatype.INT32_SIZE {
		rc := ego._buffer.WriteInt32(iv)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength += datatype.INT32_SIZE
		return core.MkSuccess(0)
	} else if curTurnWriteByte > 0 {
		memory.Int32IntoBytesBE(iv, &ego._temp, 0)
		rc := ego._buffer.WriteRawBytes(ego._temp, 0, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength = message_buffer.MAX_PACKET_BODY_SIZE
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0)
		reaminWLen := datatype.INT32_SIZE - curTurnWriteByte
		rc = ego._buffer.WriteRawBytes(ego._temp, curTurnWriteByte, reaminWLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += reaminWLen
	} else {
		rc := ego._buffer.WriteInt32(iv)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += datatype.INT32_SIZE
	}
	return core.MkSuccess(0)
}
func (ego *O1L15O1T15SerializationHelper) WriteInt64(iv int64) int32 {
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), datatype.INT64_SIZE)
	if curTurnWriteByte >= datatype.INT64_SIZE {
		rc := ego._buffer.WriteInt64(iv)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength += datatype.INT64_SIZE
		return core.MkSuccess(0)
	} else if curTurnWriteByte > 0 {
		memory.Int64IntoBytesBE(iv, &ego._temp, 0)
		rc := ego._buffer.WriteRawBytes(ego._temp, 0, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength = message_buffer.MAX_PACKET_BODY_SIZE
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0)
		reaminWLen := datatype.INT64_SIZE - curTurnWriteByte
		rc = ego._buffer.WriteRawBytes(ego._temp, curTurnWriteByte, reaminWLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += reaminWLen
	} else {
		rc := ego._buffer.WriteInt64(iv)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += datatype.INT64_SIZE
	}
	return core.MkSuccess(0)
}
func (ego *O1L15O1T15SerializationHelper) WriteUInt8(iv uint8) int32 {
	return ego.WriteInt8(int8(iv))
}
func (ego *O1L15O1T15SerializationHelper) WriteUInt16(iv uint16) int32 {
	return ego.WriteInt16(int16(iv))
}
func (ego *O1L15O1T15SerializationHelper) WriteUInt32(iv uint32) int32 {
	return ego.WriteInt32(int32(iv))
}
func (ego *O1L15O1T15SerializationHelper) WriteUInt64(iv uint64) int32 {
	return ego.WriteInt64(int64(iv))
}

func (ego *O1L15O1T15SerializationHelper) WriteRawBytes(bs []byte, srcOff int64, wLen int64) int32 {
	if wLen <= 0 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		if ego._extraLengthPos < 0 {
			panic("xxxx")
		}
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), wLen)
	if curTurnWriteByte >= wLen {
		rc := ego._buffer.WriteRawBytes(bs, srcOff, wLen)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength += int16(wLen)
		return core.MkSuccess(0)
	} else if curTurnWriteByte > 0 {
		rc := ego._buffer.WriteRawBytes(bs, srcOff, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength = message_buffer.MAX_PACKET_BODY_SIZE
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0)
		remainWLen := wLen - curTurnWriteByte
		rc = ego._buffer.WriteRawBytes(bs, srcOff+curTurnWriteByte, remainWLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += remainWLen
	} else {
		rc := ego._buffer.WriteRawBytes(bs, srcOff, wLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += wLen
	}

	return core.MkSuccess(0)
}

func (ego *O1L15O1T15SerializationHelper) WriteFloat32(fv float32) int32 {
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), datatype.FLOAT32_SIZE)
	if curTurnWriteByte >= datatype.FLOAT32_SIZE {
		rc := ego._buffer.WriteFloat32(fv)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength += datatype.FLOAT32_SIZE
		return core.MkSuccess(0)
	} else if curTurnWriteByte > 0 {
		memory.Float32IntoBytesBE(fv, &ego._temp, 0)
		rc := ego._buffer.WriteRawBytes(ego._temp, 0, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength = message_buffer.MAX_PACKET_BODY_SIZE
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0)
		reaminWLen := datatype.FLOAT32_SIZE - curTurnWriteByte
		rc = ego._buffer.WriteRawBytes(ego._temp, curTurnWriteByte, reaminWLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += reaminWLen
	} else {
		rc := ego._buffer.WriteFloat32(fv)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += datatype.FLOAT32_SIZE
	}
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15SerializationHelper) WriteFloat64(fv float64) int32 {
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE && ego._extraLengthPos < 0 {
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
	}
	logicLeft := message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength
	var curTurnWriteByte int64 = min(int64(logicLeft), datatype.FLOAT64_SIZE)
	if curTurnWriteByte >= datatype.INT64_SIZE {
		rc := ego._buffer.WriteFloat64(fv)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength += datatype.FLOAT64_SIZE
		return core.MkSuccess(0)
	} else if curTurnWriteByte > 0 {
		memory.Float64IntoBytesBE(fv, &ego._temp, 0)
		rc := ego._buffer.WriteRawBytes(ego._temp, 0, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		ego._logicDataLength = message_buffer.MAX_PACKET_BODY_SIZE
		ego._extraLengthPos = ego._buffer.WritePos()
		ego._buffer.WriteInt64(0)
		reaminWLen := datatype.FLOAT64_SIZE - curTurnWriteByte
		rc = ego._buffer.WriteRawBytes(ego._temp, curTurnWriteByte, reaminWLen)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += reaminWLen
	} else {
		rc := ego._buffer.WriteFloat64(fv)
		if core.Err(rc) {
			return rc
		}
		ego._extDataLength += datatype.FLOAT64_SIZE
	}
	return core.MkSuccess(0)
}
