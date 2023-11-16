package intrinsic

type IServiceGroup interface {
	Name() string
	AddService(any, IService) int32
	Initialize() int32
	Start() int32
	Stop() int32
	Finalize() int32
	FindServiceByKey(key any) IService
	FindAnyServiceByKey(key any) IService
}
