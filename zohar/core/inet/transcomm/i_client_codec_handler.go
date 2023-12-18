package transcomm

type IClientCodecHandler interface {
	OnReceive(*TCPClientConnection) (any, int32)
	OnSend(*TCPClientConnection, any, bool) int32
	Pulse(conn IConnection, nowTs int64)
	Reset()
}
