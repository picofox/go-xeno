package intrinsic

import (
	"strings"
)

type CronServiceGroupConfig struct {
	Params map[string]CronServiceConfig `json:"Params"`
}

func (ego *CronServiceGroupConfig) String() string {
	var ss strings.Builder
	for k, v := range ego.Params {
		ss.WriteString("    ")
		ss.WriteString(k)
		ss.WriteString(":\n")
		ss.WriteString("        ")
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	return ss.String()
}
