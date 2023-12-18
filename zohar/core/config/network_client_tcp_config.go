package config

import (
	"strconv"
	"strings"
	"xeno/zohar/core/config/intrinsic"
)

type NetworkClientTCPConfig struct {
	ServerEndPoints []string                  `json:"ServerEndPoints"`
	Count           int32                     `json:"Count"`
	Codec           string                    `json:"Codec"`
	AutoReconnect   bool                      `json:"AutoReconnect"`
	NoDelay         bool                      `json:"NoDelay"`
	KeepAlive       intrinsic.KeepAliveConfig `json:"KeepAlive"`
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
	ss.WriteString("AutoReconnect=")
	ss.WriteString(strconv.FormatBool(ego.AutoReconnect))
	ss.WriteString("\n")
	ss.WriteString("NoDelay=")
	ss.WriteString(strconv.FormatBool(ego.NoDelay))
	ss.WriteString("\n")
	return ss.String()
}
