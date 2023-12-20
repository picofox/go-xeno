package prof

import "fmt"

type RTTProf struct {
	_recent  int32
	_lowest  int32
	_highest int32
	_mean    int32
}

func (ego *RTTProf) OnUpdate(v int32) {
	if v < 0 {
		return
	}
	ego._recent = v
	if ego._mean < 0 {
		ego._mean = v
		ego._highest = v
		ego._lowest = v
	}

	ego._mean = (ego._mean + v) / 2
	if v < ego._lowest {
		ego._lowest = v
	} else if v > ego._highest {
		ego._highest = v
	}
}

func (ego *RTTProf) String() string {
	return fmt.Sprintf("[RTT:%d(%d):%d-%d]", ego._recent, ego._mean, ego._lowest, ego._highest)
}

func (ego *RTTProf) Recent() int32 {
	return ego._recent
}

func (ego *RTTProf) Lowest() int32 {
	return ego._lowest
}

func (ego *RTTProf) Highest() int32 {
	return ego._highest
}

func (ego *RTTProf) Mean() int32 {
	return ego._mean
}

func (ego *RTTProf) Reset() {
	ego._recent = -1
	ego._highest = -1
	ego._lowest = -1
	ego._mean = -1
}

func NeoRTTProf() *RTTProf {
	return &RTTProf{
		_recent:  -1,
		_lowest:  -1,
		_highest: -1,
		_mean:    -1,
	}
}
