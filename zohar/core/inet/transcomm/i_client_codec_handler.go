package transcomm

type IClientCodecHandler interface {
	OnReceive(*TCPClientConnection, any, int64, any) (int32, any, int64, any)
	//Inbound([]IServerCodecHandler, int, *TCPServerConnection, any, any) int32

	OnSend(*TCPClientConnection, any, int64, bool) (int32, any, int64, bool)

	Clear()
}
