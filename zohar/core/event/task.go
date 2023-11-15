package event

import (
	"xeno/zohar/core/concurrent"
	"xeno/zohar/core/datatype"
)

type Task struct {
	_id       int64
	_executor uint8
	_function datatype.TaskFuncType
	_arg      any
}

func (ego *Task) Executor() uint8 {
	return ego._executor
}

func (ego *Task) Function() datatype.TaskFuncType {
	return ego._function
}

func (ego *Task) Arg() any {
	return ego._arg
}

func (ego *Task) SetExecutor(e uint8) {
	ego._executor = e
}

func (ego *Task) SetArg(a any) {
	ego._arg = a
}

var sTimerExecMethodsArr = [3]func(*Task){
	func(t *Task) {
		concurrent.GetDefaultGoExecutorPool().PostTask(t._function, t._arg)
	},
	func(t *Task) {
		t._function(t._arg)
	},
	func(t *Task) {
		go t._function(t._arg)
	},
}

func (ego *Task) Execute() {
	sTimerExecMethodsArr[ego._executor](ego)
}

func (ego *Task) ExecuteBy(e uint8) {
	sTimerExecMethodsArr[e](ego)
}

func (ego *Task) SetFunction(e datatype.TaskFuncType) {
	ego._function = e
}

//func (ego *Task) SetCancel() {
//	ego._flag = ego._flag | datatype.TaskCancel
//}
//
//func (ego *Task) IsCancelling() bool {
//	if ego._flag&datatype.TaskCancel != 0 {
//		return true
//	}
//	return false
//}

func NeoTask(id int64, e uint8, f datatype.TaskFuncType, a any) *Task {
	t := Task{
		_id:       id,
		_executor: e,
		_function: f,
		_arg:      a,
	}
	return &t
}
