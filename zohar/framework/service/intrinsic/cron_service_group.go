package intrinsic

import (
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
)

type CronServiceGroup struct {
	_name     string
	_services map[string]IService
	_config   *intrinsic.CronServiceGroupConfig
	_manager  *ServiceManager
}

func (ego *CronServiceGroup) AddService(key any, svc IService) int32 {
	_, ok := ego._services[key.(string)]
	if ok {
		return core.MkErr(core.EC_ELEMENT_EXIST, 0)
	}
	ego._services[key.(string)] = svc
	return core.MkSuccess(0)
}

func (ego *CronServiceGroup) Name() string {
	return ego._name
}

func (ego *CronServiceGroup) Initialize() int32 {
	for k, v := range ego._config.Params {
		loc := time.FixedZone(k, int(v.Offset))
		svc := NeoCronService(ego._manager, loc)
		if svc == nil {
			return core.MkErr(core.EC_NULL_VALUE, 1)
		}
		rc := svc.Initialize()
		if core.Err(rc) {
			return rc
		}
		ego.AddService(k, svc)
	}

	return core.MkSuccess(0)
}

func (ego *CronServiceGroup) Start() int32 {
	for _, s := range ego._services {
		rc := s.Start()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *CronServiceGroup) Stop() int32 {
	for _, s := range ego._services {
		rc := s.Stop()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *CronServiceGroup) Finalize() int32 {
	for _, s := range ego._services {
		rc := s.Finalize()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *CronServiceGroup) FindServiceByKey(key any) IService {
	svc, ok := ego._services[key.(string)]
	if ok {
		return svc
	}
	return nil
}

func NeoCronServiceGroup(sm *ServiceManager) *CronServiceGroup {
	sg := CronServiceGroup{
		_name:     "Cron",
		_services: make(map[string]IService),
		_config:   &intrinsic.GetIntrinsicConfig().IntrinsicService.Cron,
		_manager:  sm,
	}
	return &sg
}
