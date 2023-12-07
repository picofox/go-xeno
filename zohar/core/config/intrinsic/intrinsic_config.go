package intrinsic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/io"
	"xeno/zohar/core/process"
)

type LoggerConfig struct {
	BaseOnCWD       bool     `json:"BaseOnCWD"`
	Dir             string   `json:"Dir"`
	BackupDir       string   `json:"BackupDir"`
	ZipFile         bool     `json:"ZipFile"`
	ToConsole       bool     `json:"ToConsole"`
	LineLimit       int64    `json:"LineLimit"`
	SizeLimit       int64    `json:"SizeLimit"`
	VolumeLimit     int32    `json:"VolumeLimit"`
	Type            int8     `json:"Type"`
	SplitMode       int8     `json:"SplitMode"`
	Depth           int8     `json:"SplitMode"`
	DefaultLevel    int8     `json:"DefaultLevel"`
	FileNamePattern string   `json:"FileNamePattern"`
	LinePattern     []string `json:"LinePattern"`
}

type IntrinsicConfig struct {
	CWD              string                  `json:"CWD"`
	CmdSpecRestrict  bool                    `json:"CmdSpecRestrict"`
	CmdParamSpec     string                  `json:"CmdParamSpec"`
	CmdTargetSpec    string                  `json:"CmdTargetSpec"`
	Logging          map[string]LoggerConfig `json:"Logging"`
	GoExecutorPool   GoExecutorPoolConfig    `json:"GoExecutorPool"`
	Poller           PollerConfig            `json:"Poller"`
	IntrinsicService IntrinsicServiceConfig  `json:"IntrinsicService"`
}

var sIntrinsicConfig *IntrinsicConfig = nil

func (ego *LoggerConfig) GetLogDirFullPath() string {
	return process.ComposePath(ego.BaseOnCWD, ego.Dir, false)
}
func (ego *LoggerConfig) GetBackupDirFullPath() string {
	return process.ComposePath(ego.BaseOnCWD, ego.BackupDir, false)
}

func (ego *LoggerConfig) GetBackupFileFullPath(name string) string {
	str := process.ComposePath(ego.BaseOnCWD, ego.BackupDir, false)
	str = filepath.Join(str, name)
	return str
}

func (ego *LoggerConfig) ParseString(fmt string, tm *time.Time) string {
	year, month, day := tm.Date()
	hour, min, sec := tm.Clock()

	fmt = strings.Replace(fmt, "%YYYY", strconv.Itoa(year), 1)
	fmt = strings.Replace(fmt, "%MM", strconv.Itoa(int(month)), 1)
	fmt = strings.Replace(fmt, "%DD", strconv.Itoa(day), 1)
	fmt = strings.Replace(fmt, "%hh", strconv.Itoa(hour), 1)
	fmt = strings.Replace(fmt, "%mm", strconv.Itoa(min), 1)
	fmt = strings.Replace(fmt, "%ss", strconv.Itoa(sec), 1)
	fmt = strings.Replace(fmt, "%mil", strconv.Itoa(int(tm.UnixMilli()%1000)), 1)
	fmt = strings.Replace(fmt, "%PID", strconv.Itoa(os.Getpid()), 1)
	fmt = strings.Replace(fmt, "%PID", strconv.FormatInt(process.GetCurrentGoRoutineId(), 10), 1)
	return fmt
}

func (ego *LoggerConfig) String() string {
	var ss strings.Builder
	ss.WriteString("       ")
	ss.WriteString("Type:(")
	ss.WriteString(strconv.Itoa(int(ego.Type)))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("Depth:(")
	ss.WriteString(strconv.Itoa(int(ego.Depth)))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("BaseOnCWD:(")
	ss.WriteString(strconv.FormatBool(ego.BaseOnCWD))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("Dir:(")
	ss.WriteString(ego.Dir)
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("BackupDir:(")
	ss.WriteString(ego.BackupDir)
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("ZipFile:(")
	ss.WriteString(strconv.FormatBool(ego.ZipFile))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("ToConsole:(")
	ss.WriteString(strconv.FormatBool(ego.ToConsole))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("LineLimit:(")
	ss.WriteString(strconv.Itoa(int(ego.LineLimit)))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("SizeLimit:(")
	ss.WriteString(strconv.Itoa(int(ego.SizeLimit)))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("VolumeLimit:(")
	ss.WriteString(strconv.Itoa(int(ego.VolumeLimit)))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("SplitMode:(")
	ss.WriteString(strconv.Itoa(int(ego.SplitMode)))
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("FileNamePattern:(")
	ss.WriteString(ego.FileNamePattern)
	ss.WriteString(")\n")
	ss.WriteString("       ")
	ss.WriteString("LinePattern:(")
	ss.WriteString(fmt.Sprint(ego.LinePattern))
	ss.WriteString(")\n")
	return ss.String()
}

