package memory

import "sync"

var sByteBufferCahce *ObjectCache[ByteBufferNode]
var sByteBufferCacheOnce sync.Once

func ByteBufferNodeCreator() any {
	return NeoByteBufferNode(4096)
}

func GetByteBufferCache() *ObjectCache[ByteBufferNode] {
	sByteBufferCacheOnce.Do(
		func() {
			sByteBufferCahce = NeoObjectCache[ByteBufferNode](128, ByteBufferNodeCreator)
		})
	return sByteBufferCahce
}

var sBytesCache *ObjectCache[[]byte]
var sBytesCacheOnce sync.Once

func BytesCacheCreator() any {
	ba := make([]byte, 4096)
	return &ba
}

func GetBytesCache() *ObjectCache[[]byte] {
	sBytesCacheOnce.Do(
		func() {
			sBytesCache = NeoObjectCache[[]byte](128, BytesCacheCreator)
		})
	return sBytesCache
}
