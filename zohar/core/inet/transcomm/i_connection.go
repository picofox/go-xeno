package transcomm

type IConnection interface {
	OnIncomingData()
	FileDescriptor() int
}
