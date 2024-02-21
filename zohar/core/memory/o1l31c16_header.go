package memory

import (
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type O1L31C16Header struct {
	_oAndLen uint32
	_command uint16
}

func (ego *O1L31C16Header) Set(groupType int8, length int32, cmd uint16) {
	ego._oAndLen = uint32(groupType)<<31 | (uint32(length) & 0x7FFFFFFF)
	ego._command = cmd
}

func (ego *O1L31C16Header) SetRaw(oal uint32, cmd uint16) {
	ego._oAndLen = oal
	ego._command = cmd
}

func (ego *O1L31C16Header) HeaderLength() int64 {
	return 6
}

func (ego *O1L31C16Header) GroupType() int8 {
	return int8(ego._oAndLen >> 31)
}

func (ego *O1L31C16Header) SetGroupType(groupType int8) {
	ego._oAndLen = uint32(groupType)<<31 | (uint32(ego._oAndLen) & 0x7FFFFFFF)
}

func (ego *O1L31C16Header) BodyLength() int64 {
	return int64(ego._oAndLen & 0x7FFFFFFF)
}

func (ego *O1L31C16Header) Command() uint16 {
	return ego._command
}

func (ego *O1L31C16Header) SetCommand(cmd uint16) {
	ego._command = cmd
}

func (ego *O1L31C16Header) SetBodyLength(l int64) {
	ego._oAndLen = uint32(ego.GroupType()) | (uint32(l) & 0x7FFFFFFF)
}

func (ego *O1L31C16Header) BeginSerializing(buffer IByteBuffer) (int64, int64, int32) {
	sPos := buffer.WritePos()
	rc := buffer.WriteUInt32(ego._oAndLen)
	if core.Err(rc) {
		return sPos, 0, rc
	}
	rc = buffer.WriteUInt16(ego._command)
	if core.Err(rc) {
		return sPos, 0, rc
	}
	return sPos, 6, rc
}

func (ego *O1L31C16Header) EndSerializing(buffer IByteBuffer, headerPos int64, totalLength int64) (int64, int32) {
	var rc int32 = core.MkSuccess(0)
	ego.SetBodyLength(totalLength)
	rc = buffer.SetUInt32(headerPos, ego._oAndLen)
	if core.Err(rc) {
		return 0, rc
	}
	return datatype.INT64_SIZE, rc
}

func (ego *O1L31C16Header) BeginDeserializing(buffer IByteBuffer, validate bool) (int64, int32) {
	//buffer.ReaderSeek(BUFFER_SEEK_CUR, 6)
	return 0, core.MkSuccess(0)
}

func (ego *O1L31C16Header) EndDeserializing(buffer IByteBuffer) int32 {
	return core.MkSuccess(0)
}

func (ego *O1L31C16Header) String() string {
	return fmt.Sprintf("%d,%d,%d", ego.GroupType(), ego.Command(), ego.BodyLength())
}

var _ ISerializationHeader = &O1L31C16Header{}

func NeoO1L31C16Header(oal uint32, cmd uint16) *O1L31C16Header {
	return &O1L31C16Header{
		_oAndLen: oal,
		_command: cmd,
	}
}

func NeoO1L31C16HeaderFromBytes(bs []byte) *O1L31C16Header {
	return NeoO1L31C16Header(BytesToUInt32BE(&bs, 0), BytesToUInt16BE(&bs, 4))
}

func O1L31C16HeaderFromBuffer(buffer IByteBuffer) (*O1L31C16Header, int32) {
	if buffer.ReadAvailable() >= 6 {
		ol, _ := buffer.PeekUInt32(0)
		l := int32(ol & 0x7FFFFFFF)
		if buffer.ReadAvailable() >= int64(l)+6 {
			c, _ := buffer.PeekUInt16(datatype.INT32_SIZE)
			header := &O1L31C16Header{
				_oAndLen: ol,
				_command: c,
			}
			return header, core.MkSuccess(0)
		}
	}
	return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
}
