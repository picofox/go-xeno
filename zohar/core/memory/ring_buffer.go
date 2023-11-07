package memory

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type RingBuffer struct {
	_capacity int64
	_beginPos int64
	_length   int64
	_data     []byte
	_b8Cache  []byte
}

func (ego *RingBuffer) ExpandTo(neoCapacity int64) int64 {
	if ego._capacity >= neoCapacity {
		return 0
	} else {
		wp := ego.WritePos()
		neoData := make([]byte, neoCapacity)
		if wp >= ego._beginPos {
			copy(neoData, ego._data[ego._beginPos:ego._length])
		} else {
			lenToEnd := ego._capacity - ego._beginPos
			copy(neoData, ego._data[ego._beginPos:ego._capacity])
			copy(neoData[lenToEnd:], ego._data[0:wp])
		}
		oldCapa := ego._capacity
		ego._beginPos = 0
		ego._capacity = neoCapacity
		ego._data = neoData
		return neoCapacity - oldCapa
	}
}

func (ego *RingBuffer) checkSpace(extraLength int64) int64 {
	wa := ego.WriteAvailable()
	if wa >= extraLength {
		return 0
	} else {
		atLeastToAlloc := int64(0)
		if ego._capacity < 512 {
			atLeastToAlloc = ego._capacity
		} else {
			atLeastToAlloc = 512
		}
		if atLeastToAlloc < extraLength-wa {
			atLeastToAlloc = extraLength - wa
		}
		wp := ego.WritePos()
		totolLen := ego._capacity + atLeastToAlloc
		neoData := make([]byte, totolLen)
		if wp >= ego._beginPos {
			copy(neoData, ego._data[ego._beginPos:ego._length])
		} else {
			lenToEnd := ego._capacity - ego._beginPos
			copy(neoData, ego._data[ego._beginPos:ego._capacity])
			copy(neoData[lenToEnd:], ego._data[0:wp])
		}
		ego._beginPos = 0
		ego._capacity = totolLen
		ego._data = neoData
		return atLeastToAlloc
	}
}

func (ego *RingBuffer) adjustWritePos(length int64) bool {
	if ego._length+length > ego._capacity {
		return false
	}
	ego._length = ego._length + length
	return true
}

func (ego *RingBuffer) adjustReadPos(length int64) bool {
	if ego._length >= length {
		lenToEnd := ego._capacity - ego._beginPos
		if lenToEnd < length {
			part2len := length - lenToEnd
			ego._beginPos = part2len
		} else {
			ego._beginPos += length
			if ego._beginPos == ego._capacity {
				ego._beginPos = 0
			}
		}
		ego._length -= length
		return true
	}
	return false
}

func (ego *RingBuffer) fillCachePeek(lenToEnd int64, dtSize int64) (int64, int64) {
	if lenToEnd > 0 {
		copy(ego._b8Cache, ego._data[ego._beginPos:ego._beginPos+lenToEnd])
	}
	part2Len := dtSize - lenToEnd
	copy(ego._b8Cache[lenToEnd:], ego._data[0:part2Len])

	return part2Len, ego._length - dtSize
}

func (ego *RingBuffer) fillCache(lenToEnd int64, dtSize int64) {
	if lenToEnd > 0 {
		copy(ego._b8Cache, ego._data[ego._beginPos:ego._beginPos+lenToEnd])
	}
	part2Len := dtSize - lenToEnd
	copy(ego._b8Cache[lenToEnd:], ego._data[0:part2Len])
	ego._beginPos = part2Len
	ego._length -= dtSize
}

func (ego *RingBuffer) loadFromCache(wp int64, dtSize int64) {
	lenToEnd := ego._capacity - wp
	lenFromBegin := dtSize - lenToEnd
	if lenToEnd > 0 {
		copy(ego._data[wp:], ego._b8Cache[0:lenToEnd])
	}
	copy(ego._data[0:], ego._b8Cache[lenToEnd:lenToEnd+lenFromBegin])
	ego._length += dtSize
}

func (ego *RingBuffer) InternalData() []byte {
	return ego._data
}

func (ego *RingBuffer) Capacity() int64 {
	return ego._capacity
}
func (ego *RingBuffer) BytesRef() ([]byte, []byte) {
	if ego._beginPos+ego._length > ego._capacity {
		firstPartLen := ego._capacity - ego._beginPos
		remainLen := ego._length - firstPartLen
		return ego._data[ego._beginPos : ego._beginPos+firstPartLen], ego._data[0:remainLen]
	} else {
		return ego._data[ego._beginPos:ego._length], nil
	}
}

func (ego *RingBuffer) ReadAvailable() int64 {
	return ego._length
}

func (ego *RingBuffer) WriteAvailable() int64 {
	return ego._capacity - ego._length
}

