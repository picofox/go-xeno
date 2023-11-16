package intrinsic

import "strings"

type FileSystemWatcherServiceConfig struct {
	Dirs []string `json:"Dirs"`
}

func (ego *FileSystemWatcherServiceConfig) String() string {
	return strings.Join(ego.Dirs, ",")
}
