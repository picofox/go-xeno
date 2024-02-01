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
	_appDataLength   int64
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
		_appDataLength:   0,
		_extDataLength:   0,
		_headerPos:       -1,
		_extraLengthPos:  -1,
		_buffer:          nil,
		_temp:            make([]byte, 8),
	}
}

var sO1L15O1T15SerializationHelperCache *memory.ObjectCache[O1L15O1T15SerializationHelper] = memory.NeoObjectCache[O1L15O1T15SerializationHelper](16, O1L15O1T15SerializationHelperCreator)

func (ego *O1L15O1T15SerializationHelper) _init(cmd int16, buffer memory.IByteBuffer) int32 {
	ego._command = cmd
	ego._logicDataLength = 0
	ego._appDataLength = 0
	ego._extDataLength = 0
	ego._headerPos = buffer.WritePos()
	ego._extraLengthPos = -1
	ego.SetHeader(false, false, 0)
	ego._buffer = buffer
	rc := buffer.WriteRawBytes(ego._headerData[:], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15SerializationHelper) MarkHeaderAsFinished() {
	//var lenAndO0 uint16 = uint16(ego._headerData[0]&0xff)<<8 | uint16(ego._headerData[1]&0xff)
	//var cmdAndO1 uint16 = uint16(ego._headerData[2]&0xff)<<8 | uint16(ego._headerData[3]&0xff)

	ego._headerData[0] = ego._headerData[0] ^ (1 << 7)
	ego._headerData[2] |= 1 << 7
}

func (ego *O1L15O1T15SerializationHelper) SetHeaderLength(length int16) {
	u0 := int16(ego._headerData[0]&0xff)<<8 | int16(ego._headerData[1]&0xff)
	if uint16(u0)&0x8000 != 0 {
		length = int16(uint16(length&0x7FFF) | uint16(1)<<15)
	}
	ego._headerData[0] = byte((length >> 8) & 0xFF)
	ego._headerData[1] = byte((length & 0xFF) & 0xFF)
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
	ego._headerData[3] = byte((cmdAndO1 >> 8) & 0xFF)
}

func (ego *O1L15O1T15SerializationHelper) String() string {
	var ss strings.Builder
	ss.WriteString(fmt.Sprintf("Header= %02x %02x %02x %02x", ego._headerData[0], ego._headerData[1], ego._headerData[2], ego._headerData[3]))
	ss.WriteString("\nAppDataLength=")
	ss.WriteString(strconv.FormatInt(ego._appDataLength, 10))
	ss.WriteString("\nLogicDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._logicDataLength)))
	ss.WriteString("\nExtDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._extDataLength)))
	ss.WriteString("\n_headerPos=")
	ss.WriteString(strconv.FormatInt(ego._headerPos, 10))
	ss.WriteString("\n_extraLengthPos=")
	ss.WriteString(strconv.FormatInt(ego._extraLengthPos, 10))
	return ss.String()
}

func InitializeSerialization(buffer memory.IByteBuffer, cmd int16) (*O1L15O1T15SerializationHelper, int32) {
	helper := sO1L15O1T15SerializationHelperCache.Get()
	return helper, helper._init(cmd, buffer)
}

func (ego *O1L15O1T15SerializationHelper) FinalizeSerialization() int32 {
	defer sO1L15O1T15SerializationHelperCache.Put(ego)
	if ego._appDataLength <= message_buffer.MAX_PACKET_BODY_SIZE {
		if ego._headerPos < 0 {
			panic("invalid ctx._headerPos")
		}
		ego.SetHeaderLength(int16(ego._appDataLength))
		rc := ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], ego._headerPos, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}
	} else {
		if ego._extraLengthPos <= 0 {
			panic("ctx._extraLengthPos error")
		} else {
			rc := ego._buffer.SetInt64(ego._extraLengthPos, ego._extDataLength)
			if core.Err(rc) {
				return rc
			}
		}
		if ego._logicDataLength > 0 {
			ego.SetHeader(false, true, ego._logicDataLength)
			rc := ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
			if core.Err(rc) {
				return rc
			}
		} else if ego._headerPos == 0 {
			ego.MarkHeaderAsFinished()
			rc := ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
			if core.Err(rc) {
				return rc
			}
		}
	}
	return 0
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
	blen := len(srcBA)
	rc := ego.WriteInt32(int32(blen))
	if core.Err(rc) {
		return rc
	}
	if blen > 0 {
		rc = ego.WriteRawBytes(srcBA, 0, int64(blen))
	}
	return rc
}

func (ego *O1L15O1T15SerializationHelper) WriteInt32(iv int32) int32 {
	memory.Int32IntoBytesBE(iv, &ego._temp, 0)
	return ego.WriteRawBytes(ego._temp, 0, datatype.INT32_SIZE)
}

func (ego *O1L15O1T15SerializationHelper) WriteRawBytes(bs []byte, srcOff int64, wLen int64) int32 {
	if wLen <= 0 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	if ego._logicDataLength == message_buffer.MAX_PACKET_BODY_SIZE {
		ego.SetHeader(true, true, ego._logicDataLength) //finish previous LP
		rc := ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], ego._headerPos, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}
		if ego._extraLengthPos < 0 {
			ego._extraLengthPos = ego._buffer.WritePos()
			ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
		}
		ego._logicDataLength = 0
		ego.SetHeader(true, true, 0)
		ego._headerPos = ego._buffer.WritePos() //start a neo LB
		rc = ego._buffer.WriteRawBytes(ego._headerData[:], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}
		if ego._extraLengthPos >= 0 {
			ego._extDataLength += message_buffer.O1L15O1T15_HEADER_SIZE
		}

	}
	//app data
	var curTurnWriteByte int64 = wLen
	if curTurnWriteByte > int64(message_buffer.MAX_PACKET_BODY_SIZE-ego._logicDataLength) {
		curTurnWriteByte = int64(message_buffer.MAX_PACKET_BODY_SIZE - ego._logicDataLength)
	}
	var idx int64 = 0
	for wLen > 0 {
		rc := ego._buffer.WriteRawBytes(bs, srcOff+idx, curTurnWriteByte)
		if core.Err(rc) {
			return rc
		}
		wLen -= curTurnWriteByte
		ego._logicDataLength += int16(curTurnWriteByte)
		ego._appDataLength += curTurnWriteByte
		if ego._extraLengthPos >= 0 {
			ego._extDataLength += curTurnWriteByte
		}
		if wLen == 0 {
			return rc
		}
		idx += curTurnWriteByte
		if ego._logicDataLength != message_buffer.MAX_PACKET_BODY_SIZE {
			panic("xxx")
		}
		//process LB header logic
		ego.SetHeader(true, true, ego._logicDataLength) //finish previous LP
		rc = ego._buffer.SetRawBytes(ego._headerPos, ego._headerData[0:message_buffer.O1L15O1T15_HEADER_SIZE], ego._headerPos, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}
		if ego._extraLengthPos < 0 {
			ego._extraLengthPos = ego._buffer.WritePos()
			ego._buffer.WriteInt64(0) //Extra Size Field follows 1st logic packet,but only once
		}
		ego._logicDataLength = 0
		ego.SetHeader(true, true, 0)
		ego._headerPos = ego._buffer.WritePos() //start a neo LB
		rc = ego._buffer.WriteRawBytes(ego._headerData[:], 0, message_buffer.O1L15O1T15_HEADER_SIZE)
		if core.Err(rc) {
			return rc
		}
		if ego._extraLengthPos >= 0 {
			ego._extDataLength += message_buffer.O1L15O1T15_HEADER_SIZE
		}

		curTurnWriteByte = min(wLen, message_buffer.MAX_PACKET_BODY_SIZE)
	}
	return core.MkSuccess(0)
}
