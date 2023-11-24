package main

import (
	"fmt"
	"gorm.io/gorm/logger"
	"math/bits"
	"reflect"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/db"
	"xeno/zohar/core/inet/server"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/logging/logger_adapter"
	"xeno/zohar/core/mp"
	"xeno/zohar/core/sched/timer"
	"xeno/zohar/framework"
	_ "xeno/zohar/framework"
	"xeno/zohar/framework/service/intrinsic"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func doSomething(a any) int32 {
	fmt.Printf("Do %s \n", a.(string))
	return core.MkSuccess(0)
}

func SetValues(f string, args ...any) {
	fmt.Println(f)
	fmt.Println(args)
}

func TimerFucnCb(s any) int32 {
	fmt.Printf("time due... <%s>\n", s.(*timer.Timer).Object().(string))
	return 0
}

func OnFileSystemChanged(a any) int32 {
	logging.Log(core.LL_DEBUG, "(%d) -> <%s>\n", a.([]any)[0], a.([]any)[1])
	return 0
}

func onQuiting() {
	fmt.Println("quiting")
}

func bsr(x int) int {
	return bits.Len(uint(x)) - 1
}

func CronCB(a any) int32 {
	fmt.Printf("cron tast was triggered...... content is <%ss>\n", a.(string))
	return 0
}

func Task(a any) int32 {
	fmt.Printf("task being executed")
	return 0
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	framework.Initialize()
	intrinsic.GetServiceManager().RegisterFileSystemWatcherHandler(0, OnFileSystemChanged)

	rc, errString := config.LoadKadamonyConfig()
	if core.Err(rc) {
		logging.LogFixedWidth(core.LL_SYS, 70, false, errString, "Kadamony Application Initializing ...")
	}

	var output []reflect.Value = make([]reflect.Value, 0, 10)
	mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "NeoO1L15COT15DecodeServerHandler")
	hdl := output[0].Interface().(*server.O1L15COT15DecodeServerHandler)
	if hdl == nil {
		fmt.Printf("shabile")
	}

	server := server.NeoTcpServer("0.0.0.0", 10000)
	if server == nil {
		logging.Log(core.LL_ERR, "failed")
	}
	server.Listen()
	server.Start()

	logging.Log(core.LL_SYS, "start")
	var iface logger.Interface = logger_adapter.NeoGORMLoggerAdapter(logging.GetLoggerManager().GetDefaultLogger())
	orm, err := gorm.Open(mysql.Open(config.GetKadamonyConfig().DB.Pools["DBP0"].DSN.String()), &gorm.Config{
		Logger: iface,
	})
	if err != nil {
		panic("error")
	}
	orm.AutoMigrate(&Product{})

	cfg := &config.GetKadamonyConfig().DB
	db.GetPoolManager().Initialize(cfg)
	db.GetPoolManager().ConnectDatabase()
	//
	//intrinsic.GetServiceManager().AddCronTask("default", "*/5 * * * * *", CronCB, "dislike you", datatype.TASK_EXEC_EXECUTOR_POOL)
	////
	////concurrent.GetDefaultGoExecutorPool().PostTask(doSomething, "add")
	//
	////timer.GetDefaultTimerManager().AddAbsTimerMilli(3000, 7, 5000, datatype.TASK_EXEC_EXECUTOR_POOL, TimerFucnCb, "yagami")
	//
	////concurrent.GetDefaultGoExecutorPool().PostTask(Task, nil)

	framework.WaitAll()
}
