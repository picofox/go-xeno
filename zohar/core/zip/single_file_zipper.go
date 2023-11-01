package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"xeno/zohar/core"
)

type SingleFileZipper struct {
	SrcFile string
	ZipFile string
}

func (ego *SingleFileZipper) Zip() int32 {
	zipFile, err := os.Create(ego.ZipFile)
	if err != nil {
		return core.MkErr(core.EC_FILE_OPEN_FAILED, 1)
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	w, err := zipWriter.Create(filepath.Base(ego.SrcFile))
	if err != nil {
		return core.MkErr(core.EC_FILE_OPEN_FAILED, 2)
	}
	f, err := os.Open(ego.SrcFile)
	if err != nil {
		return core.MkErr(core.EC_FILE_OPEN_FAILED, 3)
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		return core.MkErr(core.EC_FILE_WRITE_FAILED, 3)
	}
	// 第四步，关闭 zip writer，将所有数据写入指向基础 zip 文件的数据流
	zipWriter.Close()

	return core.MkSuccess(0)
}

func NeoSingleFileZipper(src string, zipFile string) *SingleFileZipper {
	return &SingleFileZipper{
		SrcFile: src,
		ZipFile: zipFile,
	}
}
