package io

import (
	"io/fs"
	"os"
	"xeno/zohar/core"
	"xeno/zohar/core/xplatform"
)

type TextFile struct {
	_file      *File
	_convertNL bool
	_nlBs      []byte
	_nlLen     int32
	_nlBSDfl   []byte
}

const STR_DEFAULT_NEWLINE = "\n"

func NeoTextFile(baseOnCWD bool, name string, flags int32, convertNL bool, customNL string) *TextFile {
	file := NeoFile(baseOnCWD, name, flags)
	var nlBs []byte
	var nlLen int32
	if len(customNL) < 1 {
		nlBs = []byte(xplatform.XPF_DEF_LineBreak)
		nlLen = int32(len(nlBs))
	} else {
		nlBs = []byte(customNL)
		nlLen = int32(len(nlBs))
	}
	nlBSDfl := []byte(STR_DEFAULT_NEWLINE)

	tf := TextFile{
		_file:      file,
		_convertNL: convertNL,
		_nlBs:      nlBs,
		_nlLen:     nlLen,
		_nlBSDfl:   nlBSDfl,
	}

	return &tf
}
func (ego *TextFile) Open(openMode int, permMode int, perm fs.FileMode) int32 {
	return ego._file.Open(openMode, permMode, perm)
}

func (ego *TextFile) Close() {
	ego._file.Close()
}

func (ego *TextFile) FullPathName() string {
	return ego._file.FullPathName()

}
func (ego *TextFile) WriteLine(str string) (int, int32) {
	ego._file.PLockExclusive()
	ego._file.TLockExclusive()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockExclusive()
	return ego.WriteLine(str)
}
func (ego *TextFile) WriteLineBared(str string) (int, int32) {
	bs, rc := ego._file.WriteStringBared(str)
	if core.Err(rc) {
		return bs, rc
	}
	bs += ego.WriteNewLineBared()
	return bs, core.MkSuccess(0)
}

func (ego *TextFile) WriteLinePartial(str string, idx int, wlen int) (int, int32) {
	ego._file.PLockExclusive()
	ego._file.TLockExclusive()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockExclusive()
	bs, rc := ego._file.WriteStringPartialBared(str, idx, wlen)
	if core.Err(rc) {
		return bs, rc
	}
	bs += ego.WriteNewLineBared()
	return bs, core.MkSuccess(0)
}

func (ego *TextFile) WriteLines(lines []string) int32 {
	ego._file.PLockExclusive()
	ego._file.TLockExclusive()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockExclusive()
	return ego.WriteLinesBared(lines)
}

func (ego *TextFile) WriteLinesBared(lines []string) int32 {
	for i := 0; i < len(lines); i++ {
		_, rc := ego._file.WriteStringBared(lines[i])
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *TextFile) WriteNewLineBared() int {
	if !ego._convertNL {
		ego._file.WriteBytesBared(ego._nlBs)
		return len(ego._nlBs)
	} else {
		ego._file.WriteBytesBared(ego._nlBSDfl)
		return len(ego._nlBSDfl)
	}
}

func (ego *TextFile) ReadLineConvNewLine(bufSz int) (string, int32) {
	ego._file.PLockShare()
	ego._file.TLockShare()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockShare()
	return ego.ReadLineConvNewLineBared(bufSz)
}

func (ego *TextFile) ReadLines(bufSz int, maxCount int) ([]string, int32) {
	ego._file.PLockShare()
	ego._file.TLockShare()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockShare()
	return ego.ReadLinesBared(bufSz, maxCount)
}
func (ego *TextFile) ReadLinesBared(bufSz int, maxCount int) ([]string, int32) {
	var r = make([]string, 0)
	if maxCount < 0 {
		for {
			str, rc := ego.ReadLine(bufSz)
			if core.Err(rc) {
				if core.IsErrType(rc, core.EC_EOF) {
					return r, core.MkErr(core.EC_EOF, 0)
				}
				return r, core.MkErr(core.EC_FILE_READ_FAILED, 0)
			}
			r = append(r, str)

		}
	} else {
		for i := 0; i < maxCount; i++ {
			str, rc := ego.ReadLine(bufSz)
			if core.Err(rc) {
				if core.IsErrType(rc, core.EC_EOF) {
					return r, core.MkErr(core.EC_EOF, 0)
				}
				return r, core.MkErr(core.EC_FILE_READ_FAILED, 0)
			}
			r = append(r, str)
		}
	}

	return r, core.MkSuccess(0)
}

func (ego *TextFile) ReadLine(bufSz int) (string, int32) {
	ego._file.PLockShare()
	ego._file.TLockShare()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockShare()
	return ego.ReadLineBared(bufSz)
}

func (ego *TextFile) ReadLineBared(bufSz int) (string, int32) {
	if ego._convertNL {
		return ego.ReadLineConvNewLine(bufSz)
	} else {
		return ego.ReadLineNoConvNewLineBared(bufSz)
	}
}

func (ego *TextFile) ReadLineNoConvNewLineBared(bufSz int) (string, int32) {
	lineBs := make([]byte, bufSz)
	idx := 0

	for {
		by, rc := ego._file.ReadByteBared()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_EOF) {
				return string(lineBs), core.MkErr(core.EC_EOF, 0)
			} else {
				return string(lineBs), core.MkErr(core.EC_FILE_READ_FAILED, 0)
			}
		} else {
			if idx >= len(lineBs) {
				lineBs = append(lineBs, by)
				idx++
			} else {
				lineBs[idx] = by
				idx++
			}

			if by == ego._nlBs[0] {
				if ego._nlLen > 1 {
					by2, rc := ego._file.ReadByteBared()
					if core.Err(rc) {
						if core.IsErrType(rc, core.EC_EOF) {
							return string(lineBs), core.MkErr(core.EC_EOF, 0)
						} else {
							return string(lineBs), core.MkErr(core.EC_FILE_READ_FAILED, 0)
						}
					} else {
						if by2 == ego._nlBs[1] {
							if idx >= len(lineBs) {
								lineBs = append(lineBs, by2)
								idx++
							} else {
								lineBs[idx] = by2
								idx++
							}
							return string(lineBs), core.MkSuccess(0)
						}
					}
				} else {
					return string(lineBs), core.MkSuccess(0)
				}
			}
		}
	}

	return string(lineBs), core.MkSuccess(0)
}

