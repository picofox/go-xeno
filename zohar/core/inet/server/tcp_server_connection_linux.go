package server

type TcpServerConnection struct {
	_buffer   []byte
	_length   int64
	_pipeline []IServerHandler
}

func (ego *TcpServerConnection) BufferLength() int64 {
	return ego._length
}

func (ego *TcpServerConnection) BufferCapacity() int64 {
	return 32768 + 4
}
