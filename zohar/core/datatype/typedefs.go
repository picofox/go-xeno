package datatype

const (
	INT32_MAX = 0x7FFFFFFF
)

const (
	TASK_EXEC_EXECUTOR_POOL   = uint8(0)
	TASK_EXEC_CURRENT_ROUTINE = uint8(1)
	TASK_EXEC_NEO_ROUTINE     = uint8(2)
)

type TaskFuncType = func(any) int32
