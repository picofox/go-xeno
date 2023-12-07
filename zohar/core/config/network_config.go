package config

import (
	"strings"
)

type NetworkConfig struct {
	Server NetworkServerConfig `json:"Server"`
	Client NetworkClientConfig `json:"Client"`
}

func (ego *NetworkConfig) String() string {
	var ss strings.Builder
	ss.WriteString("Server:\n")
	ss.WriteString(ego.Server.String())
	ss.WriteString("\n")
	ss.WriteString("Client:\n")
	ss.WriteString(ego.Client.String())
	ss.WriteString("\n")
	return ss.String()
}
