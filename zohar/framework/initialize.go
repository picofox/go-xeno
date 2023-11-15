package framework

import (
	"fmt"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/initialization"
	"xeno/zohar/core/logging"
	"xeno/zohar/framework/service/intrinsic"
)

var initFrameworkOnce sync.Once
var waitFrameworkOnce sync.Once

func init() {
	Initialize()
}

func StopAll(a any) {
	logging.Log(core.LL_SYS, "Stopping Intrinsic Services")
	intrinsic.GetServiceManager().Stop()

}

func WaitAll() {
	waitFrameworkOnce.Do(func() {
		intrinsic.GetServiceManager().Wait()
		logging.LogFixedWidth(core.LL_SYS, 70, true, "[Stopped]", "Intrinsic Services ... ")
	},
	)

	initialization.WaitAll()
}

func Initialize() {
	initialization.Initialize()

	initFrameworkOnce.Do(
		func() {
			logging.Log(core.LL_SYS, "Zohar Framework Initializing ... ")
			rc := intrinsic.GetServiceManager().Initialize()
			if core.Err(rc) {
				errStr := fmt.Sprintf("Initialize Intrinsic Services ... \t\t\t[Failed:%s]", core.ErrStr(rc))
				logging.LogFixedWidth(core.LL_FATAL, 70, false, core.ErrStr(rc), "Starting Intrinsic Services ...")
				panic(errStr)
			} else {
				logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Initialize Intrinsic Services ...")

			}
			rc = intrinsic.GetServiceManager().Start()
			if core.Err(rc) {
				errStr := fmt.Sprintf("Starting Intrinsic Services ... \t\t\t[Failed:%s]", core.ErrStr(rc))
				logging.LogFixedWidth(core.LL_FATAL, 70, false, core.ErrStr(rc), "Starting Intrinsic Services ...")
				panic(errStr)
			} else {
				logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Starting Intrinsic Services ...")

			}
			initialization.RegisterStopHandler("IntrinsicServicesStopHandler", nil, StopAll)
		},
	)

}
