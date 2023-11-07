package memory_test

import (
	"strings"
	"testing"
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
)

func Test_RingByteBuffer_Functional_Basic(t *testing.T) {
	srcBa := make([]byte, 128, 128)
	for i := 0; i < 10; i++ {
		srcBa[i] = byte(i)
	}
	dstBa := make([]byte, 128)

	var buf *memory.RingBuffer
	buf = memory.NeoByteBuffer(10)
	buf.WriteBytes(srcBa, 0, 10)

	r0, r1 := buf.BytesRef()
	if r1 != nil {
		t.Errorf("Simple 1st Write 10bs Failed")
	}
	if r0 == nil {
		t.Errorf("Simple 1st Write 10bs Failed")
	}

	nbRead := buf.ReadBytes(dstBa, 0, 5, true)
	if nbRead != 5 {
		t.Errorf("Simple 1st Read 5bs Failed got %d", nbRead)
	}

	for i := 0; i < 5; i++ {
		if dstBa[i] != byte(i) {
			t.Errorf("Simple 1st Read wrong data %d", dstBa[i])
		}
	}

	nbRead = buf.ReadBytes(dstBa, 0, 6, true)
	if nbRead != 0 {
		t.Errorf("Strict mode should not do this")
	}
	nbRead = buf.ReadBytes(dstBa, 0, 6, false)
	if nbRead != 5 {
		t.Errorf("Simple 1st Read 5bs Failed got %d", nbRead)
	}
	for i := 0; i < 5; i++ {
		if dstBa[i] != byte(i+5) {
			t.Errorf("Simple 1st Read wrong data %d", dstBa[i])
		}
	}

	buf.WriteBytes(srcBa, 0, 7)
	buf.ReadBytes(dstBa, 0, 5, true)
	lenForW := buf.WriteAvailable()
	if lenForW != 8 {
		t.Errorf("WriteAvai should be 8, but got %d", lenForW)
	}

	buf.WriteBytes(srcBa, 2, 8)
	r0, r1 = buf.BytesRef()
	if r1 == nil {
		t.Errorf("Simple 2st Write 10bs Failed")
	}
	if r0 == nil {
		t.Errorf("Simple 2st Write 10bs Failed")
	}

	if len(r0) != 5 {
		t.Errorf("Simple 2st Write r0 len should be 5, but got %d", len(r0))
	}
	if len(r1) != 5 {
		t.Errorf("Simple 2st Write r0 len should be 5, but got %d", len(r1))
	}

}

func Test_RingByteBuffer_Functional_ReSpace(t *testing.T) {
	dstBa := make([]byte, 1024)
	var buf *memory.RingBuffer
	buf = memory.NeoByteBuffer(10)
	rc := buf.WriteBytes([]byte("01234567890123456789"), 0, 20)
	if core.Err(rc) {
		t.Errorf("Write 20 bytes to 10 bytes ringbuffer failed")
	}
	nbRead := buf.ReadBytes(dstBa, 0, int64(len(dstBa)), false)
	if nbRead != 20 {
		t.Errorf("ReadBytes Failed: should 20, but got %d", rc)
	}
	str := string(dstBa[0:20])
	if strings.Compare(str, "01234567890123456789") != 0 {
		t.Errorf("Data Validate Failed")
	}

	rc = buf.WriteBytes([]byte("abcdefghijklmnopqrst"), 0, 20)
	if core.Err(rc) {
		t.Errorf("wp < rp case : Write 20 bytes to 10 bytes ringbuffer failed")
	}

	rc = buf.WriteBytes([]byte("ABCDEFGHIJKLMNOPARST"), 0, 20)
	if core.Err(rc) {
		t.Errorf("wp < rp case : ReSpace by 20 bytes failed")
	}

	nbRead = buf.ReadBytes(dstBa, 0, int64(len(dstBa)), false)
	if nbRead != 40 {
		t.Errorf("\"wp < rp case :ReadBytes Failed: should 20, but got %d", rc)
	}
	str = string(dstBa[0:40])

	if strings.Compare(str, "abcdefghijklmnopqrstABCDEFGHIJKLMNOPARST") != 0 {
		t.Errorf("Data Validate Failed")
	}

}

