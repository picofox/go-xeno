package intrinsic

type IService interface {
	Initialize() int32
	Start() int32
	Stop() int32
	Finalize() int32
}
