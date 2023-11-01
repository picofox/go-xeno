package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/config"
	"xeno/zohar/core/io"
	"xeno/zohar/core/process"
)

type LocalSyncTextLogger struct {
	_logFile    *io.TextFile
	_config     *config.LoggerConfig
	_lineBuffer strings.Builder
	_lock       sync.RWMutex

	_level          int
	_filename       string
	_name           string
	_header         int32
	_volumeNo       int32
	_lineCount      int64
	_size           int64
	_nextShiftEpoch int64
}

func (ego *LocalSyncTextLogger) SetLevel(lv int) {
	//TODO implement me
	ego._level = lv
}

func (ego *LocalSyncTextLogger) LockShared() {
	ego._lock.RLock()
}

func (ego *LocalSyncTextLogger) LockExclusive() {
	ego._lock.Lock()
}

func (ego *LocalSyncTextLogger) UnLockShared() {
	ego._lock.RUnlock()
}

func (ego *LocalSyncTextLogger) UnLockExclusive() {
	ego._lock.Unlock()
}

func (ego *LocalSyncTextLogger) CalcFileName(tm *time.Time, checkDirExist bool) {
	logFileDir := ego._config.GetLogDirFullPath()
	if checkDirExist {
		os.MkdirAll(logFileDir, 0777)
	}
	if len(ego._config.BackupDir) > 0 {
		logBackupDir := ego._config.GetBackupDirFullPath()
		os.MkdirAll(logBackupDir, 0777)
	}

	fileName := ego._config.ParseString(ego._config.FileNamePattern, tm)
	if ego._config.VolumeLimit > 0 {
		fileName += fmt.Sprintf(".%d", ego._volumeNo)
	}
	fileName += ".log"
	fileName = filepath.Join(logFileDir, fileName)

	ego._filename = fileName
}

func NeoLocalSyncTextLogger(name string, cfg *config.LoggerConfig) *LocalSyncTextLogger {

	var elemes int32 = 0
	for i := 0; i < len(cfg.LinePattern); i++ {
		if strings.ToLower(cfg.LinePattern[i]) == "date" {
			elemes |= LINE_HEADER_ELEM_DATE
		} else if strings.ToLower(cfg.LinePattern[i]) == "time" {
			elemes |= LINE_HEADER_ELEM_TIME
		} else if strings.ToLower(cfg.LinePattern[i]) == "nano" {
			elemes |= LINE_HEADER_ELEM_NANO
		} else if strings.ToLower(cfg.LinePattern[i]) == "ts" {
			elemes |= LINE_HEADER_ELEM_TS
		} else if strings.ToLower(cfg.LinePattern[i]) == "lv" {
			elemes |= LINE_HEADER_ELEM_LV
		} else if strings.ToLower(cfg.LinePattern[i]) == "pid" {
			elemes |= LINE_HEADER_ELEM_PID
		} else if strings.ToLower(cfg.LinePattern[i]) == "goid" {
			elemes |= LINE_HEADER_ELEM_GOID
		}
	}

	tm := time.Now()
	var nextEP int64 = 0
	if cfg.SplitMode == core.L_SPLIT_BY_DAY {
		nextEP = chrono.GetDayBeginMilliStampByTMOffset(&tm, 1)
	} else if cfg.SplitMode == core.L_SPLIT_BY_HOUR {
		nextEP = chrono.GetHourBeginMilliStampByTMOffset(&tm, 1)
	}

	l := LocalSyncTextLogger{

		_level:          core.LL_DEBUG,
		_name:           name,
		_header:         elemes,
		_logFile:        nil,
		_config:         cfg,
		_volumeNo:       0,
		_lineCount:      0,
		_size:           0,
		_nextShiftEpoch: nextEP,
	}

	l.CalcFileName(&tm, true)
	f := io.NeoTextFile(false, l._filename, 0, false, "\n")
	rc := f.Open(io.FILEOPEN_MODE_OPEN_OR_CREATE, io.FILEOPEM_PREM_READ_APPEND, 0755)
	if core.Err(rc) {
		return nil
	}
	l._logFile = f
	return &l
}

