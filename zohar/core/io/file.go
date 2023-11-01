package io

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"unicode/utf8"
	"xeno/zohar/core"
	"xeno/zohar/core/process"
)

const (
	FILEFLAG_THREAD_SAFE  = 0x1
	FILEFLAG_PROCESS_SAFE = 0x2
)

type File struct {
	_flags        int32
	_fileName     string
	_fileDir      string
	_fileFullPath string
	_fileHandler  *os.File
	_pLock        *FileLock
	_tLock        *sync.RWMutex
}

func NeoFile(baseOnCWD bool, name string, flags int32) *File {
	fullPath := process.ComposePath(baseOnCWD, name, false)
	fullDir := filepath.Dir(fullPath)
	fileName := filepath.Base(fullPath)
	fileDir := filepath.Base(fullDir)
	os.MkdirAll(fullDir, 0755)
	f := File{
		_flags:        flags,
		_fileName:     fileName,
		_fileDir:      fileDir,
		_fileFullPath: fullPath,
		_fileHandler:  nil,
		_pLock:        nil,
		_tLock:        nil,
	}
	return &f
}

const (
	FILEOPEN_MODE_OPEN_EXIST     = 0
	FILEOPEN_MODE_OPEN_TRUNC     = 1
	FILEOPEN_MODE_OPEN_OR_CREATE = 2
	FILEOPEN_MODE_CREATE         = 3
	FILEOPEN_MODE_CREATE_DEL     = 4
)

const (
	FILEOPEN_PERM_READ        = 0
	FILEOPEM_PREM_WRITE       = 1
	FILEOPEM_PREM_READWRITE   = 2
	FILEOPEM_PREM_APPEND      = 3
	FILEOPEM_PREM_READ_APPEND = 4
)

var sFileOpenParams [5][5]int = [5][5]int{
	{os.O_RDONLY, os.O_WRONLY, os.O_RDWR, os.O_WRONLY | os.O_APPEND, os.O_RDWR | os.O_APPEND},
	{os.O_WRONLY | os.O_TRUNC, os.O_RDWR | os.O_TRUNC, os.O_WRONLY | os.O_APPEND | os.O_TRUNC, os.O_RDWR | os.O_APPEND | os.O_TRUNC},
	{0, os.O_WRONLY | os.O_CREATE, os.O_RDWR | os.O_CREATE, os.O_WRONLY | os.O_APPEND | os.O_CREATE, os.O_RDWR | os.O_APPEND | os.O_CREATE},
	{0, os.O_WRONLY | os.O_CREATE | os.O_EXCL, os.O_RDWR | os.O_CREATE | os.O_EXCL, os.O_WRONLY | os.O_APPEND | os.O_CREATE | os.O_EXCL, os.O_RDWR | os.O_APPEND | os.O_CREATE | os.O_EXCL},
	{0, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, os.O_RDWR | os.O_TRUNC | os.O_CREATE, os.O_WRONLY | os.O_APPEND | os.O_TRUNC | os.O_CREATE, os.O_RDWR | os.O_APPEND | os.O_TRUNC | os.O_CREATE},
}

func (ego *File) FullPathName() string {
	return ego._fileFullPath
}

func (ego *File) Open(openMode int, permMode int, perm fs.FileMode) int32 {
	if ego._fileHandler != nil {
		ego._fileHandler.Close()
	}
	ego._tLock = nil
	ego._pLock = nil
	mode := sFileOpenParams[openMode][permMode]
	var err error = nil
	ego._fileHandler, err = os.OpenFile(ego._fileFullPath, mode, perm)
	if err != nil {
		return core.MkErr(core.EC_FILE_OPEN_FAILED, 1)
	}

	if ego._flags&FILEFLAG_THREAD_SAFE != 0 {
		ego._tLock = &sync.RWMutex{}
	} else {
		ego._tLock = nil
	}

	if ego._flags&FILEFLAG_PROCESS_SAFE != 0 {
		ego._pLock = CreateFileLock(ego._fileHandler)
	} else {
		ego._pLock = nil
	}

	return core.MkSuccess(0)
}

func (ego *File) Close() {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()
	ego.CloseBared()
}

func (ego *File) CloseBared() {
	if ego._fileHandler != nil {
		ego._fileHandler.Close()
		ego._fileHandler = nil
	}
	ego._tLock = nil
	ego._pLock = nil
}

func (ego *File) Seek(offset int64, whence int) {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()

	if ego._fileHandler != nil {
		ego._fileHandler.Seek(offset, whence)
	}
}

