package transcomm

import (
	"xeno/zohar/core/config/intrinsic"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/message_buffer"
)

type IConnection interface {
	OnWritable() int32
	OnIncomingData() int32
	OnPeerClosed() int32
	OnDisconnected() int32
	OnConnectingFailed() int32
	Identifier() int64
	String() string
	PreStop()
	RemoteEndPoint() *inet.IPV4EndPoint
	LocalEndPoint() *inet.IPV4EndPoint
	Type() int8
	ReactorIndex() uint32
	SetReactorIndex(uint32)
	SendMessage(msg message_buffer.INetMessage, bFlush bool) int32
	KeepAliveConfig() *intrinsic.KeepAliveConfig
}

const (
	CONNTYPE_TCP_SERVER = int8(0)
	CONNTYPE_TCP_CLIENT = int8(1)
)
