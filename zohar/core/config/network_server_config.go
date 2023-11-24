package config

import "strings"

type NetworkServerConfig struct {
	TCP map[string]NetworkServerTCPConfig `json:"TCP"`
}

func (ego *NetworkServerConfig) String() string {
	var ss strings.Builder
	for k, v := range ego.TCP {
		ss.WriteString(k)
		ss.WriteString(":\n")
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	return ss.String()
}
