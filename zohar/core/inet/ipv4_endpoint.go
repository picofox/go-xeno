package inet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
	"xeno/zohar/core"
	"xeno/zohar/core/memory"
	"xeno/zohar/core/strs"
)

const (
	EP_PROTO_IP  = 0
	EP_PROTO_TCP = 1
	EP_PROTO_UDP = 2
)

var sProtoName [3]string = [3]string{"ip", "tcp", "udp"}

type IPV4EndPoint struct {
	_identifier int64
}

func (ego *IPV4EndPoint) Identifier() int64 {
	return ego._identifier
}

func (ego *IPV4EndPoint) SetIdentifier(d int64) bool {
	if d < 0 {
		return false
	}
	ego._identifier = d
	return true
}

func (ego *IPV4EndPoint) Valid() bool {
	return ego._identifier > 0
}

func (ego *IPV4EndPoint) SetInvalid() {
	ego._identifier = -1
}

func (ego *IPV4EndPoint) Proto() int8 {
	return int8((ego._identifier >> 61) & 0x3)
}

func (ego *IPV4EndPoint) ProtoName() string {
	return sProtoName[ego.Proto()]
}

func (ego *IPV4EndPoint) EndPointString() string {
	return fmt.Sprintf("%s:%d", ego.IPV4Str(), ego.Port())
}

func (ego *IPV4EndPoint) SetProto(p int8) {
	ego._identifier = ego._identifier & 0x1FFFFFFFFFFFFFFF
	ego._identifier = ego._identifier | (int64(p) << 61)
}

func (ego *IPV4EndPoint) Mask() int8 {
	return int8((ego._identifier >> 56) & 0x1F)
}

func (ego *IPV4EndPoint) SetMask(m int8) {
	ego._identifier = ego._identifier & 0x60FFFFFFFFFFFFFF
	ego._identifier = ego._identifier | (int64(m) << 56)
}

func (ego *IPV4EndPoint) Extra() uint8 {
	return uint8((ego._identifier >> 48) & 0xFF)
}

func (ego *IPV4EndPoint) SetExtra(ex uint8) {
	ego._identifier = ego._identifier & 0x7F00FFFFFFFFFFFF
	ego._identifier = ego._identifier | (int64(ex) << 48)
}

func (ego *IPV4EndPoint) Port() uint16 {
	return uint16((ego._identifier >> 32) & 0xFFFF)
}

func (ego *IPV4EndPoint) SetPort(port uint16) {
	ego._identifier = ego._identifier & 0x7FFF0000FFFFFFFF
	ego._identifier = ego._identifier | (int64(port) << 32)
}

func (ego *IPV4EndPoint) IPV4() uint32 {
	return uint32((ego._identifier) & 0xFFFFFFFF)
}

func (ego *IPV4EndPoint) IPV4Str() string {
	return strs.IPV4UIntToString(ego.IPV4())
}

func (ego *IPV4EndPoint) SetIPV4Str(ips string) bool {
	ipv4addr, rc := strs.IPV4Addr2UIntBE(ips)
	if core.Err(rc) {
		return false
	}
	ego._identifier = ego._identifier & 0x7FFFFFFF00000000
	ego._identifier = ego._identifier | (int64(ipv4addr))
	return true
}

func (ego *IPV4EndPoint) ToSockAddr() syscall.Sockaddr {
	if ego.Proto() != EP_PROTO_TCP {
		return nil
	}

	sa := syscall.SockaddrInet4{
		Port: int(ego.Port()),
	}
	ba := sa.Addr[0:4]
	memory.UInt32IntoBytesBE(ego.IPV4(), &ba, 0)

	return &sa
}

func (ego *IPV4EndPoint) ToTCPAddr() *net.TCPAddr {
	if ego.Proto() != EP_PROTO_TCP {
		return nil
	}

	ret, err := net.ResolveTCPAddr("tcp", ego.EndPointString())
	if err != nil {
		return nil
	}
	return ret
}

func (ego *IPV4EndPoint) String() string {
	return fmt.Sprintf("%s://%s:%d (mask:%d extra:%d)", sProtoName[ego.Proto()], ego.IPV4Str(), ego.Port(), ego.Mask(), ego.Extra())
}

