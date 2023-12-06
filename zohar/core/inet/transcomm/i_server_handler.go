package transcomm

const MAX_BUFFER_MAX_CAPACITY = 32 * 1024
const MAX_PACKET_SIZE = MAX_BUFFER_MAX_CAPACITY - O1L15O1T15_HEADER_SIZE

type IServerHandler interface {
	OnReceive(*TCPServerConnection, any, int64, any) (int32, any, int64, any)
	//Inbound([]IServerHandler, int, *TCPServerConnection, any, any) int32
}
