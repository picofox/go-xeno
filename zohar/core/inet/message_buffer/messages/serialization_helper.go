package messages

import (
	"xeno/zohar/core"
	"xeno/zohar/core/algorithm"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

func AllocByteBufferBlock() *memory.ByteBufferNode {
	n := memory.GetByteBuffer4KCache().Get()
	n.Clear()
	return n
}

func CheckByteBufferListNode(bufferList *memory.ByteBufferList) (*memory.ByteBufferNode, int32) {
	bufNode := bufferList.Back()
	if bufNode == nil {
		bufNode = AllocByteBufferBlock()
		if bufNode == nil {
			return nil, core.MkErr(core.EC_NULL_VALUE, 1)
		}
		bufferList.PushBack(bufNode)
		return bufNode, core.MkSuccess(0)

	} else {
		if bufNode.WriteAvailable() <= 0 {
			bufNode = AllocByteBufferBlock()
			if bufNode == nil {
				return nil, core.MkErr(core.EC_NULL_VALUE, 2)
			}
			bufferList.PushBack(bufNode)
			return bufNode, core.MkSuccess(0)
		}
		return bufNode, core.MkSuccess(0)
	}
}

func FreeHeaders(headers []*message_buffer.MessageHeader) {
	if headers != nil {
		for i := 0; i < len(headers); i++ {
			if headers[i] != nil {
				GetHeaderCache().Put(headers[i])
			}
		}
	}
}

func AllocHeaders(logicPacketCount int64, lastPacketLength int64, cmd int16) []*message_buffer.MessageHeader {
	var ret []*message_buffer.MessageHeader = make([]*message_buffer.MessageHeader, logicPacketCount)
	if lastPacketLength > message_buffer.MAX_PACKET_BODY_SIZE {
		return nil
	}
	for i := int64(0); i < logicPacketCount; i++ {
		if i == 0 {
			if logicPacketCount == 1 {
				ret[i] = GetHeaderCache().Get()
				ret[i].Set(false, false, int16(lastPacketLength), cmd)
				return ret
			} else {
				ret[i] = GetHeaderCache().Get()
				ret[i].Set(true, false, int16(message_buffer.MAX_PACKET_BODY_SIZE), cmd)
			}
		} else if i == logicPacketCount-1 {
			ret[i] = GetHeaderCache().Get()
			ret[i].Set(false, true, int16(lastPacketLength), cmd)
		} else {
			ret[i] = GetHeaderCache().Get()
			ret[i].Set(true, true, int16(message_buffer.MAX_PACKET_BODY_SIZE), cmd)
		}
	}

	return ret
}

func ReadHeader(bufList *memory.ByteBufferList) (int32, *memory.ByteBufferNode, int64, int16) {
	var logicPacketRemain int16
	var cmd int16

	curNode := GetAvailableBufferNode(bufList)
	if curNode == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1), nil, -1, -1
	}
	if curNode.ReadAvailable() >= message_buffer.O1L15O1T15_HEADER_SIZE {
		lenAndO1, _ := curNode.ReadInt16()
		cmdAndO2, _ := curNode.ReadInt16()
		logicPacketRemain = lenAndO1 & 0x7FFF
		cmd = cmdAndO2 & 0x7FFF
		return core.MkSuccess(0), curNode, int64(logicPacketRemain), cmd

	} else if curNode.ReadAvailable() == 0 {
		curNode = GetAvailableBufferNode(bufList)
		if curNode == nil {
			return core.MkErr(core.EC_NULL_VALUE, 1), nil, -1, -1
		}
		lenAndO1, _ := curNode.ReadInt16()
		cmdAndO2, _ := curNode.ReadInt16()
		logicPacketRemain = lenAndO1 & 0x7FFF
		cmd = cmdAndO2 & 0x7FFF
		return core.MkSuccess(0), curNode, int64(logicPacketRemain), cmd

	} else {
		var ba []byte = make([]byte, message_buffer.O1L15O1T15_HEADER_SIZE)
		var lenToRead = curNode.ReadAvailable()
		curNode.ReadRawBytes(ba, 0, lenToRead, true)
		curNode = GetAvailableBufferNode(bufList)
		if curNode == nil {
			return core.MkErr(core.EC_NULL_VALUE, 1), nil, -1, -1
		}
		curNode.ReadRawBytes(ba, lenToRead, message_buffer.O1L15O1T15_HEADER_SIZE-lenToRead, true)
		lenAndO1 := memory.BytesToInt16BE(&ba, 0)
		cmdAndO2 := memory.BytesToInt16BE(&ba, 2)
		logicPacketRemain = lenAndO1 & 0x7FFF
		cmd = cmdAndO2 & 0x7FFF
		return core.MkSuccess(0), curNode, int64(logicPacketRemain), cmd
	}
}

