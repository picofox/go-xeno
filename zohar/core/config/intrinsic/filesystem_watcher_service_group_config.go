package intrinsic

import (
	"strings"
)

type FileSystemWatcherServiceGroupConfig struct {
	Params map[string]FileSystemWatcherServiceConfig `json:"Params"`
}

func (ego *FileSystemWatcherServiceGroupConfig) String() string {
	var ss strings.Builder
	for k, v := range ego.Params {
		ss.WriteString("        ")
		ss.WriteString(k)
		ss.WriteString(":\n")
		ss.WriteString("            ")
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	return ss.String()
}
