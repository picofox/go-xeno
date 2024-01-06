package messages

import (
	"xeno/zohar/core"
	"xeno/zohar/core/algorithm"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

func DebugCheckHeader(hdr []byte) {
	if hdr[0] == 255 && hdr[1] == 252 && hdr[2] == 255 && hdr[3] == 255 {
		//fmt.Printf("%d %d %d %d\n", (hdr)[0], (hdr)[1], hdr[2], (hdr)[3])
	} else if hdr[0] == 2 && hdr[1] == 7 && hdr[2] == 255 && hdr[3] == 255 {
		//fmt.Printf("%d %d %d %d\n", (hdr)[0], (hdr)[1], hdr[2], (hdr)[3])
	} else if hdr[0] == 255 && hdr[1] == 252 && hdr[2] == 127 && hdr[3] == 255 {
		//fmt.Printf("%d %d %d %d\n", (hdr)[0], (hdr)[1], hdr[2], (hdr)[3])
	} else {
		//fmt.Printf("%d %d %d %d\n", (hdr)[0], (hdr)[1], hdr[2], (hdr)[3])
	}
}

func GetAvailableBufferNode(bufList *memory.ByteBufferList) *memory.ByteBufferNode {
	curNode := bufList.Front()
	if curNode == nil {
		return nil
	} else if curNode.ReadAvailable() <= 0 {
		bufList.PopFront()
		memory.GetByteBuffer4KCache().Put(curNode)
		if bufList.Front() == nil {
			return nil
		}
		return bufList.Front()
	} else {
		return curNode
	}
}

func extractValuesFromHeader(hdrBS []byte) (int64, int16, int8) {
	u0 := memory.BytesToInt16BE(&hdrBS, 0)
	u1 := memory.BytesToInt16BE(&hdrBS, 2)
	o1 := int8(u0 >> 15 & 0x1)
	o2 := int8(u1 >> 15 & 0x1)
	l := int64(u0 & 0x7FFF)
	cmd := u1 & 0x7FFF
	var st = (o1 << 1) | o2
	return l, cmd, st
}

func SkipHeader(bufList *memory.ByteBufferList) int32 {
	byteBuf := bufList.Front()
	if byteBuf == nil {
		return core.MkErr(core.EC_TRY_AGAIN, 0)
	}

	if byteBuf.ReadAvailable() <= 0 {
		memory.GetByteBuffer4KCache().Put(byteBuf)
		bufList.PopFront()
		byteBuf = bufList.Front()

		if byteBuf == nil {
			return core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		byteBuf.ReadInt32()
		return core.MkSuccess(0)

	} else if byteBuf.ReadAvailable() <= message_buffer.O1L15O1T15_HEADER_SIZE {
		part2Len := message_buffer.O1L15O1T15_HEADER_SIZE - byteBuf.ReadAvailable()
		memory.GetByteBuffer4KCache().Put(byteBuf)
		bufList.PopFront()
		byteBuf = bufList.Front()
		if byteBuf == nil {
			return core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		byteBuf.ReaderSeek(memory.BUFFER_SEEK_SET, part2Len)
		return core.MkSuccess(0)

	} else {
		byteBuf.ReadInt32()
		return core.MkSuccess(0)
	}
}

//func PeekHeader(hdrCache []byte, byteBuf *memory.ByteBufferNode, begin int64) (*memory.ByteBufferNode, int64, int64, int16, int8, int32) {
//
//}

func PeekHeaderContent(hdrCache []byte, byteBuf *memory.ByteBufferNode, physicalIndex int64) (*memory.ByteBufferNode, int64, int64, int16, int8, int32) {
	if physicalIndex < byteBuf.ReadPos() {
		panic("pos error")
	}
	var delta int64 = physicalIndex - byteBuf.ReadPos()
	var avail int64 = byteBuf.ReadAvailable() - delta

	if avail <= 0 {
		byteBuf = byteBuf.Next()
		if byteBuf == nil {
			return nil, -1, 0, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		srcBA := byteBuf.InternalData()
		l, c, t := extractValuesFromHeader((*srcBA)[0:message_buffer.O1L15O1T15_HEADER_SIZE])
		DebugCheckHeader(*srcBA)

		return byteBuf, message_buffer.O1L15O1T15_HEADER_SIZE, l, c, t, core.MkSuccess(0)

	} else if avail <= message_buffer.O1L15O1T15_HEADER_SIZE {
		srcBA := byteBuf.InternalData()
		copy(hdrCache, (*srcBA)[physicalIndex:physicalIndex+avail])
		byteBuf = byteBuf.Next()
		if byteBuf == nil {
			return nil, -1, 0, -1, -1, core.MkErr(core.EC_TRY_AGAIN, 0)
		}
		part2Len := message_buffer.O1L15O1T15_HEADER_SIZE - avail
		srcBA = byteBuf.InternalData()
		copy(hdrCache[avail:], (*srcBA)[0:part2Len])

		DebugCheckHeader(hdrCache)

		l, c, t := extractValuesFromHeader(hdrCache)
		return byteBuf, part2Len, l, c, t, core.MkSuccess(0)

	} else {
		srcBA := byteBuf.InternalData()
		DebugCheckHeader(*srcBA)
		l, c, t := extractValuesFromHeader((*srcBA)[physicalIndex : physicalIndex+message_buffer.O1L15O1T15_HEADER_SIZE])
		return byteBuf, physicalIndex + message_buffer.O1L15O1T15_HEADER_SIZE, l, c, t, core.MkSuccess(0)
	}

}

func DeserializeStringsType(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) ([]string, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return nil, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}
	var baCount int32
	var rc int32
	baCount, logicPacketRemain, bodyLength, rc = DeserializeI32Type(bufList, logicPacketRemain, bodyLength)
	if core.Err(rc) {
		return nil, logicPacketRemain, bodyLength, rc
	}

	if baCount < 0 {
		return nil, logicPacketRemain, bodyLength, core.MkSuccess(0)
	} else if baCount == 0 {
		return make([]string, 0), logicPacketRemain, bodyLength, core.MkSuccess(0)
	}
	var retBAS []string = make([]string, baCount)
	for i := 0; i < int(baCount); i++ {
		retBAS[i], logicPacketRemain, bodyLength, rc = DeserializeStringType(bufList, logicPacketRemain, bodyLength)
		if core.Err(rc) {
			return retBAS, logicPacketRemain, bodyLength, rc
		}
	}
	return retBAS, logicPacketRemain, bodyLength, core.MkSuccess(0)
}

func DeserializeBytesSliceType(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) ([][]byte, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return nil, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}
	var baCount int32
	var rc int32
	baCount, logicPacketRemain, bodyLength, rc = DeserializeI32Type(bufList, logicPacketRemain, bodyLength)
	if core.Err(rc) {
		return nil, logicPacketRemain, bodyLength, rc
	}

	if baCount < 0 {
		return nil, logicPacketRemain, bodyLength, core.MkSuccess(0)
	} else if baCount == 0 {
		return make([][]byte, 0), logicPacketRemain, bodyLength, core.MkSuccess(0)
	}
	var retBAS [][]byte = make([][]byte, baCount)
	for i := 0; i < int(baCount); i++ {
		retBAS[i], logicPacketRemain, bodyLength, rc = DeserializeBytesType(bufList, logicPacketRemain, bodyLength)
		if core.Err(rc) {
			return retBAS, logicPacketRemain, bodyLength, rc
		}
	}
	return retBAS, logicPacketRemain, bodyLength, core.MkSuccess(0)
}

func DeserializeStringType(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (string, int64, int64, int32) {
	ba, ll, bl, rc := DeserializeBytesType(bufList, logicPacketRemain, bodyLength)
	if core.Err(rc) {
		return "", ll, bl, rc
	}
	if ba == nil || len(ba) == 0 {
		return "", ll, bl, rc
	}
	return memory.StringRef(ba), ll, bl, rc
}

func DeserializeBytesType(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) ([]byte, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return nil, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}
	var baLen int32
	var rc int32
	baLen, logicPacketRemain, bodyLength, rc = DeserializeI32Type(bufList, logicPacketRemain, bodyLength)
	if core.Err(rc) {
		return nil, logicPacketRemain, bodyLength, rc
	}
	var remainLength int64 = int64(baLen)
	var retValue []byte = make([]byte, baLen)
	var tmpBAIndex int64 = 0
	for remainLength > 0 {
		rc, idx, minV := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_TYPE_MISMATCH, 1)
		}
		if idx == 2 {
			bufNode.ReadRawBytes(retValue, tmpBAIndex, minV, true)
			logicPacketRemain -= minV
			break
		} else if idx == 1 {
			bufNode.ReadRawBytes(retValue, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return nil, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}
		} else if idx == 0 {
			bufNode.ReadRawBytes(retValue, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV

			if tmpBAIndex >= int64(baLen) {
				break
			}
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}

		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}
	return retValue, logicPacketRemain, bodyLength + int64(baLen), core.MkSuccess(0)
}

func DeserializeBoolType(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (bool, int64, int64, int32) {
	v, ll, bl, rc := DeserializeI8Type(bufList, logicPacketRemain, bodyLength)
	if core.Err(rc) {
		return false, ll, bl, rc
	}
	if v == 0 {
		return false, ll, bl, rc
	} else {
		return true, ll, bl, rc
	}
}

func DeserializeU8Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (uint8, int64, int64, int32) {
	v, ll, bl, rc := DeserializeI8Type(bufList, logicPacketRemain, bodyLength)
	return uint8(v), ll, bl, rc
}

func DeserializeI8Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (int8, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}
	var retValue int8 = 0
	var remainLength int64 = 1
	for remainLength > 0 {
		rc, idx, _ := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_TYPE_MISMATCH, 1)
		}
		if idx == 2 {
			retValue, rc = bufNode.ReadInt8()
			logicPacketRemain--
			return retValue, logicPacketRemain, bodyLength, rc

		} else if idx == 1 {
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}
		} else if idx == 0 {
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}
		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}

	return retValue, logicPacketRemain, bodyLength + datatype.INT8_SIZE, core.MkErr(core.EC_REACH_LIMIT, 0)
}

