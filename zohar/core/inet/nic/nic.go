package nic

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/strs"
)

//inet: add FlagRunning to the Flags of struct Interface, to exactly reflect the states of an interface or NIC.
//
//And a new flag(FlagRunning), and correctly set this flag while parsing the syscall result.
//
//The FlagUp flag can not distinguish the following situations:
//1. interface is plugged, automatically up, and in running(UP) state
//2. interface is not plugged, administratively or manually set to up, but in DOWN state
//
//So, We can't distinguish the state of a NIC through the FlagUp flag only.

type NIC struct {
	_if  net.Interface
	_ips []IINetAddress
}

func (ego *NIC) FindNetInfoByIPV4Address(ipU32 uint32) IINetAddress {
	for _, ip := range ego._ips {
		fmt.Printf("compare %s vs %s\n", ip.AddressString(), strs.IPV4UIntToString(ipU32))
		if ip.Type() == INET_ADDRESS_TYPE_IPV4 {
			ipv4, rc := strs.IPV4Addr2UIntBE(ip.AddressString())
			if core.Err(rc) {
				return nil
			}

			if ip.AddressString() == "192.168.0.100" {
				fmt.Printf("x")
			}

			if ipv4 == ipU32 {
				return ip
			}
		}
	}
	return nil
}

func (ego *NIC) String() string {
	var ss strings.Builder
	ss.WriteString(ego.Name())
	ss.WriteString("@")
	ss.WriteString(strconv.Itoa(ego.index()))
	ss.WriteString("\n")
	ss.WriteString("    ")
	if ego.IsRunning() && ego.IsUp() {
		ss.WriteString("Up & Running")
	} else if ego.IsRunning() {
		ss.WriteString("Running")
	} else if ego.IsUp() {
		ss.WriteString("Up")
	} else {
		ss.WriteString("Inactive")
	}
	ss.WriteString("\n")

	for _, e := range ego._ips {
		ss.WriteString("    ")
		ss.WriteString(e.String())
		ss.WriteString("\n")
	}
	return ss.String()
}

func (ego *NIC) index() int {
	return ego._if.Index
}

func (ego *NIC) MAC() []byte {
	return ego._if.HardwareAddr
}

func (ego *NIC) Name() string {
	return ego._if.Name
}

func (ego *NIC) IsUp() bool {
	return ego._if.Flags&net.FlagUp > 0
}

func (ego *NIC) IsLoopback() bool {
	return ego._if.Flags&net.FlagLoopback > 0
}

func (ego *NIC) Broadcastable() bool {
	return ego._if.Flags&net.FlagBroadcast > 0
}

func (ego *NIC) Multicastable() bool {
	return ego._if.Flags&net.FlagMulticast > 0
}

func (ego *NIC) IsP2P() bool {
	return ego._if.Flags&net.FlagPointToPoint > 0
}

func (ego *NIC) IsRunning() bool {
	return ego._if.Flags&net.FlagRunning > 0
}

func (ego *NIC) AddIPNet(ip IINetAddress) {
	ego._ips = append(ego._ips, ip)
}

func NeoNIC(i net.Interface) *NIC {
	return &NIC{
		_if:  i,
		_ips: make([]IINetAddress, 0),
	}
}
