package cms

type Finalize struct {
	_cmsid uint32
}

func (ego *Finalize) Id() uint32 {
	return ego._cmsid
}

func NeoFinalize() *Finalize {
	return &Finalize{
		_cmsid: CMSID_FINALIZE,
	}
}
