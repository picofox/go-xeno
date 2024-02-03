package memory

import (
	"container/list"
	"fmt"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/algorithm"
	"xeno/zohar/core/datatype"
)

var sEmptyBuffer []byte = make([]byte, 0)
var sEmptyByteArr [][]byte = make([][]byte, 0)
var sEmptyStringArr []string = make([]string, 0)
var sEmptyInt8Arr []int8 = make([]int8, 0)
var sEmptyUInt8Arr []uint8 = make([]uint8, 0)
var sEmptyInt16Arr []int16 = make([]int16, 0)
var sEmptyUInt16Arr []uint16 = make([]uint16, 0)
var sEmptyInt32Arr []int32 = make([]int32, 0)
var sEmptyUInt32Arr []uint32 = make([]uint32, 0)
var sEmptyInt64Arr []int64 = make([]int64, 0)
var sEmptyUInt64Arr []uint64 = make([]uint64, 0)
var sEmptyFloat64Arr []float64 = make([]float64, 0)
var sEmptyFloat32Arr []float32 = make([]float32, 0)
var sEmptyBoolArr []bool = make([]bool, 0)

func ConstEmptyBytes() []byte {
	return sEmptyBuffer
}

func ConstEmptyBoolArray() []bool {
	return sEmptyBoolArr
}

func ConstEmptyBytesArray() [][]byte {
	return sEmptyByteArr
}

func ConstEmptyStringArr() []string {
	return sEmptyStringArr
}

func ConstEmptyInt8Arr() []int8 {
	return sEmptyInt8Arr
}

func ConstEmptyUInt8Arr() []uint8 {
	return sEmptyUInt8Arr
}

func ConstEmptyInt16Arr() []int16 {
	return sEmptyInt16Arr
}

func ConstEmptyUInt16Arr() []uint16 {
	return sEmptyUInt16Arr
}

func ConstEmptyInt32Arr() []int32 {
	return sEmptyInt32Arr
}

func ConstEmptyUInt32Arr() []uint32 {
	return sEmptyUInt32Arr
}

func ConstEmptyInt64Arr() []int64 {
	return sEmptyInt64Arr
}

func ConstEmptyUInt64Arr() []uint64 {
	return sEmptyUInt64Arr
}

func ConstEmptyFloat32Arr() []float32 {
	return sEmptyFloat32Arr
}

func ConstEmptyFloat64Arr() []float64 {
	return sEmptyFloat64Arr
}

type LinkedListByteBuffer struct {
	_pieceSize int64
	_capacity  int64
	_beginPos  int64
	_length    int64
	_cache     []byte
	_list      *list.List
}

