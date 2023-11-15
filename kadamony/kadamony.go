package main

import (
	"fmt"
	"time"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/sched/timer"
	"xeno/zohar/framework"
	_ "xeno/zohar/framework"
	"xeno/zohar/framework/service/intrinsic"
)

func SetValues(f string, args ...any) {
	fmt.Println(f)
	fmt.Println(args)
}

func TimerCB(a any) {
	fmt.Printf("%s -> TimerCB : (%s)\n", time.Now().String(), a.(*timer.Timer))
}

func CronCB(a any) int32 {
	fmt.Printf("CRonCB (%s)\t", a.(string))
	return 0
}

func onQuiting() {
	fmt.Println("quiting")
}

func main() {
	framework.Initialize()
	rc, errString := config.LoadKadamonyConfig()
	if core.Err(rc) {
		logging.LogFixedWidth(core.LL_SYS, 70, false, errString, "Kadamony Application Initializing ...")
	}
	logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Kadamony Application Initializing ...")

	intrinsic.GetServiceManager().AddCronTask("default", "*/5 * * * * *", CronCB, "cron ->", datatype.TASK_EXEC_NEO_ROUTINE)
	framework.WaitAll()
}
