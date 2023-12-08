package transcomm

import (
	"xeno/zohar/core/inet"
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
}

const (
	CONNTYPE_TCP_SERVER = int8(0)
	CONNTYPE_TCP_CLIENT = int8(1)
)

const O1L15O1T15_HEADER_SIZE = 4
const MAX_BUFFER_MAX_CAPACITY = 32 * 1024
const MAX_PACKET_BODY_SIZE = MAX_BUFFER_MAX_CAPACITY - O1L15O1T15_HEADER_SIZE
