package intrinsic_test

import (
	"fmt"
	"testing"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/datetime"
	"xeno/zohar/core/process"
	"xeno/zohar/core/sched/timer"
	"xeno/zohar/framework/service/intrinsic"
)

func cbTimer0(timer *timer.Timer) int32 {
	fmt.Printf("timer [%d] due %d \n", timer.Id(), datetime.GetMonotonicMilli())
	if timer.RemainCount() == 2 {
		timer.Cancel()
	}
	return 0
}

func Test_TimerService_Functional_Basic(t *testing.T) {
	*process.GetTimestampBase() = time.Now()

	svc := intrinsic.NeoTimerService(intrinsic.GetServiceManager())
	rc := intrinsic.GetServiceManager().AddService("TimerManager", svc)
	if core.Err(rc) {
		t.Errorf("AddService Failed")
	}

	go func() {
		intrinsic.GetServiceManager().Start()
	}()

	svc.AddSecondTimer(3, 6, cbTimer0, nil)

	intrinsic.GetServiceManager().Wait()

	//intrinsic.GetServiceManager().Stop()

}
