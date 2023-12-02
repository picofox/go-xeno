package server

const MAX_PACKET_SIZE = 32*1024 - O1L15O1T15_HEADER_SIZE
const MAX_BUFFER_MAX_CAPACITY = 32 * 1024

type IServerHandler interface {
	OnReceive(*TcpServerConnection, any, any) (int32, any, any)
}
