package transcomm

import (
	"xeno/zohar/core/inet"
	"xeno/zohar/core/xplatform"
)

type TCPClientConnection struct {
}

func (T TCPClientConnection) OnIncomingData() int32 {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) Identifier() int64 {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) FileDescriptor() xplatform.FileDescriptor {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) String() string {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) PreStop() {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) RemoteEndPoint() *inet.IPV4EndPoint {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) LocalEndPoint() *inet.IPV4EndPoint {
	//TODO implement me
	panic("implement me")
}

var _ IConnection = TCPClientConnection{}
