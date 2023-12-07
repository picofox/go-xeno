package config

import (
	"strconv"
	"strings"
)

type NetworkClientTCPConfig struct {
	ServerEndPoints []string                  `json:"ServerEndPoints"`
	Count           int32                     `json:"Count"`
	Handlers        []NetworkTCPHandlerConfig `json:"Handlers"`
}

func (ego *NetworkClientTCPConfig) String() string {
	var ss strings.Builder
	ss.WriteString("ListenerEndPoints=")
	ss.WriteString(strings.Join(ego.ServerEndPoints, ","))
	ss.WriteString("\n")
	ss.WriteString(strconv.Itoa(int(ego.Count)))
	ss.WriteString("\n")
	for _, elem := range ego.Handlers {
		ss.WriteString("Handlers=")
		ss.WriteString(elem.String())
		ss.WriteString("\n")
	}
	return ss.String()
}