func (ego *IntrinsicConfig) String() string {
	var ss strings.Builder
	ss.WriteString("+++++++++++++++++++++++++++++++ IntrinsicConfig +++++++++++++++++++++++++++++++\n")
	ss.WriteString("CWD:(")
	ss.WriteString(ego.CWD)
	ss.WriteString(")\n")
	ss.WriteString("CmdSpecRestrict:(")
	ss.WriteString(strconv.FormatBool(ego.CmdSpecRestrict))
	ss.WriteString(")\n")
	ss.WriteString("CmdSpec:(")
	ss.WriteString(ego.CmdParamSpec)
	ss.WriteString(")\n")
	ss.WriteString("CmdTargetSpec:(")
	ss.WriteString(ego.CmdTargetSpec)
	ss.WriteString(")\n")
	ss.WriteString("Logging:\n")
	for k, v := range ego.Logging {
		ss.WriteString("    ")
		ss.WriteString(k)
		ss.WriteString(":\n")
		ss.WriteString(v.String())
	}
	ss.WriteString(ego.GoExecutorPool.String())
	ss.WriteString(ego.Poller.String())
	ss.WriteString(ego.IntrinsicService.String())

	ss.WriteString("------------------------------- IntrinsicConfig -------------------------------\n")
	return ss.String()
}

func GetIntrinsicConfig() *IntrinsicConfig {
	return sIntrinsicConfig
}

func LoadConfig(file *io.File) int32 {
	defer file.Close()
	bs, rc := file.ReadAll()
	if core.Err(rc) {
		return rc
	}
	var IntrinConfig IntrinsicConfig
	err := json.Unmarshal(bs, &IntrinConfig)
	if err != nil {
		fmt.Println(err)
		return core.MkErr(core.EC_JSON_UNMARSHAL_FAILED, 1)
	}

	sIntrinsicConfig = &IntrinConfig
	return core.MkSuccess(0)
}

func MakeDefaultIntrinsicConfig(file *io.File) int32 {
	path, _ := os.Executable()
	exec := filepath.Base(path)
	exec = strings.Trim(exec, filepath.Ext(exec))
	fileNamePattern := exec + "-%YYYY%MM%DD"

	var IntrinConfig IntrinsicConfig
	IntrinConfig.CWD = ""
	var loggerConfig LoggerConfig = LoggerConfig{
		BaseOnCWD:       false,
		Dir:             "log",
		BackupDir:       "log/bak",
		ZipFile:         false,
		ToConsole:       true,
		LineLimit:       10000000,
		SizeLimit:       2000 * 1024 * 1024,
		VolumeLimit:     1000,
		Type:            0,
		SplitMode:       0,
		Depth:           8,
		DefaultLevel:    core.LL_DEBUG,
		FileNamePattern: fileNamePattern,
		LinePattern:     []string{"date", "time", "nano", "ts", "lv", "pid", "goid"},
	}

	IntrinConfig.Logging = make(map[string]LoggerConfig)
	IntrinConfig.Logging["default"] = loggerConfig
	IntrinConfig.GoExecutorPool.Name = "DFL-WP"
	IntrinConfig.GoExecutorPool.InitialCount = 2
	IntrinConfig.GoExecutorPool.MaxCount = 20
	IntrinConfig.GoExecutorPool.MinCount = 2
	IntrinConfig.GoExecutorPool.QueueSize = 1024
	IntrinConfig.GoExecutorPool.HighWaterMark = 512
	IntrinConfig.GoExecutorPool.LowWaterMark = 0
	IntrinConfig.GoExecutorPool.CheckInterval = 1000

	IntrinConfig.Poller.SubReactorCount = -1

	cronCfgDfl := CronServiceConfig{
		Offset: 0,
	}
	cronCfgUtc0 := CronServiceConfig{
		Offset: 0,
	}
	cronGrpCfg := CronServiceGroupConfig{
		Params: make(map[string]CronServiceConfig),
	}
	cronGrpCfg.Params["default"] = cronCfgDfl
	cronGrpCfg.Params["utc0"] = cronCfgUtc0

	fswCfg := FileSystemWatcherServiceConfig{
		Dirs: []string{"conf", "bin"},
	}
	fswGrpCfg := FileSystemWatcherServiceGroupConfig{
		Params: make(map[string]FileSystemWatcherServiceConfig),
	}
	fswGrpCfg.Params["default"] = fswCfg

	intrinSvcConfig := IntrinsicServiceConfig{
		Cron:              cronGrpCfg,
		FileSystemWatcher: fswGrpCfg,
	}
	IntrinConfig.IntrinsicService = intrinSvcConfig

	jsonBytes, err := json.Marshal(IntrinConfig)
	if err != nil {
		return core.MkErr(core.EC_JSON_MARSHAL_FAILED, 1)
	}

	fmt.Println(string(jsonBytes))

	sIntrinsicConfig = &IntrinConfig

	_, rc := file.WriteBytes(jsonBytes)

	return rc
}
