package timer

import "sync"

var sDefaultTimerManager *TimerManager = nil
var sDefaultTimerManagerOnce sync.Once

func GetDefaultTimerManager() *TimerManager {
	sDefaultTimerManagerOnce.Do(func() {
		sDefaultTimerManager = NeoTimerManager()
	})
	return sDefaultTimerManager
}
