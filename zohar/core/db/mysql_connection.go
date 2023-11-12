package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/memory"
)

type MysqlConnection struct {
	_index               int
	_connTestFailedCount int
	_lastConnTestTs      int64
	_conn                *sql.DB
	_goContext           context.Context
	_cancelFunc          context.CancelFunc
	_transcation         *sql.Tx

	_cfgPool    *config.DBConnectionPoolConfig
	_tableNames []string
}

func (ego *MysqlConnection) TruncateTable(name string) int32 {
	_, rc := ego.Delete("TRUNCATE " + name)
	return rc
}

func (ego *MysqlConnection) _Insert(sqlString string, arg ...any) (int64, int32) {
	var err error
	var result sql.Result
	if ego._goContext == nil {
		if ego._transcation == nil {
			result, err = ego._conn.Exec(sqlString, arg...)
		} else {
			result, err = ego._transcation.Exec(sqlString, arg...)
		}
	} else {
		if ego._transcation == nil {
			result, err = ego._conn.ExecContext(ego._goContext, sqlString, arg...)
		} else {
			result, err = ego._transcation.ExecContext(ego._goContext, sqlString, arg...)
		}
	}
	if err != nil {
		return 0, core.MkErr(core.EC_DIR_ALREADY_EXIST, 1)
	}
	ra, err := result.RowsAffected()
	if err != nil {
		ra = -1
	}
	return ra, core.MkSuccess(0)
}

func (ego *MysqlConnection) Create(sqlTable *SQLTable, sqlString string, arg ...any) (int64, int32) {
	if sqlTable == nil {
		return ego._Insert(sqlString, arg)
	}

	var stmt *sql.Stmt
	var err error
	if ego._goContext == nil {
		if ego._transcation == nil {
			stmt, err = ego._conn.Prepare(sqlString)
		} else {
			stmt, err = ego._transcation.Prepare(sqlString)
		}
	} else {
		if ego._transcation == nil {
			stmt, err = ego._conn.PrepareContext(ego._goContext, sqlString)
		} else {
			stmt, err = ego._transcation.PrepareContext(ego._goContext, sqlString)
		}
	}
	if err != nil {
		return 0, core.MkErr(core.EC_DB_PREPARE_FAILED, 1)
	}

	defer stmt.Close()

	var totalCount int64 = 0
	for rowIdx := int64(0); rowIdx < int64(len(sqlTable.Data())); rowIdx++ {
		result, err := stmt.Exec(sqlTable.RowAt(rowIdx)...)
		if err != nil {
			return 0, core.MkErr(core.EC_DB_INSERT_FAILED, 1)
		}
		cnt, _ := result.RowsAffected()
		totalCount = totalCount + cnt
	}

	return totalCount, core.MkSuccess(0)
}

func (ego *MysqlConnection) Update(sqlString string, arg ...any) (int64, int32) {
	var err error
	var result sql.Result
	if ego._goContext == nil {
		if ego._transcation == nil {
			result, err = ego._conn.Exec(sqlString, arg...)
		} else {
			result, err = ego._transcation.Exec(sqlString, arg...)
		}

	} else {
		if ego._transcation == nil {
			result, err = ego._conn.ExecContext(ego._goContext, sqlString, arg...)
		} else {
			result, err = ego._transcation.ExecContext(ego._goContext, sqlString, arg...)
		}
	}

	if err != nil {
		return 0, core.MkErr(core.EC_DIR_ALREADY_EXIST, 1)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		ra = -1
	}

	return ra, core.MkSuccess(0)
}

func (ego *MysqlConnection) Delete(sqlString string, arg ...any) (int64, int32) {
	var err error
	var result sql.Result
	if ego._goContext == nil {
		if ego._transcation == nil {
			result, err = ego._conn.Exec(sqlString, arg...)
		} else {
			result, err = ego._transcation.Exec(sqlString, arg...)
		}

	} else {
		if ego._transcation == nil {
			result, err = ego._conn.ExecContext(ego._goContext, sqlString, arg...)
		} else {
			result, err = ego._transcation.ExecContext(ego._goContext, sqlString, arg...)
		}
	}

	if err != nil {
		return 0, core.MkErr(core.EC_DIR_ALREADY_EXIST, 1)
	}

	ra, err := result.RowsAffected()
	if err != nil {
		ra = -1
	}

	return ra, core.MkSuccess(0)
}

