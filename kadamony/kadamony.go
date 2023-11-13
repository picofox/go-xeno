package main

import (
	"fmt"
	"time"
	"xeno/kadamony/config"
	"xeno/zohar/core"
	"xeno/zohar/core/concurrent"
	"xeno/zohar/core/db"
	"xeno/zohar/core/finalization"
	_ "xeno/zohar/core/initialization"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/process"
	"xeno/zohar/core/sched"
)

func SetValues(f string, args ...any) {
	fmt.Println(f)
	fmt.Println(args)
}

func TimerCB(a any) {
	fmt.Printf("%s -> TimerCB : (%s)\n", time.Now().String(), a.(*sched.Timer))

}

func main() {
	defer finalization.GetGlobalFinalizer().Finalize()

	fmt.Println("Kadamony Intrinsic Initializing ... \t\t\t\t[Done]")
	kinfo := process.ProcessInfoString()
	fmt.Print(kinfo)

	netconfig := memory.CreateTLV(memory.DT_DICT, memory.T_TLV, memory.T_STR, nil)

	//netconfig.PathSet("  play[].nic[].ipv4.address ", "192.168.0.1", datatype.DT_SINGLE, datatype.T_STR, datatype.T_NULL)

	netconfig.PathSet("  default.nic[].ipv4.address ", "192.168.0.1", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[0].ipv4.address ", "192.168.0.2", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[].ipv4.address ", "192.168.0.3", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.gateway ", "192.168.0.1", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[1].id ", int32(998), memory.DT_SINGLE, memory.T_I32, memory.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.dns[] ", "202.230.129.230", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.dns[0] ", "8.8.8.8", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.dns[] ", "111.11.11.111", memory.DT_SINGLE, memory.T_STR, memory.T_NULL)
	netconfig.PathSet("  default.nic[1].ipv4.metric ", 2, memory.DT_SINGLE, memory.T_I32, memory.T_NULL)

	//_, rc := netconfig.PathGet("  default.nic[1].ipv4.metric ")

	fmt.Print("Kadamony Application Initializing ... ")
	rc, errString := config.LoadKadamonyConfig()
	if core.Err(rc) {
		fmt.Println(fmt.Sprintf("... \t\t\t\t[Failed] (%s)\n", errString))
	}

	db.GetPoolManager().Initialize(&config.GetKadamonyConfig().DB)

	fmt.Print("\t\t\t\t[Done]\n")

	concurrent.GetDefaultGoExecutorPool().Start()
	sched.GetDefaultTimerManager().Start()

	sched.GetDefaultTimerManager().AddAbsTimerSecond(5, 10, 1, sched.TIMER_EXEC_EXECUTOR_POOL, TimerCB, nil)

	sched.GetDefaultTimerManager().Wait()
	concurrent.GetDefaultGoExecutorPool().Wait()

}
