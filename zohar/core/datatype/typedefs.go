package datatype

var EmptyByteSlice []byte = make([]byte, 0)

const (
	INT32_MAX  = 0x7FFFFFFF
	INT64_MAX  = 1<<63 - 1
	UINT16_MAX = 0xFFFF

	UINT16_CAPACITY = 0x10000

	INT64_SIZE    = 8
	INT32_SIZE    = 4
	INT16_SIZE    = 2
	INT8_SIZE     = 1
	FLOAT32_SIZE  = 4
	FLOAT64_SIZE  = 8
	BYTEBOOL_SIZE = 1
)

const (
	TASK_EXEC_EXECUTOR_POOL   = uint8(0)
	TASK_EXEC_CURRENT_ROUTINE = uint8(1)
	TASK_EXEC_NEO_ROUTINE     = uint8(2)
	TASK_EXEC_OVERRIDE        = uint8(3)
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

// Boolean to int.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
