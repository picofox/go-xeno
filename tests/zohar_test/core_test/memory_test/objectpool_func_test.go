package memory_test

import (
	"testing"
	"xeno/zohar/core/memory"
)

func Test_ObjectPool_Functional_Basic(t *testing.T) {
	p0 := memory.NeoObjectPool[memory.LinearBuffer](1024, func() *memory.LinearBuffer {
		return memory.NeoLinearBuffer(128)
	})

	lb := p0.Get()
	if lb == nil {
		t.Errorf("Get Pool Element Failed")
	}

}
