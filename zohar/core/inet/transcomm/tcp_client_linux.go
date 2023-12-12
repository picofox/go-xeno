package transcomm

import (
	"fmt"
	"sync/atomic"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/logging"
)

type TCPClient struct {
	_name        string
	_config      *config.NetworkClientTCPConfig
	_connections []*TCPClientConnection
	_logger      logging.ILogger
	_poller      *Poller
	_index       atomic.Int32
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

func (ego *TCPClient) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	idx := ego._index.Add(1)
	idx = idx % int32(len(ego._connections))
	return ego._connections[idx].SendMessage(msg, bFlush)

}

func (ego *TCPClient) OnIncomingMessage(conn *TCPClientConnection, message message_buffer.INetMessage) int32 {
	fmt.Printf("Got msg [%v] \n", message)
	return core.MkSuccess(0)
}

func (ego *TCPClient) Start() int32 {
	for _, c := range ego._connections {
		rc := c.Connect()
		if core.Err(rc) {
			return rc
		}
	}

	var allReady bool = true
	for {
		allReady = true
		for _, c := range ego._connections {
			if !c._isConnected {
				allReady = false
			}
		}
		if allReady {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	ego.Log(core.LL_SYS, "All Connection is Connected.")

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
	c := &TCPClient{
		_name:   name,
		_config: config,
		_logger: logger,
		_poller: poller,
	}

	c._index.Store(0)

	return c
}
