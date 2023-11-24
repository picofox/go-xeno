package config

import (
	"strings"
)

type NetworkServerTCPHandlerConfig struct {
	Name  string            `json:"Name"`
	Param map[string]string `json:"Param"`
}

func (ego *NetworkServerTCPHandlerConfig) String() string {
	var ss strings.Builder
	ss.WriteString(ego.Name)
	ss.WriteString(": ")
	for k, v := range ego.Param {
		ss.WriteString(k)
		ss.WriteString("=")
		ss.WriteString(v)
		ss.WriteString(",")
	}
	return ss.String()
}