func WriteHeader(curNode *memory.ByteBufferNode, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList) (int32, *memory.ByteBufferNode, int) {
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
	}
	if core.Err(rc) {
		return rc, nil, -1
	}
	cnwa := curNode.WriteAvailable()
	if cnwa >= message_buffer.O1L15O1T15_HEADER_SIZE { //write header
		rc = curNode.WriteRawBytes(headers[headerIdx].Data(), 0, 4)
	} else {
		rc = curNode.WriteRawBytes(headers[headerIdx].Data(), 0, cnwa)
		curNode, rc = CheckByteBufferListNode(bufferList)
		if core.Err(rc) {
			return rc, nil, -1
		}
		rc = curNode.WriteRawBytes(headers[headerIdx].Data(), cnwa, message_buffer.O1L15O1T15_HEADER_SIZE-cnwa)
	}
	if core.Err(rc) {
		return rc, nil, -1
	}

	GetHeaderCache().Put(headers[headerIdx])
	headers[headerIdx] = nil
	return core.MkSuccess(0), curNode, headerIdx + 1
}

func SerializeBoolType(b bool, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	if b {
		return SerializeI8Type(int8(1), logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	} else {
		return SerializeI8Type(int8(0), logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	}
}

func SerializeU8Type(v uint8, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	return SerializeI8Type(int8(v), logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
}

func SerializeI8Type(v int8, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
		if curNode == nil {
			return rc, nil, -1, -1, -1, -1
		}
	}
	if logicPacketRemain < 1 {
		rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList) //last logical packet is finished
		if core.Err(rc) {
			return rc, nil, -1, -1, -1, -1
		}
		totalIndex += 4
		logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
	}

	if curNode.WriteAvailable() < 1 { //remain physical block can hold remain bytes of value
		curNode, rc = CheckByteBufferListNode(bufferList)
		if core.Err(rc) {
			return rc, nil, -1, -1, -1, -1
		}
	}
	curNode.WriteInt8(v)
	totalIndex += 1
	logicPacketRemain -= 1
	bodyLenCheck += 1
	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeU16Type(v uint16, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	return SerializeI16Type(int16(v), logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
}

func SerializeI16Type(v int16, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var tmpWriteLen int64 = 0
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
		if curNode == nil {
			return rc, nil, -1, -1, -1, -1
		}
	}
	if logicPacketRemain < 2 { //need split packet logically
		if curNode.WriteAvailable() >= logicPacketRemain { //current block can finish current packet
			last1stPartIdx := logicPacketRemain
			//not at very beginning or really occasionally, not just at begin of a physical block
			curNode.WriteInt16Begin(v, logicPacketRemain) //finish current block
			totalIndex += logicPacketRemain
			logicPacketRemain -= logicPacketRemain
			bodyLenCheck += logicPacketRemain
			rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList) //last logical packet is finished
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += 4
			logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
			//
			if curNode.WriteAvailable() >= 2-last1stPartIdx { //remain physical block can hold remain bytes of value
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 2)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			} else { //at the boundary between two physical blocks
				middlePartIdx := curNode.WriteAvailable()
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += middlePartIdx
				logicPacketRemain -= middlePartIdx
				bodyLenCheck += middlePartIdx
				curNode, rc = CheckByteBufferListNode(bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				middlePartIdx += last1stPartIdx
				tmpWriteLen, rc = curNode.WriteTrivialEnd(middlePartIdx, 2)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		} else { // (physicalBlockWriteAvailable < logicPacketRemain) or current block can not finish current packet
			last1stPartIdx := curNode.WriteAvailable()
			curNode.WriteInt16Begin(v, last1stPartIdx) //finish current block
			totalIndex += last1stPartIdx
			logicPacketRemain -= last1stPartIdx
			bodyLenCheck += last1stPartIdx
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			if logicPacketRemain >= 2-last1stPartIdx {
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 2)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			} else {
				middlePartIdx := logicPacketRemain
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += logicPacketRemain
				logicPacketRemain = 0
				bodyLenCheck += logicPacketRemain
				rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				totalIndex += 4

				logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx+middlePartIdx, 2)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		}

	} else {
		if curNode.WriteAvailable() >= 2 {
			rc = curNode.WriteInt16(v)
			totalIndex += 2
			logicPacketRemain -= 2
			bodyLenCheck += 2
		} else {
			curNode.WriteInt16Begin(v, curNode.WriteAvailable())
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			curNode.WriteTrivialEnd(curNode.WriteAvailable(), 2)
			totalIndex += 2
			logicPacketRemain -= 2
			bodyLenCheck += 2
		}
	}

	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeU32Type(v uint32, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	return SerializeI32Type(int32(v), logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
}

func SerializeI32Type(v int32, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var tmpWriteLen int64 = 0
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
		if curNode == nil {
			return rc, nil, -1, -1, -1, -1
		}
	}

	if logicPacketRemain < 4 { //need split packet logically
		if curNode.WriteAvailable() >= logicPacketRemain { //current block can finish current packet
			last1stPartIdx := logicPacketRemain
			//not at very beginning or really occasionally, not just at begin of a physical block
			curNode.WriteInt32Begin(v, logicPacketRemain) //finish current block
			totalIndex += logicPacketRemain
			logicPacketRemain -= logicPacketRemain
			bodyLenCheck += logicPacketRemain
			rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList) //last logical packet is finished
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += 4
			logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
			//

			if curNode.WriteAvailable() >= 4-last1stPartIdx { //remain physical block can hold remain bytes of value
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			} else { //at the boundary between two physical blocks
				middlePartIdx := curNode.WriteAvailable()
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += middlePartIdx
				logicPacketRemain -= middlePartIdx
				bodyLenCheck += middlePartIdx
				curNode, rc = CheckByteBufferListNode(bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				middlePartIdx += last1stPartIdx
				tmpWriteLen, rc = curNode.WriteTrivialEnd(middlePartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		} else { // (physicalBlockWriteAvailable < logicPacketRemain) or current block can not finish current packet
			last1stPartIdx := curNode.WriteAvailable()
			curNode.WriteInt32Begin(v, last1stPartIdx) //finish current block
			totalIndex += last1stPartIdx
			logicPacketRemain -= last1stPartIdx
			bodyLenCheck += last1stPartIdx
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			if logicPacketRemain >= 4-last1stPartIdx {
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen

			} else {
				middlePartIdx := logicPacketRemain
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += logicPacketRemain
				logicPacketRemain = 0
				bodyLenCheck += logicPacketRemain
				rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				totalIndex += 4

				logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx+middlePartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		}

	} else {
		if curNode.WriteAvailable() >= 4 {
			rc = curNode.WriteInt32(v)
			totalIndex += 4
			logicPacketRemain -= 4
			bodyLenCheck += 4
		} else {
			curNode.WriteInt32Begin(v, curNode.WriteAvailable())
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			curNode.WriteTrivialEnd(curNode.WriteAvailable(), 4)
			totalIndex += 4
			logicPacketRemain -= 4
			bodyLenCheck += 4
		}
	}

	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeU64Type(v uint64, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	return SerializeI64Type(int64(v), logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
}

func SerializeI64Type(v int64, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var tmpFieldLength int64 = 0
	var tmpWriteLen int64 = 0
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
		if curNode == nil {
			return rc, nil, -1, -1, -1, -1
		}
	}
	tmpFieldLength = 8
	if logicPacketRemain < tmpFieldLength { //need split packet logically
		if curNode.WriteAvailable() >= logicPacketRemain { //current block can finish current packet
			last1stPartIdx := logicPacketRemain
			//not at very beginning or really occasionally, not just at begin of a physical block
			curNode.WriteInt64Begin(v, logicPacketRemain) //finish current block
			totalIndex += logicPacketRemain
			logicPacketRemain -= logicPacketRemain
			bodyLenCheck += logicPacketRemain
			rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList) //last logical packet is finished
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += 4
			logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
			//
			if curNode.WriteAvailable() >= tmpFieldLength-last1stPartIdx { //remain physical block can hold remain bytes of value
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			} else { //at the boundary between two physical blocks
				middlePartIdx := curNode.WriteAvailable()
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += middlePartIdx
				logicPacketRemain -= middlePartIdx
				bodyLenCheck += middlePartIdx
				curNode, rc = CheckByteBufferListNode(bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				middlePartIdx += last1stPartIdx
				tmpWriteLen, rc = curNode.WriteTrivialEnd(middlePartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		} else { // (physicalBlockWriteAvailable < logicPacketRemain) or current block can not finish current packet
			last1stPartIdx := curNode.WriteAvailable()
			curNode.WriteInt64Begin(v, last1stPartIdx) //finish current block
			totalIndex += last1stPartIdx
			logicPacketRemain -= last1stPartIdx
			bodyLenCheck += last1stPartIdx
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			if logicPacketRemain >= tmpFieldLength-last1stPartIdx {
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen

			} else {
				middlePartIdx := logicPacketRemain
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += logicPacketRemain
				logicPacketRemain = 0
				bodyLenCheck += logicPacketRemain
				rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				totalIndex += 4

				logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx+middlePartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		}

	} else {
		if curNode.WriteAvailable() >= tmpFieldLength {
			rc = curNode.WriteInt64(v)
			totalIndex += 8
			logicPacketRemain -= 8
			bodyLenCheck += 8
		} else {
			curNode.WriteInt64Begin(v, curNode.WriteAvailable())
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			curNode.WriteTrivialEnd(curNode.WriteAvailable(), 8)
			totalIndex += 8
			logicPacketRemain -= 8
			bodyLenCheck += 8
		}
	}

	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeF32Type(v float32, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var tmpWriteLen int64 = 0
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
		if curNode == nil {
			return rc, nil, -1, -1, -1, -1
		}
	}
	if logicPacketRemain < 4 { //need split packet logically
		if curNode.WriteAvailable() >= logicPacketRemain { //current block can finish current packet
			last1stPartIdx := logicPacketRemain
			//not at very beginning or really occasionally, not just at begin of a physical block
			curNode.WriteFloat32Begin(v, logicPacketRemain) //finish current block
			totalIndex += logicPacketRemain
			logicPacketRemain -= logicPacketRemain
			bodyLenCheck += logicPacketRemain
			rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList) //last logical packet is finished
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += 4
			logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
			//
			if curNode.WriteAvailable() >= 4-last1stPartIdx { //remain physical block can hold remain bytes of value
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			} else { //at the boundary between two physical blocks
				middlePartIdx := curNode.WriteAvailable()
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += middlePartIdx
				logicPacketRemain -= middlePartIdx
				bodyLenCheck += middlePartIdx
				curNode, rc = CheckByteBufferListNode(bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				middlePartIdx += last1stPartIdx
				tmpWriteLen, rc = curNode.WriteTrivialEnd(middlePartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		} else { // (physicalBlockWriteAvailable < logicPacketRemain) or current block can not finish current packet
			last1stPartIdx := curNode.WriteAvailable()
			curNode.WriteFloat32Begin(v, last1stPartIdx) //finish current block
			totalIndex += last1stPartIdx
			logicPacketRemain -= last1stPartIdx
			bodyLenCheck += last1stPartIdx
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			if logicPacketRemain >= 4-last1stPartIdx {
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen

			} else {
				middlePartIdx := logicPacketRemain
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += logicPacketRemain
				logicPacketRemain = 0
				bodyLenCheck += logicPacketRemain
				rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				totalIndex += 4

				logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx+middlePartIdx, 4)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		}

	} else {
		if curNode.WriteAvailable() >= 4 {
			rc = curNode.WriteFloat32(v)
			totalIndex += 4
			logicPacketRemain -= 4
			bodyLenCheck += 4
		} else {
			curNode.WriteFloat32Begin(v, curNode.WriteAvailable())
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			curNode.WriteTrivialEnd(curNode.WriteAvailable(), 4)
			totalIndex += 4
			logicPacketRemain -= 4
			bodyLenCheck += 4
		}
	}

	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeF64Type(v float64, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var tmpFieldLength int64 = 0
	var tmpWriteLen int64 = 0
	var rc int32 = 0
	if curNode == nil {
		curNode, rc = CheckByteBufferListNode(bufferList)
		if curNode == nil {
			return rc, nil, -1, -1, -1, -1
		}
	}
	tmpFieldLength = 8

	if logicPacketRemain < tmpFieldLength { //need split packet logically
		if curNode.WriteAvailable() >= logicPacketRemain { //current block can finish current packet
			last1stPartIdx := logicPacketRemain
			//not at very beginning or really occasionally, not just at begin of a physical block
			curNode.WriteFloat64Begin(v, logicPacketRemain) //finish current block
			totalIndex += logicPacketRemain
			logicPacketRemain -= logicPacketRemain
			bodyLenCheck += logicPacketRemain
			rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList) //last logical packet is finished
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += 4
			logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
			//
			if curNode.WriteAvailable() >= tmpFieldLength-last1stPartIdx { //remain physical block can hold remain bytes of value
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			} else { //at the boundary between two physical blocks
				middlePartIdx := curNode.WriteAvailable()
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += middlePartIdx
				logicPacketRemain -= middlePartIdx
				bodyLenCheck += middlePartIdx
				curNode, rc = CheckByteBufferListNode(bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				middlePartIdx += last1stPartIdx
				tmpWriteLen, rc = curNode.WriteTrivialEnd(middlePartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		} else { // (physicalBlockWriteAvailable < logicPacketRemain) or current block can not finish current packet
			last1stPartIdx := curNode.WriteAvailable()
			curNode.WriteFloat64Begin(v, last1stPartIdx) //finish current block
			totalIndex += last1stPartIdx
			logicPacketRemain -= last1stPartIdx
			bodyLenCheck += last1stPartIdx
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			if logicPacketRemain >= tmpFieldLength-last1stPartIdx {
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen

			} else {
				middlePartIdx := logicPacketRemain
				rc = curNode.WriteTrivialMiddle(last1stPartIdx, middlePartIdx)
				totalIndex += logicPacketRemain
				logicPacketRemain = 0
				bodyLenCheck += logicPacketRemain
				rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList)
				if core.Err(rc) {
					return rc, nil, -1, -1, -1, -1
				}
				totalIndex += 4

				logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
				tmpWriteLen, rc = curNode.WriteTrivialEnd(last1stPartIdx+middlePartIdx, 8)
				totalIndex += tmpWriteLen
				logicPacketRemain -= tmpWriteLen
				bodyLenCheck += tmpWriteLen
			}
		}

	} else {
		if curNode.WriteAvailable() >= tmpFieldLength {
			rc = curNode.WriteFloat64(v)
			totalIndex += 8
			logicPacketRemain -= 8
			bodyLenCheck += 8
		} else {
			curNode.WriteFloat64Begin(v, curNode.WriteAvailable())
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			curNode.WriteTrivialEnd(curNode.WriteAvailable(), 8)
			totalIndex += 8
			logicPacketRemain -= 8
			bodyLenCheck += 8
		}
	}

	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeBytesType(bs []byte, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var bsLenCheck int = -1
	if bs != nil {
		bsLenCheck = len(bs)
	}

	if bsLenCheck > datatype.INT32_MAX {
		return core.MkErr(core.EC_REACH_LIMIT, 0), nil, -1, -1, -1, -1
	}
	var bsLen int32 = int32(bsLenCheck)
	var rc int32 = 0
	rc, curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck = SerializeI32Type(bsLen, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return rc, nil, -1, -1, -1, -1
	}
	if bsLen <= 0 {
		return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
	}
	var fieldRemainLen int64 = int64(bsLen)
	var curIdx int64 = 0
	var sel int = 0
	var currentSerializeLen int64
	var debugIdx = 0
	for fieldRemainLen > 0 {
		debugIdx++
		rc, sel, currentSerializeLen = algorithm.MinValue[int64](logicPacketRemain, curNode.WriteAvailable(), fieldRemainLen)
		if sel == 2 {
			rc = curNode.WriteRawBytes(bs, curIdx, currentSerializeLen)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += currentSerializeLen
			logicPacketRemain -= currentSerializeLen
			bodyLenCheck += currentSerializeLen
			return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck

		} else if sel == 1 {
			curNode, rc = CheckByteBufferListNode(bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			rc = curNode.WriteRawBytes(bs, curIdx, currentSerializeLen)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += currentSerializeLen
			logicPacketRemain -= currentSerializeLen
			bodyLenCheck += currentSerializeLen
			curIdx += currentSerializeLen
			fieldRemainLen -= currentSerializeLen

		} else if sel == 0 {
			rc = curNode.WriteRawBytes(bs, curIdx, currentSerializeLen)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += currentSerializeLen
			logicPacketRemain -= currentSerializeLen
			bodyLenCheck += currentSerializeLen
			curIdx += currentSerializeLen
			fieldRemainLen -= currentSerializeLen

			rc, curNode, headerIdx = WriteHeader(curNode, headers, headerIdx, bufferList)
			if core.Err(rc) {
				return rc, nil, -1, -1, -1, -1
			}
			totalIndex += 4
			logicPacketRemain = message_buffer.MAX_PACKET_BODY_SIZE
		} else {
			panic("min value has problem")
		}
	}
	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeStringType(str string, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	if len(str) == 0 {
		return SerializeBytesType(datatype.EmptyByteSlice, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	}
	ba := memory.ByteRef(str, 0, int(len(str)))
	return SerializeBytesType(ba, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
}

func SerializeBytesSliceType(ba [][]byte, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var l int32 = -1
	if ba != nil {
		l = int32(len(ba))
	}

	var rc int32
	rc, curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck = SerializeI32Type(l, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return rc, nil, -1, -1, -1, -1
	}
	if l <= 0 {
		return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
	}
	for i := int32(0); i < l; i++ {
		rc, curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck = SerializeBytesType(ba[i], logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	}
	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}

func SerializeStringsType(str []string, logicPacketRemain int64, totalIndex int64, bodyLenCheck int64, headers []*message_buffer.MessageHeader, headerIdx int, bufferList *memory.ByteBufferList, curNode *memory.ByteBufferNode) (int32, *memory.ByteBufferNode, int, int64, int64, int64) {
	var l int32 = -1
	if str != nil {
		l = int32(len(str))
	}

	var rc int32
	rc, curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck = SerializeI32Type(l, logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	if core.Err(rc) {
		return rc, nil, -1, -1, -1, -1
	}
	if l == 0 {
		return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
	} else if l < 0 {
		return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
	}
	for i := int32(0); i < l; i++ {
		rc, curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck = SerializeStringType(str[i], logicPacketRemain, totalIndex, bodyLenCheck, headers, headerIdx, bufferList, curNode)
	}
	return core.MkSuccess(0), curNode, headerIdx, totalIndex, logicPacketRemain, bodyLenCheck
}
