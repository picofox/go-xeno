package chrono

import (
	"time"
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
