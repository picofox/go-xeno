package config

import (
	"fmt"
	"strconv"
	"strings"
)

type GoExecutorPoolConfig struct {
	Name          string `json:"Name"`
	InitialCount  int32  `json:"InitialCount"`
	MaxCount      int32  `json:"MaxCount"`
	MinCount      int32  `json:"MinCount"`
	QueueSize     int    `json:"QueueSize"`
	HighWaterMark int    `json:"HighWaterMark"`
	LowWaterMark  int    `json:"LowWaterMark"`
}

func (ego *GoExecutorPoolConfig) String() string {
	var ss strings.Builder
	ss.WriteString("GoExecutorPoolConfig:\n")
	ss.WriteString("\tInitialCount: ")
	ss.WriteString(strconv.Itoa(int(ego.InitialCount)))
	ss.WriteString("\n")
	ss.WriteString("\tQueueSize: ")
	ss.WriteString(strconv.Itoa(int(ego.QueueSize)))
	ss.WriteString("\n")
	ss.WriteString(fmt.Sprintf("\tCount: %d - %d\n", ego.MinCount, ego.MaxCount))
	ss.WriteString(fmt.Sprintf("\tWaterMark: %d - %d\n", ego.LowWaterMark, ego.HighWaterMark))

	return ss.String()
}
