package sched_test

import (
	"fmt"
	"testing"
	"time"
	"xeno/zohar/core/concurrent"
	"xeno/zohar/core/config/intrinsic"
)

func cbpStart(worker *concurrent.GoWorker) {
	startFlag = 1
}

func cbpRun(worker *concurrent.GoWorker) {
	fmt.Printf("hello form %s\n", worker.String())
}

func cbpStop(worker *concurrent.GoWorker) {

}

func Test_GoWorkerPool_Functional_Basic(t *testing.T) {
	pool := concurrent.NeoGoWorkerPool(cbpStart, cbpRun, cbpStop, nil, &intrinsic.GoWorkerPoolConfig{Name: "WPool", InitialCount: 10, PulseInterval: 1000})
	if pool == nil {
		t.Errorf("Create GoWorkerPool Failed.")
	}

	pool.Wait()

	pool.SetWorkerCount(10)
	time.Sleep(10000 * time.Millisecond)

	pool.Stop()
	pool.Wait()

}
