package memory_test

import (
	"fmt"
	"testing"
	"xeno/zohar/core/inet/transcomm/prof"
	"xeno/zohar/core/memory"
)

func Test_ObjectPool_Functional_Basic(t *testing.T) {
	p0 := memory.NeoObjectPool[prof.ConnectionProfiler](128, nil)
	lb := p0.Alloc()
	if lb == nil {
		t.Errorf("Get Pool Element Failed")
	}

	p0.Free(&lb)
	if lb != nil {
		t.Errorf("lb should be nil")
	}

	fmt.Printf("P0 : %s\n", p0.String())

}
