package config

import (
	"strconv"
	"strings"
	"xeno/zohar/core/config/intrinsic"
)

type NetworkServerTCPConfig struct {
	ListenerEndPoints          []string                  `json:"ListenerEndPoints"`
	PublicConnectingEndPoints  []string                  `json:"PublicConnectingEndPoints"`
	PrivateConnectingEndPoints []string                  `json:"PrivateConnectingEndPoints"`
	Codec                      string                    `json:"Codec"`
	NoDelay                    bool                      `json:"NoDelay"`
	KeepAlive                  intrinsic.KeepAliveConfig `json:"KeepAlive"`
}

func (ego *NetworkServerTCPConfig) String() string {
	var ss strings.Builder
	ss.WriteString("ListenerEndPoints=")
	ss.WriteString(strings.Join(ego.ListenerEndPoints, ","))
	ss.WriteString("\n")
	ss.WriteString("PublicConnectingEndPoints=")
	ss.WriteString(strings.Join(ego.PublicConnectingEndPoints, ","))
	ss.WriteString("\n")
	ss.WriteString("PrivateConnectingEndPoints=")
	ss.WriteString(strings.Join(ego.PrivateConnectingEndPoints, ","))
	ss.WriteString("\n")
	ss.WriteString("Codec=")
	ss.WriteString(ego.Codec)
	ss.WriteString("\n")
	ss.WriteString("NoDelay=")
	ss.WriteString(strconv.FormatBool(ego.NoDelay))
	ss.WriteString("\n")
	ss.WriteString("KeepAlive=")
	ss.WriteString(ego.KeepAlive.String())
	ss.WriteString("\n")
	return ss.String()
}
