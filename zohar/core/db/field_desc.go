package db

import (
	"xeno/zohar/core/memory"
)

const (
	DBF_NULLABLE  = uint8(0x01)
	DBF_UNSIGNNED = uint8(0x02)
)

const (
	DBK_NONE  = 0
	DBK_PK    = 1
	DBK_INDEX = 2
)

const (
	DBF_TYPE_NA         = uint8(0)
	DBF_TYPE_TINYINT    = uint8(1)
	DBF_TYPE_SMALLINT   = uint8(2)
	DBF_TYPE_MEDIUMINT  = uint8(3)
	DBF_TYPE_INT        = uint8(4)
	DBF_TYPE_BIGINT     = uint8(5)
	DBF_TYPE_FLOAT      = uint8(6)
	DBF_TYPE_DOUBLE     = uint8(7)
	DBF_TYPE_CHAR       = uint8(8)
	DBF_TYPE_VARCHAR    = uint8(9)
	DBF_TYPE_TEXT       = uint8(10)
	DBF_TYPE_TINYBLOB   = uint8(11)
	DBF_TYPE_MEDIUMBLOB = uint8(12)
	DBF_TYPE_BLOB       = uint8(13)
	DBF_TYPE_LONGBLOB   = uint8(14)
	DBF_TYPE_DATETIME   = uint8(15)
	DBF_TYPE_DATE       = uint8(16)
	DBF_TYPE_TIME       = uint8(17)
	DBF_TYPE_TIMESTAMP  = uint8(18)
	DBF_TYPE_COUNT      = uint8(19)
)

var DBType2LocalType [DBF_TYPE_COUNT]func(u bool) uint8 = [DBF_TYPE_COUNT]func(u bool) uint8{
	func(bool) uint8 {
		return memory.T_NULL
	}, //0
	func(u bool) uint8 {
		if u {
			return memory.T_U8
		} else {
			return memory.T_I8
		}
	}, //1
	func(u bool) uint8 {
		if u {
			return memory.T_U16
		} else {
			return memory.T_I16
		}
	}, //2
	func(u bool) uint8 {
		if u {
			return memory.T_U32
		} else {
			return memory.T_I32
		}
	}, //3
	func(u bool) uint8 {
		if u {
			return memory.T_U32
		} else {
			return memory.T_I32
		}
	}, //4
	func(u bool) uint8 {
		if u {
			return memory.T_U64
		} else {
			return memory.T_I64
		}
	}, //5
	func(u bool) uint8 {
		return memory.T_F32
	}, //6
	func(u bool) uint8 {
		return memory.T_F64
	}, //7
	func(u bool) uint8 {
		return memory.T_STR
	}, //8
	func(u bool) uint8 {
		return memory.T_STR
	}, //9
	func(u bool) uint8 {
		return memory.T_STR
	}, //10
	func(u bool) uint8 {
		return memory.T_BYTES
	}, //11
	func(u bool) uint8 {
		return memory.T_BYTES
	}, //12
	func(u bool) uint8 {
		return memory.T_BYTES
	}, //13
	func(u bool) uint8 {
		return memory.T_BYTES
	}, //14
	func(u bool) uint8 {
		return memory.T_I64
	}, //15
	func(u bool) uint8 {
		return memory.T_I64
	}, //16
	func(u bool) uint8 {
		return memory.T_STR
	}, //17
	func(u bool) uint8 {
		return memory.T_I64
	}, //18
}

type FieldDesc struct {
	_name       string
	_dbType     uint8
	_localType  uint8
	_flags      uint8
	_keyType    uint8
	_pos        uint16
	_posInTable uint16
}

func NeoFieldDesc(name string, pos uint16, posInTable uint16, dbType uint8, isUnsigned bool, nullable bool, keyType uint8) *FieldDesc {
	f := uint8(0)
	if isUnsigned {
		f = f | DBF_UNSIGNNED
	}
	if nullable {
		f = f | DBF_NULLABLE
	}

	lt := DBType2LocalType[dbType](isUnsigned)
	fd := &FieldDesc{
		_name:       name,
		_pos:        pos,
		_posInTable: posInTable,
		_dbType:     dbType,
		_localType:  lt,
		_flags:      f,
		_keyType:    keyType,
	}
	return fd
}

func (ego *FieldDesc) IsPK() bool {
	if ego._keyType == DBK_PK {
		return true
	}
	return false
}

func (ego *FieldDesc) IsIndex() bool {
	if ego._keyType == DBK_INDEX {
		return true
	}
	return false
}

func (ego *FieldDesc) Name() string {
	return ego._name
}

func (ego *FieldDesc) Pos() uint16 {
	return ego._pos
}

func (ego *FieldDesc) PosInTable() uint16 {
	return ego._posInTable
}

func (ego *FieldDesc) SetPosInTable(i uint16) {
	ego._posInTable = i
}

func (ego *FieldDesc) TableIndex() uint16 {
	return ego._posInTable
}

func (ego *FieldDesc) DataBaseFieldType() uint8 {
	return ego._dbType
}

func (ego *FieldDesc) LocalType() uint8 {
	return ego._localType
}

func (ego *FieldDesc) IsUnsigned() bool {
	return ego._flags&DBF_UNSIGNNED != 0
}

func (ego *FieldDesc) Nullable() bool {
	return ego._flags&DBF_NULLABLE != 0
}

func (ego *FieldDesc) SetUnsigned(b bool) {
	if b {
		ego._flags = ego._flags | DBF_UNSIGNNED
	} else {
		ego._flags = ego._flags & (^DBF_UNSIGNNED)
	}
}

func (ego *FieldDesc) SetNullable(b bool) {
	if b {
		ego._flags = ego._flags | DBF_NULLABLE
	} else {
		ego._flags = ego._flags & (^DBF_NULLABLE)
	}
}
