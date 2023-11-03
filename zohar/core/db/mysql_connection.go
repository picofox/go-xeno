package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/logging"
)

type MysqlConnection struct {
	_index               int
	_connTestFailedCount int
	_lastConnTestTs      int64
	_conn                *sql.DB

	_cfgPool    *config.DBConnectionPoolConfig
	_tableNames []string
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
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", ego._cfgPool.Username, ego._cfgPool.Password, ego._cfgPool.IPV4Addr, ego._cfgPool.TcpPort, ego._cfgPool.DB, ego._cfgPool.ConnParam)
		var err error = nil
		ego._conn, err = sql.Open("mysql", connStr)
		if err == nil {
			rc := ego.ConnectionTest()
			if core.Err(rc) {
				logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Failed:(%s)]", ego._index, connStr)
			} else {
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
		_index:      idx,
		_cfgPool:    cfgPool,
		_tableNames: make([]string, 0),
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
