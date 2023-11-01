package fs

import (
	"os"
	"xeno/zohar/core"
)

func EnsureDir(dir_path string, perm os.FileMode, force bool) int32 {
	s, err := os.Stat(dir_path)
	if err != nil {
		err = os.MkdirAll(dir_path, perm)
		if err != nil {
			return core.MkErr(core.EC_CREATE_DIR_FAILED, 1)
		}
		return core.MkSuccess(0)
	} else {
		if s == nil {
			err = os.MkdirAll(dir_path, perm)
			if err != nil {
				return core.MkErr(core.EC_CREATE_DIR_FAILED, 3)
			}
			return core.MkSuccess(3)
		}
		if s.IsDir() {
			return core.MkSuccess(1)
		} else {
			if !force {
				return core.EC_DIR_ALREADY_EXIST
			}
			err = os.Remove(dir_path)
			if err != nil {
				return core.EC_DELETE_DIR_FAILED
			}
			err = os.MkdirAll(dir_path, perm)
			if err != nil {
				return core.MkErr(core.EC_CREATE_DIR_FAILED, 2)
			}
			return core.MkSuccess(2)
		}
	}
}