func (ego *MysqlConnection) CommitTransaction() int32 {
	if ego._transcation == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	var err error
	err = ego._transcation.Commit()
	if err != nil {
		return core.MkErr(core.EC_DB_COMMIT_TRANS_FAILED, 1)
	}
	ego._transcation = nil
	return core.MkSuccess(0)
}

func (ego *MysqlConnection) RollbackTransaction() int32 {
	if ego._transcation == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	var err error
	err = ego._transcation.Rollback()
	if err != nil {
		return core.MkErr(core.EC_DB_ROLLBACL_TRANS_FAILED, 1)
	}
	ego._transcation = nil
	return core.MkSuccess(0)
}

func (ego *MysqlConnection) SetAutoCommit(ac bool) int32 {
	var err error
	if ac {
		_, err = ego._conn.Exec("SET AUTOCOMMIT = 1")
	} else {
		_, err = ego._conn.Exec("SET AUTOCOMMIT = 0")
	}
	if err != nil {
		return core.MkErr(core.EC_DB_SET_AUTOCOMMIT_FAILED, 1)
	}

	return core.MkSuccess(0)
}

func (ego *MysqlConnection) BeginTransaction() int32 {
	var err error
	if ego._goContext == nil {
		ego._transcation, err = ego._conn.Begin()
	} else {
		ego._transcation, err = ego._conn.BeginTx(ego._goContext, nil)
	}

	if err != nil {
		return core.MkErr(core.EC_DB_SET_AUTOCOMMIT_FAILED, 1)
	}

	return core.MkSuccess(0)
}

func isBinType(name string) bool {
	name = strings.ToUpper(name)
	if strings.Contains(name, "BLOB") {
		return true
	}

	return false
}

var mysqldatebaseTypeNames2LocalType = map[string]func(int32) (any, uint8){
	"TINYINT": func(int32) (any, uint8) {
		var p sql.NullByte
		return &p, memory.T_I8
	},
	"UNSIGNED TINYINT": func(int32) (any, uint8) {
		var p sql.NullByte
		return &p, memory.T_U8
	},
	"SMALLINT": func(int32) (any, uint8) {
		var p sql.NullInt16
		return &p, memory.T_I16
	},
	"UNSIGNED SMALLINT": func(int32) (any, uint8) {
		var p sql.NullInt16
		return &p, memory.T_U16
	},
	"MEDIUNINT": func(int32) (any, uint8) {
		var p sql.NullInt32
		return &p, memory.T_I32
	},
	"UNSIGNED MEDIUNINT": func(int32) (any, uint8) {
		var p sql.NullInt32
		return &p, memory.T_U32
	},
	"INT": func(int32) (any, uint8) {
		var p sql.NullInt32
		return &p, memory.T_I32
	},
	"UNSIGNED INT": func(int32) (any, uint8) {
		var p sql.NullInt32
		return &p, memory.T_U32
	},
	"BIGINT": func(int32) (any, uint8) {
		var p sql.NullInt64
		return &p, memory.T_I64
	},
	"UNSIGNED BIGINT": func(int32) (any, uint8) {
		var p sql.NullInt64
		return &p, memory.T_U64
	},

	"VARCHAR": func(int32) (any, uint8) {
		var p sql.NullString
		return &p, memory.T_STR
	},
	"CHAR": func(int32) (any, uint8) {
		var p sql.NullString
		return &p, memory.T_STR
	},
	"TEXT": func(int32) (any, uint8) {
		var p sql.NullString
		return &p, memory.T_STR
	},
	"SMALLTEXT": func(int32) (any, uint8) {
		var p sql.NullString
		return &p, memory.T_STR
	},
	"LONGTEXT": func(int32) (any, uint8) {
		var p sql.NullString
		return &p, memory.T_STR
	},
	"DATETIME": func(flags int32) (any, uint8) {
		if flags == 0 {
			var p time.Time
			return &p, memory.T_I64
		} else {
			var p sql.NullTime
			return &p, memory.T_I64
		}
	},
	"TIMESTAMP": func(flags int32) (any, uint8) {
		if flags == 0 {
			var p time.Time
			return &p, memory.T_I64
		} else {
			var p sql.NullTime
			return &p, memory.T_I64
		}
	},
	"DATE": func(flags int32) (any, uint8) {
		if flags == 0 {
			var p time.Time
			return &p, memory.T_I64
		} else {
			var p sql.NullTime
			return &p, memory.T_I64
		}
	},
	"FLOAT": func(int32) (any, uint8) {
		var p sql.NullFloat64
		return &p, memory.T_F32
	},
	"DOUBLE": func(int32) (any, uint8) {
		var p sql.NullFloat64
		return &p, memory.T_F64
	},
}

