package transcomm

import (
	"xeno/zohar/core"
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/inet/message_buffer/messages"
)

type KeepAlive struct {
	_config            *intrinsic.KeepAliveConfig
	_lastSendTimestamp int64
	_lastRecvTimeStamp int64
	_currentTries      int32
	_recentRTT         uint32
	_message           *messages.KeepAliveMessage
}

func NeoKeepAlive(config *intrinsic.KeepAliveConfig, isServer bool) *KeepAlive {
	ka := KeepAlive{
		_config:            config,
		_lastSendTimestamp: 0,
		_lastRecvTimeStamp: 0,
		_currentTries:      0,
		_message:           messages.NeoKeepAliveMessage(isServer),
	}
	return &ka
}

func (ego *KeepAlive) Reset() {
	ego._lastRecvTimeStamp = 0
	ego._lastSendTimestamp = 0
	ego._currentTries = 0
}

func (ego *KeepAlive) OnRoundTripBack(nowTs int64) int32 {
	ego._lastRecvTimeStamp = nowTs
	ego._currentTries = 0
	return int32(nowTs - ego._lastSendTimestamp)
}

func (ego *KeepAlive) Pulse(conn IConnection, nowTs int64) int32 {
	if ego._lastRecvTimeStamp != 0 {
		if nowTs-ego._lastRecvTimeStamp > int64(ego._config.IntervalMillis) {
			ego._message.SetTimeStamp(nowTs)
			ego._lastRecvTimeStamp = 0
			ego._lastSendTimestamp = nowTs
			ego._currentTries = 0
			//conn.Log(core.LL_DEBUG, "Send keepalive on conn %s", conn.String())
			rc := conn.SendMessage(ego._message, true)
			if core.Err(rc) {
				conn.Logger().Log(core.LL_ERR, "Send keepalive failed on conn %s", conn.String())
				return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
			} else {
				return core.MkSuccess(0)
			}
		} else {
			return core.MkSuccess(0)
		}
	} else {
		v := ego.isTimeout(nowTs)
		if v {
			if ego._currentTries >= ego._config.MaxTries {
				conn.Logger().Log(core.LL_ERR, "Keepalive timeout on conn %s", conn.String())
				return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
			}
			ego._currentTries++
			ego._message.SetTimeStamp(nowTs)
			ego._lastRecvTimeStamp = 0
			ego._lastSendTimestamp = nowTs
			//conn.Log(core.LL_DEBUG, "Send keepalive on conn %s", conn.String())
			rc := conn.SendMessage(ego._message, true)
			if core.Err(rc) {
				conn.Logger().Log(core.LL_ERR, "Re-Send keepalive failed (%s) on conn %s", core.ErrStr(rc), conn.String())
				return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
			} else {
				return core.MkErr(core.EC_TRY_AGAIN, 1)
			}
		} else {
			return core.MkSuccess(0)
		}
	}
}

func (ego *KeepAlive) IsValid(nowTs int64) int32 {
	v := ego.isTimeout(nowTs)
	if v {
		if ego._currentTries >= ego._config.MaxTries {
			return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
		}
		ego._currentTries++
		return core.MkErr(core.EC_TRY_AGAIN, 1)
	}
	return core.MkSuccess(0)
}

func (ego *KeepAlive) isTimeout(nowTs int64) bool {
	to := nowTs - ego._lastSendTimestamp
	if to >= int64(ego._config.TimeoutMillis) {
		return true
	}
	return false
}
