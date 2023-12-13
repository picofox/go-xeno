package transcomm

type IClientHandler interface {
	OnReceive(*TCPClientConnection) (any, int32)

	OnSend(*TCPClientConnection, any, bool) int32

	Clear()
}