func (ego *MysqlConnection) SetContextCancelable() {
	ego._goContext, ego._cancelFunc = context.WithCancel(context.Background())
}

func (ego *MysqlConnection) CancelOperation() {
	if ego._cancelFunc != nil {
		ego._cancelFunc()
		ego._cancelFunc = nil
		ego._goContext = nil
	}
}

func calcColTypes(rows *sql.Rows) ([]any, []uint8) {
	var err error
	var colsTypes []*sql.ColumnType = nil
	colsTypes, err = rows.ColumnTypes()
	if err != nil {
		return nil, nil
	}

	colCount := len(colsTypes)
	var values []any = make([]any, colCount)
	var dirtyField []uint8 = make([]uint8, colCount)

	for i, _ := range values {
		switch colsTypes[i].ScanType().Kind() {
		case reflect.Int8:
			var p int8
			values[i] = &p
		case reflect.Uint8:
			var p uint8
			values[i] = &p
		case reflect.Int16:
			var p int16
			values[i] = &p
		case reflect.Uint16:
			var p uint16
			values[i] = &p
		case reflect.Int32:
			var p int32
			values[i] = &p
		case reflect.Uint32:
			var p uint32
			values[i] = &p
		case reflect.Int64:
			var p int64
			values[i] = &p
		case reflect.Uint64:
			var p uint64
			values[i] = &p
		case reflect.Float32:
			var p float32
			values[i] = &p
		case reflect.Float64:
			var p float64
			values[i] = &p
		case reflect.String:
			var p string
			values[i] = &p

		case reflect.Slice:
			if strings.Contains(colsTypes[i].DatabaseTypeName(), "CHAR") || strings.Contains(colsTypes[i].DatabaseTypeName(), "TEXT") {
				var p string
				values[i] = &p
			} else if isBinType(colsTypes[i].DatabaseTypeName()) {
				var p []byte
				values[i] = &p

			} else if "TIME" == strings.ToUpper(colsTypes[i].DatabaseTypeName()) {
				var p string
				values[i] = &p
			} else {
				var p []byte
				values[i] = &p
			}

		case reflect.Struct:
			f, ok := mysqldatebaseTypeNames2LocalType[strings.ToUpper(colsTypes[i].DatabaseTypeName())]
			if ok {
				lType := uint8(0)
				nb, ok := colsTypes[i].Nullable()
				if ok {
					if nb {
						values[i], lType = f(1)
					} else {
						values[i], lType = f(0)
					}
				} else {
					panic("get Nullable Failed")
				}

				dirtyField[i] = lType
			} else {
				logging.Log(core.LL_WARN, "Field Type Parse Failed %s", colsTypes[i].DatabaseTypeName())
				var p any
				values[i] = &p
			}

		default:
			//panic("Not implemented type " + colsTypes[i].ScanType().Kind().String())
			var p any
			values[i] = &p

		}
	}
	return values, dirtyField
}

