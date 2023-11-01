package config

import (
	"encoding/json"
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/io"
	"xeno/zohar/core/process"
)

type KadamonyConfig struct {
	DB config.DBConfig `json:"DB"`
}

var kadamonyConfigInstance *KadamonyConfig = nil

func GetKadamonyConfig() *KadamonyConfig {
	return kadamonyConfigInstance
}

func LoadKadamonyConfig() (int32, string) {
	fileName := process.ProgramConfFile("kadamony.", ".json")
	f := io.NeoFile(false, fileName, io.FILEFLAG_THREAD_SAFE)
	rc := f.Open(io.FILEOPEN_MODE_OPEN_EXIST, io.FILEOPEN_PERM_READ, 0755)
	if core.Err(rc) {
		return core.MkErr(core.EC_FILE_OPEN_FAILED, 1), fmt.Sprintf("Open file <%s> error", fileName)
	}
	defer f.Close()
	bs, rc := f.ReadAll()
	if core.Err(rc) {
		return core.MkErr(core.EC_FILE_READ_FAILED, 1), fmt.Sprintf("Readall file <%s> error: (%d)", fileName, rc)
	}
	var cfg KadamonyConfig
	err := json.Unmarshal(bs, &cfg)
	if err != nil {
		return core.MkErr(core.EC_JSON_UNMARSHAL_FAILED, 1), fmt.Sprintf("File <%s> to Json error: (%s)", fileName, err)
	}

	kadamonyConfigInstance = &cfg
	return core.MkSuccess(0), ""
}
