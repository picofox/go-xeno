package messages

import (
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

func IsMessageComplete(buffer memory.IByteBuffer) (bool, int16, int16, int64, int32) {
	lenAndO1, rc := buffer.PeekInt16(0)
	if core.Err(rc) {
		return false, -1, -1, -1, rc
	}
	cmdAndO2, rc := buffer.PeekInt16(0)
	if core.Err(rc) {
		return false, -1, -1, -1, rc
	}
	o1 := lenAndO1>>15&0x1 == 1
	isInternal := cmdAndO2>>15&0x1 == 1
	cmd := int16(cmdAndO2 & 0x7FFF)
	l := int16(lenAndO1 & 0x7FFF)
	if !o1 {
		if buffer.ReadAvailable() >= message_buffer.O1L15O1T15_HEADER_SIZE+int64(l) {
			return isInternal, cmd, l, 0, core.MkSuccess(0)
		} else {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 1)
		}
	} else {
		off := int64(message_buffer.O1L15O1T15_HEADER_SIZE + l)
		if buffer.ReadAvailable() < off+8 {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 2)
		}
		eLen, rc := buffer.PeekInt64(off)
		if core.Err(rc) {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_INCOMPLETE_DATA, 3)
		}
		if buffer.ReadAvailable() < off+8+eLen {
			return isInternal, -1, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 4)
		}
		return isInternal, cmd, l, eLen, core.MkSuccess(0)
	}
}

type O1L15O1T15DeserializationHelper struct {
	_command         int16
	_isInternal      bool
	_logicDataLength int16
	_extDataLength   int64
	_buffer          memory.IByteBuffer
	_temp            []byte
}

func O1L15O1T15DeserializationHelperCreator() any {
	return &O1L15O1T15DeserializationHelper{
		_command:         -1,
		_logicDataLength: 0,
		_extDataLength:   0,
		_buffer:          nil,
		_temp:            make([]byte, 8),
	}
}

var sO1L15O1T15DeserializationHelperCache *memory.ObjectCache[O1L15O1T15DeserializationHelper] = memory.NeoObjectCache[O1L15O1T15DeserializationHelper](16, O1L15O1T15DeserializationHelperCreator)

func (ego *O1L15O1T15DeserializationHelper) _init(buffer memory.IByteBuffer, isLarge bool, isInternal bool, cmd int16, logicLength int16, extraLength int64) int32 {
	ego._command = cmd
	ego._isInternal = isInternal
	ego._logicDataLength = logicLength
	ego._extDataLength = extraLength
	ego._buffer = buffer
	ego._buffer.ReadInt32()
	return core.MkSuccess(0)
}

func (ego *O1L15O1T15DeserializationHelper) String() string {
	var ss strings.Builder
	ss.WriteString("\nIsInternal=")
	ss.WriteString(strconv.FormatBool(ego._isInternal))
	ss.WriteString("\nLogicDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._logicDataLength)))
	ss.WriteString("\nExtDataLength=")
	ss.WriteString(strconv.Itoa(int(ego._extDataLength)))
	return ss.String()
}

func InitializeDeserialization(buffer memory.IByteBuffer, isLarge bool, isInternal bool, cmd int16, logicLength int16, extLength int64) (*O1L15O1T15DeserializationHelper, int32) {
	helper := sO1L15O1T15DeserializationHelperCache.Get()
	return helper, helper._init(buffer, isLarge, isInternal, cmd, logicLength, extLength)
}

func (ego *O1L15O1T15DeserializationHelper) FinalizeSerialization() int32 {
	defer sO1L15O1T15DeserializationHelperCache.Put(ego)
	return 0
}

func (ego *O1L15O1T15DeserializationHelper) ReadRawBytes(bs []byte, baOff int64, readLength int64) int32 {
	curTurnReadBytes := min(int64(ego._logicDataLength), readLength)
	if curTurnReadBytes > 0 {
		_, rc := ego._buffer.ReadRawBytes(bs, baOff, curTurnReadBytes, true)
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
		}
		ego._logicDataLength -= int16(curTurnReadBytes)
		if ego._logicDataLength == 0 {
			if !ego._buffer.ReaderSeek(memory.BUFFER_SEEK_CUR, 8) {
				return core.MkErr(core.EC_INCOMPLETE_DATA, 2)
			}
		}
		readLength -= curTurnReadBytes
		if readLength == 0 {
			return core.MkSuccess(0)
		}
	} else if curTurnReadBytes < 0 {
		panic("curTurnReadBytes < 0")
	}

	if readLength > 0 {
		if ego._extDataLength > 0 {
			_, rc := ego._buffer.ReadRawBytes(bs, baOff+curTurnReadBytes, readLength, true)
			if core.Err(rc) {
				return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
			}
		} else {
			panic("_extDataLength < 0")
		}
	}
	
	return core.MkSuccess(0)
}
