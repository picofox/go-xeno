package intrinsic

import (
	"xeno/zohar/core/datatype"
)

type ServiceCommon struct {
	_stateCode datatype.StateCode
}

func (ego *ServiceCommon) SetInitializeState() int32 {
	return ego._stateCode.SetInitializeState()
}

func (ego *ServiceCommon) SetInitializeStateResult(ok bool) int32 {
	return ego._stateCode.SetInitializeStateResult(ok)
}

func (ego *ServiceCommon) SetStartState() int32 {
	return ego._stateCode.SetStartState()
}

func (ego *ServiceCommon) SetStartStateResult(ok bool) int32 {
	return ego._stateCode.SetStartStateResult(ok)
}

func (ego *ServiceCommon) SetSuspendState() int32 {
	return ego._stateCode.SetSuspendState()
}

func (ego *ServiceCommon) SetSuspendStateResult(ok bool) int32 {
	return ego._stateCode.SetSuspendStateResult(ok)
}

func (ego *ServiceCommon) SetStopState() int32 {
	return ego._stateCode.SetStopState()
}

func (ego *ServiceCommon) SetStopStateResult(ok bool) int32 {
	return ego._stateCode.SetStopStateResult(ok)
}

func (ego *ServiceCommon) SetFinalizeState() int32 {
	return ego._stateCode.SetFinalizeState()
}

func (ego *ServiceCommon) SetFinalizeStateResult(ok bool) int32 {
	return ego._stateCode.SetFinalizeStateResult(ok)
}
