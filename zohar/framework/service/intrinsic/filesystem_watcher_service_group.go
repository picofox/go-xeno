package intrinsic

import (
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
)

type FileSystemWatcherServiceGroup struct {
	_name     string
	_services map[string]IService
	_config   *intrinsic.FileSystemWatcherServiceGroupConfig
	_manager  *ServiceManager
}

func (ego *FileSystemWatcherServiceGroup) FindAnyServiceByKey(key any) IService {
	svc, ok := ego._services[key.(string)]
	if ok {
		return svc
	}
	if len(ego._services) > 0 {
		for _, v := range ego._services {
			return v
		}
	}
	return nil
}

func (ego *FileSystemWatcherServiceGroup) AddService(key any, svc IService) int32 {
	_, ok := ego._services[key.(string)]
	if ok {
		return core.MkErr(core.EC_ELEMENT_EXIST, 0)
	}
	ego._services[key.(string)] = svc
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherServiceGroup) Name() string {
	return ego._name
}

func (ego *FileSystemWatcherServiceGroup) Initialize() int32 {
	for k, v := range ego._config.Params {
		svc := NeoFileSystemWatcherService(ego._manager, &v)
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

func (ego *FileSystemWatcherServiceGroup) Start() int32 {
	for _, elem := range ego._services {
		rc := elem.Start()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherServiceGroup) Stop() int32 {
	for _, s := range ego._services {
		rc := s.Stop()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherServiceGroup) Finalize() int32 {
	for _, s := range ego._services {
		rc := s.Finalize()
		if core.Err(rc) {
			return rc
		}
	}
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherServiceGroup) FindServiceByKey(key any) IService {
	svc, ok := ego._services[key.(string)]
	if ok {
		return svc
	}
	return nil
}

func NeoFileSystemWatcherGroup(sm *ServiceManager) *FileSystemWatcherServiceGroup {
	sg := FileSystemWatcherServiceGroup{
		_name:     "FileSystemWatcher",
		_services: make(map[string]IService),
		_config:   &intrinsic.GetIntrinsicConfig().IntrinsicService.FileSystemWatcher,
		_manager:  sm,
	}
	return &sg
}
