package transcomm

<<<<<<< HEAD:zohar/core/inet/transcomm/i_client_codec_handler.go
type IClientCodecHandler interface {
	OnReceive(*TCPClientConnection, any, int64, any) (int32, any, int64, any)
	//Inbound([]IServerCodecHandler, int, *TCPServerConnection, any, any) int32
=======
type IClientHandler interface {
	OnReceive(*TCPClientConnection) (any, int32)
>>>>>>> 7e43a4fc9ab7e9f565922f2bdc9631781a5da39c:zohar/core/inet/transcomm/i_client_handler.go

	OnSend(*TCPClientConnection, any, bool) int32

	Clear()
}
