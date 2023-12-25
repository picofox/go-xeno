package memory_test

import (
	"testing"
	"xeno/zohar/core/memory"
)

func Test_LinerBufferFixed_Functional_Basic(t *testing.T) {
	str := "I'm a 宇宙命中率最高 %#$#$^ dsfasfasf 齯躦齬34貰奮"
	strLen := len(str)
	strBA := memory.ByteRef(str, 0, 8)
	strBALen := len(strBA)

	if strLen != strBALen {
		t.Errorf("len mismatch")
	}

	for i := 0; i < strLen; i++ {
		if str[i] != strBA[i] {
			t.Errorf("error conteng ")
		}
	}
}
