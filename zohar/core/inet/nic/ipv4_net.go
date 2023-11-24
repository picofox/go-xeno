package nic

import (
	"fmt"
	"strings"
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

type IPV4Net struct {
	_addr uint32
	_mask uint32
}

func (ego *IPV4Net) Type() uint8 {
	return INET_ADDRESS_TYPE_IPV4
}

func (ego *IPV4Net) Address() []byte {
	return *memory.UInt32ToBytesBE(ego._addr)
}

func (ego *IPV4Net) NetMask() []byte {
	return *memory.UInt32ToBytesBE(ego._mask)
}

func (ego *IPV4Net) AddressString() string {
	return ego.AddrStr()
}

func (ego *IPV4Net) NetMaskStr() string {
	return strs.IPV4UIntToString(ego._mask)
}

func (ego *IPV4Net) String() string {
	return fmt.Sprintf("Type: %d Addr:%s Mask:%s", ego.Type(), ego.AddrStr(), ego.MaskStr())
}

func NeoIPV4NetByBytes(ipBa []byte, maskBa []byte) *IPV4Net {
	return &IPV4Net{
		_addr: memory.BytesToUInt32BE(&ipBa, 0),
		_mask: memory.BytesToUInt32BE(&maskBa, 0),
	}
}

func (ego *IPV4Net) AddrStr() string {
	return strs.IPV4UIntToString(ego._addr)
}

func (ego *IPV4Net) MaskStr() string {
	return strs.IPV4UIntToString(ego._mask)
}

func NeoIPV4NetByString(str string) *IPV4Net {
	s := strings.Split(str, "/")
	if len(s) < 2 {
		return nil
	}

	ip, rc := strs.IPV4Addr2UIntBE(s[0])
	if core.Err(rc) {
		return nil
	}
	mask, rc := strs.IPV4MaskBits2UIntBE(s[1])
	if core.Err(rc) {
		return nil
	}

	return &IPV4Net{
		_addr: ip,
		_mask: mask,
	}
}
