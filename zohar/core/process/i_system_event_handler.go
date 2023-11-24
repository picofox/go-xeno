package process

type ISystemEventHandler interface {
	OnLowMemory()
}
