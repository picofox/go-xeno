package datatype

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"
)

func Int8ToBytes(i int8) *[]byte {
	return &([]byte{byte(i)})
}

func Int16ToBytesBE(i int16) *[]byte {
	b0 := byte(i >> 8)
	b1 := byte(i & 0xFF)
	var ret = []byte{b0, b1}
	return &ret
}
func Int16ToBytesLE(i int16) *[]byte {
	b0 := byte(i >> 8)
	b1 := byte(i & 0xFF)
	var ret = []byte{b1, b0}
	return &ret
}

func Int32ToBytesBE(i int32) *[]byte {
	b0 := byte((i >> 24) & 0xFF)
	b1 := byte((i >> 16) & 0xFF)
	b2 := byte((i >> 8) & 0xFF)
	b3 := byte(i & 0xFF)
	var ret = []byte{b0, b1, b2, b3}
	return &ret
}
func Int32ToBytesLE(i int32) *[]byte {
	b0 := byte((i >> 24) & 0xFF)
	b1 := byte((i >> 16) & 0xFF)
	b2 := byte((i >> 8) & 0xFF)
	b3 := byte(i & 0xFF)
	var ret = []byte{b3, b2, b1, b0}
	return &ret
}

func Int64ToBytesBE(i int64) *[]byte {
	b0 := byte((i >> 56) & 0xFF)
	b1 := byte((i >> 48) & 0xFF)
	b2 := byte((i >> 40) & 0xFF)
	b3 := byte((i >> 32) & 0xFF)
	b4 := byte((i >> 24) & 0xFF)
	b5 := byte((i >> 16) & 0xFF)
	b6 := byte((i >> 8) & 0xFF)
	b7 := byte(i & 0xFF)
	var ret = []byte{b0, b1, b2, b3, b4, b5, b6, b7}
	return &ret
}
func Int64ToBytesLE(i int64) *[]byte {
	b0 := byte((i >> 56) & 0xFF)
	b1 := byte((i >> 48) & 0xFF)
	b2 := byte((i >> 40) & 0xFF)
	b3 := byte((i >> 32) & 0xFF)
	b4 := byte((i >> 24) & 0xFF)
	b5 := byte((i >> 16) & 0xFF)
	b6 := byte((i >> 8) & 0xFF)
	b7 := byte(i & 0xFF)
	var ret = []byte{b7, b6, b5, b4, b3, b2, b1, b0}
	return &ret
}

func UInt8ToBytes(i uint8) *[]byte {
	return &[]byte{i}
}

func UInt16ToBytesBE(i uint16) *[]byte {
	b0 := byte(i >> 8)
	b1 := byte(i & 0xFF)
	var ret = []byte{b0, b1}
	return &ret
}
func UInt16ToBytesLE(i uint16) *[]byte {
	b0 := byte(i >> 8)
	b1 := byte(i & 0xFF)
	var ret = []byte{b1, b0}
	return &ret
}

func UInt32ToBytesBE(i uint32) *[]byte {
	b0 := byte((i >> 24) & 0xFF)
	b1 := byte((i >> 16) & 0xFF)
	b2 := byte((i >> 8) & 0xFF)
	b3 := byte(i & 0xFF)
	var ret = []byte{b0, b1, b2, b3}
	return &ret
}
func UInt32ToBytesLE(i uint32) *[]byte {
	b0 := byte((i >> 24) & 0xFF)
	b1 := byte((i >> 16) & 0xFF)
	b2 := byte((i >> 8) & 0xFF)
	b3 := byte(i & 0xFF)
	var ret = []byte{b3, b2, b1, b0}
	return &ret
}

func UInt64ToBytesBE(i uint64) *[]byte {
	b0 := byte((i >> 56) & 0xFF)
	b1 := byte((i >> 48) & 0xFF)
	b2 := byte((i >> 40) & 0xFF)
	b3 := byte((i >> 32) & 0xFF)
	b4 := byte((i >> 24) & 0xFF)
	b5 := byte((i >> 16) & 0xFF)
	b6 := byte((i >> 8) & 0xFF)
	b7 := byte(i & 0xFF)
	var ret = []byte{b0, b1, b2, b3, b4, b5, b6, b7}
	return &ret
}
func UInt64ToBytesLE(i uint64) *[]byte {
	b0 := byte((i >> 56) & 0xFF)
	b1 := byte((i >> 48) & 0xFF)
	b2 := byte((i >> 40) & 0xFF)
	b3 := byte((i >> 32) & 0xFF)
	b4 := byte((i >> 24) & 0xFF)
	b5 := byte((i >> 16) & 0xFF)
	b6 := byte((i >> 8) & 0xFF)
	b7 := byte(i & 0xFF)
	var ret = []byte{b7, b6, b5, b4, b3, b2, b1, b0}
	return &ret
}

func BoolToBytes(i bool) *[]byte {
	if !i {
		return &([]byte{byte(0)})
	}
	return &([]byte{byte(1)})
}

func F32ToBytesBE(i float32) *[]byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, i)
	ret := bytebuf.Bytes()
	return &ret
}
func F32ToBytesLE(i float32) *[]byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, i)
	ret := bytebuf.Bytes()
	return &ret
}

func F64ToBytesBE(i float64) *[]byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.BigEndian, i)
	ret := bytebuf.Bytes()
	return &ret
}

func F64ToBytesLE(i float64) *[]byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, i)
	ret := bytebuf.Bytes()
	return &ret
}

func StrToBytes(i string) *[]byte {
	ret := []byte(i)
	return &ret
}

func IntToBytesBE(i int) *[]byte {

	if unsafe.Sizeof(i) == 4 {
		return Int32ToBytesBE(int32(i))
	} else {
		return Int64ToBytesBE(int64(i))
	}
}

func IntToBytesLE(i int) *[]byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, i)
	ret := bytebuf.Bytes()
	return &ret
}

func UIntToBytesBE(i uint) *[]byte {
	if unsafe.Sizeof(i) == 4 {
		return UInt32ToBytesBE(uint32(i))
	} else {
		return UInt64ToBytesBE(uint64(i))
	}
}

func UIntToBytesLE(i uint) *[]byte {
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, i)
	ret := bytebuf.Bytes()
	return &ret
}

func BytesToPrintable(b []byte, prefix bool, uppercase bool) string {
	var sb strings.Builder
	if prefix {
		sb.WriteString("0x")
	}
	if uppercase {
		for i := 0; i < len(b); i++ {
			sb.WriteString(fmt.Sprintf("%02X", b[i]))
		}
	} else {
		for i := 0; i < len(b); i++ {
			sb.WriteString(fmt.Sprintf("%02x", b[i]))
		}
	}

	return sb.String()
}
