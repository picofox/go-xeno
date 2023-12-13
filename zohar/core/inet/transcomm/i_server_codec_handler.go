package transcomm

type IServerCodecHandler interface {
	OnReceive(*TCPServerConnection) (any, int32)
	//Inbound([]IServerCodecHandler, int, *TCPServerConnection, any, any) int32

	Reset()
}
