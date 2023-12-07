package transcomm

import (
	"runtime"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/cms"
)

type MainReactor struct {
	_poller         *Poller
	_listener       *ListenWrapper
	_commandChannel chan cms.ICMS
}

func (ego *MainReactor) loop() {
	defer ego._poller._waitGroup.Done()
	for {
		conn := ego._listener.Accept()
		if conn == nil {
			time.Sleep(1 * time.Second)
		} else {
			c := NeoTCPServerConnection(conn, ego._listener._server._config)
			ego._poller.OnIncomingConnection(c)
			ego._listener._server.OnIncomingConnection(c)
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

func (ego *MainReactor) OnStart() {
	ego._poller.Log(core.LL_SYS, "Main Reactor <%s> Starting", ego._listener._bindAddress.EndPointString())
	ego._poller._waitGroup.Add(1)
	go ego.loop()
}

func (ego *MainReactor) OnStop() {
	ego._poller.Log(core.LL_SYS, "Main Reactor <%s> Stopping", ego._listener._bindAddress.EndPointString())
	finCMS := cms.NeoFinalize()
	ego._listener.PreStrop()
	ego._commandChannel <- finCMS
}