var ROneLocal2ToValueArr [memory.T_TLV]func(dbt uint8, nullable bool) any = [memory.T_TLV]func(dbt uint8, nullable bool) any{
	nil,
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullByte
			return &p
		} else {
			var p int8
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullInt16
			return &p
		} else {
			var p int16
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullInt32
			return &p
		} else {
			var p int32
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if dbt == DBF_TYPE_DATE || dbt == DBF_TYPE_DATETIME || dbt == DBF_TYPE_TIMESTAMP {
			if nullable {
				var p sql.NullTime
				return &p
			} else {
				var p time.Time
				return &p
			}
		} else {
			if nullable {
				var p sql.NullInt64
				return &p
			} else {
				var p int64
				return &p
			}
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullByte
			return &p
		} else {
			var p uint8
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullInt16
			return &p
		} else {
			var p uint16
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullInt32
			return &p
		} else {
			var p uint32
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullInt64
			return &p
		} else {
			var p uint64
			return &p
		}
	},
	nil,
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullFloat64
			return &p
		} else {
			var p float32
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p sql.NullFloat64
			return &p
		} else {
			var p float64
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		if nullable {
			var p []byte
			return &p
		} else {
			var p []byte
			return &p
		}
	},
	func(dbt uint8, nullable bool) any {
		var p sql.NullString
		return &p
	},
}

var ROneLocalResultParseArr [memory.T_TLV]func(v any, dt uint8, nullable bool) any = [memory.T_TLV]func(v any, dt uint8, nullable bool) any{
	nil,
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullByte)
			if nv.Valid {
				return int8(nv.Byte)
			} else {
				return nil
			}
		} else {
			return *(v.(*int8))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullInt16)
			if nv.Valid {
				return int16(nv.Int16)
			} else {
				return nil
			}
		} else {
			return *(v.(*int16))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullInt32)
			if nv.Valid {
				return int32(nv.Int32)
			} else {
				return nil
			}
		} else {
			return *(v.(*int32))
		}
	},
	func(v any, dbt uint8, nullable bool) any {
		if dbt == DBF_TYPE_DATE || dbt == DBF_TYPE_DATETIME || dbt == DBF_TYPE_TIMESTAMP {
			if nullable { //here
				nv := *v.(*sql.NullTime)
				if nv.Valid {
					return int64(nv.Time.UnixMilli())
				} else {
					return nil
				}
			} else {
				return (*v.(*time.Time)).UnixMilli()
			}
		} else {
			if nullable { //here
				nv := *v.(*sql.NullInt64)
				if nv.Valid {
					return uint64(nv.Int64)
				} else {
					return nil
				}
			} else {
				return *(v.(*uint64))
			}
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullByte)
			if nv.Valid {
				return uint8(nv.Byte)
			} else {
				return nil
			}
		} else {
			return *(v.(*uint8))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullInt16)
			if nv.Valid {
				return uint16(nv.Int16)
			} else {
				return nil
			}
		} else {
			return *(v.(*uint16))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullInt32)
			if nv.Valid {
				return uint32(nv.Int32)
			} else {
				return nil
			}
		} else {
			return *(v.(*uint32))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullInt64)
			if nv.Valid {
				return uint64(nv.Int64)
			} else {
				return nil
			}
		} else {
			return *(v.(*uint64))
		}
	},
	nil,
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullFloat64)
			if nv.Valid {
				return float32(nv.Float64)
			} else {
				return nil
			}
		} else {
			return *(v.(*float32))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		if nullable {
			nv := *v.(*sql.NullFloat64)
			if nv.Valid {
				return nv.Float64
			} else {
				return nil
			}
		} else {
			return *(v.(*float64))
		}
	},
	func(v any, dt uint8, nullable bool) any {
		bs := *(v.(*[]byte))
		if bs == nil {
			return nil
		} else {
			return bs
		}
	},
	func(v any, dt uint8, nullable bool) any {
		nv := *v.(*sql.NullString)
		if nv.Valid {
			return nv.String
		} else {
			return nil
		}
	},
}

