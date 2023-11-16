package main

import (
	"fmt"
	"path/filepath"
	"time"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/fs"
	"xeno/zohar/core/logging"
	"xeno/zohar/framework"
	_ "xeno/zohar/framework"
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

	fsw := fs.NeoFileSystemWatcher()

	path, _ := filepath.Abs(".")
	fmt.Println("wathc +" + path)

	fsw.AddDir(path)
	fsw.RegisterHandler(0, FSWCB)

	fsw.Start()
	fsw.Stop()
	framework.WaitAll()
}
