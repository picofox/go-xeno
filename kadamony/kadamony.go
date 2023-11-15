package main

import (
	"fmt"
	"time"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/event"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/sched/timer"
	"xeno/zohar/framework"
	_ "xeno/zohar/framework"
)

func SetValues(f string, args ...any) {
	fmt.Println(f)
	fmt.Println(args)
}

func TimerCB(a any) {
	fmt.Printf("%s -> TimerCB : (%s)\n", time.Now().String(), a.(*timer.Timer))
}

func EventCB(a any) int32 {
	fmt.Printf("%s -> EventCB : (%s)\n", time.Now().String(), a.(string))
	return 0
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

	go func() {
		event.GetDefaultEventManager().Register("down", datatype.TASK_EXEC_EXECUTOR_POOL, EventCB, "event trgiiger")
		for {
			time.Sleep(1000 * time.Second)
		}
	}()

	time.Sleep(1 * time.Second)
	event.GetDefaultEventManager().Fire("down", 0xff)

	framework.WaitAll()
}
