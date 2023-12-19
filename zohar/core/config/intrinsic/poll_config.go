package intrinsic

import "fmt"

type PollerConfig struct {
	SubReactorCount         int32 `json:"SubReactorCount"`
	SubReactorPulseInterval int   `json:"SubReactorPulseInterval"`
}

func (ego *PollerConfig) String() string {
	return fmt.Sprintf("SubReactorCount: %d, SubReactorPulseInterval %d", ego.SubReactorCount, ego.SubReactorPulseInterval)
}

func NeoPollerConfig(subReactorCount int32) *PollerConfig {
	return &PollerConfig{
		SubReactorCount:         subReactorCount,
		SubReactorPulseInterval: 1000,
	}
}
