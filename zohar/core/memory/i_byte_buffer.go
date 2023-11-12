package memory

type IByteBuffer interface {
	Capacity() int64
	ResizeTo(int64) int64
	ReadAvailable() int64
	WriteAvailable() int64
	Clear()
	PeekRawBytes([]byte, int64, int64, bool) (int64, int64, int64)
	ReadRawBytes([]byte, int64, int64, bool) int64
	WriteRawBytes([]byte, int64, int64) int32
	PeekFloat32() (float32, int32, int64, int64)
	ReadFloat32() (float32, int32)
	WriteFloat32(float32) int32
	PeekFloat64() (float64, int32, int64, int64)
	ReadFloat64() (float64, int32)
	WriteFloat64(float64) int32
	PeekBool() (bool, int32, int64, int64)
	ReadBool() (bool, int32)
	WriteBool(b bool) int32
	PeekInt8() (int8, int32, int64, int64)
	ReadInt8() (int8, int32)
	WriteInt8(int8) int32
	PeekUInt8() (uint8, int32, int64, int64)
	ReadUInt8() (uint8, int32)
	WriteUInt8(uint8) int32
	PeekInt16() (int16, int32, int64, int64)
	ReadInt16() (int16, int32)
	WriteInt16(int16) int32
	PeekUInt16() (uint16, int32, int64, int64)
	ReadUInt16() (uint16, int32)
	WriteUInt16(uint16) int32
	PeekInt32() (int32, int32, int64, int64)
	ReadInt32() (int32, int32)
	WriteInt32(int32) int32
	PeekUInt32() (uint32, int32, int64, int64)
	ReadUInt32() (uint32, int32)
	WriteUInt32(uint32) int32
	PeekInt64() (int64, int32, int64, int64)
	ReadInt64() (int64, int32)
	WriteInt64(int64) int32
	PeekUInt64() (uint64, int32, int64, int64)
	ReadUInt64() (uint64, int32)
	WriteUInt64(uint64) int32
	PeekBytes() ([]byte, int32, int64, int64)
	ReadBytes() ([]byte, int32)
	WriteBytes([]byte) int32
	PeekString() (string, int32, int64, int64)
	ReadString() (string, int32)
	WriteString(string) int32
	BytesRef() ([]byte, []byte)
}
