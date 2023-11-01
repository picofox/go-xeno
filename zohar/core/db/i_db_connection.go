package db

import (
	"fmt"
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/strs"
)

const (
	DB_TYPE_MYSQL = iota
	DB_TYPE_COUNT
)

type IDBConnection interface {
	Index() int
	TableNames() []string
	Connect() int32
}

func ParseTableNameConfig(tableStr string) ([]string, int32) {
	var ret []string = make([]string, 0)
	m, s, ok := strs.ExtractPinchString(tableStr, "[", "]")
	if !ok {
		ret = append(ret, m)
	} else {
		if len(s) < 2 {
			return nil, core.MkErr(core.EC_ELEMENT_NOT_FOUND, 1)
		}

		idx := strings.Index(s, "-")
		if idx >= 0 {
			left := strings.Trim(s[0:idx], " \t\r\n")
			right := strings.Trim(s[idx+1:], " \t\r\n")
			l, err := strconv.Atoi(left)
			if err != nil {
				return nil, core.MkErr(core.EC_TYPE_MISMATCH, 1)
			}
			r, err := strconv.Atoi(right)
			if err != nil {
				return nil, core.MkErr(core.EC_TYPE_MISMATCH, 2)
			}
			if l >= r {
				return nil, core.MkErr(core.EC_INDEX_OOB, 2)
			}
			for i := l; i < r; i++ {
				ret = append(ret, fmt.Sprintf("%s_%05d", m, i))
			}
			return ret, core.MkSuccess(0)
		} else {
			arr := strings.Split(s, ",")
			if len(arr) < 1 {
				return nil, core.MkErr(core.EC_TYPE_MISMATCH, 5)
			}
			for i := 0; i < len(arr); i++ {
				arr[i] = strings.Trim(arr[i], " \t\r\n")
				idx, err := strconv.Atoi(arr[i])
				if err != nil {
					return nil, core.MkErr(core.EC_TYPE_MISMATCH, 6)
				}
				ret = append(ret, fmt.Sprintf("%s_%05d", m, idx))
			}

		}

	}
	return ret, core.MkSuccess(0)
}
