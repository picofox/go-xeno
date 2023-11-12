package datetime

import (
	"time"
	"xeno/zohar/core/process"
)

func GetMonotonicNano() int64 {
	return time.Now().Sub(*process.GetTimestampBase()).Nanoseconds()
}
func GetMonotonicMilli() int64 {
	t := time.Now()
	t1 := t.Sub(*process.GetTimestampBase())
	return t1.Milliseconds()
}

func GetRealTimeNano() int64 {
	return time.Now().UnixNano()
}
