package messages

import (
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

func IsMessageComplete(buffer memory.IByteBuffer) (bool, int16, int32) {
	lenAndO1, rc := buffer.PeekInt16(0)
	if core.Err(rc) {
		return false, -1, rc
	}
	cmdAndO2, rc := buffer.PeekInt16(0)
	if core.Err(rc) {
		return false, -1, rc
	}
	o1 := lenAndO1>>15&0x1 == 1
	//o2 := cmdAndO2>>15&0x1 == 1
	cmd := int16(cmdAndO2 & 0x7FFF)
	l := int64(int16(lenAndO1 & 0x7FFF))
	if !o1 {
		if buffer.ReadAvailable() >= message_buffer.O1L15O1T15_HEADER_SIZE+l {
			return false, cmd, core.MkSuccess(0)
		} else {
			return false, -1, core.MkErr(core.EC_TRY_AGAIN, 1)
		}
	} else {
		off := int64(message_buffer.O1L15O1T15_HEADER_SIZE + l)
		if buffer.ReadAvailable() < off+8 {
			return false, -1, core.MkErr(core.EC_TRY_AGAIN, 2)
		}
		eLen, rc := buffer.PeekInt64(off)
		if core.Err(rc) {
			return false, -1, core.MkErr(core.EC_TRY_AGAIN, 3)
		}
		if buffer.ReadAvailable() < off+8+eLen {
			return false, -1, core.MkErr(core.EC_TRY_AGAIN, 4)
		}
		return true, cmd, core.MkSuccess(0)

	}
}

type O1L15O1T15DeserializationHelper struct {
	_command         int16
	_isMulti         bool
	_logicDataLength int16
	_appDataLength   int64
	_extDataLength   int64
	_lastHeaderPos   int64
	_extraLengthPos  int64
	_buffer          memory.IByteBuffer
	_temp            []byte
}

func O1L15O1T15DeserializationHelperCreator() any {
	return &O1L15O1T15DeserializationHelper{
		_command:         -1,
		_isMulti:         false,
		_logicDataLength: 0,
		_appDataLength:   0,
		_extDataLength:   0,
		_lastHeaderPos:   -1,
		_extraLengthPos:  -1,
		_buffer:          nil,
		_temp:            make([]byte, 8),
	}
}

var sO1L15O1T15DeserializationHelperCache *memory.ObjectCache[O1L15O1T15DeserializationHelper] = memory.NeoObjectCache[O1L15O1T15DeserializationHelper](16, O1L15O1T15DeserializationHelperCreator)

func (ego *O1L15O1T15DeserializationHelper) _init(buffer memory.IByteBuffer, cmd int16, isMulti bool) int32 {
	ego._command = cmd
	ego._isMulti = isMulti
	ego._logicDataLength = 0
	ego._appDataLength = 0
	ego._extDataLength = 0
	ego._lastHeaderPos = buffer.WritePos()
	ego._extraLengthPos = -1
	ego._buffer = buffer
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15DeserializationHelper) String() string {
	var ss strings.Builder
	ss.WriteString("\nAppDataLength=")
	ss.WriteString(strconv.FormatInt(ego._appDataLength, 10))
	ss.WriteString("\nLogicDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._logicDataLength)))
	ss.WriteString("\nExtDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._extDataLength)))
	ss.WriteString("\n_lastHeaderPos=")
	ss.WriteString(strconv.FormatInt(ego._lastHeaderPos, 10))
	ss.WriteString("\n_extraLengthPos=")
	ss.WriteString(strconv.FormatInt(ego._extraLengthPos, 10))
	return ss.String()
}

func InitializeDeserialization(buffer memory.IByteBuffer, cmd int16, isMulti bool) (*O1L15O1T15DeserializationHelper, int32) {
	helper := sO1L15O1T15DeserializationHelperCache.Get()
	return helper, helper._init(buffer, cmd, isMulti)
}

func (ego *O1L15O1T15DeserializationHelper) FinalizeSerialization() int32 {
	defer sO1L15O1T15DeserializationHelperCache.Put(ego)
	return 0
}

func (ego *O1L15O1T15DeserializationHelper) ReadRawBytes(bs []byte, baOff int64, readLength int64) int32 {
	ego._buffer.ReadInt32()
	if !ego._isMulti {
		_, rc := ego._buffer.ReadRawBytes(bs, baOff, readLength, true)
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		return core.MkSuccess(0)
	} else {
		var idx int64 = 0
		for readLength > 0 {
			curTurnReadBytes := min(readLength, message_buffer.MAX_PACKET_BODY_SIZE)
			_, rc := ego._buffer.ReadRawBytes(bs, baOff+idx, curTurnReadBytes, true)
			if core.Err(rc) {
				return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
			}
			idx += curTurnReadBytes
			readLength -= curTurnReadBytes
			if readLength == 0 {
				return core.MkSuccess(0)
			}
			if ego._extDataLength <= 0 {
				ego._extDataLength, rc = ego._buffer.ReadInt64()
				if core.Err(rc) {
					return core.MkErr(core.EC_INCOMPLETE_DATA, 2)
				}
			}
			lenAndO1, rc := ego._buffer.ReadInt16()
			if core.Err(rc) {
				return core.MkErr(core.EC_INCOMPLETE_DATA, 3)
			}
			_, rc = ego._buffer.ReadInt16()
			if core.Err(rc) {
				return core.MkErr(core.EC_INCOMPLETE_DATA, 4)
			}
			//o1 := lenAndO1>>15&0x1 == 1
			//o2 := cmdAndO2>>15&0x1 == 1
			//cmd := int16(cmdAndO2 & 0x7FFF)
			l := int64(int16(lenAndO1 & 0x7FFF))
			if readLength < l {
				panic("[SNH] ")
				return core.MkErr(core.EC_INCOMPLETE_DATA, 5)
			}
		}
	}
	return core.MkSuccess(0)
}
