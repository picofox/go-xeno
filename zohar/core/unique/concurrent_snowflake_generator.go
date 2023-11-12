package unique

import "sync"

type ConcurrentSnowFlakeGenerator struct {
	_lock      sync.Mutex
	_generator SnowFlakeGenerator
}

func (w *ConcurrentSnowFlakeGenerator) Next() int64 {
	w._lock.Lock()
	defer w._lock.Unlock()
	return w._generator.Next()
}

func NeoConcurrentSnowFlakeGenerator(dcId int16, wid int16) *ConcurrentSnowFlakeGenerator {
	g := &ConcurrentSnowFlakeGenerator{
		_generator: SnowFlakeGenerator{
			_lastStamp:    0,
			_workerId:     wid,
			_dataCenterId: dcId,
			_seq:          0,
		},
	}
	return g
}
