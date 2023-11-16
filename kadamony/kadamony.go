package main

import (
	"fmt"
	"time"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/logging"
	"xeno/zohar/framework"
	_ "xeno/zohar/framework"
	"xeno/zohar/framework/service/intrinsic"
)

func SetValues(f string, args ...any) {
	fmt.Println(f)
	fmt.Println(args)
}

func FSWCB(a any) int32 {
	fmt.Printf("%s -> FSWCB : (%d) (%s)\n", time.Now().String(), a.([]any)[0], a.([]any)[1])
	return 0
}

func EventCB(a any) int32 {
	fmt.Printf("%s -> Trigger Event : (%s)\n", time.Now().String(), a.(string))
	return 0
}

func CronCB(a any) int32 {
	fmt.Printf("Trigger Event (%s)\t", a.(string))
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

	intrinsic.GetServiceManager().RegisterFileSystemWatcherHandler(0, FSWCB)

	framework.WaitAll()
}
