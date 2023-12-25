package messages

import (
	"sync"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/memory"
)

var sHeaderCache *memory.ObjectCache[message_buffer.MessageHeader]
var sHeaderCacheOnce sync.Once

func HeaderCacheCreator() any {
	h := message_buffer.NeoMessageHeader()
	return &h
}

func GetHeaderCache() *memory.ObjectCache[message_buffer.MessageHeader] {
	sHeaderCacheOnce.Do(
		func() {
			sHeaderCache = memory.NeoObjectCache[message_buffer.MessageHeader](128, HeaderCacheCreator)
		})
	return sHeaderCache
}
