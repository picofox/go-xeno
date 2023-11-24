package config

import (
	"strconv"
	"strings"
)

type NetworkServerTCPConfig struct {
	BindAddr string                          `json:"BindAddr"`
	Port     int                             `json:"Port"`
	PubIP    string                          `json:"PubIP"`
	PriIP    string                          `json:"PriIP"`
	Handlers []NetworkServerTCPHandlerConfig `json:"Handlers"`
}

func (ego *NetworkServerTCPConfig) String() string {
	var ss strings.Builder
	ss.WriteString("BindAddr=")
	ss.WriteString(ego.BindAddr)
	ss.WriteString("\n")
	ss.WriteString("Port=")
	ss.WriteString(strconv.Itoa(int(ego.Port)))
	ss.WriteString("\n")
	ss.WriteString("PubIP=")
	ss.WriteString(ego.PubIP)
	ss.WriteString("\n")
	ss.WriteString("PriIP=")
	ss.WriteString(ego.PriIP)
	ss.WriteString("\n")
	for _, elem := range ego.Handlers {
		ss.WriteString("Handlers=")
		ss.WriteString(elem.String())
		ss.WriteString("\n")
	}
	return ss.String()
}
