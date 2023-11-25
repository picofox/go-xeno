package server

import (
	"net"
	"os"
	"syscall"
	"xeno/zohar/core"
)

type ListenerWrapper struct {
	fd    int
	addr  net.Addr       // listener's local addr
	ln    net.Listener   // tcp|unix listener
	pconn net.PacketConn // udp listener
	file  *os.File
}

var _ Listener = &ListenerWrapper{}

type Listener interface {
	net.Listener

	// Fd return listener's fd, used by poll.
	Fd() (fd int)
}

// Accept implements Listener.
func (ego *ListenerWrapper) Accept() (net.Conn, int32) {
	// udp
	if ego.pconn != nil {
		return ego.UDPAccept()
	}
	// tcp
	var fd, sa, err = syscall.Accept(ego.fd)
	if err != nil {
		if err == syscall.EAGAIN {
			return nil, core.MkErr(core.EC_TRY_AGAIN, 1)
		}
		return nil, core.MkErr(core.EC_NULL_VALUE, 2)
	}
	var nfd = &netFD{}
	nfd.fd = fd
	nfd.localAddr = ego.addr
	nfd.network = ego.addr.Network()
	nfd.remoteAddr = sockaddrToAddr(sa)
	return nfd, core.MkSuccess(0)
}

// TODO: UDPAccept Not implemented.
func (ego *ListenerWrapper) UDPAccept() (net.Conn, int32) {
	return nil, core.MkErr(core.EC_NULL_VALUE, 1)
}

// Close implements Listener.
func (ego *ListenerWrapper) Close() int32 {
	if ego.fd != 0 {
		syscall.Close(ego.fd)
	}
	if ego.file != nil {
		ego.file.Close()
	}
	if ego.ln != nil {
		ego.ln.Close()
	}
	if ego.pconn != nil {
		ego.pconn.Close()
	}
	return core.MkSuccess(0)
}

// Addr implements Listener.
func (ego *ListenerWrapper) Addr() net.Addr {
	return ego.addr
}

// Fd implements Listener.
func (ego *ListenerWrapper) Fd() (fd int) {
	return ego.fd
}

func (ego *ListenerWrapper) parseFD() (rc int32) {
	switch netln := ego.ln.(type) {
	case *net.TCPListener:
		ego.file, err = netln.File()
	case *net.UnixListener:
		ego.file, err = netln.File()
	default:
		return core.MkErr(core.EC_TYPE_MISMATCH, 1)
	}
	if err != nil {
		return core.MkErr(core.EC_NULL_VALUE, 1)
	}
	ego.fd = int(ego.file.Fd())
	return core.MkSuccess(0)
}
