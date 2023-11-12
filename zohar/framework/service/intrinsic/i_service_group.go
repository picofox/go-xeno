package intrinsic

type IServiceGroup interface {
	Name() string
	AddService(IService) int32
	Initialize() int32
	Start() int32
	Stop() int32
	Finalize() int32
}
