package db

import (
	"xeno/zohar/core"
	"xeno/zohar/core/config"
)

type ConnectionPool struct {
	_name              string
	_dbConnections     []IDBConnection
	_config            config.DBConnectionPoolConfig
	_table2Connections map[string]uint16
}

func (ego *ConnectionPool) GetConnection(idx int) IDBConnection {
	if idx >= len(ego._dbConnections) {
		return nil
	}
	return ego._dbConnections[idx]
}

func (ego *ConnectionPool) ConnectDatabase() int32 {
	for i := 0; i < len(ego._dbConnections); i++ {
		rc := ego._dbConnections[i].Connect()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func NeoConnectionPool(name string, cfgPool *config.DBConnectionPoolConfig) (*ConnectionPool, int32) {
	sz := len(cfgPool.Connections)
	cp := &ConnectionPool{
		_name:              name,
		_dbConnections:     make([]IDBConnection, sz),
		_config:            *cfgPool,
		_table2Connections: make(map[string]uint16),
	}
	if cfgPool.Type >= DB_TYPE_COUNT {
		return nil, core.MkErr(core.EC_INDEX_OOB, 1)
	}

	var rc int32 = 0
	for i := 0; i < sz; i++ {
		cp._dbConnections[i], rc = createDBConnectionHandles[cfgPool.Type](cfgPool, i)
		if core.Err(rc) {
			return nil, rc
		}

		tableNames := cp._dbConnections[i].TableNames()
		nLen := len(tableNames)
		for j := 0; j < nLen; j++ {
			tn := tableNames[j]
			_, ok := cp._table2Connections[tn]
			if ok {
				return nil, core.MkErr(core.EC_DIR_ALREADY_EXIST, 1)
			}
			cp._table2Connections[tn] = uint16(i)
		}

	}

	return cp, core.MkSuccess(0)
}

var createDBConnectionHandles = [DB_TYPE_COUNT]func(cfgPool *config.DBConnectionPoolConfig, idx int) (IDBConnection, int32){
	createMySqlConnectionHandles,
}
