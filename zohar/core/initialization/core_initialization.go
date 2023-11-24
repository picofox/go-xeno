package initialization

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/cmdline"
	"xeno/zohar/core/concurrent"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/finalization"
	"xeno/zohar/core/inet/server"
	"xeno/zohar/core/io"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/mp"
	"xeno/zohar/core/process"
	"xeno/zohar/core/sched/timer"
)

var initOnce sync.Once
var waitOnce sync.Once
var sProcessControlChannel chan os.Signal = make(chan os.Signal)
var sFinalizer *finalization.Finalizer = finalization.NeoFinalizer()

func init() {
	Initialize()
}

func RegisterStopHandler(name string, sub any, m func(any)) {
	sFinalizer.Register(name, sub, m)
}

func stopServiceHandler() {
	sFinalizer.Finalize()
	CoreStopAll()
}

func onQuiting(s *os.Signal) {
	stopServiceHandler()
}

func installSignalHandler() {
	signal.Notify(sProcessControlChannel, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range sProcessControlChannel {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
				logging.Log(core.LL_SYS, "Catch Sig %s", s.String())
				onQuiting(&s)
			default:
				logging.Log(core.LL_SYS, "Ignore Sig %s", s.String())
			}
		}
	}()
}

func CoreStopAll() {

	logging.Log(core.LL_SYS, "Stopping Default Timer Manager ... ")
	timer.GetDefaultTimerManager().Stop()
	logging.Log(core.LL_SYS, "Stopping Default Executor Pool ")
	concurrent.GetDefaultGoExecutorPool().Stop()
	logging.Log(core.LL_SYS, "Stopping Default Log Manager ")
	logging.GetLoggerManager().Stop()

	WaitAll()
	os.Exit(0)
}

func WaitAll() {
	waitOnce.Do(
		func() {
			timer.GetDefaultTimerManager().Wait()
			logging.LogFixedWidth(core.LL_SYS, 70, true, "[Stopped]", "Default Timer Manager...")

			concurrent.GetDefaultGoExecutorPool().Wait()
			logging.LogFixedWidth(core.LL_SYS, 70, true, "[Stopped]", "Default Executor Pool ...")

			logging.GetLoggerManager().Wait()
			logging.LogFixedWidth(core.LL_SYS, 70, true, "[Stopped]", "Default Log Manager ...")

		},
	)
}

func Initialize() {
	initOnce.Do(
		func() {
			fmt.Println("Zohar Core Initializing ... ")
			str, _ := os.Executable()
			fileName := filepath.Base(str)
			ext := filepath.Ext(fileName)
			fileName = strings.TrimSuffix(fileName, ext)
			fileName += ".json"
			str = filepath.Dir(str)
			fileName = filepath.Join(str, fileName)

			f := io.NeoFile(false, fileName, io.FILEFLAG_THREAD_SAFE)
			rc := f.Open(io.FILEOPEN_MODE_OPEN_OR_CREATE, io.FILEOPEM_PREM_READWRITE, 0755)
			if core.Err(rc) {
				panic("Intrinsic Initialization ... \t\t\t [Failed]")
			}

			sz := f.GetInfo().Size()
			if sz < 1 {
				rc = intrinsic.MakeDefaultIntrinsicConfig(f)
				if core.Err(rc) {
					panic("init: Make Default IntrinsicConfig Failed")
				}
			} else {
				rc = intrinsic.LoadConfig(f)
				if core.Err(rc) {
					panic("init: Load IntrinsicConfig Failed")
				}
			}
			fmt.Println(intrinsic.GetIntrinsicConfig().String())
			process.Initialize(intrinsic.GetIntrinsicConfig().CWD)

			installSignalHandler()

			logging.GetLoggerManager().Start()
			for k, v := range intrinsic.GetIntrinsicConfig().Logging {
				logger := logging.NeoLocalSyncTextLogger(k, &v)
				if logger == nil {
					fmt.Printf("[Failed: Logger Init (%s)]", k)
					panic("can't continued!")
				}
				logging.GetLoggerManager().Add(logger)
			}

			logging.LogRaw(core.LL_SYS, cmdline.GetArguments().String())
			logging.LogRaw(core.LL_SYS, concurrent.GetDefaultGoExecutorPool().String())

			logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Zohar Core Initializing ...")
			logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Starting Default Log Manager  ...")

			rc = concurrent.GetDefaultGoExecutorPool().Start()
			if core.Err(rc) {
				logging.LogFixedWidth(core.LL_SYS, 70, false, core.ErrStr(rc), "Starting Default Executor Pool ...")
				panic("Fatal Can not continue")
			} else {
				logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Starting Default Executor Pool ...")

			}

			rc = timer.GetDefaultTimerManager().Start()
			if core.Err(rc) {
				logging.LogFixedWidth(core.LL_SYS, 70, false, core.ErrStr(rc), "Starting Default Timer Manager ...")
				panic("Fatal Can not continue")
			} else {
				logging.LogFixedWidth(core.LL_SYS, 70, true, "", "Starting Default Timer Manager ...")

			}

			mp.GetDefaultObjectInvoker().RegisterClass("smh", server.GetHandlerRegistration())
		},
	)
}
