package nic

import (
	"net"
	"strings"
	"sync"
	"xeno/zohar/core"
	"xeno/zohar/core/logging"
)

type NICManager struct {
	_nics    []*NIC
	_nicsMap map[int]*NIC
}

var sNICManager NICManager
var sNICManagerOne sync.Once

func (ego *NICManager) clear() {
	for k, _ := range ego._nicsMap {
		delete(ego._nicsMap, k)
	}
	ego._nics = nil
	ego._nics = make([]*NIC, 0)
}

func (ego *NICManager) AddNic(nic *NIC) {
	_, ok := ego._nicsMap[nic.index()]
	if ok {
		return
	}
	ego._nicsMap[nic.index()] = nic
	ego._nics = append(ego._nics, nic)
}

func (ego *NICManager) Update() int32 {
	ego.clear()

	ifaces, err := net.Interfaces()
	if err != nil {
		logging.Log(core.LL_ERR, "Get Net.Interfaces Failed err:(%s)", err.Error())
		return core.MkErr(core.EC_NOOP, 1)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			logging.Log(core.LL_ERR, "Get Addrs Failed err:(%s)", err.Error())
			continue
		}
		nic := NeoNIC(i)
		for _, a := range addrs {
			switch a.(type) {
			case *net.IPNet:
				ipba := a.(*net.IPNet).IP.To4()
				if ipba != nil {
					ipv4Net := NeoIPV4NetByBytes(ipba, a.(*net.IPNet).Mask)
					nic.AddIPNet(ipv4Net)
				} else {
					ipv6Net := NeoIPV6NetByBytes(a.(*net.IPNet).IP.To16(), a.(*net.IPNet).Mask)
					nic.AddIPNet(ipv6Net)
				}

			default:
				logging.Log(core.LL_ERR, "unknow Type %v - %v", i.Name, a)
			}
		}
		ego.AddNic(nic)
	}
	return core.MkSuccess(0)
}

func (ego *NICManager) String() string {
	var ss strings.Builder
	for _, e := range ego._nics {
		ss.WriteString(e.String())
		ss.WriteString("\n")
	}
	return ss.String()
}

func GetNICManager() *NICManager {

	nm := NICManager{
		_nics:    make([]*NIC, 0),
		_nicsMap: make(map[int]*NIC),
	}

	sNICManagerOne.Do(
		func() {
			nm.Update()
		},
	)
	return &nm
}
