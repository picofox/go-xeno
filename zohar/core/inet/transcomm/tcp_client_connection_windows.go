package transcomm

import (
	"fmt"
	"io"
	"net"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
)

type TcpClientConnection struct {
	_conn          net.Conn
	_remoteAddress inet.IPV4EndPoint
	_localAddress  inet.IPV4EndPoint
}

func (ego *TcpClientConnection) Connect() int32 {
	var err error
	if ego._localAddress.Valid() {
		ego._conn, err = net.DialTCP("tcp", ego._localAddress.ToTCPAddr(), ego._remoteAddress.ToTCPAddr())
		if err != nil {
			return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
		}
	} else {
		ego._conn, err = net.DialTCP("tcp", nil, ego._remoteAddress.ToTCPAddr())
		if err != nil {
			return core.MkErr(core.EC_TCP_CONNECT_ERROR, 1)
		}
		ego._localAddress = inet.NeoIPV4EndPointByAddr(ego._conn.LocalAddr())
		fmt.Printf("local addr is %v\n", ego._localAddress.String())
	}
	return core.MkSuccess(0)
}

func (ego *TcpClientConnection) SendImmediately(ba []byte, offset int, length int) (int, int32) {
	n, err := ego._conn.Write(ba[offset:length])
	if err != nil {
		return n, core.MkErr(core.EC_FILE_WRITE_FAILED, 1)
	}
	return n, core.MkSuccess(0)
}

func (ego *TcpClientConnection) Recv(ba []byte, offset int, length int) (int, int32) {
	n, err := ego._conn.Read(ba[offset:length])
	if err != nil {
		if err == io.EOF {
			return n, core.MkErr(core.EC_EOF, 0)
		}
		return n, core.MkErr(core.EC_FILE_READ_FAILED, 1)
	}
	return n, core.MkSuccess(0)
}

func (ego *TcpClientConnection) Close() {
	ego._conn.Close()
	ego._conn = nil
}

func NeoTcpClientConnection(epStrRemote string, epStrLocal ...string) *TcpClientConnection {
	var addrLocal inet.IPV4EndPoint
	if len(epStrLocal) > 0 {
		addrLocal = inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, epStrRemote)
	} else {
		addrLocal.SetInvalid()
	}

	c := TcpClientConnection{
		_conn:          nil,
		_remoteAddress: inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, epStrRemote),
	}
	return &c
}
