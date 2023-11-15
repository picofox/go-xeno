package intrinsic

import "strings"

type IntrinsicServiceConfig struct {
	Cron CronServiceGroupConfig `json:"Cron"`
}

func (ego *IntrinsicServiceConfig) String() string {
	var ss strings.Builder
	ss.WriteString("IntrinsicServiceConfig:\n")
	ss.WriteString(ego.Cron.String())
	return ss.String()
}
