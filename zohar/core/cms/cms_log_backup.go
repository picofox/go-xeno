package cms

type LogBackUp struct {
	_cmsid           uint32
	AbsFilePath      string
	AbsBackupDirPath string
	ZipFile          bool
}

func (ego *LogBackUp) Id() uint32 {
	return ego._cmsid
}

func NeoCMSLogBackUp(fpath string, dpath string, zip bool) *LogBackUp {
	return &LogBackUp{
		_cmsid:           CMSID_LOG_BACKUP,
		AbsFilePath:      fpath,
		AbsBackupDirPath: dpath,
		ZipFile:          zip,
	}
}