func (ego *MysqlConnection) RetrieveField(dbt uint8, nullable bool, isUnsigned bool, sqlString string, arg ...any) (*memory.TLV, int32) {
	var row *sql.Row
	//	var err error
	if ego._goContext == nil {
		if ego._transcation == nil {
			row = ego._conn.QueryRow(sqlString, arg...)
		} else {
			row = ego._transcation.QueryRow(sqlString, arg...)
		}

	} else {
		if ego._transcation == nil {
			row = ego._conn.QueryRowContext(ego._goContext, sqlString, arg...)
		} else {
			row = ego._transcation.QueryRowContext(ego._goContext, sqlString, arg...)
		}

	}

	localType := DBType2LocalType[dbt](isUnsigned)
	f := ROneLocal2ToValueArr[localType]
	if f == nil {
		return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
	}
	value := f(dbt, nullable)
	err := row.Scan(value)
	if err != nil {
		return nil, core.MkErr(core.EC_DB_RETRIVE_DATA_FAILED, 1)
	}

	f2 := ROneLocalResultParseArr[localType]
	if f2 == nil {
		return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
	}
	fld := f2(value, dbt, nullable)
	var tlv *memory.TLV = nil
	if fld == nil {
		tlv = memory.CreateTLV(memory.DT_SINGLE, memory.T_NULL, memory.T_NULL, nil)
	} else {
		tlv = memory.CreateTLV(memory.DT_SINGLE, localType, memory.T_NULL, fld)
	}

	return tlv, core.MkSuccess(0)
}

func (ego *MysqlConnection) RetrieveRecord(rd *RecordDesc, sqlString string, arg ...any) (*memory.TLV, int32) {
	var row *sql.Row
	//	var err error
	if ego._goContext == nil {
		if ego._transcation == nil {
			row = ego._conn.QueryRow(sqlString, arg...)
		} else {
			row = ego._transcation.QueryRow(sqlString, arg...)
		}

	} else {
		if ego._transcation == nil {
			row = ego._conn.QueryRowContext(ego._goContext, sqlString, arg...)
		} else {
			row = ego._transcation.QueryRowContext(ego._goContext, sqlString, arg...)
		}

	}

	fcnt := rd.FieldCount()
	var values []any = make([]any, fcnt)
	for i := uint16(0); i < fcnt; i++ {
		fd := rd.FieldDesc(i)
		f := ROneLocal2ToValueArr[fd.LocalType()]
		if f == nil {
			return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
		}
		values[i] = f(fd.DataBaseFieldType(), fd.Nullable())
	}

	err := row.Scan(values...)
	if err != nil {
		return nil, core.MkErr(core.EC_DB_RETRIVE_DATA_FAILED, 1)
	}
	var recordData []any = make([]any, fcnt)

	for i := uint16(0); i < fcnt; i++ {
		fd := rd.FieldDesc(i)
		f := ROneLocalResultParseArr[fd.LocalType()]
		if f == nil {
			return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
		}
		recordData[i] = f(values[i], fd.DataBaseFieldType(), fd.Nullable())
	}

	tlv := memory.CreateTLV(memory.DT_LIST, memory.T_NULL, memory.T_NULL, recordData)

	return tlv, core.MkSuccess(0)
}

