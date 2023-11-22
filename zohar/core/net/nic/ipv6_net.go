package nic

import (
	"fmt"
	"net"
)

type IPV6Net struct {
	_addr net.IP
	_mask net.IPMask
}

func (ego *IPV6Net) Type() uint8 {
	return INET_ADDRESS_TYPE_IPV6
}

func (ego *IPV6Net) Address() []byte {
	return ego._addr
}

func (ego *IPV6Net) NetMask() []byte {
	return ego._mask
}

func (ego *IPV6Net) AddressString() string {
	return ego.AddrStr()
}

func (ego *IPV6Net) NetMaskStr() string {
	return ego.NetMaskStr()
}

func (ego *IPV6Net) String() string {
	return fmt.Sprintf("Type: %d Addr:%s Mask:%s", ego.Type(), ego.AddrStr(), ego.MaskStr())
}

func NeoIPV6NetByBytes(ipBa net.IP, maskBa net.IPMask) *IPV6Net {
	return &IPV6Net{
		_addr: ipBa,
		_mask: maskBa,
	}
}

func (ego *IPV6Net) AddrStr() string {
	return ego._addr.String()
}

func (ego *IPV6Net) MaskStr() string {
	return ego._mask.String()
}