// //////////
func (ego *TextFile) ReadLineConvNewLineBared(bufSz int) (string, int32) {
	lineBs := make([]byte, bufSz)
	idx := 0
	for {
		by, rc := ego._file.ReadByteBared()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_EOF) {
				return string(lineBs), core.MkErr(core.EC_EOF, 0)
			} else {
				return string(lineBs), core.MkErr(core.EC_FILE_READ_FAILED, 0)
			}
		}
		if by != ego._nlBs[0] {
			if idx >= len(lineBs) {
				lineBs = append(lineBs, by)
				idx++
			} else {
				lineBs[idx] = by
				idx++
			}

		} else {
			if ego._nlLen > 1 {
				by2, rc := ego._file.ReadByteBared()
				if core.Err(rc) {
					if core.IsErrType(rc, core.EC_EOF) {
						return string(lineBs), core.MkErr(core.EC_EOF, 0)
					} else {
						return string(lineBs), core.MkErr(core.EC_FILE_READ_FAILED, 0)
					}
				} else {
					if by2 == ego._nlBs[1] {
						if idx >= len(lineBs) {
							lineBs = append(lineBs, ego._nlBSDfl[0])
							idx++
						} else {
							lineBs[idx] = ego._nlBSDfl[0]
							idx++
						}
						return string(lineBs), core.MkSuccess(0)
					}
				}
			} else {
				lineBs = append(lineBs, ego._nlBSDfl[0])
				if idx >= len(lineBs) {
					lineBs = append(lineBs, ego._nlBSDfl[0])
					idx++
				} else {
					lineBs[idx] = ego._nlBSDfl[0]
					idx++
				}
			}
		}
	}

	return string(lineBs), core.MkSuccess(0)

}

func (ego *TextFile) ReadStringAll() (string, int32) {
	ego._file.PLockShare()
	ego._file.TLockShare()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockShare()
	return ego.ReadStringAllBared()
}

func (ego *TextFile) ReadStringAllBared() (string, int32) {
	bs, rc := ego._file.ReadAllBared()
	if core.Err(rc) {
		return "", rc
	}
	return string(bs), core.MkSuccess(0)
}

func (ego *TextFile) ReadAll() ([]byte, int32) {
	ego._file.PLockShare()
	ego._file.TLockShare()
	defer ego._file.PUnlock()
	defer ego._file.TUnLockShare()
	return ego.ReadAllBared()
}

func (ego *TextFile) ReadAllBared() ([]byte, int32) {
	return ego._file.ReadAll()
}

func (ego *TextFile) GetInfo() os.FileInfo {
	return ego._file.GetInfo()
}

func (ego *TextFile) TLockShare() {
	ego._file.TLockShare()

}
func (ego *TextFile) TLockExclusive() {
	ego._file.TLockExclusive()
}
func (ego *TextFile) TTryLockShare() bool {
	return ego._file.TTryLockShare()
}
func (ego *TextFile) TTryLockExclusive() bool {
	return ego._file.TTryLockShare()
}
func (ego *TextFile) TUnLockShare() {
	ego._file.TTryLockShare()
}
func (ego *TextFile) TUnLockExclusive() {
	ego._file.TUnLockExclusive()
}

func (ego *TextFile) PLockShare() int32 {
	return ego._file.PLockShare()
}
func (ego *TextFile) PLockExclusive() int32 {
	return ego._file.PLockExclusive()
}
func (ego *TextFile) PTryLockShare() int32 {
	return ego._file.PTryLockShare()
}
func (ego *TextFile) PTryLockExclusive() int32 {
	return ego._file.PTryLockExclusive()
}
func (ego *TextFile) PUnlock() int32 {
	return ego._file.PUnlock()
}
