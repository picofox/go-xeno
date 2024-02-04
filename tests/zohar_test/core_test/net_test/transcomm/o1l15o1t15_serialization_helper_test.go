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
var tCount int64 = 100000000

func genOne(idx int64) {
	tmplock.Lock()
	defer tmplock.Unlock()
	//fmt.Printf("%d\n", idx)

	strL := rand.Intn(sAllSizeLen)
	sHelper, rc := messages.InitializeSerialization(buffer, 0, 1)
	if core.Err(rc) {
		panic("InitializeSerialization Failed")
	}
	defer sHelper.Finalize()

	ss, _ := proRandomStrings()
	rc = sHelper.WriteStrings(ss)
	if core.Err(rc) {
		panic("ser ss Failed")
	}

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

	i8arr, _ := proRandomInt8Arr()
	rc = sHelper.WriteInt8s(i8arr)
	if core.Err(rc) {
		panic("ser i8arr Failed")
	}

	u8arr, _ := proRandomUInt8Arr()
	rc = sHelper.WriteUInt8s(u8arr)
	if core.Err(rc) {
		panic("ser ui8arr Failed")
	}

	i16arr, _ := proRandomInt16Arr()
	rc = sHelper.WriteInt16s(i16arr)
	if core.Err(rc) {
		panic("ser i16arr Failed")
	}

	u16arr, _ := proRandomUInt16Arr()
	rc = sHelper.WriteUInt16s(u16arr)
	if core.Err(rc) {
		panic("ser u168arr Failed")
	}

	i32arr, _ := proRandomInt32Arr()
	rc = sHelper.WriteInt32s(i32arr)
	if core.Err(rc) {
		panic("ser i32arr Failed")
	}

	u32arr, _ := proRandomUInt32Arr()
	rc = sHelper.WriteUInt32s(u32arr)
	if core.Err(rc) {
		panic("ser u32arr Failed")
	}

	i64arr, _ := proRandomInt64Arr()
	rc = sHelper.WriteInt64s(i64arr)
	if core.Err(rc) {
		panic("ser i64arr Failed")
	}

	u64arr, _ := proRandomUInt64Arr()
	rc = sHelper.WriteUInt64s(u64arr)
	if core.Err(rc) {
		panic("ser u64arr Failed")
	}

	boolArr, _ := proRandomBoolArr()
	rc = sHelper.WriteBools(boolArr)
	if core.Err(rc) {
		panic("ser boolarr Failed")
	}

	f64arr, _ := proRandomFloat64Arr()
	rc = sHelper.WriteFloat64s(f64arr)
	if core.Err(rc) {
		panic("ser u64arr Failed")
	}

	f32arr, _ := proRandomFloat32Arr()
	rc = sHelper.WriteFloat32s(f32arr)
	if core.Err(rc) {
		panic("ser u32arr Failed")
	}
}

func gen() {
	for i := int64(0); i < tCount; i++ {
		genOne(i)
	}
}

var ci int64 = 0

