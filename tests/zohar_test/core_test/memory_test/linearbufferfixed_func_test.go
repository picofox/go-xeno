package memory_test

import (
	"fmt"
	"testing"
	"xeno/zohar/core/memory"
)

func Test_LinerBufferFixed_Functional_Basic(t *testing.T) {

	buf := memory.NeoLinearBufferFixed(8)
	buf.WriteInt8(int8(1))
	buf.WriteInt8(int8(2))
	buf.WriteInt8(int8(3))
	buf.WriteInt8(int8(4))
	buf.WriteInt8(int8(5))
	buf.WriteInt8(int8(6))
	buf.WriteInt8(int8(7))
	if buf.WriteAvailable() != 1 {
		t.Errorf("Write avialable failed\n")
	}
	buf.ReadInt8()
	if buf.WriteAvailable() != 1 {
		t.Errorf("Write avialable failed\n")
	}
	buf.WriteInt8(int8(8))

	if buf.WriteAvailable() != 0 {
		t.Errorf("Write avialable failed 3\n")
	}
	buf.ReadInt8()
	buf.ReadInt8()
	buf.ReadInt8()
	buf.ReadInt8()
	buf.ReadInt8()
	buf.ReadInt8()
	buf.ReadInt8()
	if buf.ReadAvailable() != 0 {
		t.Errorf("Write avialable failed 4\n")
	}

	fmt.Printf("wa=%d\n", buf.WriteAvailable())

	//str := "I'm a 宇宙命中率最高 %#$#$^ dsfasfasf 齯躦齬34貰奮"
	//strLen := len(str)
	//strBA := memory.ByteRef(str, 0, 8)
	//strBALen := len(strBA)
	//
	//if strLen != strBALen {
	//	t.Errorf("len mismatch")
	//}
	//
	//for i := 0; i < strLen; i++ {
	//	if str[i] != strBA[i] {
	//		t.Errorf("error conteng ")
	//	}
	//}
}