func (ego *LinkedListByteBuffer) SetInt64(pos int64, iv int64) int32 {
	if pos+datatype.INT64_SIZE > ego._beginPos+ego._length {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	node, begPos := ego.findNode(pos)
	if node == nil {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	return ego.SetInt64ByNode(node, begPos, iv)
}

func (ego *LinkedListByteBuffer) SetInt64ByNode(node *list.Element, pos int64, iv int64) int32 {
	remainSpaceInCurBlock := ego._pieceSize - pos
	if remainSpaceInCurBlock >= datatype.INT64_SIZE {
		Int64IntoBytesBE(iv, node.Value.(*[]byte), pos)
		return core.MkSuccess(0)
	} else {
		Int64IntoBytesBE(iv, &ego._cache, 0)
		return ego.SetRawBytesByNode(node, pos, ego._cache, 0, datatype.INT64_SIZE)
	}
}

func (ego *LinkedListByteBuffer) SetInt32(pos int64, iv int32) int32 {
	node, begPos := ego.findNode(pos)
	if node == nil {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	return ego.SetInt32ByNode(node, begPos, iv)
}

func (ego *LinkedListByteBuffer) SetInt32ByNode(node *list.Element, pos int64, iv int32) int32 {
	remainSpaceInCurBlock := ego._pieceSize - pos
	if remainSpaceInCurBlock >= datatype.INT32_SIZE {
		Int32IntoBytesBE(iv, node.Value.(*[]byte), pos)
		return core.MkSuccess(0)
	} else {
		Int32IntoBytesBE(iv, &ego._cache, 0)
		return ego.SetRawBytesByNode(node, pos, ego._cache, 0, datatype.INT32_SIZE)
	}
}

func (ego *LinkedListByteBuffer) SetRawBytesByNode(node *list.Element, pos int64, bs []byte, offset int64, length int64) int32 {
	remainSpaceInCurBlock := ego._pieceSize - pos
	if remainSpaceInCurBlock >= length {
		copy((*(node.Value.(*[]byte)))[pos:pos+length], bs[offset:offset+length])
		return core.MkSuccess(0)
	} else {
		curTurnToWrite := remainSpaceInCurBlock
		off := int64(0)
		for length > 0 {
			copy((*(node.Value.(*[]byte)))[pos:pos+curTurnToWrite], bs[offset+off:offset+off+curTurnToWrite])
			off += curTurnToWrite
			length -= curTurnToWrite
			if length < 0 {
				panic("[SNH] src Len <0")
			} else if length == 0 {
				break
			}
			if node.Next() == nil {
				return core.MkErr(core.EC_NULL_VALUE, 1)
			}
			node = node.Next()
			pos = 0
			curTurnToWrite = min(ego._pieceSize, length)
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) SetRawBytes(pos int64, bs []byte, offset int64, length int64) int32 {
	node, begPos := ego.findNode(pos)
	if node == nil {
		return core.MkErr(core.EC_REACH_LIMIT, 1)
	}
	return ego.SetRawBytesByNode(node, begPos, bs, offset, length)
}

func (ego *LinkedListByteBuffer) SetLength(ll int64) {
	ego._length = ll
}

func (ego *LinkedListByteBuffer) String() string {
	var ss strings.Builder
	ss.WriteString(fmt.Sprintf("cap:%d, beg:%d, len:%d lcnt:%d bal:%d", ego._capacity, ego._beginPos, ego._length, ego._list.Len(), GetDefaultBufferCacheManager().GetCache(ego._pieceSize).Balance()))
	return ss.String()
}

func (ego *LinkedListByteBuffer) findNode(beginPos int64) (*list.Element, int64) {
	if beginPos >= ego._beginPos+ego._length {
		return nil, -1
	}
	begNode := ego._list.Front()
	if begNode == nil {
		return nil, -1
	}
	skipNodes := beginPos / ego._pieceSize
	begPos := beginPos % ego._pieceSize
	for skipNodes > 0 {
		if begNode.Next() == nil {
			return nil, -1
		}
		begNode = begNode.Next()
		skipNodes--
	}
	if begNode == nil {
		return nil, -1
	}
	return begNode, begPos
}

func (ego *LinkedListByteBuffer) preRead() int32 {
	for ego._beginPos >= ego._pieceSize {
		rc := ego._clearFront()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) postRead(adjustSize int64) int32 {
	ego._beginPos += adjustSize
	ego._length -= adjustSize
	if ego._beginPos >= ego._pieceSize {
		rc := ego._clearFront()
		if core.Err(rc) {
			return core.MkErr(core.EC_INCOMPLETE_DATA, 102)
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) _clearFront() int32 {

	if ego._list != nil {
		if ego._list.Front() != nil {
			buf := ego._list.Front().Value.(*[]byte)
			GetDefaultBufferCacheManager().GetCache(ego._pieceSize).Put(buf)
			ego._beginPos -= ego._pieceSize
			ego._capacity -= ego._pieceSize
			ego._list.Remove(ego._list.Front())
			if ego._beginPos < 0 {
				panic("[SNH] begin pos become negative")
			}

			if ego._beginPos > 4095 {
				fmt.Printf("beg %d\n", ego._beginPos)
			}

			return core.MkSuccess(0)
		} else {
			fmt.Printf("list head\n")
		}
	} else {
		fmt.Printf("list null\n")
	}
	return core.MkErr(core.EC_NULL_VALUE, 1)
}

func (ego *LinkedListByteBuffer) bufferForWriting() (*[]byte, int64) {
	var pBuf *[]byte = nil
	wp := ego._beginPos + ego._length
	curNodeCount := wp / ego._pieceSize
	remainBytesInNode := wp % ego._pieceSize
	if ego._list.Back() == nil {
		pBuf = ego.addNode()
		if pBuf == nil {
			return nil, -1
		}
	} else {
		if int64(ego._list.Len()) < curNodeCount+1 {
			pBuf = ego.addNode()
			if pBuf == nil {
				return nil, -1
			}
		} else if int64(ego._list.Len()) == curNodeCount+1 {
			pBuf = ego._list.Back().Value.(*[]byte)
		} else {
			panic("[SNH] redundant buffer found")
		}
	}
	return pBuf, remainBytesInNode
}

func (ego *LinkedListByteBuffer) addNode() *[]byte {
	b := GetDefaultBufferCacheManager().Get(ego._pieceSize)
	if b == nil {
		panic("GetBytesCache().Get() Failed")
	}
	ego._list.PushBack(b)
	ego._capacity += int64(cap(*b))
	return b
}

func (ego *LinkedListByteBuffer) MaxUsedPiece() int64 {
	beglen := ego._beginPos + ego._length
	if beglen == 0 {
		return 0
	}

	n := algorithm.AlignSize(beglen, ego._pieceSize)
	return n / ego._pieceSize
}

func (ego *LinkedListByteBuffer) PieceSize() int64 {
	return ego._pieceSize
}

func (ego *LinkedListByteBuffer) Capacity() int64 {
	return ego._capacity
}

func (ego *LinkedListByteBuffer) ReadAvailable() int64 {
	return ego._length
}

func (ego *LinkedListByteBuffer) WriteAvailable() int64 {
	return ego._capacity - ego._beginPos - ego._length
}

func (ego *LinkedListByteBuffer) Clear() {
	ego._beginPos = 0
	ego._length = 0
	ego._capacity = 0
	for e := ego._list.Front(); e != nil; e = e.Next() {
		b := e.Value.(*[]byte)
		GetDefaultBufferCacheManager().Put(ego._pieceSize, b)
	}
	ego._list.Init()
}

func (ego *LinkedListByteBuffer) WriteInt8(i int8) int32 {
	buf, beg := ego.bufferForWriting()
	curTurnWrite := ego._pieceSize - beg
	if curTurnWrite > 0 {
		(*buf)[beg] = byte(i)
		ego._length += datatype.INT8_SIZE
	} else {
		panic("[SNH] should have at least 1 byte space")
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) WriteUInt8(u uint8) int32 {
	return ego.WriteInt8(int8(u))
}

func (ego *LinkedListByteBuffer) ReadUInt8() (uint8, int32) {
	v, r := ego.ReadInt8()
	return uint8(v), r
}

func (ego *LinkedListByteBuffer) ReadInt8() (int8, int32) {
	if ego._length < datatype.INT8_SIZE || ego._list.Front() == nil {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curReadAvail := ego._pieceSize - ego._beginPos
	if curReadAvail > 0 {
		var retV int8 = int8((*(ego._list.Front().Value.(*[]byte)))[ego._beginPos])
		ego.postRead(datatype.INT8_SIZE)
		return retV, core.MkSuccess(0)
	} else {
		panic("[SNH] should have at least 1 byte space")
	}
}

func (ego *LinkedListByteBuffer) PeekInt8(srcOff int64) (int8, int32) {
	if ego._beginPos+ego._length-srcOff < datatype.INT8_SIZE {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	begNode, segBeginPos := ego.findNode(ego._beginPos + srcOff)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	buf := begNode.Value.(*[]byte)
	leftSpace := ego._pieceSize - segBeginPos
	if leftSpace >= datatype.INT8_SIZE {
		return int8((*buf)[segBeginPos]), core.MkSuccess(0)
	} else {
		panic("[SNH] 1 byte even not has")
	}
}

func (ego *LinkedListByteBuffer) PeekUInt8(srcOff int64) (uint8, int32) {
	v, rc := ego.PeekInt8(srcOff)
	return uint8(v), rc
}

func (ego *LinkedListByteBuffer) WriteRawBytes(bs []byte, srcOff int64, srcLength int64) int32 {
	if srcLength < 0 {
		srcLength = int64(len(bs))
	}
	if srcLength == 0 {
		return core.MkSuccess(0)
	}

	bs = bs[srcOff : srcOff+srcLength]
	off := int64(0)

	if ego._beginPos >= ego._pieceSize {
		panic("[SNH] begin pos too large")
	}

	buf, beg := ego.bufferForWriting()
	curTurnToWrite := ego._pieceSize - beg

	if curTurnToWrite >= srcLength {
		copy((*buf)[beg:beg+curTurnToWrite], bs[off:off+srcLength])
		ego._length += srcLength
	} else {
		for srcLength > 0 {
			copy((*buf)[beg:beg+curTurnToWrite], bs[off:off+curTurnToWrite])
			ego._length += curTurnToWrite
			off += curTurnToWrite
			srcLength -= curTurnToWrite
			if srcLength < 0 {
				panic("[SNH] src Len <0")
			} else if srcLength == 0 {
				break
			}
			buf = ego.addNode()
			if buf == nil {
				return core.MkErr(core.EC_NULL_VALUE, 2)
			}
			beg = 0
			curTurnToWrite = min(ego._pieceSize, srcLength)
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadRawBytes(bs []byte, baOff int64, readLength int64, isStrict bool) (int64, int32) {
	if readLength < 0 {
		readLength = int64(cap(bs)) - baOff
	}
	if ego._length < readLength {
		if isStrict {
			return 0, core.MkErr(core.EC_REACH_LIMIT, 1)
		} else {
			readLength = ego._length
		}
	}
	if ego._length < readLength {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curTurnReadLength := ego._pieceSize - ego._beginPos
	if curTurnReadLength > readLength {
		curTurnReadLength = readLength
	}
	idx := int64(0)
	var origReadLength int64 = readLength
	for readLength > 0 {
		if ego._list.Front() == nil {
			return origReadLength - readLength, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		copy(bs[baOff+idx:baOff+idx+curTurnReadLength], (*(ego._list.Front().Value.(*[]byte)))[ego._beginPos:ego._beginPos+curTurnReadLength])
		idx += curTurnReadLength
		ego._beginPos += curTurnReadLength
		ego._length -= curTurnReadLength
		readLength -= curTurnReadLength
		if ego._beginPos >= ego._pieceSize {
			rc := ego._clearFront()
			if core.Err(rc) {
				return origReadLength - readLength, rc
			}

		}
		if readLength == 0 {
			break
		} else if readLength < 0 {
			panic("Read length < 0")
		}
		if readLength > ego._pieceSize {
			curTurnReadLength = ego._pieceSize
		} else {
			curTurnReadLength = readLength
		}
	}
	return origReadLength, core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekRawBytes(srcOff int64, dstBuf []byte, dstOff int64, readLength int64, isStrict bool) (int64, int32) {
	if readLength < 0 {
		readLength = int64(cap(dstBuf)) - dstOff
	}
	var simuBeginPos int64 = ego._beginPos + srcOff
	var simuLength int64 = ego._length - srcOff
	if simuBeginPos >= ego._beginPos+ego._length {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	var curTurnReadLength int64 = 0
	if simuLength < readLength {
		if isStrict {
			return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
		} else {
			readLength = simuLength
		}
	}

	if readLength == 0 {
		return 0, core.MkSuccess(0)
	}
	begNode, segBeginPos := ego.findNode(simuBeginPos)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	readAvailOfThisNode := ego._pieceSize - segBeginPos
	if readAvailOfThisNode >= readLength {
		curTurnReadLength = readLength
	} else {
		curTurnReadLength = readAvailOfThisNode
	}

	idx := int64(0)
	buf := begNode.Value.(*[]byte)
	for idx < readLength {
		copy(dstBuf[dstOff+idx:dstOff+idx+curTurnReadLength], (*buf)[segBeginPos:segBeginPos+curTurnReadLength])
		idx += curTurnReadLength
		simuBeginPos += curTurnReadLength
		simuLength -= curTurnReadLength
		if idx == readLength {
			break
		} else if idx > readLength {
			panic("[SNH] excessive read!")
		}

		if simuBeginPos >= ego._pieceSize {
			begNode = begNode.Next()
			if begNode == nil {
				return 0, core.MkErr(core.EC_TRY_AGAIN, 3)
			}
			buf = begNode.Value.(*[]byte)
			segBeginPos = 0
		}

		if readLength-idx > ego._pieceSize {
			curTurnReadLength = ego._pieceSize
		} else {
			curTurnReadLength = readLength - idx
		}
	}
	return readLength, core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) WriteInt32(i int32) int32 {
	buf, beg := ego.bufferForWriting()
	curTurnWrite := ego._pieceSize - beg

	if curTurnWrite >= datatype.INT32_SIZE {
		Int32IntoBytesBE(i, buf, beg)
		ego._length += datatype.INT32_SIZE
	} else {
		Int32IntoBytesBE(i, &ego._cache, 0)
		rc := ego.WriteRawBytes(ego._cache, 0, datatype.INT32_SIZE)
		if core.Err(rc) {
			return rc
		}
	}

	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadInt32() (int32, int32) {
	if ego._length < datatype.INT32_SIZE || ego._list.Front() == nil {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curReadAvail := ego._pieceSize - ego._beginPos
	if curReadAvail >= datatype.INT32_SIZE {
		v := BytesToInt32BE(ego._list.Front().Value.(*[]byte), ego._beginPos)
		rc := ego.postRead(datatype.INT32_SIZE)
		if core.Err(rc) {
			return 0, rc
		}
		return v, core.MkSuccess(0)
	} else {
		_, rc := ego.ReadRawBytes(ego._cache, 0, datatype.INT32_SIZE, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToInt32BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) PeekInt32(srcOff int64) (int32, int32) {
	if ego._beginPos+ego._length-srcOff < datatype.INT32_SIZE {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	begNode, segBeginPos := ego.findNode(ego._beginPos + srcOff)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	buf := begNode.Value.(*[]byte)
	leftSpace := ego._pieceSize - segBeginPos
	if leftSpace >= datatype.INT32_SIZE {
		v := BytesToInt32BE(buf, segBeginPos)
		return v, core.MkSuccess(0)
	} else {
		rd, rc := ego.PeekRawBytes(srcOff, ego._cache, 0, datatype.INT32_SIZE, true)
		if core.Err(rc) {
			return 0, rc
		}
		if rd != datatype.INT32_SIZE {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToInt32BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) WriteBytes(srcBA []byte) int32 {
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

func (ego *LinkedListByteBuffer) PeekBytesArray() ([][]byte, int32) {
	return nil, 0
}

func (ego *LinkedListByteBuffer) ReadBytesArray() ([][]byte, int32) {
	rc := ego.isByteArrayDataReadyToRead()
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyByteArr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret [][]byte = make([][]byte, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadBytes()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) ReadStrings() ([]string, int32) {
	rc := ego.isByteArrayDataReadyToRead()
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyStringArr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []string = make([]string, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadString()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) ReadBytes() ([]byte, int32) {
	if ego._length < datatype.INT32_SIZE {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	bLen, rc := ego.PeekInt32(0)
	if core.Err(rc) {
		return nil, rc
	}
	if bLen < 0 {
		ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
		return nil, core.MkSuccess(0)
	} else if bLen == 0 {
		ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
		return sEmptyBuffer, core.MkSuccess(0)
	}
	if ego._length < int64(bLen+datatype.INT32_SIZE) {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	ego.ReadInt32()
	rBA := make([]byte, bLen)
	_, rc = ego.ReadRawBytes(rBA, 0, int64(bLen), true)
	if core.Err(rc) {
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
	}
	return rBA, core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekBytes(srcOff int64) ([]byte, int32) {
	if ego._length < datatype.INT32_SIZE {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	bLen, rc := ego.PeekInt32(0)
	if core.Err(rc) {
		return nil, rc
	}
	if bLen < 0 {
		return nil, core.MkSuccess(0)
	} else if bLen == 0 {
		return sEmptyBuffer, core.MkSuccess(0)
	}
	if ego._length < int64(bLen+datatype.INT32_SIZE) {
		return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
	}

	rBA := make([]byte, bLen)
	rd, rc := ego.PeekRawBytes(srcOff+datatype.INT32_SIZE, rBA, 0, int64(bLen), true)
	if core.Err(rc) || rd != int64(bLen) {
		return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
	}

	return rBA, core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) WriteString(s string) int32 {
	bLen := len(s)
	if bLen == 0 {
		return ego.WriteInt32(0)
	}
	ba := ByteRef(s, 0, int(bLen))
	rc := ego.WriteBytes(ba)
	if core.Err(rc) {
		return rc
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadString() (string, int32) {
	rBA, rc := ego.ReadBytes()
	if core.Err(rc) {
		return "", rc
	}

	if rBA == nil {
		return "", core.MkErr(core.EC_NULL_VALUE, 1)
	} else if len(rBA) == 0 {
		return "", core.MkSuccess(0)
	}

	return StringRef(rBA), core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekString(srcOff int64) (string, int32) {
	rBA, rc := ego.PeekBytes(srcOff)
	if core.Err(rc) {
		return "", rc
	}
	if rBA == nil {
		return "", core.MkErr(core.EC_NULL_VALUE, 1)
	} else if len(rBA) == 0 {
		return "", core.MkSuccess(0)
	}
	return StringRef(rBA), core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) WriteBytesArray(bya [][]byte) int32 {
	if bya == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(bya)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteBytes(bya[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) WriteStrings(strs []string) int32 {
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

func (ego *LinkedListByteBuffer) isIntArrayReadyToRead(intSize int32) int32 {
	arrCount, rc := ego.PeekInt32(0)
	if core.Err(rc) {
		return rc
	}
	if arrCount < -1 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	} else if arrCount == -1 {
		return core.MkSuccess(1)
	} else if arrCount == 0 {
		return core.MkSuccess(2)
	}
	var curLength int64 = ego._length - datatype.INT32_SIZE
	if curLength < int64(intSize*arrCount) {
		return core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) isByteArrayDataReadyToRead() int32 {
	arrCount, rc := ego.PeekInt32(0)
	if core.Err(rc) {
		return rc
	}
	if arrCount < -1 {
		return core.MkErr(core.EC_INCOMPLETE_DATA, 1)
	} else if arrCount == -1 {
		return core.MkSuccess(1)
	} else if arrCount == 0 {
		return core.MkSuccess(2)
	}
	var curLength int64 = ego._length - datatype.INT32_SIZE
	var idx int64 = datatype.INT32_SIZE
	var bLen int32 = 0
	for i := int32(0); i < arrCount; i++ {
		bLen, rc = ego.PeekInt32(idx)
		if core.Err(rc) {
			return rc
		}
		idx += datatype.INT32_SIZE
		curLength -= datatype.INT32_SIZE
		if curLength < int64(bLen) {
			return core.MkErr(core.EC_TRY_AGAIN, 1)
		}
		idx += int64(bLen)
		curLength -= int64(bLen)
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekFloat32(srcOff int64) (float32, int32) {
	if ego._beginPos+ego._length-srcOff < datatype.INT32_SIZE {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	begNode, segBeginPos := ego.findNode(ego._beginPos + srcOff)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	buf := begNode.Value.(*[]byte)
	leftSpace := ego._pieceSize - segBeginPos
	if leftSpace >= datatype.FLOAT32_SIZE {
		v := BytesToFloat32BE(buf, segBeginPos)
		return v, core.MkSuccess(0)
	} else {
		rd, rc := ego.PeekRawBytes(srcOff, ego._cache, 0, datatype.FLOAT32_SIZE, true)
		if core.Err(rc) {
			return 0, rc
		}
		if rd != datatype.FLOAT32_SIZE {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToFloat32BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) ReadFloat32() (float32, int32) {
	if ego._length < datatype.FLOAT32_SIZE || ego._list.Front() == nil {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curReadAvail := ego._pieceSize - ego._beginPos
	if curReadAvail >= datatype.FLOAT32_SIZE {
		v := BytesToFloat32BE(ego._list.Front().Value.(*[]byte), ego._beginPos)

		rc := ego.postRead(datatype.FLOAT32_SIZE)
		if core.Err(rc) {
			return 0, rc
		}
		return v, core.MkSuccess(0)
	} else {
		_, rc := ego.ReadRawBytes(ego._cache, 0, datatype.FLOAT32_SIZE, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToFloat32BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) WriteFloat32(f float32) int32 {
	buf, beg := ego.bufferForWriting()
	curTurnWrite := ego._pieceSize - beg

	if curTurnWrite >= datatype.FLOAT32_SIZE {
		Float32IntoBytesBE(f, buf, beg)
		ego._length += datatype.FLOAT32_SIZE
	} else {
		Float32IntoBytesBE(f, &ego._cache, 0)
		rc := ego.WriteRawBytes(ego._cache, 0, datatype.FLOAT32_SIZE)
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekFloat64(srcOff int64) (float64, int32) {
	if ego._beginPos+ego._length-srcOff < datatype.INT64_SIZE {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	begNode, segBeginPos := ego.findNode(ego._beginPos + srcOff)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	buf := begNode.Value.(*[]byte)
	leftSpace := ego._pieceSize - segBeginPos
	if leftSpace >= datatype.FLOAT64_SIZE {
		v := BytesToFloat64BE(buf, segBeginPos)
		return v, core.MkSuccess(0)
	} else {
		rd, rc := ego.PeekRawBytes(srcOff, ego._cache, 0, datatype.FLOAT64_SIZE, true)
		if core.Err(rc) {
			return 0, rc
		}
		if rd != datatype.FLOAT64_SIZE {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToFloat64BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) ReadFloat64() (float64, int32) {
	if ego._length < datatype.INT32_SIZE || ego._list.Front() == nil {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curReadAvail := ego._pieceSize - ego._beginPos
	if curReadAvail > datatype.FLOAT64_SIZE {
		v := BytesToFloat64BE(ego._list.Front().Value.(*[]byte), ego._beginPos)
		rc := ego.postRead(datatype.FLOAT64_SIZE)
		if core.Err(rc) {
			return 0, rc
		}
		return v, core.MkSuccess(0)
	} else {
		_, rc := ego.ReadRawBytes(ego._cache, 0, datatype.FLOAT64_SIZE, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToFloat64BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) WriteFloat64(f float64) int32 {
	buf, beg := ego.bufferForWriting()
	curTurnWrite := ego._pieceSize - beg
	if curTurnWrite >= datatype.FLOAT64_SIZE {
		Float64IntoBytesBE(f, buf, beg)
		ego._length += datatype.FLOAT64_SIZE
	} else {
		Float64IntoBytesBE(f, &ego._cache, 0)
		rc := ego.WriteRawBytes(ego._cache, 0, datatype.FLOAT64_SIZE)
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekBool(srcOff int64) (bool, int32) {
	iv, rc := ego.PeekInt8(srcOff)
	if core.Err(rc) {
		return false, rc
	}
	if iv != 0 {
		return true, rc
	}
	return false, rc
}

func (ego *LinkedListByteBuffer) ReadBool() (bool, int32) {
	iv, rc := ego.ReadInt8()
	if core.Err(rc) {
		return false, rc
	}
	if iv != 0 {
		return true, rc
	}
	return false, rc
}

func (ego *LinkedListByteBuffer) WriteBool(b bool) int32 {
	if b {
		return ego.WriteInt8(1)
	} else {
		return ego.WriteInt8(0)
	}
}

func (ego *LinkedListByteBuffer) PeekInt16(srcOff int64) (int16, int32) {
	if ego._beginPos+ego._length-srcOff < datatype.INT16_SIZE {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	begNode, segBeginPos := ego.findNode(ego._beginPos + srcOff)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	buf := begNode.Value.(*[]byte)
	leftSpace := ego._pieceSize - segBeginPos
	if leftSpace >= datatype.INT16_SIZE {
		v := BytesToInt16BE(buf, segBeginPos)
		return v, core.MkSuccess(0)
	} else {
		rd, rc := ego.PeekRawBytes(srcOff, ego._cache, 0, datatype.INT16_SIZE, true)
		if core.Err(rc) {
			return 0, rc
		}
		if rd != datatype.INT16_SIZE {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToInt16BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) ReadInt16() (int16, int32) {
	if ego._length < datatype.INT16_SIZE || ego._list.Front() == nil {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curReadAvail := ego._pieceSize - ego._beginPos
	if curReadAvail >= datatype.INT16_SIZE {
		v := BytesToInt16BE(ego._list.Front().Value.(*[]byte), ego._beginPos)
		rc := ego.postRead(datatype.INT16_SIZE)
		if core.Err(rc) {
			return 0, rc
		}
		return v, core.MkSuccess(0)
	} else {
		_, rc := ego.ReadRawBytes(ego._cache, 0, datatype.INT16_SIZE, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToInt16BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) WriteInt16(i int16) int32 {
	buf, beg := ego.bufferForWriting()
	curTurnWrite := ego._pieceSize - beg
	if curTurnWrite >= datatype.INT16_SIZE {
		Int16IntoBytesBE(i, buf, beg)
		ego._length += datatype.INT16_SIZE
	} else {
		Int16IntoBytesBE(i, &ego._cache, 0)
		rc := ego.WriteRawBytes(ego._cache, 0, datatype.INT16_SIZE)
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekUInt16(srcOff int64) (uint16, int32) {
	v, r := ego.PeekInt16(srcOff)
	return uint16(v), r
}

func (ego *LinkedListByteBuffer) ReadUInt16() (uint16, int32) {
	v, r := ego.ReadInt16()
	return uint16(v), r
}

func (ego *LinkedListByteBuffer) WriteUInt16(u uint16) int32 {
	return ego.WriteInt16(int16(u))
}

func (ego *LinkedListByteBuffer) PeekUInt32(srcOff int64) (uint32, int32) {
	v, r := ego.PeekInt32(srcOff)
	return uint32(v), r
}

func (ego *LinkedListByteBuffer) ReadUInt32() (uint32, int32) {
	v, r := ego.ReadInt32()
	return uint32(v), r
}

func (ego *LinkedListByteBuffer) WriteUInt32(u uint32) int32 {
	return ego.WriteInt32(int32(u))
}

func (ego *LinkedListByteBuffer) PeekInt64(srcOff int64) (int64, int32) {
	if ego._length-srcOff < datatype.INT64_SIZE {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	begNode, segBeginPos := ego.findNode(ego._beginPos + srcOff)
	if begNode == nil || segBeginPos < 0 {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	buf := begNode.Value.(*[]byte)
	leftSpace := ego._pieceSize - segBeginPos
	if leftSpace >= datatype.INT64_SIZE {
		v := BytesToInt64BE(buf, segBeginPos)
		return v, core.MkSuccess(0)
	} else {
		rd, rc := ego.PeekRawBytes(srcOff, ego._cache, 0, datatype.INT64_SIZE, true)
		if core.Err(rc) {
			return 0, rc
		}
		if rd != datatype.INT64_SIZE {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToInt64BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) ReadInt64() (int64, int32) {
	if ego._length < datatype.INT32_SIZE || ego._list.Front() == nil {
		return 0, core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	curReadAvail := ego._pieceSize - ego._beginPos
	if curReadAvail >= datatype.INT64_SIZE {
		v := BytesToInt64BE(ego._list.Front().Value.(*[]byte), ego._beginPos)

		rc := ego.postRead(datatype.INT64_SIZE)
		if core.Err(rc) {
			return 0, rc
		}
		return v, core.MkSuccess(0)
	} else {
		_, rc := ego.ReadRawBytes(ego._cache, 0, datatype.INT64_SIZE, true)
		if core.Err(rc) {
			return 0, core.MkErr(core.EC_INCOMPLETE_DATA, 2)
		}
		v := BytesToInt64BE(&ego._cache, 0)
		return v, core.MkSuccess(0)
	}
}

func (ego *LinkedListByteBuffer) WriteInt64(i int64) int32 {
	buf, beg := ego.bufferForWriting()
	curTurnWrite := ego._pieceSize - beg

	if curTurnWrite >= datatype.INT64_SIZE {
		Int64IntoBytesBE(i, buf, beg)
		ego._length += datatype.INT64_SIZE
	} else {
		Int64IntoBytesBE(i, &ego._cache, 0)
		rc := ego.WriteRawBytes(ego._cache, 0, datatype.INT64_SIZE)
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) PeekUInt64(srcOff int64) (uint64, int32) {
	v, r := ego.PeekInt32(srcOff)
	return uint64(v), r
}

func (ego *LinkedListByteBuffer) ReadUInt64() (uint64, int32) {
	v, r := ego.ReadInt64()
	return uint64(v), r
}

func (ego *LinkedListByteBuffer) WriteUInt64(u uint64) int32 {
	return ego.WriteInt64(int64(u))
}

func (ego *LinkedListByteBuffer) ReadPos() int64 {
	return ego._beginPos
}

func (ego *LinkedListByteBuffer) WritePos() int64 {
	return ego._beginPos + ego._length
}

func (ego *LinkedListByteBuffer) WriterSeek(whence int, offset int64) int32 {
	if whence == BUFFER_SEEK_CUR {
		if ego._list == nil {
			return core.MkErr(core.EC_NULL_VALUE, 1)
		}
		if offset < 0 {
			return core.MkErr(core.EC_REACH_LIMIT, 1)
		} else if offset == 0 {
			return core.MkSuccess(0)
		}

		avail := ego._pieceSize - ego._beginPos - ego._length
		if offset < avail {
			ego._length += offset
			if ego._list.Len() == 0 {
				ego.addNode()
			}
		} else {
			curTurnBytes := avail
			for offset > 0 {
				ego.addNode()
				offset -= curTurnBytes
				curTurnBytes = ego._pieceSize
			}
			ego._length += offset
		}
		return core.MkSuccess(0)
	}
	return core.MkErr(core.EC_INVALID_STATE, 1)
}

func (ego *LinkedListByteBuffer) ReaderSeek(whence int, offset int64) bool {
	if whence == BUFFER_SEEK_CUR {
		if ego._list == nil || ego._list.Front() == nil {
			return false
		}
		if offset < 0 {
			return false
		} else if offset == 0 {
			return true
		}
		if ego._length < offset {
			return false
		}

		var inOneBuf bool = false
		curTurnReadLength := ego._pieceSize - ego._beginPos
		//curTurnReadLength = min(curTurnReadLength, ego._length)
		if curTurnReadLength >= offset {
			curTurnReadLength = offset
			inOneBuf = true
		}
		if inOneBuf {
			rc := ego.postRead(offset)
			if core.Err(rc) {
				return false
			}
		} else {
			idx := int64(0)
			for offset > 0 {
				idx += curTurnReadLength
				ego._beginPos += curTurnReadLength
				ego._length -= curTurnReadLength
				offset -= curTurnReadLength
				if ego._beginPos >= ego._pieceSize {
					rc := ego._clearFront()
					if core.Err(rc) {
						return false
					}
				}
				if offset == 0 {
					break
				} else if offset < 0 {
					panic("remain offset < 0")
				} else {
					if offset > ego._pieceSize {
						curTurnReadLength = ego._pieceSize
					} else {
						curTurnReadLength = offset
					}
				}
			}
		}
		return true
	}
	return false
}

func (ego *LinkedListByteBuffer) WriteInt8Array(ia []int8) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	if l > 0 {
		rc := ego.WriteRawBytes(Int8ArrayToByteArrayRef(ia), 0, -1)
		if core.Err(rc) {
			return rc
		}
	}

	return core.MkSuccess(0)
}

//func (ego *LinkedListByteBuffer) WriteInt8Array(ia []int8) int32 {
//	if ia == nil {
//		ego.WriteInt32(-1)
//		return core.MkSuccess(0)
//	}
//	l := len(ia)
//	ego.WriteInt32(int32(l))
//	for i := 0; i < l; i++ {
//		rc := ego.WriteInt8(ia[i])
//		if core.Err(rc) {
//			return rc
//		}
//	}
//	return core.MkSuccess(0)
//}

func (ego *LinkedListByteBuffer) ReadInt8Array() ([]int8, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT8_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyInt8Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []int8 = make([]int8, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadInt8()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteUInt8Array(ia []uint8) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	if l > 0 {
		rc := ego.WriteRawBytes(UInt8ArrayToByteArrayRef(ia), 0, -1)
		if core.Err(rc) {
			return rc
		}
	}

	return core.MkSuccess(0)
}

//func (ego *LinkedListByteBuffer) WriteUInt8Array(ia []uint8) int32 {
//	if ia == nil {
//		ego.WriteInt32(-1)
//		return core.MkSuccess(0)
//	}
//	l := len(ia)
//	ego.WriteInt32(int32(l))
//	for i := 0; i < l; i++ {
//		rc := ego.WriteUInt8(ia[i])
//		if core.Err(rc) {
//			return rc
//		}
//	}
//	return core.MkSuccess(0)
//}

func (ego *LinkedListByteBuffer) ReadUInt8Array() ([]uint8, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT8_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyUInt8Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []uint8 = make([]uint8, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadUInt8()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteInt16Array(ia []int16) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteInt16(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadInt16Array() ([]int16, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT16_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyInt16Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []int16 = make([]int16, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadInt16()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteUInt16Array(ia []uint16) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteUInt16(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadUInt16Array() ([]uint16, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT16_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyUInt16Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []uint16 = make([]uint16, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadUInt16()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteInt32Array(ia []int32) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteInt32(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadInt32Array() ([]int32, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT32_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyInt32Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []int32 = make([]int32, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadInt32()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteUInt32Array(ia []uint32) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteUInt32(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadUInt32Array() ([]uint32, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT32_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyUInt32Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []uint32 = make([]uint32, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadUInt32()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteInt64Array(ia []int64) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteInt64(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadInt64Array() ([]int64, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT64_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyInt64Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []int64 = make([]int64, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadInt64()
				if core.Err(rc) {
					et, em := core.ExErr(rc)
					fmt.Printf("%d, %d\n", et, em)
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteUInt64Array(ia []uint64) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteUInt64(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *LinkedListByteBuffer) ReadUInt64Array() ([]uint64, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.INT64_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyUInt64Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []uint64 = make([]uint64, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadUInt64()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteBoolArray(ia []bool) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteBool(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}
func (ego *LinkedListByteBuffer) ReadBoolArray() ([]bool, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.BYTEBOOL_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyBoolArr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []bool = make([]bool, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadBool()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteFloat32Array(ia []float32) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteFloat32(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}
func (ego *LinkedListByteBuffer) ReadFloat32Array() ([]float32, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.FLOAT32_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyFloat32Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []float32 = make([]float32, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadFloat32()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) WriteFloat64Array(ia []float64) int32 {
	if ia == nil {
		ego.WriteInt32(-1)
		return core.MkSuccess(0)
	}
	l := len(ia)
	ego.WriteInt32(int32(l))
	for i := 0; i < l; i++ {
		rc := ego.WriteFloat64(ia[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}
func (ego *LinkedListByteBuffer) ReadFloat64Array() ([]float64, int32) {
	rc := ego.isIntArrayReadyToRead(datatype.FLOAT64_SIZE)
	et, em := core.ExErr(rc)
	if et == core.EC_TRY_AGAIN {
		return nil, rc
	} else if et == core.EC_OK {
		if em == 1 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return nil, core.MkSuccess(0)
		} else if em == 2 {
			ego.ReaderSeek(BUFFER_SEEK_CUR, datatype.INT32_SIZE)
			return sEmptyFloat64Arr, core.MkSuccess(0)
		} else {
			var cnt int32 = 0
			cnt, rc = ego.ReadInt32()
			var ret []float64 = make([]float64, cnt)
			for i := int32(0); i < cnt; i++ {
				ret[i], rc = ego.ReadFloat64()
				if core.Err(rc) {
					return nil, core.MkErr(core.EC_INCOMPLETE_DATA, 1)
				}
			}
			return ret, core.MkSuccess(0)
		}
	} else {
		return nil, rc
	}
}

func (ego *LinkedListByteBuffer) PieceCount() int64 {
	return int64(ego._list.Len())
}

func NeoLinkedListByteBuffer(pieceSize int64) *LinkedListByteBuffer {
	lb := &LinkedListByteBuffer{
		_pieceSize: pieceSize,
		_list:      list.New(),
		_cache:     make([]byte, 8),
		_length:    0,
		_beginPos:  0,
	}
	return lb
}

var _ IByteBuffer = &LinkedListByteBuffer{}
