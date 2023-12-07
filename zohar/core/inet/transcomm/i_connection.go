package transcomm

import (
	"xeno/zohar/core/inet"
	"xeno/zohar/core/xplatform"
)

type IConnection interface {
	OnWritable() int32
	OnIncomingData() int32
	Identifier() int64
	FileDescriptor() xplatform.FileDescriptor
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
