package transcomm

type IServerCodecHandler interface {
	OnReceive(*TCPServerConnection) (any, int32)
	OnSend(*TCPServerConnection, any, bool) int32
	Pulse(conn IConnection, nowTs int64)
	OnKeepAlive(ts int64, delta int32)
	Reset()
}