func (ego *File) Flush() {
	if ego._fileHandler != nil {
		ego._fileHandler.Sync()
	}
}
func (ego *File) WriteString(str string) (int, int32) {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()
	return ego.WriteStringBared(str)
}

func (ego *File) WriteStringBared(str string) (int, int32) {
	return ego.WriteBytesBared([]byte(str))
}

func (ego *File) WriteStringPartial(str string, idx int, wlen int) (int, int32) {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()
	return ego.WriteStringPartialBared(str, idx, wlen)
}

func (ego *File) WriteStringPartialBared(str string, idx int, wlen int) (int, int32) {
	rc := core.EC_NOOP
	if ego._fileHandler != nil {
		rs := []rune(str)[idx : idx+wlen]
		var buf = make([]byte, utf8.UTFMax)
		var bsWrite int = 0
		for _, r := range rs {
			nBs := utf8.EncodeRune(buf, r)
			_, rc = ego.WriteBared(buf, 0, nBs)
			if core.Err(rc) {
				return bsWrite, core.MkErr(core.EC_FILE_WRITE_FAILED, 1)
			}
			bsWrite += nBs
		}
		return bsWrite, core.MkSuccess(0)
	}
	return 0, rc
}

func (ego *File) Write(bs []byte, bufPos int, maxLen int) (int, int32) {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()
	return ego.WriteBared(bs, bufPos, maxLen)
}

func (ego *File) WriteBared(bs []byte, bufPos int, maxLen int) (int, int32) {
	var writeLen int = maxLen
	if writeLen > len(bs)-bufPos {
		writeLen = len(bs) - bufPos
	}

	if ego._fileHandler != nil {
		nbs, err := ego._fileHandler.Write(bs[bufPos : bufPos+writeLen])
		if err != nil {
			return nbs, core.MkErr(core.EC_FILE_WRITE_FAILED, 1)
		}
		return nbs, core.MkSuccess(0)
	}

	return 0, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) WriteBytesBared(bs []byte) (int, int32) {
	if ego._fileHandler != nil {
		nbs, err := ego._fileHandler.Write(bs)
		if err != nil {
			return nbs, core.MkErr(core.EC_FILE_WRITE_FAILED, 1)
		}
		return nbs, core.MkSuccess(0)
	}

	return 0, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) WriteBytes(bs []byte) (int, int32) {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()
	return ego.WriteBytesBared(bs)
}

func (ego *File) ReadByte() (bool, byte, int32) {
	ego.PLockExclusive()
	ego.TLockExclusive()
	defer ego.PUnlock()
	defer ego.TUnLockExclusive()
	return ego.ReadByte()
}

func (ego *File) ReadByteBared() (byte, int32) {
	if ego._fileHandler != nil {
		rdata := make([]byte, 1)
		_, err := io.ReadFull(ego._fileHandler, rdata)
		if err != nil {
			if err == io.EOF {
				return 0, core.MkErr(core.EC_EOF, 0)
			} else {
				return 0, core.MkErr(core.EC_FILE_READ_FAILED, 2)
			}
		} else {
			return rdata[0], core.MkSuccess(0)
		}
	}
	return 0, core.MkErr(core.EC_NOOP, 1)
}

func (ego *File) ReadN(n int) ([]byte, int32) {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()
	return ego.ReadNBared(n)
}

func (ego *File) ReadNBared(n int) ([]byte, int32) {
	byteSlice := make([]byte, n)

	if ego._fileHandler != nil {
		numBytesRead, err := io.ReadFull(ego._fileHandler, byteSlice)
		if numBytesRead < n {
			if err == io.EOF {
				return byteSlice, core.MkErr(core.EC_EOF, 0)
			} else {
				return nil, core.MkErr(core.EC_FILE_READ_FAILED, 1)
			}
		}
		return byteSlice, core.MkSuccess(0)
	}
	return nil, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) ReadNAt(rdata []byte, off int64, whence int) (int, int32) {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()
	return ego.ReadNAtBared(rdata, off, whence)
}

func (ego *File) ReadNAtBared(rdata []byte, off int64, whence int) (int, int32) {
	if ego._fileHandler != nil {
		ego._fileHandler.Seek(off, whence)
		numBytesRead, err := io.ReadFull(ego._fileHandler, rdata)
		if numBytesRead < len(rdata) {
			if err == io.EOF {
				return numBytesRead, core.MkErr(core.EC_EOF, 0)
			} else {
				return numBytesRead, core.MkErr(core.EC_FILE_READ_FAILED, 1)
			}
		}
		return numBytesRead, core.MkSuccess(0)
	}
	return 0, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) ReadAt(rdata []byte, off int64) (int, int32) {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()
	return ego.ReadAtBared(rdata, off)
}

