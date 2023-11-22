package fs

import (
	"github.com/fsnotify/fsnotify"
	"runtime"
	"sync"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
	"xeno/zohar/core/event"
	"xeno/zohar/core/logging"
)

const (
	FS_WATCH_HAS_CREATE = 0x1
	FS_WATCH_HAS_WRITE  = 0x2
	FS_WATCH_HAS_REMOVE = 0x4
	FS_WATCH_HAS_RENAME = 0x8
	FS_WATCH_HAS_CHMOD  = 0x10
)

type FileSystemWatcher struct {
	_watcher *fsnotify.Watcher
	_handler *event.Task
	_lock    sync.RWMutex
}

func (ego *FileSystemWatcher) RegisterHandler(e uint8, f datatype.TaskFuncType) {
	arg := make([]any, 2)
	ego._handler = event.NeoTask(e, f, arg)
}

func (ego *FileSystemWatcher) AddWatch(path string) {
	ego._watcher.Add(path)
}

func (ego *FileSystemWatcher) RemoveWatch(path string) {
	ego._watcher.Remove(path)
}

func (ego *FileSystemWatcher) loop() {
	for {
		select {
		case ev := <-ego._watcher.Events:
			{
				if ev.Op&fsnotify.Create == fsnotify.Create {
					//logging.Log(core.LL_DEBUG, "file created <%s> (%s)", ev.Name, ev.String())
					if ego._handler.Function() != nil {
						ego._handler.Arg().([]any)[0] = FS_WATCH_HAS_CREATE
						ego._handler.Arg().([]any)[1] = ev.Name
						ego._handler.Execute()
					}
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					if ego._handler.Function() != nil {
						ego._handler.Arg().([]any)[0] = FS_WATCH_HAS_WRITE
						ego._handler.Arg().([]any)[1] = ev.Name
						ego._handler.Execute()
					}
					//logging.Log(core.LL_DEBUG, "file wrote <%s> (%s)", ev.Name, ev.String())
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					if ego._handler.Function() != nil {
						ego._handler.Arg().([]any)[0] = FS_WATCH_HAS_REMOVE
						ego._handler.Arg().([]any)[1] = ev.Name
						ego._handler.Execute()
					}

					//logging.Log(core.LL_DEBUG, "file removed <%s> (%s)", ev.Name, ev.String())
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					if ego._handler.Function() != nil {
						ego._handler.Arg().([]any)[0] = FS_WATCH_HAS_RENAME
						ego._handler.Arg().([]any)[1] = ev.Name
						ego._handler.Execute()
					}

					//logging.Log(core.LL_DEBUG, "file renamed <%s> (%s)", ev.Name, ev.String())
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					if ego._handler.Function() != nil {
						ego._handler.Arg().([]any)[0] = FS_WATCH_HAS_CHMOD
						ego._handler.Arg().([]any)[1] = ev.Name
						ego._handler.Execute()
					}
					//logging.Log(core.LL_DEBUG, "file chmod <%s> (%s)", ev.Name, ev.String())
				}

			}

		case err := <-ego._watcher.Errors:
			{
				if err != nil {
					logging.Log(core.LL_ERR, "FSW error: (%s)", err.Error())
					time.Sleep(3 * time.Second)
				} else {
					logging.Log(core.LL_SYS, "File Watcher Stopped")
					runtime.Goexit()
				}
			}
		}
	}
}

func (ego *FileSystemWatcher) Start() int32 {
	go ego.loop()
	return core.MkSuccess(0)
}

func (ego *FileSystemWatcher) Stop() int32 {
	ego._watcher.Close()
	return core.MkSuccess(0)
}

func NeoFileSystemWatcher() *FileSystemWatcher {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}
	fsw := FileSystemWatcher{
		_watcher: w,
		_handler: nil,
	}

	return &fsw
}
