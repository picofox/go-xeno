package chrono

type IntervalCounter struct {
	_interval int64
	_lastTs   int64
}

func (ego *IntervalCounter) UpdateNow() bool {
	nowTs := GetRealTimeMilli()
	return ego.Update(nowTs)
}

func (ego *IntervalCounter) Update(nowTs int64) bool {
	if nowTs >= ego._lastTs {
		ego._lastTs = nowTs
		return true
	}
	return false
}

func NeoIntervalCounter(interval int64, lastTs ...int64) *IntervalCounter {
	var ts int64
	if len(lastTs) > 0 {
		ts = lastTs[0]
	} else {
		ts = GetRealTimeMilli()
	}

	ic := IntervalCounter{
		_interval: interval,
		_lastTs:   ts,
	}
	return &ic
}
