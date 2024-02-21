package transcomm

import (
	"reflect"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/chrono"
	"xeno/zohar/core/config"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
	"xeno/zohar/core/inet/message_buffer/messages"
	"xeno/zohar/core/logging"
	"xeno/zohar/core/mp"
)

type TCPClient struct {
	_name            string
	_config          *config.NetworkClientTCPConfig
	_connections     []*TCPClientConnection
	_logger          logging.ILogger
	_poller          *Poller
	_router          IClientMessageRouter
	_lastSendConnIdx int
}

func (ego *TCPClient) SendMessageWithConnection(idxConn int, msg message_buffer.INetMessage, bFlush bool) int32 {
	if ego._connections[idxConn] != nil {
		return ego._connections[idxConn].SendMessage(msg, bFlush)
	} else {
		ego.Log(core.LL_SYS, "Connection <idx:%d> Invalid", idxConn)
	}
	return core.MkErr(core.EC_NULL_VALUE, 1)
}

func (ego *TCPClient) SendMessage(msg message_buffer.INetMessage, bFlush bool) int32 {
	if len(ego._connections) == 1 {
		return ego._connections[0].SendMessage(msg, bFlush)
	}
	if ego._connections[ego._lastSendConnIdx] != nil {
		rc := ego._connections[ego._lastSendConnIdx].SendMessage(msg, bFlush)
		ego._lastSendConnIdx++
		if ego._lastSendConnIdx >= len(ego._connections) {
			ego._lastSendConnIdx = 0
		}
		return rc
	} else {
		startIdx := ego._lastSendConnIdx
		for ego._connections[ego._lastSendConnIdx] == nil {
			ego._lastSendConnIdx++
			if ego._lastSendConnIdx >= len(ego._connections) {
				ego._lastSendConnIdx = 0
			}
			if ego._lastSendConnIdx == startIdx {
				return core.MkErr(core.EC_NULL_VALUE, 1)
			}
		}
		rc := ego._connections[ego._lastSendConnIdx].SendMessage(msg, bFlush)
		if ego._lastSendConnIdx >= len(ego._connections) {
			ego._lastSendConnIdx = 0
		}
		return rc
	}
}

func (ego *TCPClient) OnIncomingMessage(conn *TCPClientConnection, message message_buffer.INetMessage) int32 {
	return ego._router.OnIncomingMessage(conn, message)
}

func (ego *TCPClient) OnPeerClosed(connection *TCPClientConnection) int32 {
	ego.Log(core.LL_SYS, "Connection Peer <%s> Closed.", connection.String())
	ego._connections[connection._index] = nil
	ego._poller.SubReactorEnded(connection.Identifier())
	return core.MkSuccess(0)
}

func (ego *TCPClient) OnDisconnected(connection *TCPClientConnection) int32 {
	ego.Log(core.LL_SYS, "Connection Peer <%s> Disconnected.", connection.String())
	ego._connections[connection._index] = nil
	ego._poller.SubReactorEnded(connection.Identifier())
	return core.MkSuccess(0)
}

func (ego *TCPClient) OnIOError(connection *TCPClientConnection) int32 {
	ego.Log(core.LL_SYS, "Connection IO <%s> Error.", connection.String())
	ego._connections[connection._index] = nil
	ego._poller.SubReactorEnded(connection.Identifier())
	return core.MkSuccess(0)
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

func (ego *TCPClient) Stop() int32 {
	for i := 0; i < len(ego._connections); i++ {
		ego._connections[i] = nil
		ego._poller.SubReactorEnded(ego._connections[i].Identifier())
		ego._connections[i].Close()
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

func (ego *TCPClient) OnKeepAliveMessage(conn *TCPClientConnection, message message_buffer.INetMessage) int32 {
	var pkam *messages.KeepAliveMessage = message.(*messages.KeepAliveMessage)
	if pkam.IsServer() {
		conn.SendMessage(message, true)
	} else {
		ts := chrono.GetRealTimeMilli()
		delta := ts - pkam.TimeStamp()
		conn.OnKeepAlive(ts, int32(delta))
	}
	return core.MkSuccess(0)
}

func (ego *TCPClient) OnProcTestMessage(conn *TCPClientConnection, message message_buffer.INetMessage) int32 {
	var m *messages.ProcTestMessage = message.(*messages.ProcTestMessage)
	if m.IsServer {
		conn.SendMessage(message, true)
	} else {
		if core.Err(m.Validate()) {
			panic("invalid msg")
		}
	}
	return core.MkSuccess(0)
}

func NeoTCPClient(name string, poller *Poller, config *config.NetworkClientTCPConfig, logger logging.ILogger) *TCPClient {
	c := &TCPClient{
		_name:   name,
		_config: config,
		_logger: logger,
		_poller: poller,
		_router: nil,
	}

	var output = make([]reflect.Value, 0, 1)
	rc := mp.GetDefaultObjectInvoker().Invoke(&output, "smh", "Reflect"+config.Codec, c)
	if core.Err(rc) {
		return nil
	}
	h := output[0].Interface().(IClientMessageRouter)
	c._router = h

	c._router.RegisterHandler(messages.INTERNAL_MSG_GRP_TYPE, messages.KEEP_ALIVE_MESSAGE_ID, c.OnKeepAliveMessage)
	c._router.RegisterHandler(messages.INTERNAL_MSG_GRP_TYPE, messages.PROC_TEST_MESSAGE_ID, c.OnProcTestMessage)

	return c
}
