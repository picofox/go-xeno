package server

type TcpServerConnection struct {
	_pipeline []IServerHandler
}
