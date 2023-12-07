package transcomm

import (
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/logging"
)

type TCPClient struct {
	_name        string
	_config      *config.NetworkClientTCPConfig
	_connections []*TCPClientConnection
	_logger      logging.ILogger
	_poller      *Poller
}

func (ego *TCPClient) Initialize() int32 {
	for idx, targetStr := range ego._config.ServerEndPoints {
		rAddr := inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, targetStr)
		for i := int32(0); i < ego._config.Count; i++ {
			c := NeoTCPClientConnection(idx, ego, rAddr)
			ego._connections = append(ego._connections, c)
		}
	}

	return core.MkSuccess(0)
}

func (ego *TCPClient) Start() int32 {
	for _, c := range ego._connections {
		for {
			rc := c.Connect()
			if core.Err(rc) {
				ego.Log(core.LL_ERR, "Connecting to %s Failed, Will Retry", c._remoteEndPoint.EndPointString())
			} else {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}

		ego._poller.OnIncomingConnection(c)
	}

	return core.MkSuccess(0)
}

func (ego *TCPClient) Log(lv int, fmt string, arg ...any) {
	if ego._logger != nil {
		ego._logger.Log(lv, fmt, arg...)
	}
}

func (ego *TCPClient) LogFixedWidth(lv int, leftLen int, ok bool, failStr string, format string, arg ...any) {
	if ego._logger != nil {
		ego._logger.LogFixedWidth(lv, leftLen, ok, failStr, format, arg...)
	}
}

func NeoTCPClient(name string, poller *Poller, config *config.NetworkClientTCPConfig, logger logging.ILogger) *TCPClient {
	return &TCPClient{
		_name:   name,
		_config: config,
		_logger: logger,
		_poller: poller,
	}
}