func (ego *MysqlConnection) Retrieve(sqlString string, arg ...any) ([]*memory.TLV, int32) {
	var rData []*memory.TLV = make([]*memory.TLV, 0)

	var rows *sql.Rows
	var err error
	if ego._goContext == nil {
		if ego._transcation == nil {
			rows, err = ego._conn.Query(sqlString, arg...)
		} else {
			rows, err = ego._transcation.Query(sqlString, arg...)
		}
	} else {
		if ego._transcation == nil {
			rows, err = ego._conn.QueryContext(ego._goContext, sqlString, arg...)
		} else {
			rows, err = ego._transcation.QueryContext(ego._goContext, sqlString, arg...)
		}
	}
	if err != nil {
		return nil, core.MkErr(core.EC_DB_RETRIVE_DATA_FAILED, 1)
	}

	var colCount int = 0
	values, dirtyFields := calcColTypes(rows)
	if values == nil {
		return nil, core.MkErr(core.EC_DB_GET_COL_INFO_FAILED, 1)
	}

	colCount = len(values)
	for rows.Next() {

		err = rows.Scan(values...)
		if err != nil {
			return nil, core.MkErr(core.EC_DB_RETRIVE_DATA_FAILED, 1)
		}

		var recordData []any = make([]any, colCount)

		for i := 0; i < colCount; i++ {
			//var rawValue = *values[i].(*interface{})
			switch values[i].(type) {
			case *int8:
				recordData[i] = *values[i].(*int8)
			case *uint8:
				recordData[i] = *values[i].(*uint8)
			case *int16:
				recordData[i] = *values[i].(*int16)
			case *uint16:
				recordData[i] = *values[i].(*uint16)
			case *int32:
				recordData[i] = *values[i].(*int32)
			case *uint32:
				recordData[i] = *values[i].(*uint32)
			case *int64:
				recordData[i] = *values[i].(*int64)
			case *uint64:
				recordData[i] = *values[i].(*uint64)
			case *string:
				recordData[i] = *values[i].(*string)
			case *float32:
				recordData[i] = *values[i].(*float32)
			case *float64:
				recordData[i] = *values[i].(*float64)
			case *[]byte:
				recordData[i] = *values[i].(*[]byte)
				var bs = recordData[i].([]byte)
				if bs == nil {
					recordData[i] = nil
				} else {
					recordData[i] = bs
				}

			case *sql.NullTime:
				v := *values[i].(*sql.NullTime)
				if v.Valid {
					recordData[i] = int64(v.Time.UnixMilli())
				} else {
					recordData[i] = nil
				}
			case *sql.NullInt64:
				v := *values[i].(*sql.NullInt64)
				if v.Valid {
					if dirtyFields[i] == 0 {
						recordData[i] = v.Int64
					} else {
						if dirtyFields[i] == memory.T_I64 {
							recordData[i] = int64(v.Int64)
						} else if dirtyFields[i] == memory.T_U64 {
							recordData[i] = uint32(v.Int64)
						} else {
							panic("invalid byte sub type")
						}
					}
				} else {
					recordData[i] = nil
				}
			case *sql.NullInt32:
				v := *values[i].(*sql.NullInt32)
				if v.Valid {
					if dirtyFields[i] == 0 {
						recordData[i] = v.Int32
					} else {
						if dirtyFields[i] == memory.T_I32 {
							recordData[i] = int32(v.Int32)
						} else if dirtyFields[i] == memory.T_U32 {
							recordData[i] = uint32(v.Int32)
						} else {
							panic("invalid byte sub type")
						}
					}
				} else {
					recordData[i] = nil
				}
			case *sql.NullInt16:
				v := *values[i].(*sql.NullInt16)
				if v.Valid {
					if dirtyFields[i] == 0 {
						recordData[i] = v.Int16
					} else {
						if dirtyFields[i] == memory.T_I16 {
							recordData[i] = int16(v.Int16)
						} else if dirtyFields[i] == memory.T_U16 {
							recordData[i] = uint16(v.Int16)
						} else {
							panic("invalid byte sub type")
						}
					}
				} else {
					recordData[i] = nil
				}
			case *sql.NullFloat64:
				v := *values[i].(*sql.NullFloat64)
				if v.Valid {
					if dirtyFields[i] == memory.T_F32 {
						recordData[i] = float32(v.Float64)
					} else if dirtyFields[i] == memory.T_F64 {
						recordData[i] = v.Float64
					}
				} else {
					recordData[i] = nil
				}

			case *time.Time:
				recordData[i] = (*values[i].(*time.Time)).UnixMilli()

			case *sql.NullByte:
				v := *values[i].(*sql.NullByte)
				if v.Valid {
					if dirtyFields[i] == 0 {
						recordData[i] = v.Byte
					} else {
						if dirtyFields[i] == memory.T_I8 {
							recordData[i] = int8(v.Byte)
						} else if dirtyFields[i] == memory.T_U8 {
							recordData[i] = uint8(v.Byte)
						} else {
							panic("invalid byte sub type")
						}
					}

				} else {
					recordData[i] = nil
				}

			case *any:
				panic(fmt.Sprintf("Unrecgnised scan type %T", values[i]))

			}
		}

		tlv := memory.CreateTLV(memory.DT_LIST, memory.T_NULL, memory.T_NULL, recordData)
		rData = append(rData, tlv)

		//var rawValue = *(values[0].(*interface{}))

	}

	return rData, core.MkSuccess(0)
}

