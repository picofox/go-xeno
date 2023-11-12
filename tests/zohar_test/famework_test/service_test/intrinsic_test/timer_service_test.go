package intrinsic_test

import (
	"fmt"
	"testing"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/datetime"
	"xeno/zohar/core/process"
	"xeno/zohar/framework/service/intrinsic"
)

func cbTimer0(obj any) int32 {
	fmt.Printf("timer due %d\n", datetime.GetMonotonicMilli())
	return 0
}

func Test_TimerService_Functional_Basic(t *testing.T) {
	*process.GetTimestampBase() = time.Now()

	svc := intrinsic.NeoTimerService(intrinsic.GetServiceManager())

	rc := intrinsic.GetServiceManager().AddService("TimerService", svc)
	if core.Err(rc) {
		t.Errorf("AddService Failed")
	}

	go func() {
		intrinsic.GetServiceManager().Start()
	}()

	svc.AddTimer(100, 6, cbTimer0, nil)

	intrinsic.GetServiceManager().Wait()

	//intrinsic.GetServiceManager().Stop()

}
