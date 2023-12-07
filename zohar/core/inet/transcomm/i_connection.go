package transcomm

import (
	"xeno/zohar/core/inet"
	"xeno/zohar/core/xplatform"
)

type IConnection interface {
	OnIncomingData() int32
	Identifier() int64
	FileDescriptor() xplatform.FileDescriptor
	String() string
	PreStop()
	RemoteEndPoint() *inet.IPV4EndPoint
	LocalEndPoint() *inet.IPV4EndPoint
}
