package prof

import "strings"

type ConnectionProfiler struct {
	_rtt *RTTProf
}

func (ego *ConnectionProfiler) GetRTTProf() *RTTProf {
	return ego._rtt
}

func NeoConnectionProfiler() *ConnectionProfiler {
	return &ConnectionProfiler{
		_rtt: NeoRTTProf(),
	}
}
func (ego *ConnectionProfiler) String() string {
	var ss strings.Builder
	ss.WriteString(ego._rtt.String())
	return ss.String()
}
