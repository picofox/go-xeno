package intrinsic

import (
	"fmt"
)

type KeepAliveConfig struct {
	Enable         bool  `json:"Enable"`
	TimeoutMillis  int32 `json:"TimeoutMillis"`
	MaxTries       int32 `json:"MaxTries"`
	IntervalMillis int32 `json:"IntervalMillis"`
}

func (ego *KeepAliveConfig) String() string {
	return fmt.Sprintf("e:%t t:%d, n:%d i:%d", ego.Enable, ego.TimeoutMillis, ego.MaxTries, ego.IntervalMillis)
}
