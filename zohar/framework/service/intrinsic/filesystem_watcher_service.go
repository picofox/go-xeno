package intrinsic

import (
	"path/filepath"
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/fs"
	"xeno/zohar/core/process"
)

type FileSystemWatcherService struct {
	ServiceCommon
	_fsw            *fs.FileSystemWatcher
	_config         *intrinsic.FileSystemWatcherServiceConfig
	_serviceManager *ServiceManager
	_handler        datatype.TaskFuncType
}

func (ego *FileSystemWatcherService) RegisterHandler(e uint8, f datatype.TaskFuncType) {
	ego._fsw.RegisterHandler(e, f)
}

func (ego *FileSystemWatcherService) AddWatch(dir string) {
	if filepath.IsAbs(dir) {
		ego._fsw.AddWatch(dir)
	} else {
		dstr := process.ComposePath(false, dir, true)
		ego._fsw.AddWatch(dstr)
	}

	ego._fsw.AddWatch(dir)
}

func (ego *FileSystemWatcherService) RemoveWatch(dir string) {
	if filepath.IsAbs(dir) {
		ego._fsw.RemoveWatch(dir)
	} else {
		dstr := process.ComposePath(false, dir, true)
		ego._fsw.RemoveWatch(dstr)
	}

	ego._fsw.RemoveWatch(dir)
}

func (ego *FileSystemWatcherService) Initialize() int32 {
	rc := ego.SetInitializeState()
	if core.Err(rc) {
		return rc
	}

	for _, d := range ego._config.Dirs {
		if filepath.IsAbs(d) {
			ego._fsw.AddWatch(d)
		} else {
			dstr := process.ComposePath(false, d, true)
			ego._fsw.AddWatch(dstr)
		}

	}

	ego.SetInitializeStateResult(true)
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherService) Finalize() int32 {
	rc := ego.SetFinalizeState()
	if core.Err(rc) {
		return rc
	}

	for _, d := range ego._config.Dirs {
		ego._fsw.RemoveWatch(d)
	}

	ego.SetFinalizeStateResult(true)
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherService) Start() int32 {
	rc := ego.SetStartState()
	if core.Err(rc) {
		return rc
	}

	ego._fsw.Start()

	ego.SetStartStateResult(true)
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcherService) Stop() int32 {
	rc := ego.SetStopState()
	if core.Err(rc) {
		return rc
	}
	ego._fsw.Stop()
	ego.SetStopStateResult(true)
	return core.MkSuccess(0)
}

func NeoFileSystemWatcherService(sm *ServiceManager, cfg *intrinsic.FileSystemWatcherServiceConfig) *FileSystemWatcherService {
	s := FileSystemWatcherService{
		ServiceCommon: ServiceCommon{
			_stateCode: datatype.StateCode(0),
		},
		_fsw:            fs.NeoFileSystemWatcher(),
		_serviceManager: sm,
		_config:         cfg,
	}

	return &s
}
