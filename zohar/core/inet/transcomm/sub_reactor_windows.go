package transcomm

import (
	"runtime"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
)

type SubReactor struct {
	_poller         *Poller
	_commandChannel chan cms.ICMS
	_connection     IConnection
}

func (ego *SubReactor) loop() int32 {
	defer ego._poller._waitGroup.Done()
	for {
		ego._poller.Log(core.LL_DEBUG, "Reading socket.... <%s>", ego._connection.String())
		rc := ego._connection.OnIncomingData()
		if core.Err(rc) {
			if core.IsErrType(rc, core.EC_EOF) {
				ego._poller.Log(core.LL_DEBUG, "Connection Closed <%s>: remain %d", ego._connection.String(), ego._poller.SubReactorCount())
				ego._poller.SubReactorEnded(ego)
				runtime.Goexit()
			}
		}

		select {
		case m := <-ego._commandChannel:
			if m.Id() == cms.CMSID_FINALIZE {
				runtime.Goexit()
			}
		default:
		}
	}
}

func (ego *SubReactor) OnStart() {
	ego._poller.Log(core.LL_SYS, "Sub Reactor Starting")
	ego._poller._waitGroup.Add(1)
	go ego.loop()
}

func (ego *SubReactor) OnStop() {
	ego._poller.Log(core.LL_SYS, "Sub Reactor <%s> Stopping", ego._connection.String())
	ego._connection.PreStop()
	finCMS := cms.NeoFinalize()
	ego._commandChannel <- finCMS
}
