package intrinsic

import (
	"xeno/zohar/core"
	"xeno/zohar/core/datatype"
)

type ServiceCommon struct {
	_state uint8
}

func (ego *ServiceCommon) BeginInitializing() int32 {
	if ego._state != datatype.Uninitialized {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndInitializing() int32 {
	ego._state = datatype.Initializing
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginInitialized() int32 {
	if ego._state != datatype.Initializing {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndInitialized() int32 {
	ego._state = datatype.Initialized
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginStarting() int32 {
	if ego._state != datatype.Initialized {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndStarting() int32 {
	ego._state = datatype.Starting
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginStarted() int32 {
	if ego._state != datatype.Starting {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndStarted() int32 {
	ego._state = datatype.Started
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginSuspending() int32 {
	if ego._state != datatype.Started {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndSuspending() int32 {
	ego._state = datatype.Suspending
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginSuspended() int32 {
	if ego._state != datatype.Suspending {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndSuspended() int32 {
	ego._state = datatype.Suspended
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginStopping() int32 {
	if ego._state != datatype.Started && ego._state != datatype.Suspended {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndStopping() int32 {
	ego._state = datatype.Stopping
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginStopped() int32 {
	if ego._state != datatype.Stopping {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndStopped() int32 {
	ego._state = datatype.Stopped
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginFinalizing() int32 {
	if ego._state != datatype.Stopped {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndFinalizing() int32 {
	ego._state = datatype.Finalizing
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) BeginUninitialized() int32 {
	if ego._state != datatype.Finalizing {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	return core.MkSuccess(0)
}

func (ego *ServiceCommon) EndUninitialized() int32 {
	ego._state = datatype.Uninitialized
	return core.MkSuccess(0)
}
