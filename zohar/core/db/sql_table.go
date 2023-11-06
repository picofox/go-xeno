package db

import "xeno/zohar/core"

type SQLTable struct {
	RecordDesc
	_dbName    string
	_tableName string
	_dataGrid  [][]any
}

func (ego *SQLTable) Data() [][]any {
	return ego._dataGrid
}

func (ego *SQLTable) RowAt(idx int64) []any {
	return ego._dataGrid[idx]
}

func NeoSQLTable(dbName string, tableName string, fieldCnt int, bOrig bool) *SQLTable {

	return &SQLTable{
		_dbName:    dbName,
		_tableName: tableName,
		RecordDesc: RecordDesc{
			_fieldDesc:       make([]*FieldDesc, fieldCnt),
			_filedDescByName: make(map[string]*FieldDesc),
			_original:        bOrig,
		},
		_dataGrid: make([][]any, 0),
	}
}

func (ego *SQLTable) InitRow() {
	cnt := len(ego._fieldDesc)
	rowRaw := make([]any, cnt)
	ego._dataGrid = append(ego._dataGrid, rowRaw)
}

func (ego *SQLTable) SetRow(idx uint16, recVal []any) int32 {
	cnt := len(ego._fieldDesc)
	if len(recVal) != cnt {
		return core.MkErr(core.EC_INDEX_OOB, 1)
	}

	if int(idx) < len(ego._dataGrid) {
		ego._dataGrid[idx] = recVal
	} else {
		if idx == 0xFFFF {
			ego._dataGrid = append(ego._dataGrid, recVal)
		} else {
			return core.MkErr(core.EC_INDEX_OOB, 1)
		}
	}

	return core.MkSuccess(0)
}

func (ego *SQLTable) AddRow(recVal []any) int32 {
	return ego.SetRow(uint16(0xFFFF), recVal)
}
