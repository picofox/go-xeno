package server

type IServerHandler interface {
	OnReceive(connection *TcpServerConnection, obj any) int32
}
