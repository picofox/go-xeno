package transcomm

import (
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/logging"
)

type Poller struct {
	_logger logging.ILogger
}

func (ego *Poller) OnIncomingConnection(connection IConnection) {

}

func (ego *Poller) RegisterTCPServer(svr *TCPServer) {

}

func (ego *Poller) OnServerStart(svr *TCPServer) int32 {

	return core.MkSuccess(0)
}

func (ego *Poller) Start() int32 {

	return core.MkSuccess(0)
}

func (ego *Poller) Wait() int32 {

	return core.MkSuccess(0)
}

func (ego *Poller) Stop() int32 {
	return core.MkSuccess(0)
}

func (ego *Poller) NeoSubReactor() *SubReactor {
	sr := SubReactor{
		_poller: ego,
	}
	return &sr
}

func (ego *Poller) neoMainReactor() *MainReactor {
	mr := MainReactor{
		_poller: ego,
	}

	return &mr
}

func (ego *Poller) Initialize() int32 {

	return core.MkSuccess(0)
}

func (ego *Poller) Log(lv int, fmt string, arg ...any) {
	if ego._logger != nil {
		ego._logger.Log(lv, fmt, arg...)
	}
}

func (ego *Poller) LogFixedWidth(lv int, leftLen int, ok bool, failStr string, format string, arg ...any) {
	if ego._logger != nil {
		ego._logger.LogFixedWidth(lv, leftLen, ok, failStr, format, arg...)
	}
}

func NeoPoller() *Poller {
	p := Poller{
		_logger: logging.GetLoggerManager().GetDefaultLogger(),
	}

	return &p
}

var sDefaultPollerInstance *Poller
var sDefaultPollerInstanceOnce sync.Once

func GetDefaultPoller() *Poller {
	sDefaultPollerInstanceOnce.Do(
		func() {
			sDefaultPollerInstance = NeoPoller()
		})
	return sDefaultPollerInstance
}
