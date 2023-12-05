package datatype

import "xeno/zohar/core"

type StateCode int8

const (
	Uninitialized = int8(0)
	Initializing  = int8(1)
	Initialized   = int8(2)
	Starting      = int8(3)
	Started       = int8(4)
	Suspending    = int8(5)
	Suspended     = int8(6)
	Stopping      = int8(7)
	Stopped       = int8(8)
	Finalizing    = int8(9)
)

func (ego *StateCode) HasError() bool {
	if *ego < 0 || int8(*ego) > Finalizing {
		return true
	}
	return false
}

func MarkError(sc *StateCode) {
	v := sc.Code()
	uv := uint8(1 << 7)
	v |= int8(uv)
	*sc = StateCode(v)
}

func Reset(sc *StateCode) {
	*sc = StateCode(0)
}

func (ego *StateCode) SetCode(code int8) {
	*ego = StateCode(code)
}

func (ego *StateCode) Code() int8 {
	return int8(*ego) & 0x7f
}

func (ego *StateCode) String() string {
	return StateCodeToString(int8(*ego))
}

func (ego *StateCode) SetInitializeState() int32 {
	if ego.HasError() {
		return core.MkErr(core.EC_INVALID_STATE, 0)
	}

	if ego.Code() != Uninitialized {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}

	ego.SetCode(Initializing)

	return core.MkSuccess(0)
}

func (ego *StateCode) SetInitializeStateResult(ok bool) int32 {
	if ego.Code() != Initializing {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ok {
		ego.SetCode(Initialized)
	} else {
		MarkError(ego)
	}

	return core.MkSuccess(0)
}

func (ego *StateCode) SetStartState() int32 {
	if ego.HasError() {
		return core.MkErr(core.EC_INVALID_STATE, 0)
	}
	if ego.Code() != Initialized && ego.Code() != Suspended && ego.Code() != Stopped {
		if ego.Code() == Started {
			return core.MkErr(core.EC_NOOP, 0)
		}
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	ego.SetCode(Starting)
	return core.MkSuccess(0)
}

func (ego *StateCode) SetStartStateResult(ok bool) int32 {
	if ego.Code() != Starting {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ok {
		ego.SetCode(Started)
	} else {
		MarkError(ego)
	}

	return core.MkSuccess(0)
}

func (ego *StateCode) SetSuspendState() int32 {

	if ego.HasError() {
		return core.MkErr(core.EC_INVALID_STATE, 0)
	}
	if ego.Code() != Started {
		if ego.Code() == Suspended {
			return core.MkErr(core.EC_NOOP, 0)
		}
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	ego.SetCode(Suspending)
	return core.MkSuccess(0)
}

func (ego *StateCode) SetSuspendStateResult(ok bool) int32 {

	if ego.Code() != Suspending {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ok {
		ego.SetCode(Suspended)
	} else {
		MarkError(ego)
	}

	return core.MkSuccess(0)
}

func (ego *StateCode) SetStopState() int32 {
	if ego.HasError() {
		return core.MkErr(core.EC_INVALID_STATE, 0)
	}
	if ego.Code() != Started && ego.Code() != Suspended {
		if ego.Code() == Stopped {
			return core.MkErr(core.EC_NOOP, 0)
		}
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	ego.SetCode(Stopping)
	return core.MkSuccess(0)
}

func (ego *StateCode) SetStopStateResult(ok bool) int32 {
	if ego.Code() != Stopping {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ok {
		ego.SetCode(Stopped)
	} else {
		MarkError(ego)
	}

	return core.MkSuccess(0)
}

func (ego *StateCode) SetFinalizeState() int32 {
	if ego.HasError() {
		return core.MkErr(core.EC_INVALID_STATE, 0)
	}
	if ego.Code() != Stopped && ego.Code() != Initialized {
		if ego.Code() == Uninitialized {
			return core.MkErr(core.EC_NOOP, 0)
		}
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	ego.SetCode(Finalizing)
	return core.MkSuccess(0)
}

func (ego *StateCode) SetFinalizeStateResult(ok bool) int32 {
	if ego.Code() != Finalizing {
		return core.MkErr(core.EC_INVALID_STATE, 1)
	}
	if ok {
		ego.SetCode(Uninitialized)
	} else {
		MarkError(ego)
	}

	return core.MkSuccess(0)
}
