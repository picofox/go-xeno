package chrono

import (
	"time"
	"xeno/zohar/core/process"
)

func GetDayBeginMilliStampByTM(t time.Time) int64 {
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return addTime.UnixMilli()
}

func GetDayBeginMilliStampByTMOffset(t *time.Time, offsetDays int64) int64 {
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return addTime.UnixMilli() + 86400000*offsetDays
}

func GetHourBeginMilliStampByTMOffset(t *time.Time, offSetHours int64) int64 {
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return addTime.UnixMilli() + 3600000*offSetHours
}

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

func GetRealTimeMilli() int64 {
	return time.Now().UnixMilli()
}
