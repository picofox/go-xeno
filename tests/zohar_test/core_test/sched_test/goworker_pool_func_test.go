package sched_test

import (
	"fmt"
	"testing"
	"time"
	"xeno/zohar/core/config"
	"xeno/zohar/core/sched"
)

func cbpStart(worker *sched.GoWorker) {
	startFlag = 1
}

func cbpRun(worker *sched.GoWorker) {
	fmt.Printf("hello form %s\n", worker.String())
}

func cbpStop(worker *sched.GoWorker) {

}

func Test_GoWorkerPool_Functional_Basic(t *testing.T) {
	pool := sched.NeoGoWorkerPool(cbpStart, cbpRun, cbpStop, nil, &config.GoWorkerPoolConfig{Name: "WPool", InitialCount: 10, PulseInterval: 1000})
	if pool == nil {
		t.Errorf("Create GoWorkerPool Failed.")
	}

	pool.Wait()

	pool.SetWorkerCount(10)
	time.Sleep(10000 * time.Millisecond)

	pool.Stop()
	pool.Wait()

}
