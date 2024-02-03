package chrono

import (
	"fmt"
	"strings"
)

type StopWatch struct {
	_record []int64
	_descs  []string
}

func (ego *StopWatch) String() string {
	var ss strings.Builder
	if len(ego._record) != len(ego._descs) {
		panic("data broken between record and descs")
	}
	if len(ego._record) == 1 {
		return fmt.Sprintf("Start @ %d", ego._record[0])
	} else if len(ego._record) == 0 {
		return fmt.Sprintf("Not Started StopWatch")
	}

	for i := 1; i < len(ego._record); i++ {
		diff := ego._record[i] - ego._record[i-1]
		ss.WriteString("[")
		ss.WriteString(ego._descs[i-1])
		ss.WriteString("-")
		ss.WriteString(ego._descs[i])
		ss.WriteString("]")
		ss.WriteString(fmt.Sprintf(":(%.6f)", float64(diff)/1000000.0))
		if i < len(ego._record)-1 {
			ss.WriteString(", ")
		} else {
			diff := ego._record[i] - ego._record[0]
			ss.WriteString(fmt.Sprintf(" | (%.5f)", float64(diff)/1000000.0))
		}
	}
	return ss.String()
}

func (ego *StopWatch) Begin(begStr string) {
	b := GetRealTimeNano()
	ego._record = append(ego._record, b)
	ego._descs = append(ego._descs, begStr)
}

func (ego *StopWatch) Stop(endStr string) int64 {
	m := GetRealTimeNano()
	ego._record = append(ego._record, m)
	ego._descs = append(ego._descs, endStr)
	return m - ego._record[0]
}

func (ego *StopWatch) Mark(markStr string) {
	e := GetRealTimeNano()
	ego._record = append(ego._record, e)
	ego._descs = append(ego._descs, markStr)

}

func (ego *StopWatch) GetRecord(times int) int64 {
	return ego._record[times] - ego._record[0]
}

func (ego *StopWatch) GetRecordRaw(times int) int64 {
	return ego._record[times]
}

func (ego *StopWatch) GetRecordRel(times int) int64 {
	return ego._record[times] - ego._record[times-1]
}

func (ego *StopWatch) GetRecordRecent() int64 {
	return ego._record[len(ego._record)-1] - ego._record[0]
}

func (ego *StopWatch) Clear() {
	ego._record = make([]int64, 0)
	ego._descs = make([]string, 0)
}

func NeoStopWatch() *StopWatch {
	return &StopWatch{
		_record: make([]int64, 0),
		_descs:  make([]string, 0),
	}
}
