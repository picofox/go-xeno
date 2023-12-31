package logging

import "xeno/zohar/core"

const (
	LINE_HEADER_ELEM_DATE  = 0x1
	LINE_HEADER_ELEM_TIME  = 0x2
	LINE_HEADER_ELEM_NANO  = 0x4
	LINE_HEADER_ELEM_MILLI = 0x8
	LINE_HEADER_ELEM_MICRO = 0x10
	LINE_HEADER_ELEM_TS    = 0x20
	LINE_HEADER_ELEM_LV    = 0x40
	LINE_HEADER_ELEM_PID   = 0x80
	LINE_HEADER_ELEM_GOID  = 0x100
	LINE_TAILER_ELEM_LPOS  = 0x200
	LINE_TAILER_ELEM_SPOS  = 0x400
)

var LEVEL_NAMES []string = []string{
	"[SYS]", "[FTL]", "[ERR]", "[WRN]", "[INF]", "[DBG]",
}

func GetLogLevelName(idx int) string {
	if idx < 0 || idx > core.LL_DEBUG {
		return "[NA]"
	}
	return LEVEL_NAMES[idx]
}

type ILogger interface {
	Name() string
	SetLevel(lv int)
	Log(int, string, ...any)
	LogRaw(int, string, ...any)
	LogFixedWidth(int, int, bool, string, string, ...any)
}
