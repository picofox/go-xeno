package excel

import (
	"github.com/xuri/excelize/v2"
	_ "github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"xeno/zohar/core"
	"xeno/zohar/core/process"
)

type ExcelFile struct {
	_fileName        string
	_fileDir         string
	_fileFullPath    string
	_excelFileHandle *excelize.File
}

func (ego *ExcelFile) Open() int32 {
	var err error
	ego._excelFileHandle, err = excelize.OpenFile(ego._fileFullPath)
	if err != nil {
		return core.MkErr(core.EC_INDEX_OOB, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ExcelFile) Close() {
	if ego._excelFileHandle != nil {
		ego._excelFileHandle.Close()
	}
}

func (ego *ExcelFile) ReadAll() {

}

func NeoExcelFile(baseOnCWD bool, name string) *ExcelFile {
	fullPath := process.ComposePath(baseOnCWD, name, false)
	fullDir := filepath.Dir(fullPath)
	fileName := filepath.Base(fullPath)
	fileDir := filepath.Base(fullDir)
	os.MkdirAll(fullDir, 0755)

	f := ExcelFile{
		_fileName:        fileName,
		_fileDir:         fileDir,
		_fileFullPath:    fullPath,
		_excelFileHandle: nil,
	}
	return &f
}
