package transcomm

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

var sAllSize []int64 = []int64{32768, 2, 11, 1024*32 + 1, 3}

// var sAllSize []int64 = []int64{11, 2, 234, 4096, 12344, 32760, 32755, 32765, 565475, 1024 * 32, 1024 * 33, 1024*32 + 1, 1024*1024*1 + 1}
var sAllSizeLen int = len(sAllSize)
var buffer *memory.LinkedListByteBuffer = memory.NeoLinkedListByteBuffer(datatype.SIZE_4K)
var totalBs int64 = 0
var tmplock sync.Mutex
var tCount int64 = 1000000000

func genOne(idx int64) {
	tmplock.Lock()
	defer tmplock.Unlock()
	//fmt.Printf("%d\n", idx)

	strL := rand.Intn(sAllSizeLen)
	sHelper, rc := messages.InitializeSerialization(buffer, false, 1)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}
	defer sHelper.FinalizeSerialization()

	s := strs.CreateSampleString(int(sAllSize[strL]), "@", "$")
	if len(s) != int(sAllSize[strL]) {
		panic("create string failed")
	}
	rc = sHelper.WriteString(s)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}
	rc = sHelper.WriteString(s)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}
	rc = sHelper.WriteInt64(idx)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}

	rc = sHelper.WriteBytes(nil)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}

	rc = sHelper.WriteInt16(int16(idx % 32767))
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}

	rc = sHelper.WriteInt8(int8(idx % 127))
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}

	rc = sHelper.WriteUInt8(uint8(idx % 255))
	if core.Err(rc) {
		panic("ser uint8 Failed")
	}

	rc = sHelper.WriteUInt16(uint16(idx % 65535))
	if core.Err(rc) {
		panic("ser uint16 Failed")
	}

	rc = sHelper.WriteUInt32(uint32(idx % 0xFFFFFFFF))
	if core.Err(rc) {
		panic("ser uint32 Failed")
	}

	rc = sHelper.WriteUInt64(uint64(idx))
	if core.Err(rc) {
		panic("ser uint64 Failed")
	}

	var bv bool
	if idx%2 == 0 {
		bv = false
	} else {
		bv = true
	}
	rc = sHelper.WriteBool(bv)
	if core.Err(rc) {
		panic("ser bool Failed")
	}

	rc = sHelper.WriteFloat32(3.14)
	if core.Err(rc) {
		panic("ser f32 Failed")
	}

	rc = sHelper.WriteFloat64(2.71828)
	if core.Err(rc) {
		panic("ser f64 Failed")
	}

}

func gen() {
	for i := int64(0); i < tCount; i++ {
		genOne(i)
	}
}

var ci int64 = 0

func baseTestOne() int64 {
	var isInternal bool
	var cmd int16
	var ll int16
	var el int64
	var rc int32
	var sz int64
	tmplock.Lock()
	defer tmplock.Unlock()
	isInternal, cmd, ll, el, rc = messages.IsMessageComplete(buffer)
	if core.Err(rc) {
		if !core.IsErrType(rc, core.EC_TRY_AGAIN) {
			panic("comple failed")
		} else {
			return 0
		}
	}

	sz = int64(ll)
	sz += el
	if el > 0 {
		sz += 8
	}

	dHelper, r := messages.InitializeDeserialization(buffer, isInternal, cmd, ll, el)
	if core.Err(r) {
		panic("create dHelper failed")
	}

	defer dHelper.FinalizeDeserialization()

	//fmt.Printf("%t, %d, %d, %d\n", isInternal, cmd, ll, el)
	var rs string
	rs, rc = dHelper.ReadString()
	if core.Err(rc) {
		panic("read string failed")
	}
	if len(rs) > 2 && (rs[0] != '@' || rs[len(rs)-1] != '$') {
		panic("validate failed")
	}
	rs, rc = dHelper.ReadString()
	if core.Err(rc) {
		panic("read string failed")
	}
	if len(rs) > 2 && (rs[0] != '@' || rs[len(rs)-1] != '$') {
		panic("validate failed")
	}
	var i64v int64 = 0
	i64v, rc = dHelper.ReadInt64()
	if core.Err(rc) {
		panic("read i64 failed")
	}
	if i64v != ci {
		panic("validate i64 failed")
	}

	var rba []byte
	rba, rc = dHelper.ReadBytes()
	if core.Err(rc) {
		panic("read i64 failed")
	}
	if rba != nil {
		panic("validate rba failed")
	}

	var i16v int16 = 0
	i16v, rc = dHelper.ReadInt16()
	if core.Err(rc) {
		panic("read i16 failed")
	}
	if i16v != int16(ci%32767) {
		panic("validate i16 failed")
	}
	var i8v int8 = 0
	i8v, rc = dHelper.ReadInt8()
	if core.Err(rc) {
		panic("read i16 failed")
	}
	if i8v != int8(ci%127) {
		panic("validate i8 failed")
	}

	var u8v uint8 = 0
	u8v, rc = dHelper.ReadUInt8()
	if core.Err(rc) {
		panic("read u8 failed")
	}
	if u8v != uint8(ci%255) {
		panic("validate u8 failed")
	}

	var u16v uint16 = 0
	u16v, rc = dHelper.ReadUInt16()
	if core.Err(rc) {
		panic("read u16 failed")
	}
	if u16v != uint16(ci%65535) {
		panic("validate u16 failed")
	}

	var u32v uint32 = 0
	u32v, rc = dHelper.ReadUInt32()
	if core.Err(rc) {
		panic("read u32 failed")
	}
	if u32v != uint32(ci%0xFFFFFFFF) {
		panic("validate u32 failed")
	}

	var u64v uint64 = 0
	u64v, rc = dHelper.ReadUInt64()
	if core.Err(rc) {
		panic("read u64 failed")
	}
	if u64v != uint64(ci) {
		panic("validate u64 failed")
	}

	var bv bool
	bv, rc = dHelper.ReadBool()
	if core.Err(rc) {
		panic("read b failed")
	}
	var cbv bool
	if ci%2 == 0 {
		cbv = false
	} else {
		cbv = true
	}
	if bv != cbv {
		panic("validate bool failed")
	}

	var f32v float32 = 0
	f32v, rc = dHelper.ReadFloat32()
	if core.Err(rc) {
		panic("read f32 failed")
	}
	if f32v != 3.14 {
		panic("validate u32 failed")
	}

	var f64v float64 = 0
	f64v, rc = dHelper.ReadFloat64()
	if core.Err(rc) {
		panic("read f64 failed")
	}
	if f64v != 2.71828 {
		panic("validate u64 failed")
	}

	ci++
	return sz
}

var psw *chrono.StopWatch = chrono.NeoStopWatch()

func Test_O1L15O1C15Serializer_Functional_Basic(t *testing.T) {
	var idx int64 = 0
	rand.Seed(84950729087)
	var bsPrintNext int64 = 1
	go gen()
	psw.Begin(".")
	for {

		sz := baseTestOne()
		if sz == 0 {
			continue
		}
		//fmt.Printf("get %d msg\n", sz)
		totalBs += sz
		idx++
		if idx >= tCount {
			return
		}
		if idx%100000 == 0 {
			fmt.Printf("Total: %d count\n", idx)
		}
		if totalBs > bsPrintNext*1048576*1000 {
			psw.Mark(".")
			t := psw.GetRecordRecent()
			avgspd := float64(totalBs) / 1048576.0 / (float64(t) / 1000000000.0)
			fmt.Printf("Total: %.2f MB MeanSPD = %.2f MBPS\n", float64(totalBs)/1048576.0, avgspd)
			bsPrintNext++
		}
	}

}
