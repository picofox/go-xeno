package memory

import "container/list"

const (
	BUFFER_SEEK_CUR = 0
	BUFFER_SEEK_SET = 1
)

type IByteBuffer interface {
	Capacity() int64
	ReadAvailable() int64
	WriteAvailable() int64
	Clear()
	PeekRawBytes(int64, []byte, int64, int64, bool) (int64, int32)
	ReadRawBytes([]byte, int64, int64, bool) (int64, int32)
	WriteRawBytes([]byte, int64, int64) int32
	PeekFloat32(int64) (float32, int32)
	ReadFloat32() (float32, int32)
	WriteFloat32(float32) int32
	PeekFloat64(int64) (float64, int32)
	ReadFloat64() (float64, int32)
	WriteFloat64(float64) int32
	PeekBool(int64) (bool, int32)
	ReadBool() (bool, int32)
	WriteBool(bool) int32
	PeekInt8(int64) (int8, int32)
	ReadInt8() (int8, int32)
	WriteInt8(int8) int32
	PeekUInt8(int64) (uint8, int32)
	ReadUInt8() (uint8, int32)
	WriteUInt8(uint8) int32
	PeekInt16(int64) (int16, int32)
	ReadInt16() (int16, int32)
	WriteInt16(int16) int32
	PeekUInt16(int64) (uint16, int32)
	ReadUInt16() (uint16, int32)
	WriteUInt16(uint16) int32
	PeekInt32(int64) (int32, int32)
	ReadInt32() (int32, int32)
	WriteInt32(int32) int32
	PeekUInt32(int64) (uint32, int32)
	ReadUInt32() (uint32, int32)
	WriteUInt32(uint32) int32
	PeekInt64(int64) (int64, int32)
	ReadInt64() (int64, int32)
	WriteInt64(int64) int32
	PeekUInt64(int64) (uint64, int32)
	ReadUInt64() (uint64, int32)
	WriteUInt64(uint64) int32
	PeekBytes(int64) ([]byte, int32)
	ReadBytes() ([]byte, int32)
	WriteBytes([]byte) int32
	PeekString(int64) (string, int32)
	ReadString() (string, int32)
	WriteString(string) int32
	ReadPos() int64
	WritePos() int64
	ReaderSeek(int, int64) bool
	WriterSeek(int, int64) int32

	WriteStrings([]string) int32
	ReadStrings() ([]string, int32)
	WriteBytess([][]byte) int32
	ReadBytess() ([][]byte, int32)

	SetRawBytes(int64, []byte, int64, int64) int32
	SetRawBytesByNode(*list.Element, int64, []byte, int64, int64) int32

	SetInt32(int64, int32) int32
	SetInt32ByNode(*list.Element, int64, int32) int32

	SetInt64(int64, int64) int32
	SetInt64ByNode(*list.Element, int64, int64) int32
}
