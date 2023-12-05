package transcomm

const MAX_BUFFER_MAX_CAPACITY = 32*1024 + O1L15O1T15_HEADER_SIZE
const MAX_PACKET_SIZE = MAX_BUFFER_MAX_CAPACITY - O1L15O1T15_HEADER_SIZE

type IServerHandler interface {
	OnReceive(*TcpServerConnection, any, any) (int32, any, any)
}
