package cmdline

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/config"
	"xeno/zohar/core/memory"
)

type Arguments struct {
	_isStrict    bool
	_program     string
	_targetSpec  *ArgumentSpec
	_shortSpecs  [128]*ArgumentSpec
	_shortValues [128]*memory.TLV
	_targets     []*memory.TLV
	_specs       map[string]*ArgumentSpec
	_values      map[string]*memory.TLV
}

func (ego *Arguments) GetLongParam(lcmd string) *memory.TLV {
	tlv, ok := ego._values[lcmd]
	if ok {
		return tlv
	}
	return nil
}

func CmdArgGetLongParamValue[T any](lcmd string, dfl T) (T, int32) {
	tlv := GetArguments().GetLongParam(lcmd)
	if tlv == nil {
		return dfl, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	return tlv.Value().(T), core.MkSuccess(0)
}

func CmdArgGetLongParamAsDictElement[KT any, T any](lcmd string, key KT, dfl T) (T, int32) {
	tlv := GetArguments().GetLongParam(lcmd)
	if tlv == nil {
		return dfl, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	v, rc := tlv.GetDictValue(key)
	if core.Err(rc) {
		return dfl, rc
	}
	return v.(T), rc
}

func CmdArgGetLongParamAsListElement[T any](lcmd string, idx uint32, dfl T) (T, int32) {
	tlv := GetArguments().GetLongParam(lcmd)
	if tlv == nil {
		return dfl, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	v, rc := tlv.GetListValue(idx)
	if core.Err(rc) {
		return dfl, rc
	}
	return v.(T), rc
}

func (ego *Arguments) GetLongParamValue(lcmd string) any {
	tlv, ok := ego._values[lcmd]
	if ok {
		return tlv.Value()
	}
	return nil
}

func (ego *Arguments) GetShortParam(scmd uint8) *memory.TLV {
	if scmd < 128 {
		return ego._shortValues[scmd]
	}
	return nil
}

func CmdArgGetShortParamValue[T any](scmd uint8, dfl T) (T, int32) {
	tlv := GetArguments().GetShortParam(scmd)
	if tlv == nil {
		return dfl, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	return tlv.Value().(T), core.MkSuccess(0)
}

func CmdArgGetShortParamAsDictElement[KT any, T any](scmd uint8, key KT, dfl T) (T, int32) {
	tlv := GetArguments().GetShortParam(scmd)
	if tlv == nil {
		return dfl, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}

	v, rc := tlv.GetDictValue(key)
	if core.Err(rc) {
		return dfl, rc
	}
	return v.(T), rc
}

func CmdArgGetShortParamAsListElement[T any](arg *Arguments, scmd uint8, idx uint32, dfl T) (T, int32) {
	tlv := arg.GetShortParam(scmd)
	if tlv == nil {
		return dfl, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
	}
	v, rc := tlv.GetListValue(idx)
	if core.Err(rc) {
		return dfl, rc
	}
	return v.(T), rc
}

func (ego *Arguments) GetShortParamValue(scmd uint8) any {
	if scmd < 128 {
		return ego._shortValues[scmd].Value()
	}
	return nil
}

func (ego *Arguments) string() string {
	var ss strings.Builder
	ss.WriteString("IsStrict:")
	ss.WriteString(strconv.FormatBool(ego._isStrict))

	return ss.String()
}

var cmdArgInstance *Arguments
var once sync.Once

func GetArguments() *Arguments {
	once.Do(func() {
		errString := ""
		cmdArgInstance, errString = Initialize()
		if cmdArgInstance == nil {
			panic(fmt.Sprintf("Init Command Line Arguments \t\t\t[Failed:%s]", errString))
		}
	})
	return cmdArgInstance
}

func Initialize() (*Arguments, string) {
	var specRawPart []string
	a := &Arguments{
		_isStrict:    config.GetIntrinsicConfig().CmdSpecRestrict,
		_program:     os.Args[0],
		_targetSpec:  nil,
		_shortSpecs:  [128]*ArgumentSpec{},
		_shortValues: [128]*memory.TLV{},
		_targets:     make([]*memory.TLV, 0),
		_specs:       make(map[string]*ArgumentSpec),
		_values:      make(map[string]*memory.TLV),
	}

	if len(config.GetIntrinsicConfig().CmdTargetSpec) > 1 {
		targetSubPart := strings.Split(config.GetIntrinsicConfig().CmdTargetSpec, ":")
		if len(targetSubPart) != 3 {
			return nil, "Invalid Target Spec of " + config.GetIntrinsicConfig().CmdTargetSpec
		}
		minS, _ := strconv.Atoi(targetSubPart[1])
		maxS, _ := strconv.Atoi(targetSubPart[2])
		a._targetSpec = NeoArgumentSpec(0, targetSubPart[0], true, minS, maxS, "")
	}

	if len(config.GetIntrinsicConfig().CmdParamSpec) > 1 {
		specRawPart = strings.Split(config.GetIntrinsicConfig().CmdParamSpec, ",")
		for i := 0; i < len(specRawPart); i++ {
			specRawPart[i] = strings.Trim(specRawPart[i], " \t\r\n")
		}
	}
	specLen := len(specRawPart)
	if specLen > 0 {
		for i := 0; i < len(specRawPart); i++ {
			specSubPart := strings.Split(specRawPart[i], ":")
			if len(specSubPart) < 4 {
				return nil, "Invalid Argument Spec of " + specRawPart[i]
			}
			for j := 0; j < len(specSubPart); j++ {
				specSubPart[j] = strings.Trim(specSubPart[j], " \t\r\n")
			}

			if len(specSubPart) != 6 {
				return nil, "Invalid Sub Count of Argument Spec of " + specRawPart[i]
			}

			var sn uint8 = 0
			if len(specSubPart[0]) > 0 {
				sn = specSubPart[0][0]
			}
			minS, _ := strconv.Atoi(specSubPart[3])
			maxS, _ := strconv.Atoi(specSubPart[4])
			opt, _ := strconv.ParseBool(specSubPart[5])

			asp := NeoArgumentSpec(uint8(sn), specSubPart[2], opt, minS, maxS, specSubPart[1])
			if !asp.IsValid() {
				return nil, "Invalid Argument Spec of " + asp.String()
			}

			if asp.HasLong() {
				a._specs[asp._longName] = asp
				if asp.HasShort() {
					a._shortSpecs[asp._shortName] = asp
				}
			} else if asp.HasShort() {
				a._shortSpecs[asp._shortName] = asp
				a._specs[asp.ShortCommandString()] = asp
			}
		}
	}

	errStr := a.ParseArgs(os.Args)
	if len(errStr) > 0 {
		panic(errStr)
	}

	return a, ""
}

func (ego *Arguments) checkDone() string {
	for _, v := range ego._specs {
		if v.IsSingle() {
			rc, errString := ego.IsSingleValueValid(v)
			if len(errString) > 0 {
				return errString
			}

			if !rc {
				return fmt.Sprintf("Arg ParseFailed Single Value Validation Failed: (%s)\n", v.String())
			}
		}
	}

	for k, _ := range ego._values {
		v, ok := ego._specs[k]
		if !ok {
			return fmt.Sprintf("Arg Value (%s) is not specific.", v.String())
		}
	}
	return ""
}

func (ego *Arguments) IsSingleValueValid(asp *ArgumentSpec) (bool, string) {
	if asp.Optional() {
		return true, ""
	} else if asp.IsFlag() {
		return true, ""
	}

	if asp.HasShort() {
		if ego._shortValues[asp._shortName] == nil {
			return false, fmt.Sprintf("No value is given to mandatory ARG: %c\n", asp._shortName)
		}
	} else if asp.HasLong() {
		_, ok := ego._values[asp._longName]
		if !ok {
			return false, fmt.Sprintf("No value is given to mandatory ARG: %s\n", asp._longName)
		}
	}
	return true, ""
}

func (ego *Arguments) IsMulFlags(arg string) bool {
	for i := 1; i < len(arg); i++ {
		if ego._shortSpecs[rune(arg[i])] != nil {
			return false
		}

		if !ego._shortSpecs[rune(arg[i])].IsFlag() {
			return false
		}
	}
	return true
}

func (ego *Arguments) String() string {
	var ss strings.Builder

	ss.WriteString(ego._program)
	ss.WriteString("full spec:\n")
	for k, v := range ego._specs {
		ss.WriteString("\t")
		ss.WriteString(k)
		ss.WriteString(" -> ")
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	ss.WriteString("Short Args Specs:\n")
	for i := 0; i < 128; i++ {
		if ego._shortSpecs[i] != nil && ego._shortSpecs[i].IsValid() {
			ss.WriteString("\t")
			ss.WriteString(fmt.Sprintf("%c", i))
			ss.WriteString(" -> ")
			ss.WriteString(ego._shortSpecs[i].String())
			ss.WriteString("\n")
		}
	}
	ss.WriteString("Target Spec:\n")
	ss.WriteString("\t")
	ss.WriteString(ego._targetSpec.String())
	ss.WriteString("\n")
	ss.WriteString("Full Input:\n")
	for k, v := range ego._values {
		ss.WriteString("\t")
		ss.WriteString(k)
		ss.WriteString(" -> ")
		ss.WriteString(v.String())
		ss.WriteString("\n")
	}
	ss.WriteString("Targets List:\n")
	for i := 0; i < len(ego._targets); i++ {
		ss.WriteString("\t")
		ss.WriteString(ego._targets[i].String())
		ss.WriteString("\n")
	}
	return ss.String()
}

func (ego *Arguments) findSpecByShort(scmd uint8) *ArgumentSpec {
	if scmd < 0 || scmd > 127 {
		return nil
	}
	if ego._shortSpecs[scmd] == nil {
		return nil
	}
	if ego._shortSpecs[scmd].HasShort() {
		return ego._shortSpecs[scmd]
	}
	return nil
}

func (ego *Arguments) findSpecByLong(lcmd string) *ArgumentSpec {
	if len(lcmd) < 1 {
		return nil
	}
	v, ok := ego._specs[lcmd]
	if ok {
		if v.HasLong() {
			return v
		}
		return nil
	}
	return nil
}

func (ego *Arguments) findSpec(scmd uint8, lcmd string) *ArgumentSpec {
	sp := ego.findSpecByShort(scmd)
	if sp != nil {
		return sp
	}
	return ego.findSpecByLong(lcmd)
}

func (ego *Arguments) addNeoTarget(target string) string {
	curCount := len(ego._targets)
	if curCount >= ego._targetSpec._maxSACount {
		return fmt.Sprintf("Can't add Target (%s), Too Many Targets: %u", target, curCount)
	}

	tlv := memory.CreateTLV(memory.DT_SINGLE, ego._targetSpec.SingleType(), memory.T_NULL, target)
	ego._targets = append(ego._targets, tlv)
	return ""
}

func (ego *Arguments) updateNeoValueBase(scmd uint8, lcmd string) (*memory.TLV, string) {
	var tlv *memory.TLV = nil
	if len(lcmd) > 0 {
		asp := ego.findSpecByLong(lcmd)
		if asp == nil {
			return nil, fmt.Sprintf("No Spec is found for ARG: %s\n", lcmd)
		}

		tlv = memory.CreateTLV(asp.ContainerType(), asp.SingleType(), asp.keyType(), nil)
		if tlv == nil {
			return nil, fmt.Sprintf("Value allocate failed for ARG: %s\n", lcmd)
		}
		ego._values[lcmd] = tlv

		if asp.HasShort() {
			ego._shortValues[asp._shortName] = tlv
		}
	} else if scmd != 0 {
		asp := ego.findSpecByShort(scmd)
		if asp == nil {
			return nil, fmt.Sprintf("No Spec is found for ARG: %c\n", scmd)
		}
		if !asp.HasLong() {
			str := asp.ShortCommandString()
			tlv = memory.CreateTLV(asp.ContainerType(), asp.SingleType(), asp.keyType(), nil)
			if tlv == nil {
				return nil, fmt.Sprintf("Value allocate failed for ARG: %s\n", lcmd)
			}
			if ego._shortValues[scmd] != nil {
				delete(ego._values, str)
				ego._shortValues[scmd] = nil
			}
			ego._shortValues[scmd] = tlv
			ego._values[str] = tlv

		} else {
			tlv = memory.CreateTLV(asp.ContainerType(), asp.SingleType(), asp.keyType(), nil)
			if tlv == nil {
				return nil, fmt.Sprintf("Value allocate failed for ARG: %s\n", lcmd)
			}
			_, ok := ego._values[asp._longName]
			if ok {
				delete(ego._values, asp._longName)
				ego._shortValues[scmd] = nil
			}

			ego._shortValues[scmd] = tlv
			ego._values[asp._longName] = tlv
		}
	}

	return tlv, ""
}

func (ego *Arguments) updateNeoValue(scmd uint8, lcmd string, val string) (*memory.TLV, string) {
	asp := ego.findSpec(scmd, lcmd)
	if asp == nil {
		return nil, fmt.Sprintf("No value is given to ARG: %c or %s\n", scmd, lcmd)
	}
	if asp._maxSACount != 1 || asp._minSACount != 1 {
		return nil, fmt.Sprintf("Only 1 sub arg is allowd for ARG: %c or %s", scmd, lcmd)
	}

	tlv, rcs := ego.updateNeoValueBase(scmd, lcmd)
	if tlv == nil {
		return nil, rcs
	}

	tlv.SetSingleValue(val)
	return tlv, ""
}

func (ego *Arguments) addMultiParams(scmd uint8, lcmd string, av []string, offset int, maxCount int) string {
	asp := ego.findSpec(scmd, lcmd)
	if asp == nil {
		return fmt.Sprintf("No value is given to ARG: %c or %s\n", scmd, lcmd)
	}

	tlv, rcs := ego.updateNeoValueBase(scmd, lcmd)
	if tlv == nil {
		return rcs
	}

	if asp.IsList() {
		for i := 0; i < maxCount; i++ {
			tlv.PushBack(av[i+offset])
		}
	} else if asp.IsDict() {
		for i := 0; i < maxCount; i++ {
			kvp := strings.SplitN(av[i+offset], "=", 2)
			if len(kvp) != 2 {
				if i < asp._minSACount {
					return fmt.Sprintf("Dicket ELEM strs format error %s", av)
				} else {
					break
				}

			}
			tlv.SetDictValue(kvp[0], kvp[1])
		}
	} else {
		return fmt.Sprintf("Arg %s is multi-value type, but it's container type is %d\n", asp.name(), asp.ContainerType())
	}

	return ""
}

func (ego *Arguments) ParseArgs(av []string) string {
	errStr := ego.checkDone()
	if len(errStr) > 0 {
		return errStr
	}

	ac := len(av)
	var tlv *memory.TLV = nil

	for i := 1; i < len(av); i++ {
		oneArg := strings.Trim(av[i], " \t\r\n")
		argLen := len(oneArg)
		if oneArg[0] == '-' {
			if argLen > 1 && oneArg[1] == '-' { //long
				if argLen == 2 {
					return fmt.Sprintf("[%s] is not a valid arg", oneArg)
				}
				asp := ego.findSpecByLong(oneArg[2:])
				if asp == nil {
					return fmt.Sprintf("No Spec is Found for ARG: %c", oneArg[1])
				}
				if asp.IsFlag() {
					tlv, errStr = ego.updateNeoValueBase(uint8(0), asp._longName)
					if tlv == nil {
						return errStr
					}
				} else {
					if asp._minSACount == 1 && asp._maxSACount == 1 {
						if i >= len(av)-1 {
							return fmt.Sprintf("Arg %s need 1 Param, but no Param Found", asp._longName)
						}
						tlv, errStr = ego.updateNeoValue(uint8(0), asp._longName, av[i+1])
						if tlv == nil {
							return errStr
						}
						i++
					} else {
						var subCnt int = 0
						var maxSub = ac - i - 1
						if maxSub > asp._maxSACount {
							maxSub = asp._maxSACount
						}
						for subIdx := i + 1; subIdx < i+maxSub; subIdx++ {
							if av[subIdx][0] != '-' {
								if asp.IsDict() {
									kvp := strings.Split(av[subIdx], "=")
									if len(kvp) != 2 {
										if subCnt < asp._minSACount {
											return fmt.Sprintf("No meet dict min SA Requirement (%s)", asp.String())
										} else {
											break
										}
									}
								}
								subCnt++
							} else {
								break
							}
						}
						if subCnt < asp._minSACount {
							return fmt.Sprintf("Arg %s need at Least %d Param, but %d is given\n", oneArg, asp._minSACount, subCnt)
						}

						errStr = ego.addMultiParams(0, asp._longName, av, i+1, subCnt)
						if len(errStr) > 0 {
							return errStr
						}

						i = i + subCnt
					}
				}
			} else { //short
				if argLen > 2 {
					if ego.IsMulFlags(oneArg) {
						for j := 1; j < argLen; j++ {
							tlv, errStr = ego.updateNeoValueBase(oneArg[j], "")
							if tlv == nil {
								return errStr
							}
						}
					} else {
						tlv, errStr = ego.updateNeoValue(oneArg[1], "", oneArg[2:])
						if tlv == nil {
							return errStr
						}
					}
				} else {
					asp := ego.findSpecByShort(oneArg[1])
					if asp == nil {
						return fmt.Sprintf("No Spec is Found for ARG: %c", oneArg[1])
					}
					if asp.IsFlag() {
						tlv, errStr = ego.updateNeoValueBase(asp._shortName, "")
						if tlv == nil {
							return errStr
						}
					} else {
						if asp._minSACount == 1 && asp._maxSACount == 1 {
							if i > ac-1 {
								return fmt.Sprintf("Arg %c need 1 Param, but no Param Found", oneArg[1])
							}
							tlv, errStr = ego.updateNeoValue(asp._shortName, "", av[i+1])
							if tlv == nil {
								return errStr
							}
							i++
						} else {
							var subCnt int = 0
							var maxSub = ac - i - 1
							if maxSub > asp._maxSACount {
								maxSub = asp._maxSACount
							}
							for subIdx := i + 1; subIdx <= i+maxSub; subIdx++ {
								if av[subIdx][0] != '-' {
									if asp.IsDict() {
										kvp := strings.Split(av[subIdx], "=")
										if len(kvp) != 2 {
											if subCnt < asp._minSACount {
												return fmt.Sprintf("No meet dict min SA Requirement (%s)", asp.String())
											} else {
												break
											}
										}
									}
									subCnt++
								} else {
									break
								}
							}

							if subCnt < asp._minSACount {
								return fmt.Sprintf("Arg %c need at Least %d Param, but %d is given\n", oneArg[0], asp._minSACount, subCnt)
							}

							errStr = ego.addMultiParams(asp._shortName, "", av, i+1, subCnt)
							if len(errStr) > 0 {
								return errStr
							}
							i = i + subCnt
						}
					}
				}
			}
		} else {
			errStr = ego.addNeoTarget(oneArg)
			if len(errStr) > 0 {
				return errStr
			}
		}
	}

	return ""
}
