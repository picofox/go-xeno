package net

import (
	"fmt"
	"xeno/zohar/core"
	"xeno/zohar/core/strs"
)

const (
	EP_PROTO_IP  = 0
	EP_PROTO_TCP = 1
	EP_PROTO_UDP = 2
)

var sProtoName [3]string = [3]string{"ip", "tcp", "udp"}

type IPV4EndPoint struct {
	_data int64
}

func (ego *IPV4EndPoint) Data() int64 {
	return ego._data
}

func (ego *IPV4EndPoint) Valid() bool {
	return ego._data > 0
}

func (ego *IPV4EndPoint) Proto() int8 {
	return int8((ego._data >> 61) & 0x3)
}

func (ego *IPV4EndPoint) SetProto(p int8) {
	ego._data = ego._data & 0x1FFFFFFFFFFFFFFF
	ego._data = ego._data | (int64(p) << 61)
}

func (ego *IPV4EndPoint) Mask() int8 {
	return int8((ego._data >> 56) & 0x1F)
}

func (ego *IPV4EndPoint) SetMask(m int8) {
	ego._data = ego._data & 0x60FFFFFFFFFFFFFF
	ego._data = ego._data | (int64(m) << 56)
}

func (ego *IPV4EndPoint) Extra() uint8 {
	return uint8((ego._data >> 48) & 0xFF)
}

func (ego *IPV4EndPoint) SetExtra(ex uint8) {
	ego._data = ego._data & 0x7F00FFFFFFFFFFFF
	ego._data = ego._data | (int64(ex) << 48)
}

func (ego *IPV4EndPoint) Port() uint16 {
	return uint16((ego._data >> 32) & 0xFFFF)
}

func (ego *IPV4EndPoint) SetPort(port uint16) {
	ego._data = ego._data & 0x7FFF0000FFFFFFFF
	ego._data = ego._data | (int64(port) << 32)
}

func (ego *IPV4EndPoint) IPV4() uint32 {
	return uint32((ego._data) & 0xFFFFFFFF)
}

func (ego *IPV4EndPoint) IPV4Str() string {
	return strs.IPV4UIntToString(ego.IPV4())
}

func (ego *IPV4EndPoint) SetIPV4Str(ips string) bool {
	ipv4addr, rc := strs.IPV4Addr2UIntBE(ips)
	if core.Err(rc) {
		return false
	}
	ego._data = ego._data & 0x7FFFFFFF00000000
	ego._data = ego._data | (int64(ipv4addr))
	return true
}

func (ego *IPV4EndPoint) String() string {
	return fmt.Sprintf("%s://%s:%d (mask:%d extra:%d)", sProtoName[ego.Proto()], ego.IPV4Str(), ego.Port(), ego.Mask(), ego.Extra())
}

func NeoIPV4EndPoint(proto int8, mask int8, extra uint8, ipv4Addr uint32, port uint16) IPV4EndPoint {
	var d int64 = 0
	d = d | (int64(proto&0x7) << 61)
	d = d | (int64(mask&0x1F) << 56)
	d = d | (int64(extra&0xFF) << 48)
	d = d | (int64(port) << 32)
	d = d | (int64(ipv4Addr))
	return IPV4EndPoint{
		_data: d,
	}
}

func NeoIPV4EndPointByStrIP(proto int8, mask int8, extra uint8, ipv4AddrStr string, port uint16) IPV4EndPoint {
	var d int64 = 0
	ipv4Addr, rc := strs.IPV4Addr2UIntBE(ipv4AddrStr)
	if core.Err(rc) {
		return IPV4EndPoint{
			_data: -1,
		}
	}
	d = d | (int64(proto&0x7) << 61)
	d = d | (int64(mask&0x1F) << 56)
	d = d | (int64(extra&0xFF) << 48)
	d = d | (int64(port) << 32)
	d = d | (int64(ipv4Addr))
	return IPV4EndPoint{
		_data: d,
	}
}
