package memory

import (
	"unsafe"
)

//// string转ytes
//func Str2sbyte(s string) (b []byte) {
//	*(*string)(unsafe.Pointer(&b)) = s	// 把s的地址付给b
//	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2*unsafe.Sizeof(&b))) = len(s)	// 修改容量为长度
//	return
//}
//
//// []byte转string
//func Sbyte2str(b []byte) string {
//	return *(*string)(unsafe.Pointer(&b))
//}

func ByteRef(str string, off int64, length int) (ba []byte) {
	ss := str[off:length]
	*(*string)(unsafe.Pointer(&ba)) = ss
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&ba)) + 2*unsafe.Sizeof(&ba))) = len(str)
	return ba
}

func StringRef(ba []byte) string {
	return *(*string)(unsafe.Pointer(&ba))
}

//func StringSerialize(str string, list *ByteBufferList, blockSize int64) (*ByteBufferList, int32, int32) {
//	l := len(str)
//	if l > datatype.INT32_MAX {
//		return nil, -1, core.MkErr(core.EC_REACH_LIMIT, 1)
//	}
//	var rc int32 = 0
//	if list == nil {
//		list = NeoByteBufferList()
//	}
//
//	var remainLen int64 = int64(l)
//	for remainLen >= blockSize {
//		ba := make([]byte, blockSize)
//
//		bf := AttachLinearBufferFixed(ba, blockSize)
//		n := AdoptByteBufferNode(bf)
//		list.PushBack(n)
//		remainLen -= blockSize
//	}
//
//	return list, int32(l), rc
//}