func (ego *RingBuffer) Clear() {
	ego._beginPos = 0
	ego._length = 0
}

func (ego *RingBuffer) WritePos() int64 {
	wp := ego._beginPos + ego._length
	if wp >= ego._capacity {
		wp -= ego._capacity
	}
	return wp
}

func (ego *RingBuffer) PeekBytes(ba []byte, baOff int64, peekLength int64, isStrict bool) int64 {
	if ego._length < peekLength {
		if isStrict {
			return 0
		} else {
			peekLength = ego._length
		}
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < peekLength {
		if lenToEnd > 0 {
			copy(ba[baOff:], ego._data[ego._beginPos:ego._beginPos+lenToEnd])
		}
		part2len := peekLength - lenToEnd
		copy(ba[baOff+lenToEnd:], ego._data[0:part2len])
		return peekLength
	} else {
		copy(ba[baOff:], ego._data[ego._beginPos:ego._beginPos+peekLength])
		return peekLength
	}
}

func (ego *RingBuffer) ReadBytes(ba []byte, baOff int64, readLength int64, isStrict bool) int64 {
	if ego._length < readLength {
		if isStrict {
			return 0
		} else {
			readLength = ego._length
		}
	}

	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < readLength {
		if lenToEnd > 0 {
			copy(ba[baOff:], ego._data[ego._beginPos:ego._beginPos+lenToEnd])
		}
		part2len := readLength - lenToEnd
		copy(ba[baOff+lenToEnd:], ego._data[0:part2len])
		ego._beginPos = part2len
		ego._length -= readLength
		return readLength
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

func (ego *RingBuffer) WriteBytes(ba []byte, srcOff int64, srcLength int64) int32 {
	if srcLength < 0 {
		srcLength = int64(len(ba))
	}
	if ego.checkSpace(srcLength) < 0 {
		return core.MkErr(core.EC_NULL_VALUE, 0)
	}

	wp := ego.WritePos()
	if wp+srcLength <= ego._capacity {
		copy(ego._data[wp:], ba[srcOff:srcOff+srcLength])
		ego._length += srcLength
	} else {
		lenToEnd := ego._capacity - wp
		if lenToEnd > 0 {
			copy(ego._data[wp:], ba[srcOff:srcOff+lenToEnd])
		}
		copy(ego._data[0:], ba[srcOff+lenToEnd:srcOff+srcLength])
		ego._length += srcLength
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekBool() (bool, int32, int64, int64) {
	iv, rc, beg, rLen := ego.PeekInt8()
	if core.Err(rc) {
		return false, rc, beg, rLen
	}
	if iv != 0 {
		return true, rc, beg, rLen
	}
	return false, rc, beg, rLen
}

func (ego *RingBuffer) ReadBool() (bool, int32) {
	iv, rc := ego.ReadInt8()
	if core.Err(rc) {
		return false, rc
	}
	if iv != 0 {
		return true, rc
	}
	return false, rc
}

func (ego *RingBuffer) WriteBool(b bool) int32 {
	if ego.checkSpace(1) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
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

func (ego *RingBuffer) PeekInt8() (int8, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := int8(ego._data[ego._beginPos])
	beg := ego._beginPos + 1
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 1
}

func (ego *RingBuffer) ReadInt8() (int8, int32) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := int8(ego._data[ego._beginPos])
	ego._beginPos++
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	ego._length--
	return rc, core.MkSuccess(0)
}

func (ego *RingBuffer) WriteInt8(iv int8) int32 {
	if ego.checkSpace(1) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	ego._data[wp] = byte(iv)
	ego._length++
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekUInt8() (uint8, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	rc := uint8(ego._data[ego._beginPos])
	beg := ego._beginPos + 1
	if beg == ego._capacity {
		beg = 0
	}
	return rc, core.MkSuccess(0), beg, ego._length - 1
}

func (ego *RingBuffer) ReadUInt8() (uint8, int32) {
	readable := ego.ReadAvailable()
	if readable < 1 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	rc := uint8(ego._data[ego._beginPos])
	ego._beginPos++
	if ego._beginPos == ego._capacity {
		ego._beginPos = 0
	}
	return rc, core.MkSuccess(0)
}

func (ego *RingBuffer) WriteUInt8(iv uint8) int32 {
	if ego.checkSpace(1) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	ego._data[wp] = byte(iv)
	ego._length++
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekInt16() (int16, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 2 {
		beg, rlen := ego.fillCachePeek(lenToEnd, 2)
		return datatype.BytesToInt16BE(&ego._b8Cache, 0), core.MkSuccess(0), beg, rlen
	} else {
		rc := datatype.BytesToInt16BE(&ego._data, ego._beginPos)
		beg := ego._beginPos + 2
		if beg == ego._capacity {
			beg = 0
		}
		return rc, core.MkSuccess(0), beg, ego._length - 2
	}
}

func (ego *RingBuffer) ReadInt16() (int16, int32) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 2 {
		ego.fillCache(lenToEnd, 2)
		return datatype.BytesToInt16BE(&ego._b8Cache, 0), core.MkSuccess(0)
	} else {
		rc := datatype.BytesToInt16BE(&ego._data, ego._beginPos)
		ego._beginPos += 2
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= 2
		return rc, core.MkSuccess(0)
	}
}

func (ego *RingBuffer) WriteInt16(iv int16) int32 {
	if ego.checkSpace(2) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	if wp+2 <= ego._capacity {
		datatype.Int16IntoBytesBE(iv, &ego._data, wp)
		ego._length += 2
	} else {
		datatype.Int16IntoBytesBE(iv, &ego._b8Cache, 0)
		ego.loadFromCache(wp, 2)
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekUInt16() (uint16, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 2 {
		beg, rlen := ego.fillCachePeek(lenToEnd, 2)
		return datatype.BytesToUInt16BE(&ego._b8Cache, 0), core.MkSuccess(0), beg, rlen
	} else {
		rc := datatype.BytesToUInt16BE(&ego._data, ego._beginPos)
		beg := ego._beginPos + 2
		if beg == ego._capacity {
			beg = 0
		}
		return rc, core.MkSuccess(0), beg, ego._length - 2
	}
}

func (ego *RingBuffer) ReadUInt16() (uint16, int32) {
	readable := ego.ReadAvailable()
	if readable < 2 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 2 {
		ego.fillCache(lenToEnd, 2)
		return datatype.BytesToUInt16BE(&ego._b8Cache, 0), core.MkSuccess(0)
	} else {
		rc := datatype.BytesToUInt16BE(&ego._data, ego._beginPos)
		ego._beginPos += 2
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= 2
		return rc, core.MkSuccess(0)
	}
}

func (ego *RingBuffer) WriteUInt16(iv uint16) int32 {
	if ego.checkSpace(2) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	if wp+2 <= ego._capacity {
		datatype.UInt16IntoBytesBE(iv, &ego._data, wp)
		ego._length += 2
	} else {
		datatype.UInt16IntoBytesBE(iv, &ego._b8Cache, 0)
		ego.loadFromCache(wp, 2)
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekInt32() (int32, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 4 {
		beg, rlen := ego.fillCachePeek(lenToEnd, 4)
		return datatype.BytesToInt32BE(&ego._b8Cache, 0), core.MkSuccess(0), beg, rlen
	} else {
		rc := datatype.BytesToInt32BE(&ego._data, ego._beginPos)
		beg := ego._beginPos + 4
		if beg == ego._capacity {
			beg = 0
		}
		return rc, core.MkSuccess(0), beg, ego._length - 4
	}
}

func (ego *RingBuffer) ReadInt32() (int32, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 4 {
		ego.fillCache(lenToEnd, 4)
		return datatype.BytesToInt32BE(&ego._b8Cache, 0), core.MkSuccess(0)
	} else {
		rc := datatype.BytesToInt32BE(&ego._data, ego._beginPos)
		ego._beginPos += 4
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= 4
		return rc, core.MkSuccess(0)
	}
}

func (ego *RingBuffer) WriteInt32(iv int32) int32 {
	if ego.checkSpace(4) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	if wp+4 <= ego._capacity {
		datatype.Int32IntoBytesBE(iv, &ego._data, wp)
		ego._length += 4
	} else {
		datatype.Int32IntoBytesBE(iv, &ego._b8Cache, 0)
		ego.loadFromCache(wp, 4)
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekUInt32() (uint32, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 4 {
		beg, rlen := ego.fillCachePeek(lenToEnd, 4)
		return datatype.BytesToUInt32BE(&ego._b8Cache, 0), core.MkSuccess(0), beg, rlen
	} else {
		rc := datatype.BytesToUInt32BE(&ego._data, ego._beginPos)
		beg := ego._beginPos + 4
		if beg == ego._capacity {
			beg = 0
		}
		return rc, core.MkSuccess(0), beg, ego._length - 4
	}
}

func (ego *RingBuffer) ReadUInt32() (uint32, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 4 {
		ego.fillCache(lenToEnd, 4)
		return datatype.BytesToUInt32BE(&ego._b8Cache, 0), core.MkSuccess(0)
	} else {
		rc := datatype.BytesToUInt32BE(&ego._data, ego._beginPos)
		ego._beginPos += 4
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= 4
		return rc, core.MkSuccess(0)
	}
}

func (ego *RingBuffer) WriteUInt32(iv uint32) int32 {
	if ego.checkSpace(4) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	if wp+4 <= ego._capacity {
		datatype.UInt32IntoBytesBE(iv, &ego._data, wp)
		ego._length += 4
	} else {
		datatype.UInt32IntoBytesBE(iv, &ego._b8Cache, 0)
		ego.loadFromCache(wp, 4)
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekInt64() (int64, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 8 {
		beg, rlen := ego.fillCachePeek(lenToEnd, 8)
		return datatype.BytesToInt64BE(&ego._b8Cache, 0), core.MkSuccess(0), beg, rlen
	} else {
		rc := datatype.BytesToInt64BE(&ego._data, ego._beginPos)
		beg := ego._beginPos + 8
		if beg == ego._capacity {
			beg = 0
		}
		return rc, core.MkSuccess(0), beg, ego._length - 8
	}
}

func (ego *RingBuffer) ReadInt64() (int64, int32) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 8 {
		ego.fillCache(lenToEnd, 8)
		return datatype.BytesToInt64BE(&ego._b8Cache, 0), core.MkSuccess(0)
	} else {
		rc := datatype.BytesToInt64BE(&ego._data, ego._beginPos)
		ego._beginPos += 8
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= 8
		return rc, core.MkSuccess(0)
	}
}

func (ego *RingBuffer) WriteInt64(iv int64) int32 {
	if ego.checkSpace(8) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	if wp+8 <= ego._capacity {
		datatype.Int64IntoBytesBE(iv, &ego._data, wp)
		ego._length += 8
	} else {
		datatype.Int64IntoBytesBE(iv, &ego._b8Cache, 0)
		ego.loadFromCache(wp, 8)
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) PeekUInt64() (uint64, int32, int64, int64) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1), -1, -1
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 8 {
		beg, rlen := ego.fillCachePeek(lenToEnd, 8)
		return datatype.BytesToUInt64BE(&ego._b8Cache, 0), core.MkSuccess(0), beg, rlen
	} else {
		rc := datatype.BytesToUInt64BE(&ego._data, ego._beginPos)
		beg := ego._beginPos + 8
		if beg == ego._capacity {
			beg = 0
		}
		return rc, core.MkSuccess(0), beg, ego._length - 8
	}
}

func (ego *RingBuffer) ReadUInt64() (uint64, int32) {
	readable := ego.ReadAvailable()
	if readable < 8 {
		return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	lenToEnd := ego._capacity - ego._beginPos
	if lenToEnd < 8 {
		ego.fillCache(lenToEnd, 8)
		return datatype.BytesToUInt64BE(&ego._b8Cache, 0), core.MkSuccess(0)
	} else {
		rc := datatype.BytesToUInt64BE(&ego._data, ego._beginPos)
		ego._beginPos += 8
		if ego._beginPos == ego._capacity {
			ego._beginPos = 0
		}
		ego._length -= 8
		return rc, core.MkSuccess(0)
	}
}

func (ego *RingBuffer) WriteUInt64(iv uint64) int32 {
	if ego.checkSpace(8) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	wp := ego.WritePos()
	if wp+8 <= ego._capacity {
		datatype.UInt64IntoBytesBE(iv, &ego._data, wp)
		ego._length += 8
	} else {
		datatype.UInt64IntoBytesBE(iv, &ego._b8Cache, 0)
		ego.loadFromCache(wp, 8)
	}
	return core.MkSuccess(0)
}

func (ego *RingBuffer) ReadString() (string, int32) {
	readable := ego.ReadAvailable()
	if readable < 4 {
		return "", core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	bLen, rc, updateBeg, updateLen := ego.PeekInt32()
	if core.Err(rc) {
		return "", rc
	}
	if readable < int64(4+bLen) {
		return "", core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	}
	ego._beginPos = updateBeg
	ego._length = updateLen
	if bLen > 0 {
		rBA := make([]byte, bLen)
		if ego.ReadBytes(rBA, 0, int64(bLen), true) != int64(bLen) {
			return "", core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		return string(rBA), core.MkSuccess(0)
	}

	return "", core.MkSuccess(0)
}

func (ego *RingBuffer) WriteString(str string) int32 {
	ba := []byte(str)
	blen := len(ba)
	if blen > datatype.INT32_MAX {
		return core.MkErr(core.EC_INDEX_OOB, 0)
	}
	if ego.checkSpace(int64(4+blen)) < 0 {
		return core.MkErr(core.EC_RESPACE_FAILED, 1)
	}
	ego.WriteInt32(int32(blen))
	if blen > 0 {
		ego.WriteBytes(ba, 0, int64(blen))
	}

	return core.MkSuccess(0)
}

func NeoByteBuffer(capacity int64) *RingBuffer {
	bf := &RingBuffer{
		_capacity: capacity,
		_beginPos: 0,
		_length:   0,
		_data:     make([]byte, capacity),
		_b8Cache:  make([]byte, 8),
	}
	return bf
}
