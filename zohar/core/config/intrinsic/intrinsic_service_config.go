package intrinsic

import "strings"

type IntrinsicServiceConfig struct {
	Cron              CronServiceGroupConfig              `json:"Cron"`
	FileSystemWatcher FileSystemWatcherServiceGroupConfig `json:"FileSystemWatcher"`
}

func (ego *IntrinsicServiceConfig) String() string {
	var ss strings.Builder
	ss.WriteString("IntrinsicServiceConfig:\n")
	ss.WriteString("    Cron:\n")

	ss.WriteString(ego.Cron.String())
	ss.WriteString("    FileSystemWatcher:\n")

	ss.WriteString(ego.FileSystemWatcher.String())
	return ss.String()
}
