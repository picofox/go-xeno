package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
	"xeno/zohar/core/finalization"
	"xeno/zohar/core/zip"
)

type LoggerManager struct {
	loggers       map[string]ILogger
	defaultLogger ILogger
	channel       chan cms.ICMS
	stopped       bool
}

func (ego *LoggerManager) Start() {
	go ego.MaintenanceRoutine()
}

func (ego *LoggerManager) Stop() {
	ego.channel <- cms.NeoFinalize()
	for {
		if ego.stopped {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	close(ego.channel)
}

func (ego *LoggerManager) BackUp(filePath string, backupPath string, zip bool) {
	m := cms.NeoCMSLogBackUp(filePath, backupPath, zip)
	ego.channel <- m
}

func (ego *LoggerManager) MaintenanceRoutine() {
	finalization.GetGlobalFinalizer().Register("LoggerManager", ego, finalizeLoggerManagerMaintenanceRoutine)
	for {
		m := <-ego.channel
		if m.Id() == cms.CMSID_LOG_BACKUP {
			bakFile := filepath.Base(m.(*cms.LogBackUp).AbsFilePath)
			bakFile = filepath.Join(m.(*cms.LogBackUp).AbsBackupDirPath, bakFile)
			err := os.Rename(m.(*cms.LogBackUp).AbsFilePath, bakFile)
			if err != nil {
				fmt.Println("rename error " + bakFile)
			}

			if m.(*cms.LogBackUp).ZipFile {
				zipper := zip.NeoSingleFileZipper(bakFile, bakFile+".zip")
				rc := zipper.Zip()
				if core.Err(rc) {
					fmt.Printf("Zip file <%s> Failed!", bakFile)
				} else {
					os.Remove(bakFile)
				}

			}

		} else if m.Id() == cms.CMSID_FINALIZE {
			ego.stopped = true
			runtime.Goexit()
		}
		sz := len(ego.channel)
		fmt.Println("channel " + strconv.Itoa(sz))
	}
}

func (ego *LoggerManager) Add(logger ILogger) {
	ego.loggers[logger.Name()] = logger
	if logger.Name() == "default" {
		ego.defaultLogger = logger
	}
}

func (ego *LoggerManager) Default() ILogger {
	return ego.defaultLogger
}

func Log(lv int, format string, arg ...any) {
	GetLoggerManager().Default().Log(lv, format, arg...)
}

func LogRaw(lv int, format string, arg ...any) {
	GetLoggerManager().Default().LogRaw(lv, format, arg...)
}

func (ego *LoggerManager) Get(name string) ILogger {
	v2, ok := ego.loggers[name]
	if !ok {
		return nil
	}
	return v2
}

var lmInstance LoggerManager
var once sync.Once

func GetLoggerManager() *LoggerManager {
	once.Do(func() {
		lmInstance = LoggerManager{
			channel:       make(chan cms.ICMS, 128),
			loggers:       make(map[string]ILogger),
			stopped:       false,
			defaultLogger: nil,
		}
	})
	return &lmInstance
}

func finalizeLoggerManagerMaintenanceRoutine(subject any) {
	subject.(*LoggerManager).Stop()

}
