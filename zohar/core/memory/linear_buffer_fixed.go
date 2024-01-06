package memory

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type LinearBufferFixed struct {
	_capacity int64
	_beginPos int64
	_length   int64
	_data     []byte
	_cache    []byte
}

func (ego *LinearBufferFixed) DataRef() *[]byte {
	return &ego._data
}

func (ego *LinearBufferFixed) WriteStrings(strs []string) int32 {
	if strs == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(strs)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteString(strs[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) ReadStrings() ([]string, int32) {
	l, rc := ego.ReadInt32()
	if l < 0 {
		return nil, core.MkSuccess(0)
	} else if l == 0 {
		return make([]string, 0), core.MkSuccess(0)
	}
	r := make([]string, l)
	for i := int32(0); i < l; i++ {
		r[i], rc = ego.ReadString()
		if core.Err(rc) {
			return nil, rc
		}
	}
	return r, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) Compact() int64 {
	if ego._beginPos > 0 {
		if ego._length > 0 {
			copy(ego._data[0:], ego._data[ego._beginPos:ego._beginPos+ego._length])
		}
		ego._beginPos = 0
	}
	return ego.WriteAvailable()
}

func (ego *LinearBufferFixed) ResizeTo(newSize int64) int64 {
	return 0
}

func (ego *LinearBufferFixed) checkSpace(extraLength int64) int64 {
	wa := ego.WriteAvailable()
	if wa >= extraLength {
		return 0
	} else {
		return wa - extraLength
	}
}

func (ego *LinearBufferFixed) WritePos() int64 {
	wp := ego._beginPos + ego._length
	return wp
}

func (ego *LinearBufferFixed) ReadPos() int64 {
	return ego._beginPos
}

func (ego *LinearBufferFixed) WriterSeek(whence int, offset int64) bool {
	if whence == BUFFER_SEEK_CUR {
		absPosTest := ego.WritePos() + offset
		return ego.WriterSeek(whence, absPosTest)
	} else if whence == BUFFER_SEEK_SET {
		if offset < 0 || offset > ego._capacity {
			return false
		}

		ego._length = offset - ego._beginPos
	}
	return true
}

func (ego *LinearBufferFixed) ReaderSeek(whence int, offset int64) bool {
	if whence == BUFFER_SEEK_CUR {
		absPosTest := ego._beginPos + offset
		if absPosTest < 0 || absPosTest >= ego._beginPos+ego._length {
			return false
		}
		ego._beginPos += offset
		if ego._beginPos >= ego.Capacity() {
			ego._beginPos = 0
		}
		ego._length -= offset
	} else if whence == BUFFER_SEEK_SET {
		delta := offset - ego._beginPos
		return ego.ReaderSeek(BUFFER_SEEK_CUR, delta)

	}
	return true
}
func (ego *LinearBufferFixed) SliceOf(length int64) []byte {
	if ego._length < 1 {
		return ego._data[ego._beginPos:ego._beginPos]
	}
	ba := ego._data[ego._beginPos : ego._beginPos+length]
	ego.ReaderSeek(BUFFER_SEEK_CUR, length)
	return ba
}

func (ego *LinearBufferFixed) Capacity() int64 {
	return ego._capacity
}

func (ego *LinearBufferFixed) ReadAvailable() int64 {
	return ego._length
}

func (ego *LinearBufferFixed) WriteAvailable() int64 {
	return ego._capacity - ego._length
}

func (ego *LinearBufferFixed) Clear() {
	ego._beginPos = 0
	ego._length = 0
}

func (ego *LinearBufferFixed) Reset() {
	ego._beginPos = 0
	ego._length = 0
	ego._capacity = 0
	ego._data = make([]byte, 0)
}

func (ego *LinearBufferFixed) PeekRawBytes(ba []byte, baOff int64, peekLength int64, isStrict bool) (int64, int64, int64) {
	if ego._length < peekLength {
		if isStrict {
			return 0, -1, -1
		} else {
			peekLength = ego._length
		}
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < peekLength {
		return -1, -1, -1
	} else {
		copy(ba[baOff:], ego._data[ego._beginPos:ego._beginPos+peekLength])
		beginPos := ego._beginPos + peekLength
		if beginPos == ego._capacity {
			beginPos = 0
		}
		return peekLength, beginPos, ego._length - peekLength
	}
}

func (ego *LinearBufferFixed) ReadRawBytes(ba []byte, baOff int64, readLength int64, isStrict bool) int64 {
	if readLength == 0 {
		return 0
	}
	if ego._length < readLength {
		if isStrict {
			return 0
		} else {
			readLength = ego._length
		}
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < readLength {
		return -1
	} else {
		copy(ba[baOff:], ego._data[ego._beginPos:ego._beginPos+readLength])
		ego._beginPos += readLength
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= readLength
		return readLength
	}

}

func (ego *LinearBufferFixed) WriteRawBytes(ba []byte, srcOff int64, srcLength int64) int32 {
	if srcLength < 0 {
		srcLength = int64(len(ba))
	}
	if ego.checkSpace(srcLength) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 0)
	}
	wp := ego.WritePos()
	copy(ego._data[wp:], ba[srcOff:srcOff+srcLength])
	ego._length += srcLength
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekFloat32() (float32, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	f32 := BytesToFloat32BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 4
	if beg == ego._capacity {
		beg = 0
	}
	return f32, core.MkSuccess(0), beg, ego._length - 4
}

func (ego *LinearBufferFixed) ReadFloat32() (float32, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToFloat32BE(&ego._data, ego._beginPos)
	ego._beginPos += 4
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 4
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteFloat32(fv float32) int32 {
	if ego.checkSpace(4) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	Float32IntoBytesBE(fv, &ego._data, wp)
	ego._length += 4
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekFloat64() (float64, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	f64 := BytesToFloat64BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 8
	if beg == ego._capacity {
		beg = 0
	}
	return f64, core.MkSuccess(0), beg, ego._length - 8
}

func (ego *LinearBufferFixed) ReadFloat64() (float64, int32) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToFloat64BE(&ego._data, ego._beginPos)
	ego._beginPos += 8
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 8
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteFloat64(fv float64) int32 {
	if ego.checkSpace(8) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	Float64IntoBytesBE(fv, &ego._data, wp)
	ego._length += 8
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekBool() (bool, int32, int64, int64) {
	iv, rc, beg, rLen := ego.PeekInt8()
	if core.Err(rc) {
		return false, rc, beg, rLen
	}
	if iv != 0 {
		return true, rc, beg, rLen
	}
	return false, rc, beg, rLen
}

func (ego *LinearBufferFixed) ReadBool() (bool, int32) {
	iv, rc := ego.ReadInt8()
	if core.Err(rc) {
		return false, rc
	}
	if iv != 0 {
		return true, rc
	}
	return false, rc
}

func (ego *LinearBufferFixed) WriteBool(b bool) int32 {
	if ego.checkSpace(1) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	if b {
		ego._data[wp] = uint8(1)
	} else {
		ego._data[wp] = uint8(0)
	}
	ego._length++
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekInt8() (int8, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := int8(ego._data[ego._beginPos])
	return rc, core.MkSuccess(0), ego._beginPos + 1, ego._length - 1
}

func (ego *LinearBufferFixed) ReadInt8() (int8, int32) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := int8(ego._data[ego._beginPos])
	ego._beginPos++
	ego._length--
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt8(iv int8) int32 {
	if ego.checkSpace(1) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	ego._data[wp] = byte(iv)
	ego._length++
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekUInt8() (uint8, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := uint8(ego._data[ego._beginPos])
	return rc, core.MkSuccess(0), ego._beginPos + 1, ego._length - 1
}

func (ego *LinearBufferFixed) ReadUInt8() (uint8, int32) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := uint8(ego._data[ego._beginPos])
	ego._beginPos++
	ego._length--
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteUInt8(u uint8) int32 {
	if ego.checkSpace(1) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	ego._data[wp] = byte(u)
	ego._length++
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekInt16() (int16, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := BytesToInt16BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 2
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 2
}

func (ego *LinearBufferFixed) ReadInt16() (int16, int32) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToInt16BE(&ego._data, ego._beginPos)
	ego._beginPos += 2
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 2
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt16(iv int16) int32 {
	if ego.checkSpace(2) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	if wp+2 <= ego._capacity {
		Int16IntoBytesBE(iv, &ego._data, wp)
		ego._length += 2
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekUInt16() (uint16, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := BytesToUInt16BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 2
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 2
}

func (ego *LinearBufferFixed) ReadUInt16() (uint16, int32) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToUInt16BE(&ego._data, ego._beginPos)
	ego._beginPos += 2
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 2
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteUInt16(uv uint16) int32 {
	if ego.checkSpace(2) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	if wp+2 <= ego._capacity {
		UInt16IntoBytesBE(uv, &ego._data, wp)
		ego._length += 2
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekInt32() (int32, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := BytesToInt32BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 4
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 4
}

func (ego *LinearBufferFixed) ReadInt32() (int32, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToInt32BE(&ego._data, ego._beginPos)
	ego._beginPos += 4
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 4
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt32(iv int32) int32 {
	if ego.checkSpace(4) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	Int32IntoBytesBE(iv, &ego._data, wp)
	ego._length += 4
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekUInt32() (uint32, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := BytesToUInt32BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 4
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 4
}

func (ego *LinearBufferFixed) ReadUInt32() (uint32, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToUInt32BE(&ego._data, ego._beginPos)
	ego._beginPos += 4
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 4
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteUInt32(uv uint32) int32 {
	if ego.checkSpace(4) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	UInt32IntoBytesBE(uv, &ego._data, wp)
	ego._length += 4
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekInt64() (int64, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := BytesToInt64BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 8
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 8
}

func (ego *LinearBufferFixed) ReadInt64() (int64, int32) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToInt64BE(&ego._data, ego._beginPos)
	ego._beginPos += 8
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 8
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt64(iv int64) int32 {
	if ego.checkSpace(8) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	Int64IntoBytesBE(iv, &ego._data, wp)
	ego._length += 8
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt16Begin(iv int16, length int64) int32 {
	if ego.checkSpace(length) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	Int16IntoBytesBE(iv, &ego._cache, 0)
	if length > 0 {
		ego.WriteRawBytes(ego._cache, 0, length)
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt32Begin(iv int32, length int64) int32 {
	if ego.checkSpace(length) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	Int32IntoBytesBE(iv, &ego._cache, 0)
	if length > 0 {
		ego.WriteRawBytes(ego._cache, 0, length)
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteFloat32Begin(iv float32, length int64) int32 {
	if ego.checkSpace(length) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	Float32IntoBytesBE(iv, &ego._cache, 0)
	if length > 0 {
		ego.WriteRawBytes(ego._cache, 0, length)
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteInt64Begin(iv int64, length int64) int32 {
	if ego.checkSpace(length) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	Int64IntoBytesBE(iv, &ego._cache, 0)
	if length > 0 {
		ego.WriteRawBytes(ego._cache, 0, length)
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteFloat64Begin(iv float64, length int64) int32 {
	if ego.checkSpace(length) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	Float64IntoBytesBE(iv, &ego._cache, 0)
	if length > 0 {
		ego.WriteRawBytes(ego._cache, 0, length)
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteTrivialMiddle(off int64, l int64) int32 {
	if ego.checkSpace(l) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	if l > 0 {
		ego.WriteRawBytes(ego._cache, off, l)
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteTrivialEnd(off int64, SIZEOFTYPE int64) (int64, int32) {
	l := SIZEOFTYPE - off
	if l < 1 {
		return 0, core.MkSuccess(0)
	}
	if ego.checkSpace(l) < 0 {
		return 0, core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	if l > 0 {
		ego.WriteRawBytes(ego._cache, off, l)
	}
	return l, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekUInt64() (uint64, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := BytesToUInt64BE(&ego._data, ego._beginPos)
	beg := ego._beginPos + 8
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 8
}

func (ego *LinearBufferFixed) ReadUInt64() (uint64, int32) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := BytesToUInt64BE(&ego._data, ego._beginPos)
	ego._beginPos += 8
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length -= 8
	return rc, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteUInt64(uv uint64) int32 {
	if ego.checkSpace(8) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	wp := ego.WritePos()
	UInt64IntoBytesBE(uv, &ego._data, wp)
	ego._length += 8
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekBytes() ([]byte, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	bLen, rc, updateBeg, updateLen := ego.PeekInt32()
	if core.Err(rc) {
		return nil, rc, -1, -1
	}
	if readable < int64(4+bLen) {
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	saveBeg := ego._beginPos
	saveLen := ego._length
	ego._beginPos = updateBeg
	ego._length = updateLen
	if bLen > 0 {
		rBA := make([]byte, bLen)
		pLen, beg, rLen := ego.PeekRawBytes(rBA, 0, int64(bLen), true)
		ego._beginPos = saveBeg
		ego._length = saveLen
		if pLen != int64(bLen) {
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2), -1, -1
		}

		return rBA, core.MkSuccess(0), beg, rLen
	} else if bLen == 0 {
		ego._beginPos = saveBeg
		ego._length = saveLen
		return make([]byte, 0), core.MkSuccess(0), updateBeg, updateLen
	}

	ego._beginPos = saveBeg
	ego._length = saveLen
	return nil, core.MkSuccess(0), updateBeg, updateLen
}

func (ego *LinearBufferFixed) ReadBytes() ([]byte, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	bLen, rc, updateBeg, updateLen := ego.PeekInt32()
	if core.Err(rc) {
		return nil, rc
	}
	if readable < int64(4+bLen) {
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	ego._beginPos = updateBeg
	ego._length = updateLen
	if bLen > 0 {
		rBA := make([]byte, bLen)
		if ego.ReadRawBytes(rBA, 0, int64(bLen), true) != int64(bLen) {
			return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		return rBA, core.MkSuccess(0)
	} else if bLen == 0 {
		return make([]byte, 0), core.MkSuccess(0)
	}

	return nil, core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteBytes(srcBA []byte) int32 {
	blen := len(srcBA)
	if blen > datatype.INT32_MAX {
		return core.MkErr(core.EC_INDEX_OOB, 0)
	}
	if ego.checkSpace(int64(4+blen)) < 0 {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	ego.WriteInt32(int32(blen))
	if blen > 0 {
		ego.WriteRawBytes(srcBA, 0, int64(blen))
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) PeekString() (string, int32, int64, int64) {
	rBA, rc, beg, rLen := ego.PeekBytes()
	if core.Err(rc) {
		return "", rc, -1, -1
	}
	if rBA == nil {
		return "", core.MkSuccess(0), beg, rLen
	} else if len(rBA) == 0 {
		return "", core.MkSuccess(0), beg, rLen
	}
	return string(rBA), core.MkSuccess(0), beg, rLen

}

func (ego *LinearBufferFixed) ReadString() (string, int32) {
	rBA, rc := ego.ReadBytes()
	if core.Err(rc) {
		return "", rc
	}

	if rBA == nil {
		return "", core.MkSuccess(0)
	} else if len(rBA) == 0 {
		return "", core.MkSuccess(0)
	}
	return string(rBA), core.MkSuccess(0)
}

func (ego *LinearBufferFixed) WriteString(str string) int32 {
	ba := []byte(str)
	rc := ego.WriteBytes(ba)
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *LinearBufferFixed) BytesRef(length int64) ([]byte, []byte) {
	if ego._length < 1 {
		return nil, nil
	}
	if length < 0 {
		length = ego._length
	} else if length > ego._length {
		panic("buffer out of scope")
	}

	return ego._data[ego._beginPos : ego._beginPos+length], nil
}

func (ego *LinearBufferFixed) InternalData() *[]byte {
	return &ego._data
}

var _ IByteBuffer = &LinearBufferFixed{}

func NeoLinearBufferFixed(capacity int64) *LinearBufferFixed {
	bf := &LinearBufferFixed{
		_capacity: capacity,
		_beginPos: 0,
		_length:   0,
		_data:     make([]byte, capacity),
		_cache:    make([]byte, 8),
	}
	return bf
}

func NeoLinearBufferFixedCopy(ba []byte, capacity int64) *LinearBufferFixed {
	if ba == nil || len(ba) > int(capacity) {
		return NeoLinearBufferFixed(capacity)
	}
	l := len(ba)
	bf := &LinearBufferFixed{
		_capacity: capacity,
		_beginPos: 0,
		_length:   int64(l),
		_data:     make([]byte, capacity),
		_cache:    make([]byte, 8),
	}
	copy(bf._data, ba[0:l])
	return bf
}

func AttachLinearBufferFixed(ba []byte, capacity int64) *LinearBufferFixed {
	if ba == nil || len(ba) > int(capacity) {
		return nil
	}
	bf := &LinearBufferFixed{
		_capacity: capacity,
		_beginPos: 0,
		_length:   int64(len(ba)),
		_data:     ba,
		_cache:    make([]byte, 8),
	}
	return bf
}
