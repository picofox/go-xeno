package transcomm

import (
	"xeno/zohar/core/inet"
)

type IConnection interface {
	OnWritable() int32
	OnIncomingData() int32
	Identifier() int64
	String() string
	PreStop()
	RemoteEndPoint() *inet.IPV4EndPoint
	LocalEndPoint() *inet.IPV4EndPoint
	Type() int8
}

const (
	CONNTYPE_TCP_SERVER = int8(0)
	CONNTYPE_TCP_CLIENT = int8(1)
)
