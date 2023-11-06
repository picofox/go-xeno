package db

import (
	"xeno/zohar/core"
)

type RecordDesc struct {
	_fieldDesc       []*FieldDesc
	_filedDescByName map[string]*FieldDesc
	_original        bool
}

func NeoRecordDesc(fieldCnt int, bOrig bool) *RecordDesc {
	return &RecordDesc{
		_fieldDesc:       make([]*FieldDesc, fieldCnt),
		_filedDescByName: make(map[string]*FieldDesc),
		_original:        bOrig,
	}
}

func (ego *RecordDesc) AddFieldDesc(name string, posInTable uint16, dbType uint8, isUnsigned bool, nullable bool) int32 {
	return ego.SetFieldDesc(name, 0xFFFF, posInTable, dbType, isUnsigned, nullable)
}

func (ego *RecordDesc) SetFieldDesc(name string, pos uint16, posInTable uint16, dbType uint8, isUnsigned bool, nullable bool) int32 {
	flen := uint16(len(ego._fieldDesc))
	if pos >= flen {
		if pos == 0xFFFF {
			pFD := NeoFieldDesc(name, flen, posInTable, dbType, isUnsigned, nullable)
			ego._filedDescByName[name] = pFD
			ego._fieldDesc = append(ego._fieldDesc, pFD)
			return core.MkSuccess(0)
		} else {
			return core.MkErr(core.EC_INDEX_OOB, 1)
		}
	}

	pFD := NeoFieldDesc(name, flen, posInTable, dbType, isUnsigned, nullable)
	ego._fieldDesc[pos] = pFD
	ego._filedDescByName[name] = pFD
	return core.MkSuccess(0)
}

func (ego *RecordDesc) FieldDescByName(name string) *FieldDesc {
	fd, ok := ego._filedDescByName[name]
	if !ok {
		return nil
	}
	return fd
}

func (ego *RecordDesc) FieldDesc(pos uint16) *FieldDesc {
	return ego._fieldDesc[pos]
}

func (ego *RecordDesc) IsOriginal() bool {
	return ego._original
}

func (ego *RecordDesc) FieldCount() uint16 {
	return uint16(len(ego._fieldDesc))
}
