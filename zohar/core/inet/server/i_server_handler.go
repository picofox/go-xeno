package server

type IServerHandler interface {
	OnReceive(*TcpServerConnection, any, any) (int32, any, any)
}
