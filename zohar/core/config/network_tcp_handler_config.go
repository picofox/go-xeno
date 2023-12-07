package config

import (
	"strings"
)

type NetworkTCPHandlerConfig struct {
	Name  string            `json:"Name"`
	Param map[string]string `json:"Param"`
}

func (ego *NetworkTCPHandlerConfig) String() string {
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