func (ego *File) ReadAtBared(rdata []byte, off int64) (int, int32) {

	if ego._fileHandler != nil {
		nRead, err := ego._fileHandler.ReadAt(rdata, off)
		if err != nil {
			if err == io.EOF {
				return nRead, core.MkErr(core.EC_EOF, 0)
			} else {
				return nRead, core.MkErr(core.EC_FILE_READ_FAILED, 1)
			}
		}

		return nRead, core.MkSuccess(0)
	}
	return 0, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) Read(rdata []byte) (int, int32) {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()
	return ego.ReadBared(rdata)
}

func (ego *File) ReadBared(rdata []byte) (int, int32) {

	if ego._fileHandler != nil {
		nRead, err := ego._fileHandler.Read(rdata)
		if err != nil {
			if err == io.EOF {
				return nRead, core.MkErr(core.EC_EOF, 0)
			} else {
				return nRead, core.MkErr(core.EC_FILE_READ_FAILED, 1)
			}
		}
		return nRead, core.MkSuccess(0)
	}
	return 0, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) GetInfo() os.FileInfo {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()
	return ego.GetInfoBared()
}

func (ego *File) GetInfoBared() os.FileInfo {
	if ego._fileHandler != nil {
		fi, err := ego._fileHandler.Stat()
		if err != nil {
			return nil
		}
		return fi
	}
	return nil
}

func (ego *File) ReadAll() ([]byte, int32) {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()
	return ego.ReadAllBared()
}

func (ego *File) ReadAllBared() ([]byte, int32) {
	ego.PLockShare()
	ego.TLockShare()
	defer ego.PUnlock()
	defer ego.TUnLockShare()

	if ego._fileHandler != nil {
		fi, err := ego._fileHandler.Stat()
		if err != nil {
			return nil, core.MkErr(core.EC_FILE_STAT_FAILED, 1)
		}
		if fi.Size() < 1 {
			return []byte{}, core.MkSuccess(0)
		}
		sz := fi.Size()
		byteSlice := make([]byte, sz)
		numBytesRead, err := io.ReadFull(ego._fileHandler, byteSlice)
		if int64(numBytesRead) < sz {
			if err == io.EOF {
				return byteSlice, core.MkErr(core.EC_EOF, 0)
			} else {
				return byteSlice, core.MkErr(core.EC_FILE_READ_FAILED, 1)
			}
		}
		return byteSlice, core.MkSuccess(0)
	}
	return []byte{}, core.MkErr(core.EC_NOOP, 0)
}

func (ego *File) TLockShare() {
	if ego._tLock != nil {
		ego._tLock.RLock()
	}
}
func (ego *File) TLockExclusive() {
	if ego._tLock != nil {
		ego._tLock.Lock()
	}
}
func (ego *File) TTryLockShare() bool {
	if ego._tLock != nil {
		return ego._tLock.TryRLock()
	}
	return true
}
func (ego *File) TTryLockExclusive() bool {
	if ego._tLock != nil {
		return ego._tLock.TryLock()
	}
	return true
}
func (ego *File) TUnLockShare() {
	if ego._tLock != nil {
		ego._tLock.RUnlock()
	}
}
func (ego *File) TUnLockExclusive() {
	if ego._tLock != nil {
		ego._tLock.Unlock()
	}
}

func (ego *File) PLockShare() int32 {
	if ego._pLock != nil {
		return ego._pLock.LockShare()
	}
	return core.MkSuccess(0)
}
func (ego *File) PLockExclusive() int32 {
	if ego._pLock != nil {
		return ego._pLock.LockExclusive()
	}
	return core.MkSuccess(0)
}
func (ego *File) PTryLockShare() int32 {
	if ego._pLock != nil {
		return ego._pLock.TryLockShare()
	}
	return core.MkSuccess(0)
}
func (ego *File) PTryLockExclusive() int32 {
	if ego._pLock != nil {
		return ego._pLock.TryLockExclusive()
	}
	return core.MkSuccess(0)
}
func (ego *File) PUnlock() int32 {
	if ego._pLock != nil {
		return ego._pLock.Unlock()
	}
	return core.MkSuccess(0)
}
