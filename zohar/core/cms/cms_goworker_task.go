package cms

type GoWorkerTask struct {
	_cmsid     uint32
	_procedure func(any)
	_object    any
}

func (ego *GoWorkerTask) Id() uint32 {
	return ego._cmsid
}

func (ego *GoWorkerTask) Exec() {
	if ego._procedure != nil {
		ego._procedure(ego._object)
	}
}

func NeoCMSGoWorkerTask(proc func(any), obj any) *GoWorkerTask {
	return &GoWorkerTask{
		_cmsid:     CMSID_GOWORKER_TASK,
		_procedure: proc,
		_object:    obj,
	}
}
