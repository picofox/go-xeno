package config

import "strings"

type NetworkClientConfig struct {
	TCP map[string]NetworkClientTCPConfig `json:"TCP"`
}

func (ego *NetworkClientConfig) GetTCP(key string) *NetworkClientTCPConfig {
	cfg, ok := ego.TCP[key]
	if ok {
		return &cfg
	}
	return nil
}

func (ego *NetworkClientConfig) String() string {
	var ss strings.Builder
	for k, v := range ego.TCP {
		ss.WriteString(k)
		ss.WriteString(":\n")
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	return ss.String()
}
