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
		rc := ego._connection.OnIncomingData()
		if core.Err(rc) {
			et, _ := core.ExErr(rc)
			if et == core.EC_TRY_AGAIN {

			} else if et == core.EC_EOF {
				ego._connection.OnPeerClosed()
				runtime.Goexit()
			} else {
				ego._connection.OnDisconnected()
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
