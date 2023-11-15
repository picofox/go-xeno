package intrinsic

import (
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/logging"
)

var _serviceManagerInstance ServiceManager
var _serviceManagerOnce sync.Once

type ServiceManager struct {
	_groups    map[string]IServiceGroup
	_waitGroup sync.WaitGroup
	_started   bool
}

func (ego *ServiceManager) Start() int32 {
	for _, v := range ego._groups {
		rc := v.Start()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *ServiceManager) Stop() int32 {
	for _, v := range ego._groups {
		rc := v.Stop()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *ServiceManager) addRef() {
	ego._waitGroup.Add(1)
}
func (ego *ServiceManager) delRef() {
	ego._waitGroup.Done()
}

func (ego *ServiceManager) Initialize() int32 {
	cronCfg := &intrinsic.GetIntrinsicConfig().IntrinsicService
	if cronCfg == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	grp := NeoCronServiceGroup(ego)
	if grp == nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}

	rc := grp.Initialize()
	if core.Err(rc) {
		return rc
	}

	ego.AddGroup(grp)

	return core.MkSuccess(0)
}

func (ego *ServiceManager) Wait() {

	for {
		if !ego._started {
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

	ego._waitGroup.Wait()
	ego._started = false
}

func (ego *ServiceManager) GetGroup(name string) IServiceGroup {
	grp, ok := ego._groups[name]
	if !ok {
		return nil
	}
	return grp
}

func (ego *ServiceManager) AddGroup(grp IServiceGroup) int32 {
	g := ego.GetGroup(grp.Name())
	if g != nil {
		return core.MkErr(core.EC_NOOP, 1)
	}
	ego._groups[grp.Name()] = grp
	return core.MkSuccess(0)
}

func (ego *ServiceManager) AddService(name string, key any, svc IService) int32 {
	g := ego.GetGroup(name)
	if g == nil {
		g = NeoCronServiceGroup(ego)
		ego.AddGroup(g)
	}
	return g.AddService(key, svc)
}

func (ego *ServiceManager) AddCronTask(which string, spec string, cmd datatype.TaskFuncType, a any, executor uint8) int32 {
	grp := ego.GetGroup("Cron")
	if grp == nil {
		logging.Log(core.LL_ERR, "Cron Service Group Not Found")
		return core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}

	var cs *CronService = nil
	cs = grp.FindServiceByKey(which).(*CronService)
	return cs.AddCron(spec, cmd, a, executor)
}

func GetServiceManager() *ServiceManager {
	_serviceManagerOnce.Do(func() {
		_serviceManagerInstance = ServiceManager{
			_groups:  make(map[string]IServiceGroup),
			_started: true,
		}
	})
	return &_serviceManagerInstance
}
