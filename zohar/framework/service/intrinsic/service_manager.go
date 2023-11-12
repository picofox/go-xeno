package intrinsic

import (
	"sync"
	"time"
	"xeno/zohar/core"
)

var _serviceManagerInstance ServiceManager
var once sync.Once

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

func (ego *ServiceManager) Wait() {
	for {
		if !ego._started {
			time.Sleep(100 * time.Millisecond)
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

func (ego *ServiceManager) AddService(name string, svc IService) int32 {
	g := ego.GetGroup(name)
	if g == nil {
		g = NeoTimerServiceGroup()
		ego.AddGroup(g)
	}
	return g.AddService(svc)
}

func GetServiceManager() *ServiceManager {
	once.Do(func() {
		_serviceManagerInstance = ServiceManager{
			_groups:  make(map[string]IServiceGroup),
			_started: true,
		}
	})
	return &_serviceManagerInstance
}
