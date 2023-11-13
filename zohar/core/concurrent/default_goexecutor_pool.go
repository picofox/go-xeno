package concurrent

import "sync"

var sDefaultGoExecutorPool *GoExecutorPool = nil
var sDefaultGoExecutorPoolonce sync.Once

func GetDefaultGoExecutorPool() *GoExecutorPool {
	sDefaultGoExecutorPoolonce.Do(func() {
		sDefaultGoExecutorPool = NeoGoExecutorPool()
	})
	return sDefaultGoExecutorPool
}
