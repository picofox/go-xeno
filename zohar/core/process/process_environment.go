package process

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
	"xeno/zohar/core"
)

const (
	DIR_BASE = iota
	DIR_BIN
	DIR_LOG
	DIR_CONF
	DIR_TMP

	DIR_NA
)

var g_subdir_str = [DIR_NA]string{
	"",
	"bin",
	"log",
	"conf",
	"tmp",
}

var __g_exec_file_path string = ""
var __g_cwd = ""
var __g_main_dir_path = ""
var __g_prog_base = ""
var __g_prog_log_dir_path = ""
var __g_prog_tmp_dir_path = ""
var __g_prog_conf_dir_path = ""
var __g_prog_bin_dir_path = ""
var __g_base_unixtime time.Time

//__g_main_dir_path = os.path.dirname(__s_main_file_path)
//__s_main_file_name = os.path.basename(__s_main_file_path)
//__s_main_file_base_name = os.path.splitext(__s_main_file_name)[0]
//__s_main_file_ext_name = os.path.splitext(__s_main_file_name)[1]
//__s_log_dir = os.path.join(__s_main_dir_path, 'log')

//__s_run_mode = S_RUN_MISC

func RelToAbs(dirType int, sub string, filePath string) string {
	if filepath.IsAbs(filePath) {
		return filePath
	}
	if dirType == DIR_BASE {
		return filepath.Join(__g_prog_base, filePath)
	} else if dirType >= DIR_NA {
		return filepath.Join(__g_prog_base, sub, filePath)
	}
	return filepath.Join(__g_prog_base, g_subdir_str[dirType], filePath)
}

func ProgramBase() string {
	return __g_prog_base
}

func ProgramConfFile(baseName string, ext string) string {
	if strings.HasSuffix(baseName, ".") {
		baseName = baseName + "conf" + ext
		fileBaseName := filepath.Join(__g_prog_conf_dir_path, baseName)
		return fileBaseName
	}

	baseName = baseName + ext
	fileBaseName := filepath.Join(__g_prog_conf_dir_path, baseName)
	return fileBaseName

}

func ExecutablePath() string {
	if __g_exec_file_path == "" {
		__g_exec_file_path, _ = os.Executable()
		__g_exec_file_path = strings.TrimRight(__g_exec_file_path, "\n")
	}
	return __g_exec_file_path
}

func CWD() string {
	if __g_cwd == "" {
		__g_cwd, _ = os.Getwd()
	}
	return __g_cwd
}

func MainDirPath() string {
	if __g_main_dir_path == "" {
		dir := getCurrentAbPathByExecutable()
		tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
		if strings.Contains(dir, tmpDir) {
			__g_main_dir_path = getCurrentAbPathByCaller()
		}
		__g_main_dir_path = dir
	}
	return __g_main_dir_path

}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func GetTimestampBase() *time.Time {
	return &__g_base_unixtime
}

func Initialize(cwdPath string) int32 {
	__g_base_unixtime = time.Now()
	ExecutablePath()
	MainDirPath()
	__g_main_dir_path = MainDirPath()
	s2 := filepath.Join(__g_main_dir_path, "../")
	__g_prog_base, _ = filepath.Abs(s2)

	__g_prog_log_dir_path = filepath.Join(__g_prog_base, "log")
	__g_prog_tmp_dir_path = filepath.Join(__g_prog_base, "tmp")
	__g_prog_conf_dir_path = filepath.Join(__g_prog_base, "conf")
	__g_prog_bin_dir_path = filepath.Join(__g_prog_base, "bin")

	if core.Err(EnsureDir(__g_prog_log_dir_path, 0755, false)) {
		return core.MkErr(core.EC_ENSURE_DIR_FAILED, 1)
	}
	if core.Err(EnsureDir(__g_prog_tmp_dir_path, 0755, false)) {
		return core.MkErr(core.EC_ENSURE_DIR_FAILED, 2)
	}
	if core.Err(EnsureDir(__g_prog_conf_dir_path, 0755, false)) {
		return core.MkErr(core.EC_ENSURE_DIR_FAILED, 1)
	}

	if cwdPath != "" {
		str, err := filepath.Abs(cwdPath)
		if err != nil {
			return core.MkErr(core.EC_TO_ABS_PATH_FAILED, 1)
		}
		if os.Chdir(str) != nil {
			return core.MkErr(core.EC_SET_CWD_FAILED, 1)
		}
	} else {
		if os.Chdir(__g_prog_base) != nil {
			return core.MkErr(core.EC_SET_CWD_FAILED, 1)
		}
	}

	CWD()

	return core.MkSuccess(0)
}

func GetCurrentGoRoutineId() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	id, _ := strconv.Atoi(idField)
	return int64(id)
}

func ComposePath(basedOnCWD bool, name string, createDir bool) string {
	if filepath.IsAbs(name) {
		return name
	} else {
		if basedOnCWD {
			name = filepath.Join(CWD(), name)
		} else {
			name = filepath.Join(ProgramBase(), name)
		}
	}

	if createDir {
		os.MkdirAll(name, 0777)
	}
	return name
}

func ProcessInfoString() string {
	var ss strings.Builder
	ss.WriteString("Exe Path: <")
	ss.WriteString(__g_exec_file_path)
	ss.WriteString(">\n")

	ss.WriteString("Main DIR Path: <")
	ss.WriteString(__g_main_dir_path)
	ss.WriteString(">\n")

	ss.WriteString("Program Base: <")
	ss.WriteString(__g_prog_base)
	ss.WriteString(">\n")

	ss.WriteString("Log Dir: <")
	ss.WriteString(__g_prog_log_dir_path)
	ss.WriteString(">\n")

	ss.WriteString("Tmp Dir: <")
	ss.WriteString(__g_prog_tmp_dir_path)
	ss.WriteString(">\n")

	ss.WriteString("Bin Dir: <")
	ss.WriteString(__g_prog_bin_dir_path)
	ss.WriteString(">\n")

	ss.WriteString("Conf Dir: <")
	ss.WriteString(__g_prog_conf_dir_path)
	ss.WriteString(">\n")

	ss.WriteString("CWD: <")
	ss.WriteString(__g_cwd)
	ss.WriteString(">\n")

	return ss.String()
}
