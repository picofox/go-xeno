package intrinsic

import "xeno/zohar/core"

type TimerServiceGroup struct {
	_name     string
	_services []IService
}

func (ego *TimerServiceGroup) AddService(svc IService) int32 {
	ego._services = append(ego._services, svc)
	return core.MkSuccess(0)
}

func (ego *TimerServiceGroup) Name() string {
	return ego._name
}

func (ego *TimerServiceGroup) Initialize() int32 {
	for _, s := range ego._services {
		rc := s.Initialize()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *TimerServiceGroup) Start() int32 {
	for _, s := range ego._services {
		rc := s.Start()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *TimerServiceGroup) Stop() int32 {
	for _, s := range ego._services {
		rc := s.Stop()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *TimerServiceGroup) Finalize() int32 {
	for _, s := range ego._services {
		rc := s.Finalize()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func NeoTimerServiceGroup() *TimerServiceGroup {
	sg := TimerServiceGroup{
		_name:     "TimerService",
		_services: make([]IService, 0),
	}
	return &sg
}
