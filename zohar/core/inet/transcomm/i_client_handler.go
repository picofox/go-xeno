package transcomm

type IClientHandler interface {
	OnReceive(*TCPClientConnection, any, int64, any) (int32, any, int64, any)
	//Inbound([]IServerHandler, int, *TCPServerConnection, any, any) int32
}
