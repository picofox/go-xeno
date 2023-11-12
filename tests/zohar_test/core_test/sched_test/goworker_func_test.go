package sched_test

import (
	"sync"
	"testing"
	"time"
	"xeno/zohar/core/sched"
)

var startFlag = 0

func cbStart(worker *sched.GoWorker) {
	startFlag = 1
}

func cbRun(worker *sched.GoWorker) {

}

func cbStop(worker *sched.GoWorker) {

}

func Test_GoWorker_Functional_Basic(t *testing.T) {
	wg := sync.WaitGroup{}
	w := sched.NeoGoWorker("test", 0, cbStart, cbRun, cbStop, nil, 1000, &wg)
	if w == nil {
		t.Errorf("Neo GoWorker Failed")
	}

	w.Start()
	time.Sleep(100 * time.Millisecond)
	if startFlag != 1 {
		t.Errorf("start cb may not called")
	}
	w.Stop()
	w.Wait()
}