func (ego *LocalSyncTextLogger) writeHeader(tm *time.Time, lv int) {
	ts := tm.UnixMilli()
	if ego._header&LINE_HEADER_ELEM_DATE != 0 {
		ego._lineBuffer.WriteString(fmt.Sprintf("%04d-%02d-%02d\t", tm.Year(), tm.Month(), tm.Day()))
	}
	if ego._header&LINE_HEADER_ELEM_TIME != 0 {
		if ego._header&LINE_HEADER_ELEM_NANO != 0 {
			ego._lineBuffer.WriteString(fmt.Sprintf("%02d:%02d:%02d.%06d\t", tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond()))
		} else {
			ego._lineBuffer.WriteString(fmt.Sprintf("%02d:%02d:%02d\t", tm.Hour(), tm.Minute(), tm.Second()))
		}
	}
	if ego._header&LINE_HEADER_ELEM_TS != 0 {
		ego._lineBuffer.WriteString(strconv.FormatInt(ts, 10))
		ego._lineBuffer.WriteString("\t")
	}

	if ego._header&LINE_HEADER_ELEM_LV != 0 {
		ego._lineBuffer.WriteString(GetLogLevelName(lv))
		ego._lineBuffer.WriteString("\t")
	}

	if ego._header&LINE_HEADER_ELEM_PID != 0 {
		ego._lineBuffer.WriteString(strconv.Itoa(os.Getpid()))
		ego._lineBuffer.WriteString("\t")
	}

	if ego._header&LINE_HEADER_ELEM_GOID != 0 {
		ego._lineBuffer.WriteString(strconv.FormatInt(process.GetCurrentGoRoutineId(), 10))
		ego._lineBuffer.WriteString("\t")
	}
}

func (ego *LocalSyncTextLogger) Name() string {
	return ego._name
}

func (ego *LocalSyncTextLogger) checkNeoFile(tm *time.Time) {
	oldFileName := ego._filename
	currentMIL := tm.UnixMilli()
	changeDay := false
	changeVol := false
	if currentMIL >= ego._nextShiftEpoch {
		changeDay = true

	} else {
		if ego._config.SizeLimit > 0 {
			if ego._size > ego._config.SizeLimit {
				changeVol = true
			}
		}
		if ego._config.LineLimit > 0 {
			if ego._lineCount > ego._config.LineLimit {
				changeVol = true
			}
		}
	}

	if !changeVol && !changeDay {
		return
	}

	ego._size = 0
	ego._lineCount = 0
	if changeDay {
		if ego._config.SplitMode == core.L_SPLIT_BY_DAY {
			ego._nextShiftEpoch = chrono.GetDayBeginMilliStampByTMOffset(tm, 1)
		} else if ego._config.SplitMode == core.L_SPLIT_BY_HOUR {
			ego._nextShiftEpoch = chrono.GetHourBeginMilliStampByTMOffset(tm, 1)
		}
		ego._volumeNo = 0
	} else if changeVol {
		ego._volumeNo++
	}

	ego._logFile.Close()
	ego.CalcFileName(tm, false)

	ego._logFile = io.NeoTextFile(false, ego._filename, 0, false, "\n")
	rc := ego._logFile.Open(io.FILEOPEN_MODE_OPEN_OR_CREATE, io.FILEOPEM_PREM_READ_APPEND, 0755)
	if core.Err(rc) {
		fmt.Printf("Reopen File (%s) Failed\n", ego._filename)
	}

	if len(ego._config.BackupDir) > 0 {
		GetLoggerManager().BackUp(oldFileName, ego._config.GetBackupDirFullPath(), ego._config.ZipFile)
	}

}

func (ego *LocalSyncTextLogger) Log(lv int, format string, arg ...any) {
	if lv > ego._level {
		return
	}

	tm := time.Now()

	ego.LockExclusive()
	defer ego.UnLockExclusive()
	ego._lineBuffer.Reset()

	ego.checkNeoFile(&tm)

	ego.writeHeader(&tm, lv)
	ego._lineBuffer.WriteString(fmt.Sprintf(format, arg...))
	s := ego._lineBuffer.String()
	ego._logFile.WriteLineBared(s)
	if ego._config.ToConsole {
		fmt.Println(s)
	}
	if ego._config.SizeLimit > 0 {
		ego._size = ego._size + int64(len(s))
	}
	if ego._config.LineLimit > 0 {
		ego._lineCount++
	}
}

func (ego *LocalSyncTextLogger) LogRaw(lv int, format string, arg ...any) {
	if lv > ego._level {
		return
	}

	tm := time.Now()

	ego.LockExclusive()
	defer ego.UnLockExclusive()
	ego._lineBuffer.Reset()

	ego.checkNeoFile(&tm)
	ego._lineBuffer.WriteString(fmt.Sprintf(format, arg...))
	s := ego._lineBuffer.String()
	ego._logFile.WriteLineBared(s)
	if ego._config.ToConsole {
		fmt.Println(s)
	}
	if ego._config.SizeLimit > 0 {
		ego._size = ego._size + int64(len(s))
	}

}
