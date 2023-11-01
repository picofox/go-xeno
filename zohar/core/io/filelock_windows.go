package io

import (
	"golang.org/x/sys/windows"
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

const (
	reserved = 0
	allBytes = ^uint32(0)
)

func (ego *FileLock) LockShare() int32 {
	// Per https://golang.org/issue/19098, “Programs currently expect the Fd
	// method to return a handle that uses ordinary synchronous I/O.”
	// However, LockFileEx still requires an OVERLAPPED structure,
	// which contains the file offset of the beginning of the lock range.
	// We want to lock the entire file, so we leave the offset as zero.
	ol := new(windows.Overlapped)
	err := windows.LockFileEx(windows.Handle(ego._fileHandler.Fd()), uint32(0), reserved, allBytes, allBytes, ol)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) TryLockShare() int32 {
	ol := new(windows.Overlapped)
	err := windows.LockFileEx(windows.Handle(ego._fileHandler.Fd()), uint32(0)|windows.LOCKFILE_FAIL_IMMEDIATELY, reserved, allBytes, allBytes, ol)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) LockExclusive() int32 {
	ol := new(windows.Overlapped)
	err := windows.LockFileEx(windows.Handle(ego._fileHandler.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, reserved, allBytes, allBytes, ol)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) TryLockExclusive() int32 {
	ol := new(windows.Overlapped)
	err := windows.LockFileEx(windows.Handle(ego._fileHandler.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK|windows.LOCKFILE_FAIL_IMMEDIATELY, reserved, allBytes, allBytes, ol)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *FileLock) Unlock() int32 {
	ol := new(windows.Overlapped)
	err := windows.UnlockFileEx(windows.Handle(ego._fileHandler.Fd()), reserved, allBytes, allBytes, ol)
	if err != nil {
		return core.MkErr(core.EC_TRY_LOCK_FILE_FAILED, 1)
	}
	return core.MkSuccess(0)
}
