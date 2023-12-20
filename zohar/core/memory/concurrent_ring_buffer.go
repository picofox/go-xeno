package memory

import (
	"sync"
	"xeno/zohar/core"
)

type ConcurrentRingBuffer struct {
	_ringBuffer *RingBuffer
	_lock       sync.RWMutex
}

func (ego *ConcurrentRingBuffer) ReadStrings() ([]string, int32) {
	ego._lock.Lock()
	defer ego._lock.Unlock()

	l, rc := ego._ringBuffer.ReadInt32()
	if l < 0 {
		return nil, core.MkSuccess(0)
	} else if l == 0 {
		return make([]string, 0), core.MkSuccess(0)
	}
	r := make([]string, l)
	for i := int32(0); i < l; i++ {
		r[i], rc = ego._ringBuffer.ReadString()
		if core.Err(rc) {
			return nil, rc
		}
	}
	return r, core.MkSuccess(0)
}

func (ego *ConcurrentRingBuffer) ResizeTo(i int64) int64 {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) PeekBytes() ([]byte, int32, int64, int64) {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) PeekString() (string, int32, int64, int64) {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) ReadString() (string, int32) {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) WriteString(s string) int32 {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) BytesRef(i int64) ([]byte, []byte) {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) ReadPos() int64 {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) ReaderSeek(whence int, offset int64) bool {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) WriterSeek(whence int, offset int64) bool {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) InternalData() *[]byte {
	//TODO implement me
	panic("implement me")
}

func (ego *ConcurrentRingBuffer) WriteStrings(strs []string) int32 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	l := len(strs)
	ego._ringBuffer.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego._ringBuffer.WriteString(strs[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *ConcurrentRingBuffer) Capacity() int64 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	ret := ego._ringBuffer.Capacity()
	return ret
}

func (ego *ConcurrentRingBuffer) ReadAvailable() int64 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	ret := ego._ringBuffer.ReadAvailable()
	return ret
}

func (ego *ConcurrentRingBuffer) WriteAvailable() int64 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	ret := ego._ringBuffer.WriteAvailable()
	return ret
}

func (ego *ConcurrentRingBuffer) Clear() {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	ego._ringBuffer.Clear()
}

func (ego *ConcurrentRingBuffer) WritePos() int64 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	ret := ego._ringBuffer.WritePos()
	return ret
}

func (ego *ConcurrentRingBuffer) PeekRawBytes(ba []byte, baOff int64, peekLength int64, isStrict bool) (int64, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	return ego._ringBuffer.PeekRawBytes(ba, baOff, peekLength, isStrict)
}

func (ego *ConcurrentRingBuffer) ReadRawBytes(ba []byte, baOff int64, readLength int64, isStrict bool) int64 {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	ret := ego._ringBuffer.ReadRawBytes(ba, baOff, readLength, isStrict)
	return ret
}

func (ego *ConcurrentRingBuffer) WriteRawBytes(ba []byte, srcOff int64, srcLength int64) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	ret := ego._ringBuffer.WriteRawBytes(ba, srcOff, srcLength)
	return ret
}

func (ego *ConcurrentRingBuffer) PeekFloat32() (float32, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekFloat32()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadFloat32() (float32, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadFloat32()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteFloat32(iv float32) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteFloat32(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekFloat64() (float64, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekFloat64()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadFloat64() (float64, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadFloat64()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteFloat64(iv float64) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteFloat64(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekBool() (bool, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekBool()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadBool() (bool, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadBool()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteBool(iv bool) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteBool(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekInt8() (int8, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekInt8()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadInt8() (int8, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadInt8()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteInt8(iv int8) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteInt8(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekUInt8() (uint8, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekUInt8()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadUInt8() (uint8, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadUInt8()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteUInt8(iv uint8) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteUInt8(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekInt16() (int16, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekInt16()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadInt16() (int16, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadInt16()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteInt16(iv int16) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteInt16(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekUInt16() (uint16, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekUInt16()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadUInt16() (uint16, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadUInt16()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteUInt16(iv uint16) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteUInt16(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekInt32() (int32, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekInt32()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadInt32() (int32, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadInt32()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteInt32(iv int32) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteInt32(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) ReadBytes() ([]byte, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	return ego._ringBuffer.ReadBytes()
}

func (ego *ConcurrentRingBuffer) WriteBytes(srcBA []byte) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	return ego._ringBuffer.WriteBytes(srcBA)
}

func (ego *ConcurrentRingBuffer) PeekUInt32() (uint32, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekUInt32()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadUInt32() (uint32, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadUInt32()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteUInt32(iv uint32) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteUInt32(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekInt64() (int64, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekInt64()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadInt64() (int64, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadInt64()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteInt64(iv int64) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteInt64(iv)
	return rc
}

func (ego *ConcurrentRingBuffer) PeekUInt64() (uint64, int32, int64, int64) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc, beg, rlen := ego._ringBuffer.PeekUInt64()
	return v, rc, beg, rlen
}

func (ego *ConcurrentRingBuffer) ReadUInt64() (uint64, int32) {
	ego._lock.RLock()
	defer ego._lock.RUnlock()
	v, rc := ego._ringBuffer.ReadUInt64()
	return v, rc
}

func (ego *ConcurrentRingBuffer) WriteUInt64(iv uint64) int32 {
	ego._lock.Lock()
	defer ego._lock.Unlock()
	rc := ego._ringBuffer.WriteUInt64(iv)
	return rc
}

func NeoConcurrentRingBuffer(capacity int64) *ConcurrentRingBuffer {
	bf := &ConcurrentRingBuffer{
		_ringBuffer: &RingBuffer{
			_capacity: capacity,
			_beginPos: 0,
			_length:   0,
			_data:     make([]byte, capacity),
			_b8Cache:  make([]byte, 8),
		},
	}
	return bf
}

var _ IByteBuffer = &ConcurrentRingBuffer{}
