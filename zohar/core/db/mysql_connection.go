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
	_index int
	_conn  *sql.DB

	_cfgPool    *config.DBConnectionPoolConfig
	_tableNames []string
}

func (ego *MysqlConnection) Index() int {
	return ego._index
}

func (ego *MysqlConnection) TableNames() []string {
	return ego._tableNames

}

func (ego *MysqlConnection) Connect() int32 {
	var nTries uint16 = 0
	for {
		logging.Log(core.LL_INFO, "Connecting (%d) to db ...", ego._index)
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", ego._cfgPool.Username, ego._cfgPool.Password, ego._cfgPool.IPV4Addr, ego._cfgPool.TcpPort, ego._cfgPool.DB, ego._cfgPool.ConnParam)
		var err error = nil
		ego._conn, err = sql.Open("mysql", connStr)
		if err == nil {
			q := fmt.Sprintf("use %s", ego._cfgPool.DB)
			_, err = ego._conn.Exec(q)
			if err != nil {
				logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Failed:(%s not exist)]", ego._index, ego._cfgPool.DB)
				return core.MkErr(core.EC_CONNECT_DB_FAILED, 1)
			}
			logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Success:(%s)]", ego._index, connStr)

			return core.MkSuccess(0)
		}
		logging.Log(core.LL_INFO, "Connect_%d to db \t\t\t\t[Failed:(%s)]", ego._index, connStr)

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