func baseTestOne() int64 {
	var mGrpId int8
	var cmd int16
	var ll int16
	var el int64
	var rc int32
	var sz int64
	tmplock.Lock()
	defer tmplock.Unlock()
	mGrpId, cmd, ll, el, rc = messages.IsMessageComplete(buffer)
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

	dHelper, r := messages.InitializeDeserialization(buffer, mGrpId, cmd, ll, el)
	if core.Err(r) {
		panic("create dHelper failed")
	}

	defer dHelper.Finalize()

	var ss []string = nil
	ss, rc = dHelper.ReadStrings()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(ss) > 0 {
		for i := 0; i < len(ss); i++ {
			if len(ss[i]) > 1 {
				if ss[i][0] != '@' || ss[i][len(ss[i])-1] != '$' {
					panic("valid ss failed")
				}
			}
		}
	}

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

	var i8arr []int8 = nil
	i8arr, rc = dHelper.ReadInt8s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(i8arr) > 0 {
		for i := 0; i < len(i8arr); i++ {
			if i8arr[i] != int8(i%127) {
				panic("valid i8arr failed")
			}
		}
	}
	var u8arr []uint8 = nil
	u8arr, rc = dHelper.ReadUInt8s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(u8arr) > 0 {
		for i := 0; i < len(u8arr); i++ {
			if u8arr[i] != uint8(i%255) {
				panic("valid i8arr failed")
			}
		}
	}

	var i16arr []int16 = nil
	i16arr, rc = dHelper.ReadInt16s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(i16arr) > 0 {
		for i := 0; i < len(i16arr); i++ {
			if i16arr[i] != int16(i%32767) {
				panic("valid i8arr failed")
			}
		}
	}
	var u16arr []uint16 = nil
	u16arr, rc = dHelper.ReadUInt16s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(u16arr) > 0 {
		for i := 0; i < len(u16arr); i++ {
			if u16arr[i] != uint16(i%65535) {
				panic("valid i8arr failed")
			}
		}
	}

	var i32arr []int32 = nil
	i32arr, rc = dHelper.ReadInt32s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(i32arr) > 0 {
		for i := 0; i < len(i32arr); i++ {
			if i32arr[i] != int32(i%0x7FFFFFFF) {
				panic("valid i32arr failed")
			}
		}
	}
	var u32arr []uint32 = nil
	u32arr, rc = dHelper.ReadUInt32s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(u32arr) > 0 {
		for i := 0; i < len(u32arr); i++ {
			if u32arr[i] != uint32(i%0xFFFFFFFF) {
				panic("valid u32arr failed")
			}
		}
	}

	var i64arr []int64 = nil
	i64arr, rc = dHelper.ReadInt64s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(i64arr) > 0 {
		for i := 0; i < len(i64arr); i++ {
			if i64arr[i] != int64(i) {
				panic("valid i64arr failed")
			}
		}
	}
	var u64arr []uint64 = nil
	u64arr, rc = dHelper.ReadUInt64s()
	if core.Err(rc) {
		panic("read ss failed")
	}
	if len(u64arr) > 0 {
		for i := 0; i < len(u64arr); i++ {
			if u64arr[i] != uint64(i) {
				panic("valid i64arr failed")
			}
		}
	}

	var blarr []bool = nil
	blarr, rc = dHelper.ReadBools()
	if core.Err(rc) {
		panic("read bl arr failed")
	}
	if len(blarr) > 0 {
		for i := 0; i < len(blarr); i++ {
			if blarr[i] != true {
				panic("valid bl arr failed")
			}
		}
	}

	var f64arr []float64 = nil
	f64arr, rc = dHelper.ReadFloat64s()
	if core.Err(rc) {
		panic("read f64 failed")
	}
	if len(f64arr) > 0 {
		for i := 0; i < len(f64arr); i++ {
			if f64arr[i] != 2.71828 {
				panic("valid i64arr failed")
			}
		}
	}

	var f32arr []float32 = nil
	f32arr, rc = dHelper.ReadFloat32s()
	if core.Err(rc) {
		panic("read f32 failed")
	}
	if len(f32arr) > 0 {
		for i := 0; i < len(f32arr); i++ {
			if f32arr[i] != 3.14 {
				panic("valid i32arr failed")
			}
		}
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

func proRandomStrings() ([]string, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyStringArr(), 4
		//fmt.Printf("Prod empty Strings\n")
	} else {
		var ss []string = make([]string, cnt)
		for i := 0; i < cnt; i++ {
			ll := rand.Intn(32767)
			if ll < 2 {
				ss[i] = ""
			} else {
				ss[i] = strs.CreateSampleString(ll, "@", "$")
			}
			bsRet += int64(4 + ll)
		}
		return ss, bsRet
	}
}

func proRandomBoolArr() ([]bool, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyBoolArray(), 4
	} else {
		var ss []bool = make([]bool, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = true
			bsRet += 1
		}
		return ss, bsRet
	}
}

func proRandomInt8Arr() ([]int8, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyInt8Arr(), 4
	} else {
		var ss []int8 = make([]int8, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = int8(i % 127)
			bsRet += 1
		}
		return ss, bsRet
	}
}
func proRandomUInt8Arr() ([]uint8, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyUInt8Arr(), 4
	} else {
		var ss []uint8 = make([]uint8, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = uint8(i % 255)
			bsRet += 1
		}
		return ss, bsRet
	}
}

func proRandomInt16Arr() ([]int16, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyInt16Arr(), 4
	} else {
		var ss []int16 = make([]int16, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = int16(i % 32767)
			bsRet += 2
		}
		return ss, bsRet
	}
}
func proRandomUInt16Arr() ([]uint16, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyUInt16Arr(), 4
	} else {
		var ss []uint16 = make([]uint16, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = uint16(i % 65535)
			bsRet += 2
		}
		return ss, bsRet
	}
}

func proRandomInt32Arr() ([]int32, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyInt32Arr(), 4
	} else {
		var ss []int32 = make([]int32, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = int32(i % 0x7FFFFFFF)
			bsRet += 4
		}
		return ss, bsRet
	}
}
func proRandomUInt32Arr() ([]uint32, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyUInt32Arr(), 4
	} else {
		var ss []uint32 = make([]uint32, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = uint32(i % 0xFFFFFFFF)
			bsRet += 4
		}
		return ss, bsRet
	}
}

func proRandomInt64Arr() ([]int64, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyInt64Arr(), 4
	} else {
		var ss []int64 = make([]int64, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = int64(i)
			bsRet += 8
		}
		return ss, bsRet
	}
}
func proRandomUInt64Arr() ([]uint64, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyUInt64Arr(), 4
	} else {
		var ss []uint64 = make([]uint64, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = uint64(i)
			bsRet += 8
		}
		return ss, bsRet
	}
}

func proRandomFloat32Arr() ([]float32, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyFloat32Arr(), 4
	} else {
		var ss []float32 = make([]float32, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = 3.14
			bsRet += 4
		}
		return ss, bsRet
	}
}

func proRandomFloat64Arr() ([]float64, int64) {
	cnt := rand.Intn(18)
	cnt--
	var bsRet int64 = 0
	if cnt < 0 {
		return nil, 4
	} else if cnt == 0 {
		return memory.ConstEmptyFloat64Arr(), 4
	} else {
		var ss []float64 = make([]float64, cnt)
		for i := 0; i < cnt; i++ {
			ss[i] = 2.71828
			bsRet += 8
		}
		return ss, bsRet
	}
}
