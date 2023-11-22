package config

type DBConnectionConfig struct {
	Tables []string `json:"Tables"`
}

type DBConnectionPoolConfig struct {
	Type        uint16                `json:"Type"`
	MaxTries    uint16                `json:"MaxTries"`
	KeepAlive   uint16                `json:"KeepAlive"`
	DSN         DSNConfig             `json:"DSN"`
	Connections []*DBConnectionConfig `json:"Connections"`
}

type DBConfig struct {
	Pools map[string]*DBConnectionPoolConfig `json:"Pools"`
}
