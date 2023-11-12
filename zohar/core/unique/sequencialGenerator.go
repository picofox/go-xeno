package unique

import "sync/atomic"

type SequentialGenerator struct {
	_seq atomic.Int64
}

func (ego *SequentialGenerator) Next() int64 {
	return ego._seq.Add(1)
}
