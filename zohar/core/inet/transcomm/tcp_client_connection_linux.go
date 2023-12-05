package transcomm

type TCPClientConnection struct {
}

func (T TCPClientConnection) OnIncomingData() {
	//TODO implement me
	panic("implement me")
}

func (T TCPClientConnection) FileDescriptor() int {
	//TODO implement me
	panic("implement me")
}

var _ IConnection = TCPClientConnection{}
