package chrono

type StopWatch struct {
	_beginTS int64
	_record  []int64
}

func (ego *StopWatch) Begin() {
	ego._beginTS = GetRealTimeNano()
}

func (ego *StopWatch) Stop() int64 {
	r := GetRealTimeNano() - ego._beginTS
	return r
}

func (ego *StopWatch) Mark() int64 {
	e := GetRealTimeNano()
	ego._record = append(ego._record, e)
	return e - ego._beginTS
}

func (ego *StopWatch) GetRecord(times int) int64 {
	return ego._record[times] - ego._beginTS
}

func (ego *StopWatch) GetRecordRaw(times int) int64 {
	return ego._record[times]
}

func (ego *StopWatch) GetRecordRel(times int) int64 {
	if times > 0 {
		return ego._record[times] - ego._record[times-1]
	} else {
		return ego._record[times] - ego._beginTS
	}
}

func (ego *StopWatch) Clear() {
	ego._beginTS = 0
	ego._record = make([]int64, 0)
}

func NeoStopWatch() *StopWatch {
	return &StopWatch{
		_beginTS: 0,
		_record:  make([]int64, 0),
	}
}
