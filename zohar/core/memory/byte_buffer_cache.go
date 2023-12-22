package memory

import "sync"

var sByteBufferCahce *ObjectCache[ByteBufferNode]
var sByteBufferCacheOnce sync.Once

func ByteBufferNode4KCreator() any {
	return NeoByteBufferNode(4096)
}

func GetByteBuffer4KCache() *ObjectCache[ByteBufferNode] {
	sByteBufferCacheOnce.Do(
		func() {
			sByteBufferCahce = NeoObjectCache[ByteBufferNode](128, ByteBufferNode4KCreator)
		})
	return sByteBufferCahce
}
