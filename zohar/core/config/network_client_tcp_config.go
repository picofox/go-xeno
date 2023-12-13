package config

import (
	"strconv"
	"strings"
)

type NetworkClientTCPConfig struct {
	ServerEndPoints []string `json:"ServerEndPoints"`
	Count           int32    `json:"Count"`
	Codec           string   `json:"Codec"`
}

func (ego *NetworkClientTCPConfig) String() string {
	var ss strings.Builder
	ss.WriteString("ListenerEndPoints=")
	ss.WriteString(strings.Join(ego.ServerEndPoints, ","))
	ss.WriteString("\n")
	ss.WriteString(strconv.Itoa(int(ego.Count)))
	ss.WriteString("\n")
	ss.WriteString("Codec=")
	ss.WriteString(ego.Codec)
	ss.WriteString("\n")

	return ss.String()
}
