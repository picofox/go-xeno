package config

type DBConnectionConfig struct {
	Tables []string `json:"Tables"`
}

type DBConnectionPoolConfig struct {
	IPV4Addr    string                `json:"IPV4Addr"`
	Type        uint16                `json:"Type"`
	MaxTries    uint16                `json:"MaxTries"`
	TcpPort     uint16                `json:"TcpPort"`
	KeepAlive   uint16                `json:"KeepAlive"`
	DB          string                `json:"DB"`
	Username    string                `json:"Username"`
	Password    string                `json:"Password"`
	ConnParam   string                `json:"ConnParam"`
	Connections []*DBConnectionConfig `json:"Connections"`
}

type DBConfig struct {
	Pools map[string]*DBConnectionPoolConfig `json:"Pools"`
}
