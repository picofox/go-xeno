package memory

import (
	"fmt"
	"testing"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/memory"
)

func Test_ByteBufferManager_Functional_Basic(t *testing.T) {
	bap := memory.GetDefaultBufferCacheManager().Get(datatype.SIZE_4K)
	if bap == nil {
		t.Errorf("Get byte form buffer cache manager failed")
	}
	memory.GetDefaultBufferCacheManager().Put(datatype.SIZE_4K, bap)

	var ba []byte = make([]byte, 16)

	fmt.Printf("len is %d/%d\n", len(ba), cap(ba))

}
