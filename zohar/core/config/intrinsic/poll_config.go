package intrinsic

import "fmt"

type PollerConfig struct {
	SubReactorCount int32 `json:"SubReactorCount"`
}

func (ego *PollerConfig) String() string {
	return fmt.Sprintf("SubReactorCount: %d ", ego.SubReactorCount)
}

func NeoPollerConfig(subReactorCount int32) *PollerConfig {
	return &PollerConfig{
		SubReactorCount: subReactorCount,
	}
}
