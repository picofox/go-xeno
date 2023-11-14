package initialization

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/cmdline"
	"xeno/zohar/core/concurrent"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/io"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/process"
)

func init() {
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

	logging.GetLoggerManager().Start()
	for k, v := range intrinsic.GetIntrinsicConfig().Logging {
		logger := logging.NeoLocalSyncTextLogger(k, &v)
		if logger == nil {
			fmt.Printf("[Failed: Logger Init (%s)]", k)
			panic("can't continued!")
		}
		logging.GetLoggerManager().Add(logger)
	}

	fmt.Println(cmdline.GetArguments().String())

	fmt.Println(concurrent.GetDefaultGoExecutorPool().String())

}
