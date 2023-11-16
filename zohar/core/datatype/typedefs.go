package datatype

const (
	INT32_MAX = 0x7FFFFFFF
)

const (
	TASK_EXEC_EXECUTOR_POOL   = uint8(0)
	TASK_EXEC_CURRENT_ROUTINE = uint8(1)
	TASK_EXEC_NEO_ROUTINE     = uint8(2)
)

const TaskCancel = 1

type TaskFuncType = func(any) int32

var sStateToString []string = []string{
	"Uninitialized",
	"Initializing",
	"Initialized",
	"Starting",
	"Started",
	"Suspending",
	"Suspended",
	"Stopping",
	"Stopped",
	"Finalizing",
}

func StateCodeToString(c uint8) string {
	if c > Finalizing {
		return "UnknowState"
	}
	return sStateToString[c]
}

const (
	Uninitialized = uint8(0)
	Initializing  = uint8(1)
	Initialized   = uint8(2)
	Starting      = uint8(3)
	Started       = uint8(4)
	Suspending    = uint8(5)
	Suspended     = uint8(6)
	Stopping      = uint8(7)
	Stopped       = uint8(8)
	Finalizing    = uint8(9)
)
