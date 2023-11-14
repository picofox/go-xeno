package intrinsic

import "fmt"

type CronServiceConfig struct {
	Offset int32 `json:"Offset"`
}

func (ego *CronServiceConfig) String() string {
	return fmt.Sprintf("Off:%d", ego.Offset)
}
