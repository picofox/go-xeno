package logger_adapter

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/logging"
)

var sLevelToLevelMap [4]int = [4]int{core.LL_FATAL, core.LL_ERR, core.LL_WARN, core.LL_INFO}

type GORMLoggerAdapter struct {
	_level  int
	_logger logging.ILogger
}

func (ego *GORMLoggerAdapter) LogMode(ll logger.LogLevel) logger.Interface {
	ego._level = int(ll)
	ego._logger.SetLevel(sLevelToLevelMap[ll])
	return ego
}
func (ego *GORMLoggerAdapter) Info(ctx context.Context, str string, arg ...interface{}) {
	ego._logger.Log(core.LL_INFO, str, arg...)
}
func (ego *GORMLoggerAdapter) Warn(ctx context.Context, str string, arg ...interface{}) {
	ego._logger.Log(core.LL_WARN, str, arg...)
}
func (ego *GORMLoggerAdapter) Error(ctx context.Context, str string, arg ...interface{}) {
	ego._logger.Log(core.LL_ERR, str, arg...)
}

func (ego *GORMLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if ego._level <= int(logger.Silent) {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && ego._level >= int(logger.Error) && (!errors.Is(err, gorm.ErrRecordNotFound)):
		sql, rows := fc()
		if rows == -1 {
			ego._logger.Log(core.LL_ERR, "Line:%s err:%s %f - %s ", utils.FileWithLineNum(), err.Error(), float64(elapsed.Nanoseconds())/1e6, sql)
		} else {
			ego._logger.Log(core.LL_ERR, "Line:%s err:%s %f %d %s ", utils.FileWithLineNum(), err.Error(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > time.Duration(10*time.Millisecond) && ego._level >= int(logger.Warn):
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= 10ms")
		if rows == -1 {
			ego._logger.Log(core.LL_WARN, "Line:%s SlowLog:%s %f - %s ", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, sql)
		} else {
			ego._logger.Log(core.LL_WARN, "Line:%s SlowLog:%s %f %d %s ", utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case ego._level == int(logger.Info):
		sql, rows := fc()
		if rows == -1 {
			ego._logger.Log(core.LL_WARN, "Line:%s %f - %s ", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, sql)
		} else {
			ego._logger.Log(core.LL_WARN, "Line:%s %f %d %s ", utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func NeoGORMLoggerAdapter(logger logging.ILogger) *GORMLoggerAdapter {
	return &GORMLoggerAdapter{
		_level:  4,
		_logger: logger,
	}

}

func NNN2() {

}
