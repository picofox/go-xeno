package config

type GoWorkerPoolConfig struct {
	Name          string `json:"Name"`
	InitialCount  int32  `json:"InitialCount"`
	PulseInterval int64  `json:"PulseInterval"`
}
