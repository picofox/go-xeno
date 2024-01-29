package memory

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

var longText string = strs.CreateSampleText("@abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-", 16*1024, datatype.SIZE_1M, "#")

var lb *memory.LinkedListByteBuffer = memory.NeoLinkedListByteBuffer(datatype.SIZE_4K)
var lock sync.Mutex
var TestCount int = 10000
var SingleTestCount int = 100000000
var Int32TestCount int = 100000
var retRawBs []byte = make([]byte, datatype.SIZE_1M+1)
var readCount int = 0

var sw *chrono.StopWatch = chrono.NeoStopWatch()
var conCount atomic.Int64
var bsCount atomic.Int64

func readLongText() int64 {
	lock.Lock()
	defer lock.Unlock()
	bs, rc := lb.ReadBytes()
	if core.Err(rc) {
		return -1
	}
	return int64(len(bs))
}

// ------------------------Test bool Arr types-------------------------------------------------------------------------
func proBoolArrType() {
	lock.Lock()
	defer lock.Unlock()
	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteBoolArray(make([]bool, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteBoolArray(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]bool, cnt)
		ia[0] = false
		ia[cnt-1] = true
		lb.WriteBoolArray(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conBoolArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadBoolArray()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadBoolArray error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != true {
			fmt.Printf("ReadBoolArray validate error\n")
			return false
		}
	} else {
		if ia[0] != false || ia[len(ia)-1] != true {
			fmt.Printf("ReadBoolArray validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerBoolArr() {
	for i := 0; i < Int32TestCount; i++ {
		proBoolArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_BoolArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerBoolArr()

	sw.Begin("sglB")
	for {
		rc := conBoolArrType()
		if !rc {
			panic("conFloat64ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test f64 Arr types-------------------------------------------------------------------------
func proFloat64ArrType() {
	lock.Lock()
	defer lock.Unlock()
	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteFloat64Array(make([]float64, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteFloat64Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]float64, cnt)
		ia[0] = 3.1415926
		ia[cnt-1] = 2.71828
		lb.WriteFloat64Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conFloat64ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadFloat64Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadFloat64Array error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 2.71828 {
			fmt.Printf("ReadFloat64Array validate error\n")
			return false
		}
	} else {
		if ia[0] != 3.1415926 || ia[len(ia)-1] != 2.71828 {
			fmt.Printf("ReadFloat64Array validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerFloat64Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proFloat64ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_Float64ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerFloat64Arr()

	sw.Begin("sglB")
	for {
		rc := conFloat64ArrType()
		if !rc {
			panic("conFloat64ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test f32 Arr types-------------------------------------------------------------------------
func proFloat32ArrType() {
	lock.Lock()
	defer lock.Unlock()
	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteFloat32Array(make([]float32, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteFloat32Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]float32, cnt)
		ia[0] = 3.14
		ia[cnt-1] = 2.71
		lb.WriteFloat32Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conFloat32ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadFloat32Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadFloat32Array error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 2.71 {
			fmt.Printf("ReadFloat32Array validate error\n")
			return false
		}
	} else {
		if ia[0] != 3.14 || ia[len(ia)-1] != 2.71 {
			fmt.Printf("ReadInt32Array validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerFloat32Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proFloat32ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_Float32ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerFloat32Arr()

	sw.Begin("sglB")
	for {
		rc := conFloat32ArrType()
		if !rc {
			panic("conFloat32ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test u8 Arr types-------------------------------------------------------------------------
func proUInt8ArrType() {
	lock.Lock()
	defer lock.Unlock()

	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteUInt8Array(make([]uint8, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteUInt8Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]uint8, cnt)
		ia[0] = 1
		ia[cnt-1] = 255
		lb.WriteUInt8Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conUInt8ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadUInt8Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("conUInt8ArrType error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 255 {
			fmt.Printf("conUInt8ArrType validate error\n")
			return false
		}
	} else {
		if ia[0] != 1 || ia[len(ia)-1] != 255 {
			fmt.Printf("conUInt8ArrType validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerUInt8Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proUInt8ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_UInt8ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerUInt8Arr()

	sw.Begin("sglB")
	for {
		rc := conUInt8ArrType()
		if !rc {
			panic("conUInt8ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}
}

// ------------------------Test i8 Arr types-------------------------------------------------------------------------
func proInt8ArrType() {
	lock.Lock()
	defer lock.Unlock()

	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteInt8Array(make([]int8, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteInt8Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]int8, cnt)
		ia[0] = -128
		ia[cnt-1] = 127
		lb.WriteInt8Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conInt8ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadInt8Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("conInt8ArrType error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 127 {
			fmt.Printf("conInt8ArrType validate error\n")
			return false
		}
	} else {
		if ia[0] != -128 || ia[len(ia)-1] != 127 {
			fmt.Printf("conInt8ArrType validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerInt8Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proInt8ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_Int8ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerInt8Arr()

	sw.Begin("sglB")
	for {
		rc := conInt8ArrType()
		if !rc {
			panic("conInt8ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}
}

// ------------------------Test u16 Arr types-------------------------------------------------------------------------
func proUInt16ArrType() {
	lock.Lock()
	defer lock.Unlock()

	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteUInt16Array(make([]uint16, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteUInt16Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]uint16, cnt)
		ia[0] = 12345
		ia[cnt-1] = 65535
		lb.WriteUInt16Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conUInt16ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadUInt16Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("conUInt16ArrType error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 65535 {
			fmt.Printf("conUInt16ArrType validate error\n")
			return false
		}
	} else {
		if ia[0] != 12345 || ia[len(ia)-1] != 65535 {
			fmt.Printf("conUInt16ArrType validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerUInt16Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proUInt16ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_UInt16ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerUInt16Arr()

	sw.Begin("sglB")
	for {
		rc := conUInt16ArrType()
		if !rc {
			panic("conUInt16ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test i16 Arr types-------------------------------------------------------------------------
func proInt16ArrType() {
	lock.Lock()
	defer lock.Unlock()

	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteInt16Array(make([]int16, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteUInt16Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]int16, cnt)
		ia[0] = -32768
		ia[cnt-1] = 32767
		lb.WriteInt16Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conInt16ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadInt16Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("conUInt16ArrType error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 32767 {
			fmt.Printf("conUInt16ArrType validate error\n")
			return false
		}
	} else {
		if ia[0] != -32768 || ia[len(ia)-1] != 32767 {
			fmt.Printf("conUInt16ArrType validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerInt16Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proInt16ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_Int16ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerInt16Arr()

	sw.Begin("sglB")
	for {
		rc := conInt16ArrType()
		if !rc {
			panic("conUInt16ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test u64 Arr types-------------------------------------------------------------------------
func proUInt64ArrType() {
	lock.Lock()
	defer lock.Unlock()

	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteUInt64Array(make([]uint64, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteUInt64Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]uint64, cnt)
		ia[0] = 0x1234567890123456
		ia[cnt-1] = 0x3142750893480187
		lb.WriteUInt64Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conUInt64ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadUInt64Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadUInt64Array error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 0x3142750893480187 {
			fmt.Printf("ReadUInt64Array validate error\n")
			return false
		}
	} else {
		if ia[0] != 0x1234567890123456 || ia[len(ia)-1] != 0x3142750893480187 {
			fmt.Printf("ReadUInt64Array validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerUInt64Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proUInt64ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_UInt64ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerUInt64Arr()

	sw.Begin("sglB")
	for {
		rc := conUInt64ArrType()
		if !rc {
			panic("conInt64ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test i64 Arr types-------------------------------------------------------------------------
func proInt64ArrType() {
	lock.Lock()
	defer lock.Unlock()

	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteInt64Array(make([]int64, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteInt64Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]int64, cnt)
		ia[0] = 0x1234567890123456
		ia[cnt-1] = 0x3142750893480187
		lb.WriteInt64Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conInt64ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadInt64Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadInt64Array error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 0x3142750893480187 {
			fmt.Printf("ReadInt64Array validate error\n")
			return false
		}
	} else {
		if ia[0] != 0x1234567890123456 || ia[len(ia)-1] != 0x3142750893480187 {
			fmt.Printf("ReadInt64Array validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerInt64Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proInt64ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_Int64ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	//go ProducerInt64Arr()

	sw.Begin("sglB")
	for {
		proInt64ArrType()
		rc := conInt64ArrType()
		if !rc {
			panic("conInt64ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test u32 Arr types-------------------------------------------------------------------------
func proUInt32ArrType() {
	lock.Lock()
	defer lock.Unlock()
	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteUInt32Array(make([]uint32, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteUInt32Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]uint32, cnt)
		ia[0] = 0x12345678
		ia[cnt-1] = 0x76543210
		lb.WriteUInt32Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conUInt32ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadInt32Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadInt32Array error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 0x76543210 {
			fmt.Printf("ReadInt32Array validate error\n")
			return false
		}
	} else {
		if ia[0] != 0x12345678 || ia[len(ia)-1] != 0x76543210 {
			fmt.Printf("ReadInt32Array validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerUInt32Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proUInt32ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_UInt32ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerUInt32Arr()

	sw.Begin("sglB")
	for {
		rc := conUInt32ArrType()
		if !rc {
			panic("conInt32ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}
}

// ------------------------Test i32 Arr types-------------------------------------------------------------------------
func proInt32ArrType() {
	lock.Lock()
	defer lock.Unlock()
	rEmp := rand.Intn(16)
	if rEmp == 0 {
		lb.WriteInt32Array(make([]int32, 0))
		bsCount.Add(4)
	} else if rEmp == 1 {
		lb.WriteInt32Array(nil)
		bsCount.Add(4)
	} else {
		cnt := rand.Intn(10000) + 1
		ia := make([]int32, cnt)
		ia[0] = 0x12345678
		ia[cnt-1] = 0x76543210
		lb.WriteInt32Array(ia)
		bsCount.Add(int64(4 + cnt*4))
	}
}

func conInt32ArrType() bool {
	lock.Lock()
	defer lock.Unlock()
	ia, rc := lb.ReadInt32Array()
	if core.Err(rc) {
		if core.IsErrType(rc, core.EC_TRY_AGAIN) {
			return true
		}
		fmt.Printf("ReadInt32Array error\n")
		return false
	}
	if ia == nil {

	} else if len(ia) == 0 {

	} else if len(ia) == 1 {
		if ia[0] != 0x76543210 {
			fmt.Printf("ReadInt32Array validate error\n")
			return false
		}
	} else {
		if ia[0] != 0x12345678 || ia[len(ia)-1] != 0x76543210 {
			fmt.Printf("ReadInt32Array validate error\n")
			return false
		}
	}
	conCount.Add(1)
	return true
}
func ProducerInt32Arr() {
	for i := 0; i < Int32TestCount; i++ {
		proInt32ArrType()
	}
}

func Test_LinkedListByteBuffer_Functional_Int32ArrTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerInt32Arr()

	sw.Begin("sglB")
	for {
		rc := conInt32ArrType()
		if !rc {
			panic("conInt32ArrType error")
		}
		if conCount.Load() >= int64(Int32TestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test Single types-------------------------------------------------------------------------
func proRandomSingleType(i int) {
	lock.Lock()
	defer lock.Unlock()
	tp := i % 11
	if tp == 0 {
		lb.WriteInt8(-128)
		bsCount.Add(1)
	} else if tp == 1 {
		lb.WriteFloat32(3.14)
		bsCount.Add(4)
	} else if tp == 2 {
		lb.WriteFloat64(2.71828)
		bsCount.Add(8)
	} else if tp == 3 {
		lb.WriteInt32(0x7FFFFFFF)
		bsCount.Add(4)
	} else if tp == 4 {
		lb.WriteUInt32(0xFFFFFFFF)
		bsCount.Add(4)
	} else if tp == 5 {
		lb.WriteUInt8(255)
		bsCount.Add(1)
	} else if tp == 6 {
		lb.WriteInt16(32767)
		bsCount.Add(2)
	} else if tp == 7 {
		lb.WriteInt64(12345678987654321)
		bsCount.Add(8)
	} else if tp == 8 {
		lb.WriteBool(true)
		bsCount.Add(1)
	} else if tp == 9 {
		lb.WriteUInt16(65535)
		bsCount.Add(2)
	} else if tp == 10 {
		lb.WriteUInt64(0xFFFFFFFFFFFFFFFF)
		bsCount.Add(8)
	}
}
func conRandomSingleType(i int) int {
	lock.Lock()
	defer lock.Unlock()
	tp := i % 11
	if tp == 0 {
		i8, rc := lb.ReadInt8()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadInt8 Failed %s \n", core.ErrStr(rc)))
		}
		if i8 != -128 {
			panic("ReadInt8 val Failed")
		}
	} else if tp == 1 {
		f32, rc := lb.ReadFloat32()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}
			panic(fmt.Sprintf("ReadF32 Failed %s \n", core.ErrStr(rc)))
		}
		if f32 != 3.14 {
			panic("ReadF32 val Failed")
		}

	} else if tp == 2 {
		f64, rc := lb.ReadFloat64()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}
			panic(fmt.Sprintf("ReadF64 Failed %s \n", core.ErrStr(rc)))
		}
		if f64 != 2.71828 {
			panic("ReadF64 val Failed")
		}

	} else if tp == 3 {
		i32, rc := lb.ReadInt32()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}
			panic(fmt.Sprintf("ReadI32 Failed %s \n", core.ErrStr(rc)))
		}
		if i32 != 0x7FFFFFFF {
			panic("ReadI32 val Failed")
		}
	} else if tp == 4 {
		i32, rc := lb.ReadUInt32()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}
			panic(fmt.Sprintf("ReadU32 Failed %s \n", core.ErrStr(rc)))
		}
		if i32 != 0xFFFFFFFF {
			panic("ReadU32 val Failed")
		}

	} else if tp == 5 {
		i8, rc := lb.ReadUInt8()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadUInt8 Failed %s \n", core.ErrStr(rc)))
		}
		if i8 != 255 {
			panic("ReadUInt8 val Failed")
		}

	} else if tp == 6 {
		iv, rc := lb.ReadInt16()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadInt16 Failed %s \n", core.ErrStr(rc)))
		}
		if iv != 32767 {
			panic("ReadUInt8 val Failed")
		}

	} else if tp == 7 {
		iv, rc := lb.ReadInt64()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadInt64 Failed %s \n", core.ErrStr(rc)))
		}
		if iv != 12345678987654321 {
			panic("ReadUInt8 val Failed")
		}

	} else if tp == 8 {
		iv, rc := lb.ReadBool()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadBool Failed %s \n", core.ErrStr(rc)))
		}
		if iv != true {
			panic("ReadBool val Failed")
		}

	} else if tp == 9 {
		iv, rc := lb.ReadUInt16()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadUInt16 Failed %s \n", core.ErrStr(rc)))
		}
		if iv != 65535 {
			panic("ReadUInt16 val Failed")
		}

	} else if tp == 10 {
		iv, rc := lb.ReadUInt64()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_TRY_AGAIN) {
				return i
			}

			panic(fmt.Sprintf("ReadUInt64 Failed %s \n", core.ErrStr(rc)))
		}
		if iv != 0xFFFFFFFFFFFFFFFF {
			panic("ReadUInt64 val Failed")
		}

	} else {
		panic("^^")
	}

	conCount.Add(1)
	return i + 1
}

func ProducerSingleType() {
	for i := 0; i < SingleTestCount; i++ {
		proRandomSingleType(i)
	}
}

func Test_LinkedListByteBuffer_Functional_SingleTypes(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	go ProducerSingleType()

	sw.Begin("sglB")
	idx := 0
	for {
		idx = conRandomSingleType(idx)
		if conCount.Load() >= int64(SingleTestCount) {
			rSec := float64(sw.Stop("seqE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

// ------------------------Test String ---------------------------------------------------------------------------
func proRandomString1() int32 {
	lock.Lock()
	defer lock.Unlock()
	t := rand.Intn(20)
	if t == 0 {
		rc := lb.WriteString("")
		if core.Err(rc) {
			panic("Write 3")
		}
		bsCount.Add(4)

	} else {
		ll := rand.Intn(32767 * 32)
		ss := strs.CreateSampleString(ll, "@", "$")
		rc := lb.WriteString(ss)
		if core.Err(rc) {
			panic("Write 3")
		}
		bsCount.Add(int64(ll))
	}

	return 0
}
func conRandomString1() bool {
	lock.Lock()
	defer lock.Unlock()
	s, rc := lb.ReadString()
	if core.Err(rc) {
		et, _ := core.ExErr(rc)
		if et == core.EC_TRY_AGAIN {
			return true
		}
		fmt.Printf("Err %s\n", core.ErrStr(rc))
		panic("read failed")
		return false
	} else {

		if s == "" {

		} else if s[0] != '@' || s[len(s)-1] != '$' {
			fmt.Printf("%d -> [%s]", conCount.Load(), s)
			panic("validate failed")
			return false
		}
		conCount.Add(1)
		//fmt.Printf("%d Read str of Len %d\n", rCnt, len(s))
	}
	return true
}
func ProducerString() {
	for i := 0; i < TestCount; i++ {
		proRandomString1()
	}
}

func Test_LinkedListByteBuffer_Functional_StringCon(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	sw.Begin("conB")
	go ProducerString()
	for {
		conRandomString1()
		if conCount.Load() >= int64(TestCount) {
			rSec := float64(sw.Stop("conE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}
	}

}

func Test_LinkedListByteBuffer_Functional_StringSeq(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	sw.Begin("seqB")
	for i := 0; i < TestCount; i++ {
		proRandomString1()
		conRandomString1()
	}
	sw.Stop("seqE")
	rSec := float64(sw.Stop("seqE")) / 1000000000.0
	rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec

	fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
}

// ------------------------Test Strings ---------------------------------------------------------------------------
func conRandomString() {
	lock.Lock()
	defer lock.Unlock()
	bs, rc := lb.ReadStrings()
	if core.Err(rc) {
		et, _ := core.ExErr(rc)
		if et == core.EC_TRY_AGAIN {
			return
		} else {
			panic("ReadString Failed")
		}
	}

	conCount.Add(1)

	if bs == nil {
		//fmt.Printf("Got nil Strings\n")
	} else if len(bs) == 0 {
		//fmt.Printf("Got empty Strings\n")
	} else {
		cnt := len(bs)
		//fmt.Printf("Got %d Strings\n", cnt)
		for i := 0; i < cnt; i++ {
			if len(bs[i]) > 0 {
				if bs[i][0] != '@' || bs[i][len(bs[i])-1] != '$' {
					panic("validate error")
				}
			}
		}
	}
}

func proRandomString() int32 {
	lock.Lock()
	defer lock.Unlock()
	cnt := rand.Intn(18)
	cnt--
	var ss []string = make([]string, 0)
	if cnt < 0 {
		rc := lb.WriteStrings(nil)
		bsCount.Add(4)
		if core.Err(rc) {
			panic("Write 1")
		}
		//fmt.Printf("Prod nil Strings\n")
	} else if cnt == 0 {
		rc := lb.WriteStrings(make([]string, 0))
		if core.Err(rc) {
			panic("Write 2")
		}
		bsCount.Add(4)
		//fmt.Printf("Prod empty Strings\n")
	} else {
		for i := 0; i < cnt; i++ {
			ll := rand.Intn(32767)
			if ll < 3 {
				ss = append(ss, "")
				//fmt.Printf("----- Prod len=%d String\n", 0)
			} else {
				ss = append(ss, strs.CreateSampleString(ll, "@", "$"))
				//fmt.Printf("----- Prod len=%d String\n", ll)
			}
			bsCount.Add(int64(4 + ll))
		}
		rc := lb.WriteStrings(ss)

		if core.Err(rc) {
			panic("Write 3")
		}
		//fmt.Printf("Prod %d Strings\n", cnt)
	}
	return 0
}

func Producer() {
	for i := 0; i < TestCount; i++ {
		proRandomString()
	}
}

func Test_LinkedListByteBuffer_Functional_StringsSeq(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	sw.Begin("seqB")
	for i := 0; i < TestCount; i++ {
		proRandomString()
		conRandomString()
	}
	rSec := float64(sw.Stop("seqE")) / 1000000000.0
	rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec

	fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())

}

func Test_LinkedListByteBuffer_Functional_StringsCon(t *testing.T) {
	if lb.ReadAvailable() != 0 {
		fmt.Printf("read avial is not zero %d\n", lb.ReadAvailable())
	}
	conCount.Store(0)
	bsCount.Store(0)
	sw.Clear()
	sw.Begin("conB")
	go Producer()
	//go Producer()
	for {
		conRandomString()
		if conCount.Load() >= int64(TestCount) {
			rSec := float64(sw.Stop("conE")) / 1000000000.0
			rSpeed := float64(bsCount.Load()/1024.0/1024.0) / rSec
			fmt.Printf("cnt:%d, rSpd:%f lb: %s\n", conCount.Load(), rSpeed, lb.String())
			return
		}

	}
}
