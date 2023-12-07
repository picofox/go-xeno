package main

import (
	"fmt"
	"xeno/deus/config"
	"xeno/zohar/core"
	"xeno/zohar/core/db"
	"xeno/zohar/core/inet/transcomm"
	"xeno/zohar/core/logging"
	"xeno/zohar/framework"
	"xeno/zohar/framework/service/intrinsic"
)

func OnFileSystemChanged(a any) int32 {
	logging.Log(core.LL_DEBUG, "(%d) -> <%s>\n", a.([]any)[0], a.([]any)[1])
	return 0
}

func main() {
	logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Deus Application Initializing ...")
	framework.Initialize()
	intrinsic.GetServiceManager().RegisterFileSystemWatcherHandler(0, OnFileSystemChanged)

	rc, errString := config.LoadDeusConfig()
	if core.Err(rc) {
		logging.LogFixedWidth(core.LL_SYS, 70, false, errString, "Deus Application Initializing ...")
	}

	svr := transcomm.NeoTcpServer("Default", config.GetKadamonyConfig().Network.Server.GetTCP("Defaut"), logging.GetLoggerManager().GetDefaultLogger())
	if svr == nil {
		fmt.Printf("Failed")
	}
	transcomm.GetDefaultPoller().RegisterTCPServer(svr)

	rc = svr.Initialize()
	rc = svr.Start()
	fmt.Println(rc)

	cfg := &config.GetKadamonyConfig().DB
	db.GetPoolManager().Initialize(cfg)
	db.GetPoolManager().ConnectDatabase()

	framework.WaitAll()

}