func (ego *IPV4EndPoint) URI(path ...string) string {
	var ss strings.Builder
	ss.WriteString(fmt.Sprintf("%s://%s:%d", sProtoName[ego.Proto()], ego.IPV4Str(), ego.Port()))
	for _, e := range path {
		ss.WriteString("/")
		ss.WriteString(e)
	}
	return ss.String()
}

func NeoIPV4EndPoint(proto int8, mask int8, extra uint8, ipv4Addr uint32, port uint16) IPV4EndPoint {
	var d int64 = 0
	d = d | (int64(proto&0x7) << 61)
	d = d | (int64(mask&0x1F) << 56)
	d = d | (int64(extra&0xFF) << 48)
	d = d | (int64(port) << 32)
	d = d | (int64(ipv4Addr))
	return IPV4EndPoint{
		_identifier: d,
	}
}

func NeoIPV4EndPointByEPStr(proto int8, mask int8, extra uint8, ipv4AddrStr string) IPV4EndPoint {
	ss := strings.Split(ipv4AddrStr, ":")
	if ss == nil || len(ss) < 2 {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}

	addrsList, err := net.LookupIP(ss[0])
	if err != nil {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	ss[0] = addrsList[0].String()

	return NeoIPV4EndPointByStrIP(proto, mask, extra, ss[0], uint16(port))
}
func NeoIPV4EndPointBySockAddr(proto int8, mask int8, extra uint8, sockaddr syscall.Sockaddr) IPV4EndPoint {
	var ip uint32 = 0
	var port uint16 = uint16(0)
	ok := false
	switch sockaddr.(type) {
	case *syscall.SockaddrInet4:
		ba := sockaddr.(*syscall.SockaddrInet4).Addr[0:4]
		ip = memory.BytesToUInt32BE(&ba, 0)
		port = uint16(sockaddr.(*syscall.SockaddrInet4).Port)
		ok = true
	case *syscall.SockaddrInet6:
		ba := sockaddr.(*syscall.SockaddrInet6).Addr[12:16]
		ip = memory.BytesToUInt32BE(&ba, 0)
		port = uint16(sockaddr.(*syscall.SockaddrInet6).Port)
		ok = true
	}
	if !ok {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	var d int64 = 0
	d = d | (int64(proto&0x7) << 61)
	d = d | (int64(mask&0x1F) << 56)
	d = d | (int64(extra&0xFF) << 48)
	d = d | (int64(port) << 32)
	d = d | (int64(ip))
	return IPV4EndPoint{
		_identifier: d,
	}
}

func NeoIPV4EndPointByStrIP(proto int8, mask int8, extra uint8, ipv4AddrStr string, port uint16) IPV4EndPoint {
	if strings.ToLower(ipv4AddrStr) == "localhost" {
		ipv4AddrStr = "127.0.0.1"
	} else {
		addrsList, err := net.LookupIP(ipv4AddrStr)
		if err != nil {
			return IPV4EndPoint{
				_identifier: -1,
			}
		}
		ipv4AddrStr = addrsList[0].String()
	}

	var d int64 = 0
	ipv4Addr, rc := strs.IPV4Addr2UIntBE(ipv4AddrStr)
	if core.Err(rc) {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	d = d | (int64(proto&0x7) << 61)
	d = d | (int64(mask&0x1F) << 56)
	d = d | (int64(extra&0xFF) << 48)
	d = d | (int64(port) << 32)
	d = d | (int64(ipv4Addr))
	return IPV4EndPoint{
		_identifier: d,
	}
}

func NeoIPV4EndPointByIdentifier(id int64) IPV4EndPoint {
	if id < 0 {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	return IPV4EndPoint{
		_identifier: id,
	}
}

func NeoIPV4EndPointByAddr(a net.Addr) IPV4EndPoint {
	str := a.String()
	ss := strings.Split(str, ":")
	if ss == nil || len(ss) < 2 {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	port, err := strconv.Atoi(ss[1])
	if err != nil {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}

	var proto int8
	if strings.ToLower(a.Network()) == "tcp" {
		proto = EP_PROTO_TCP
	} else if strings.ToLower(a.Network()) == "udp" {
		proto = EP_PROTO_UDP
	} else {
		return IPV4EndPoint{
			_identifier: -1,
		}
	}
	return NeoIPV4EndPointByStrIP(proto, 0, 0, ss[0], uint16(port))
}