func DeserializeU16Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (uint16, int64, int64, int32) {
	v, ll, bl, rc := DeserializeI16Type(bufList, logicPacketRemain, bodyLength)
	return uint16(v), ll, bl, rc
}

func DeserializeI16Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (int16, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	var retValue int16 = 0
	var remainLength int64 = 2
	var tmpBA []byte = make([]byte, remainLength)
	var tmpBAIndex int64 = 0
	for remainLength > 0 {
		rc, idx, minV := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_MIN_VALUE_FIND_ERROR, 1)
		}
		if idx == 2 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
		} else if idx == 1 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}
		} else if idx == 0 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}
		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}

	retValue = memory.BytesToInt16BE(&tmpBA, 0)
	return retValue, logicPacketRemain, bodyLength + datatype.INT16_SIZE, core.MkSuccess(0)
}

func DeserializeU32Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (uint32, int64, int64, int32) {
	v, ll, bl, rc := DeserializeI32Type(bufList, logicPacketRemain, bodyLength)
	return uint32(v), ll, bl, rc
}

func DeserializeI32Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (int32, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	var retValue int32 = 0
	var remainLength int64 = 4
	var tmpBA []byte = make([]byte, remainLength)
	var tmpBAIndex int64 = 0
	for remainLength > 0 {
		rc, idx, minV := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_MIN_VALUE_FIND_ERROR, 1)
		}
		if idx == 2 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
		} else if idx == 1 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}
		} else if idx == 0 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}
		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}

	retValue = memory.BytesToInt32BE(&tmpBA, 0)
	return retValue, logicPacketRemain, bodyLength + datatype.INT32_SIZE, core.MkSuccess(0)
}

