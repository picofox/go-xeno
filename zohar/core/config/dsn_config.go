package config

import "strings"

type DSNConfig struct {
	Username    string            `json:"Username"`
	Password    string            `json:"Password"`
	Protocol    string            `json:"Protocol"`
	IP          string            `json:"IP"`
	Port        uint16            `json:"Port"`
	DB          string            `json:"DB"`
	ParamString string            `json:"ParamString"`
	Params      map[string]string `json:"Params"`
}

func (ego *DSNConfig) String() string {
	var ss strings.Builder
	if len(ego.Username) > 0 {
		ss.WriteString(ego.Username)
	}
	if len(ego.ParamString) > 0 {
		ss.WriteString(":")
		ss.WriteString(ego.Password)
	}
	if len(ego.Username) > 0 {
		ss.WriteString("@")
	}
	if len(ego.Protocol) > 0 {
		ss.WriteString(ego.Protocol)
	}
	ss.WriteString("/")
	if len(ego.DB) > 0 {
		ss.WriteString(ego.DB)
	}

	if len(ego.ParamString) > 0 && len(ego.Params) > 0 {
		ss.WriteString("?")
		ss.WriteString(ego.ParamString)
		idx := 0
		for k, v := range ego.Params {
			if idx == 0 {
				if ego.ParamString[len(ego.ParamString)-1] != '&' {
					ss.WriteString("&")
				}
			} else {
				ss.WriteString("&")
			}

			ss.WriteString(k)
			ss.WriteString("=")
			ss.WriteString(v)
			idx++
		}
	} else if len(ego.ParamString) > 0 {
		ss.WriteString("?")
		ss.WriteString(ego.ParamString)
	} else if len(ego.Params) > 0 {
		ss.WriteString("?")
		idx := 0
		for k, v := range ego.Params {
			if idx > 0 {
				ss.WriteString("&")
			}
			ss.WriteString(k)
			ss.WriteString("=")
			ss.WriteString(v)
			idx++
		}
	}
	return ss.String()
}

func (ego *DSNConfig) AddParam(k string, v string) {
	ego.Params[k] = v
}

func NeoDSN(name string, pass string, proto string, ip string, port uint16, db string, paramStr string) *DSNConfig {
	return &DSNConfig{
		Username:    name,
		Password:    pass,
		Protocol:    proto,
		DB:          db,
		ParamString: paramStr,
	}
}
