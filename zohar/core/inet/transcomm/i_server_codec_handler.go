package transcomm

type IServerCodecHandler interface {
	OnReceive(*TCPServerConnection) (any, int32)
	OnSend(*TCPServerConnection, any, bool) int32

	Reset()
}
