package config

import (
	"encoding/json"
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/io"
	"xeno/zohar/core/process"
)

type DeusConfig struct {
	DB      config.DBConfig      `json:"DB"`
	Network config.NetworkConfig `json:"Network"`
}

func (ego *DeusConfig) String() string {
	return ego.Network.String()
}

var deusConfigConfigInstance *DeusConfig = nil

func GetDeusConfig() *DeusConfig {
	return deusConfigConfigInstance
}

func LoadDeusConfig() (int32, string) {
	fileName := process.ProgramConfFile("deus.", ".json")
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
	var cfg DeusConfig
	err := json.Unmarshal(bs, &cfg)
	if err != nil {
		return core.MkErr(core.EC_JSON_UNMARSHAL_FAILED, 1), fmt.Sprintf("File <%s> to Json error: (%s)", fileName, err)
	}

	deusConfigConfigInstance = &cfg
	return core.MkSuccess(0), ""
}
