package intrinsic_test

import (
	"fmt"
	"testing"
	"xeno/zohar/core/datetime"
	"xeno/zohar/core/sched/timer"
)

func cbTimer0(timer *timer.Timer) int32 {
	fmt.Printf("timer [%d] due %d \n", timer.Id(), datetime.GetMonotonicMilli())
	if timer.RemainCount() == 2 {
		timer.Cancel()
	}
	return 0
}

func Test_TimerService_Functional_Basic(t *testing.T) {

	//intrinsic.GetServiceManager().Stop()

}
