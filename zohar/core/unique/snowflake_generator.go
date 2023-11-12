package unique

import (
	"time"
)

const (
	INVALID_SNOW_FLAKE_ID = int64(-1)

	workerIDBits     = uint64(5) // 10bit 工作机器ID中的 5bit workerID
	dataCenterIDBits = uint64(5) // 10 bit 工作机器ID中的 5bit dataCenterID
	sequenceBits     = uint64(12)

	maxWorkerID     = int64(-1) ^ (int64(-1) << workerIDBits) //节点ID的最大值 用于防止溢出
	maxDataCenterID = int64(-1) ^ (int64(-1) << dataCenterIDBits)
	maxSequence     = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits // 时间戳向左偏移量
	dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
	workLeft = uint8(12) // workLeft = sequenceBits // 节点IDx向左偏移量
	// 2020-05-20 08:00:00 +0800 CST
	twepoch = int64(1589923200000) // 常量时间戳(毫秒)
)

type SnowFlakeGenerator struct {
	_lastStamp    int64 // 记录上一次ID的时间戳
	_workerId     int16 // 该节点的ID
	_dataCenterId int16 // 该节点的 数据中心ID
	_seq          int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成4096个ID
}

func (w *SnowFlakeGenerator) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *SnowFlakeGenerator) Next() int64 {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w._lastStamp {
		return INVALID_SNOW_FLAKE_ID
	}
	if w._lastStamp == timeStamp {
		w._seq = (w._seq + 1) & maxSequence
		if w._seq == 0 {
			for timeStamp <= w._lastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w._seq = 0
	}

	w._lastStamp = timeStamp
	id := ((timeStamp - twepoch) << timeLeft) |
		(int64(w._dataCenterId) << dataLeft) |
		(int64(w._workerId) << workLeft) |
		w._seq

	return id
}

func NeoSnowFlakeGenerator(dcId int16, wid int16) *SnowFlakeGenerator {
	g := &SnowFlakeGenerator{
		_lastStamp:    0,
		_workerId:     wid,
		_dataCenterId: dcId,
		_seq:          0,
	}
	return g
}