func (ego *MysqlConnection) ConnectionTest() int32 {
	err := ego._conn.Ping()
	ego._lastConnTestTs = time.Now().UnixMilli()
	if err != nil {
		ego._connTestFailedCount++
		return core.MkErr(core.EC_PING_DB_FAILED, 1)
	} else {
		ego._connTestFailedCount = 0
		return core.MkSuccess(0)
	}
}

func (ego *MysqlConnection) Index() int {
	return ego._index
}

func (ego *MysqlConnection) TableNames() []string {
	return ego._tableNames

}

func (ego *MysqlConnection) CreateDataBase(name string, chaset string, ci string) int32 {
	ddlStr := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s, DEFAULT CHARACTER SET %s, DEFAULT COLLATE %s", name, chaset, ci)
	_, err := ego._conn.Exec(ddlStr)
	if err != nil {
		return core.MkErr(core.EC_CREATE_DB_FAILED, 1)
	}
	return core.MkSuccess(0)
}

func (ego *MysqlConnection) Connect() int32 {
	var nTries uint16 = 0
	for {
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s&parseTime=true", ego._cfgPool.Username, ego._cfgPool.Password, ego._cfgPool.IPV4Addr, ego._cfgPool.TcpPort, ego._cfgPool.DB, ego._cfgPool.ConnParam)
		var err error = nil
		ego._conn, err = sql.Open("mysql", connStr)
		if err == nil {
			rc := ego.ConnectionTest()
			if core.Err(rc) {
				logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Failed:(%s)]", ego._index, connStr)
			} else {
				ego._conn.SetMaxOpenConns(1)
				ego._conn.SetMaxIdleConns(1)
				ego._conn.SetConnMaxIdleTime(time.Second * time.Duration(ego._cfgPool.KeepAlive))

				logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Success:(%s)]", ego._index, connStr)
				return core.MkSuccess(0)
			}
		} else {
			logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Failed:(%s) Maybe Syntax?]", ego._index, connStr)
		}

		time.Sleep(1000 * time.Millisecond)
		nTries++
		if ego._cfgPool.MaxTries > 0 && nTries >= ego._cfgPool.MaxTries {
			logging.Log(core.LL_INFO, "Connect_%d to db Failed and Retry reached Max Time %d", ego._index, ego._cfgPool.MaxTries)
			return core.MkErr(core.EC_CONNECT_DB_FAILED, 1)
		}
	}

}

func createMySqlConnectionHandles(cfgPool *config.DBConnectionPoolConfig, idx int) (IDBConnection, int32) {
	mc := MysqlConnection{
		_index:               idx,
		_cfgPool:             cfgPool,
		_lastConnTestTs:      0,
		_goContext:           nil,
		_cancelFunc:          nil,
		_transcation:         nil,
		_connTestFailedCount: 0,
		_tableNames:          make([]string, 0),
	}

	for i := 0; i < len(cfgPool.Connections[idx].Tables); i++ {
		tableNames, rc := ParseTableNameConfig(cfgPool.Connections[idx].Tables[i])
		if core.Err(rc) {
			return nil, core.MkErr(rc, 1)
		}

		for j := 0; j < len(tableNames); j++ {
			mc._tableNames = append(mc._tableNames, tableNames[j])
		}
	}

	return &mc, core.MkSuccess(0)
}
