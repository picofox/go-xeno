package sched

type Task struct {
	_executor uint8
	_function TaskFuncType
	_args     any
}

func (ego *Task) Executor() uint8 {
	return ego._executor
}

func (ego *Task) Function() TaskFuncType {
	return ego._function
}

func (ego *Task) Arg() any {
	return ego._args
}

func (ego *Task) SetExecutor(exe uint8) {
	ego._executor = exe
}

func (ego *Task) SetFunction(f func(a any) int32) {
	ego._function = f
}

func (ego *Task) SetArg(a any) {
	ego._args = a
}

type TaskFuncType = func(any) int32

//var sExecMethodsArr = [3]func(task *Task) int32{
//	func(t *Task) int32 {
//		concurrent.GetDefaultGoExecutorPool().PostTask(t.Function(), t.Arg())
//		return core.MkSuccess(0)
//	},
//	func(t *Task) int32 {
//		t.Function()(t.Arg())
//		return core.MkSuccess(0)
//	},
//	func(t *Task) int32 {
//		go t.Function()(t.Arg())
//		return core.MkSuccess(0)
//	},
//}
//
//func (ego *Task) Execute(a any) {
//
//}

func NeoTask(exe uint8, f func(a any) int32, a any) *Task {
	return &Task{
		_executor: exe,
		_function: f,
		_args:     a,
	}
}