func DeserializeU64Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (uint64, int64, int64, int32) {
	v, ll, bl, rc := DeserializeI64Type(bufList, logicPacketRemain, bodyLength)
	return uint64(v), ll, bl, rc
}

func DeserializeI64Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (int64, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	var retValue int64 = 0
	var remainLength int64 = 8
	var tmpBA []byte = make([]byte, remainLength)
	var tmpBAIndex int64 = 0
	for remainLength > 0 {
		rc, idx, minV := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_MIN_VALUE_FIND_ERROR, 1)
		}
		if idx == 2 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
		} else if idx == 1 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}
		} else if idx == 0 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}
		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}

	retValue = memory.BytesToInt64BE(&tmpBA, 0)
	return retValue, logicPacketRemain, bodyLength + datatype.INT64_SIZE, core.MkSuccess(0)
}

func DeserializeF32Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (float32, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	var retValue float32 = 0
	var remainLength int64 = 4
	var tmpBA []byte = make([]byte, remainLength)
	var tmpBAIndex int64 = 0
	for remainLength > 0 {
		rc, idx, minV := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_MIN_VALUE_FIND_ERROR, 1)
		}
		if idx == 2 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
		} else if idx == 1 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}
		} else if idx == 0 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}
		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}

	retValue = memory.BytesToFloat32BE(&tmpBA, 0)
	return retValue, logicPacketRemain, bodyLength + datatype.FLOAT32_SIZE, core.MkSuccess(0)
}

func DeserializeF64Type(bufList *memory.ByteBufferList, logicPacketRemain int64, bodyLength int64) (float64, int64, int64, int32) {
	bufNode := GetAvailableBufferNode(bufList)
	if bufNode == nil {
		return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
	}

	var retValue float64 = 0
	var remainLength int64 = 8
	var tmpBA []byte = make([]byte, remainLength)
	var tmpBAIndex int64 = 0
	for remainLength > 0 {
		rc, idx, minV := algorithm.MinValue[int64](logicPacketRemain, bufNode.ReadAvailable(), remainLength)
		if core.Err(rc) {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_MIN_VALUE_FIND_ERROR, 1)
		}
		if idx == 2 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
		} else if idx == 1 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			bufNode = GetAvailableBufferNode(bufList)
			if bufNode == nil {
				return 0, logicPacketRemain, bodyLength, core.MkErr(core.EC_NULL_VALUE, 1)
			}

		} else if idx == 0 {
			bufNode.ReadRawBytes(tmpBA, tmpBAIndex, minV, true)
			tmpBAIndex += minV
			remainLength -= minV
			logicPacketRemain -= minV
			rc, bufNode, logicPacketRemain, _ = ReadHeader(bufList)
			if core.Err(rc) {
				return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_READ_HEADER_ERROR, 2)
			}
		} else {
			return retValue, logicPacketRemain, bodyLength, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 2)
		}
	}

	retValue = memory.BytesToFloat64BE(&tmpBA, 0)
	return retValue, logicPacketRemain, bodyLength + datatype.FLOAT64_SIZE, core.MkSuccess(0)
}
