package io

import (
	"golang.org/x/sys/unix"
	"os"
	"xeno/zohar/core"
)

type FileLock struct {
	_fileHandler *os.File
}

func CreateFileLock(fh *os.File) *FileLock {
	return &FileLock{
		_fileHandler: fh,
	}
}

func (ego *FileLock) LockShare() int32 {
	unix.Flock(int(ego._fileHandler.Fd()), unix.LOCK_SH)
	return core.MkSuccess(0)
}

func (ego *FileLock) TryLockShare() int32 {
	unix.Flock(int(ego._fileHandler.Fd()), unix.LOCK_SH|unix.LOCK_NB)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) LockExclusive() int32 {
	err := unix.Flock(int(ego._fileHandler.Fd()), unix.LOCK_EX)
	if err != nil {
		return core.MkErr(core.EC_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) TryLockExclusive() int32 {
	err := unix.Flock(int(ego._fileHandler.Fd()), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) Unlock() int32 {
	unix.Flock(int(ego._fileHandler.Fd()), unix.LOCK_UN)
	return core.MkSuccess(0)
}