func Test_RingByteBuffer_Functional_String(t *testing.T) {
	var buf *memory.RingBuffer
	buf = memory.NeoByteBuffer(28)

	rc := buf.WriteString("0123456789")
	if core.Err(rc) {
		t.Errorf("case0 (Normal Write): WriteSimpleString failed")
	}
	str, rc := buf.ReadString()
	if core.Err(rc) {
		t.Errorf("case0 (Normal Write): ReadSimpleString failed")
	}
	if str != "0123456789" {
		t.Errorf("case0 (Normal Write): validation Failed")
	}

	buf.Clear()

	rc = buf.WriteString("0123456789")
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): WriteString at end failed")
	}
	rc = buf.WriteString("abcdefghij")
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): WriteString at end failed")
	}
	str, rc = buf.ReadString()
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): ReadString at endfailed")
	}
	if str != "0123456789" {
		t.Errorf("case1 (Normal Write): validation Failed")
	}
	str, rc = buf.ReadString()
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): ReadString at endfailed")
	}
	if str != "abcdefghij" {
		t.Errorf("case1 (Normal Write): validation Failed")
	}

	buf.Clear()

	rc = buf.WriteString("012345")
	if core.Err(rc) {
		t.Errorf("case3 (Half Write): WriteString at half failed")
	}
	str, rc = buf.ReadString()
	if core.Err(rc) {
		t.Errorf("case3 (Half Write): ReadString at endfailed")
	}
	if str != "012345" {
		t.Errorf("case3 (Half Write): validation Failed")
	}

	rc = buf.WriteString("0123456789abc")
	if core.Err(rc) {
		t.Errorf("case3 (Half Write): WriteString at half failed")
	}

	rc = buf.WriteString("fox")
	if core.Err(rc) {
		t.Errorf("case3 (Half Write): WriteString at half failed")
	}
	str, rc = buf.ReadString()
	if core.Err(rc) {
		t.Errorf("case3 (Half Write): ReadString at endfailed")
	}
	if str != "0123456789abc" {
		t.Errorf("case1 (Half Write): validation Failed")
	}
	str, rc = buf.ReadString()
	if core.Err(rc) {
		t.Errorf("case3 (Half Write): ReadString at endfailed")
	}
	if str != "fox" {
		t.Errorf("case1 (Half Write): validation Failed")
	}

}

func Test_RingByteBuffer_Functional_Int32(t *testing.T) {
	var buf *memory.RingBuffer
	buf = memory.NeoByteBuffer(11)

	rc := buf.WriteInt32(-3242342)
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): Write  failed")
	}
	rc = buf.WriteInt32(459783498)
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): Write  failed")
	}
	rc = buf.WriteInt32(-1)
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): Write  failed")
	}

	iv, rc := buf.ReadInt32()
	if core.Err(rc) || iv != -3242342 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

	iv, rc = buf.ReadInt32()
	if core.Err(rc) || iv != 459783498 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}
	iv, rc = buf.ReadInt32()
	if core.Err(rc) || iv != -1 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

	rc = buf.WriteUInt32(0xFFFFFFFF)
	if core.Err(rc) {
		t.Errorf("case1 (UInt64 Write): Write  failed")
	}
	rc = buf.WriteUInt32(3345235)
	if core.Err(rc) {
		t.Errorf("case1 (UInt64 Write): Write  failed")
	}
	rc = buf.WriteUInt32(0)
	if core.Err(rc) {
		t.Errorf("case1 (UInt64 Write): Write  failed")
	}

	uv, rc := buf.ReadUInt32()
	if core.Err(rc) || uv != 0xFFFFFFFF {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

	uv, rc = buf.ReadUInt32()
	if core.Err(rc) || uv != 3345235 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}
	uv, rc = buf.ReadUInt32()
	if core.Err(rc) || uv != 0 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

}

func Test_RingByteBuffer_Functional_Int64(t *testing.T) {
	var buf *memory.RingBuffer
	buf = memory.NeoByteBuffer(19)

	rc := buf.WriteInt64(-3242342)
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): Write  failed")
	}
	rc = buf.WriteInt64(4597834573498)
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): Write  failed")
	}
	rc = buf.WriteInt64(0)
	if core.Err(rc) {
		t.Errorf("case1 (Normal Write): Write  failed")
	}

	iv, rc := buf.ReadInt64()
	if core.Err(rc) || iv != -3242342 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

	iv, rc = buf.ReadInt64()
	if core.Err(rc) || iv != 4597834573498 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}
	iv, rc = buf.ReadInt64()
	if core.Err(rc) || iv != 0 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

	rc = buf.WriteUInt64(0xFFFFFFFFFFFFFFFF)
	if core.Err(rc) {
		t.Errorf("case1 (UInt64 Write): Write  failed")
	}
	rc = buf.WriteUInt64(124597834573498)
	if core.Err(rc) {
		t.Errorf("case1 (UInt64 Write): Write  failed")
	}
	rc = buf.WriteUInt64(47802343242342)
	if core.Err(rc) {
		t.Errorf("case1 (UInt64 Write): Write  failed")
	}

	uv, rc := buf.ReadUInt64()
	if core.Err(rc) || uv != 0xFFFFFFFFFFFFFFFF {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

	uv, rc = buf.ReadUInt64()
	if core.Err(rc) || uv != 124597834573498 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}
	uv, rc = buf.ReadUInt64()
	if core.Err(rc) || uv != 47802343242342 {
		t.Errorf("case1 (Normal Write): Read or Validate failed")
	}

}
