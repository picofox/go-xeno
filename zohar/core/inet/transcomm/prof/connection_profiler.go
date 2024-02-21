package prof

import (
	"strconv"
	"strings"
)

type ConnectionProfiler struct {
	_rtt           *RTTProf
	_bytesSent     int64
	_bytesReceived int64
}

func (ego *ConnectionProfiler) OnBytesSent(n int64) {
	ego._bytesSent += n
}

func (ego *ConnectionProfiler) OnBytesReceived(n int64) {
	ego._bytesReceived += n
}

func (ego *ConnectionProfiler) GetRTTProf() *RTTProf {
	return ego._rtt
}

func NeoConnectionProfiler() *ConnectionProfiler {
	return &ConnectionProfiler{
		_rtt:           NeoRTTProf(),
		_bytesSent:     0,
		_bytesReceived: 0,
	}
}

func (ego *ConnectionProfiler) Reset() {
	ego._rtt.Reset()
	ego._bytesReceived = 0
	ego._bytesSent = 0
}

func (ego *ConnectionProfiler) String() string {
	var ss strings.Builder
	ss.WriteString(ego._rtt.String())
	ss.WriteString(" Bytes IO: [")
	ss.WriteString(strconv.FormatInt(ego._bytesReceived, 10))
	ss.WriteString(":")
	ss.WriteString(strconv.FormatInt(ego._bytesSent, 10))
	ss.WriteString("]")
	return ss.String()
}
