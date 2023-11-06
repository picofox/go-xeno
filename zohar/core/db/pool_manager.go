package db

import (
	"fmt"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/logging"
)

var poolManagerInstance *PoolManager
var pmOnce sync.Once

type PoolManager struct {
	_pools  map[string]*ConnectionPool
	_config config.DBConfig
}

func GetPoolManager() *PoolManager {
	pmOnce.Do(func() {
		poolManagerInstance = &PoolManager{
			_pools: make(map[string]*ConnectionPool),
		}
	})
	return poolManagerInstance
}

func (ego *PoolManager) GetPool(name string) *ConnectionPool {
	p, ok := ego._pools[name]
	if ok {
		return p
	}
	return nil
}

func (ego *PoolManager) Initialize(cfgDB *config.DBConfig) (int32, string) {
	ego._config = *cfgDB
	for k, v := range ego._config.Pools {
		pool, rc := NeoConnectionPool(k, v)
		if core.Err(rc) {
			return rc, fmt.Sprintf("Init DB Connection Pool (%s) Failed", k)
		}
		ego._pools[k] = pool
	}
	return core.MkSuccess(0), ""
}

func (ego *PoolManager) ConnectDatabase() int32 {
	for k, v := range ego._pools {
		logging.Log(core.LL_INFO, "Connecting Database For Pool: %s", k)
		rc := v.ConnectDatabase()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}
