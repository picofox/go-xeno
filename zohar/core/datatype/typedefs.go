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

func StateCodeToString(c int8) string {
	code := c & 0x7F
	isErr := (c>>7)&1 == 0
	if isErr {
		if code > Finalizing {
			return "Err:" + "NAState"
		}
		return "Err:" + sStateToString[code]
	}

	if code > Finalizing {
		return "NAState"
	}
	return sStateToString[code]

}
