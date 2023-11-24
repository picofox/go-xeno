package config

import "strings"

type NetworkConfig struct {
	Server NetworkServerConfig `json:"Server"`
}

func (ego *NetworkConfig) String() string {
	var ss strings.Builder
	ss.WriteString("Server:\n")
	ss.WriteString(ego.Server.String())
	return ss.String()
}
