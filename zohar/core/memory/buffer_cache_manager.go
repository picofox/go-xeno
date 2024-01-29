package memory

import (
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/algorithm"
	"xeno/zohar/core/datatype"
)

var sBufferCacheManager *BufferCacheManager
var sBufferCacheManagerOnce sync.Once

type BufferCacheManager struct {
	_begin  int64
	_caches []*ObjectCache[[]byte]
}

func (ego *BufferCacheManager) GetCache(size int64) *ObjectCache[[]byte] {
	n, rc := algorithm.PowerOf2(size)
	if core.Err(rc) {
		return nil
	}
	n -= ego._begin
	return ego._caches[n]
}

func (ego *BufferCacheManager) Get(size int64) *[]byte {
	n, rc := algorithm.PowerOf2(size)

	if core.Err(rc) {
		return nil
	}

	n -= ego._begin

	ba := ego._caches[n].Get()

	if size != int64(cap(*ba)) {
		panic("size error")
	}
	return ba
}

func (ego *BufferCacheManager) Put(size int64, ba *[]byte) {
	n, rc := algorithm.PowerOf2(size)
	if core.Err(rc) {
		return
	}

	n -= ego._begin
	if size != int64(cap(*ba)) {
		panic("size error")
	}

	ego._caches[n].Put(ba)
}

var sByteBufferCreators = []func() any{
	func() any { //0
		ba := make([]byte, 1)
		return &ba
	},
	func() any { //1
		ba := make([]byte, 2)
		return &ba
	},
	func() any { //2
		ba := make([]byte, 4)
		return &ba
	},
	func() any { //3
		ba := make([]byte, 8)
		return &ba
	},
	func() any { //4
		ba := make([]byte, 16)
		return &ba
	},
	func() any { //5
		ba := make([]byte, 32)
		return &ba
	},
	func() any { //6
		ba := make([]byte, 64)
		return &ba
	},
	func() any { //7
		ba := make([]byte, 128)
		return &ba
	},
	func() any { //8
		ba := make([]byte, 256)
		return &ba
	},
	func() any { //9
		ba := make([]byte, 512)
		return &ba
	},
	func() any { //10
		ba := make([]byte, datatype.SIZE_1K)
		return &ba
	},
	func() any { //1
		ba := make([]byte, datatype.SIZE_2K)
		return &ba
	},
	func() any { //12
		ba := make([]byte, datatype.SIZE_4K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_8K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_16K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_32K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_64K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_128K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_256K)
		return &ba
	},
	func() any {
		ba := make([]byte, datatype.SIZE_512K)
		return &ba
	},
	func() any { //7
		ba := make([]byte, datatype.SIZE_1M)
		return &ba
	},
}

func neoBufferCacheManager(begin int64, length int64, sizeArr []int64) *BufferCacheManager {
	bcm := BufferCacheManager{
		_begin:  begin,
		_caches: make([]*ObjectCache[[]byte], length),
	}

	for i := begin; i < begin+length; i++ {
		bcm._caches[i-begin] = NeoObjectCache[[]byte](sizeArr[i-begin], sByteBufferCreators[i])
	}

	return &bcm
}

func GetDefaultBufferCacheManager() *BufferCacheManager {
	sBufferCacheManagerOnce.Do(
		func() {
			sBufferCacheManager = neoBufferCacheManager(5, 16, []int64{-1, -1, -1, -1, -1, -1, -1, 1024, 1024, -1, -1, -1, -1, -1, -1, -1})
		})
	return sBufferCacheManager
}
