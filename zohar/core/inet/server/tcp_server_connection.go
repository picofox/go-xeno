package server

import (
	"errors"
	"io"
	"net"
	"os"
	"time"
	"xeno/zohar/core"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/logging"
)

type TcpServerConnection struct {
	_conn     net.Conn
	_endPoint inet.IPV4EndPoint
	_buffer   []byte
	_length   int64
}

func (ego *TcpServerConnection) Shutdown() {
	ego._conn.Close()
}

func (ego *TcpServerConnection) TryRead() int {
	n, err := ego._conn.Read(ego._buffer[ego._length:])

	if err != nil {
		if err == io.EOF {
			logging.Log(core.LL_ERR, "Read Conn <%s> Closed", ego.String())
			return -1
		} else if errors.Is(err, os.ErrDeadlineExceeded) {
			return 0
		} else {
			logging.Log(core.LL_ERR, "Read Conn <%s> Error: %s", ego.String(), err.Error())
			return -2
		}
	}
	ego._length = ego._length + int64(n)

	return n
}

func (ego *TcpServerConnection) BufferLength() int64 {
	return ego._length
}

func (ego *TcpServerConnection) BufferCapacity() int64 {
	return 32768 + 4
}

func (ego *TcpServerConnection) SetNextReadTimeout(t time.Time) {
	ego._conn.SetReadDeadline(t)
}

func (ego *TcpServerConnection) String() string {
	return ego._endPoint.String()
}

func (ego *TcpServerConnection) Identifier() int64 {
	return ego._endPoint.Identifier()
}

func NeoTcpServerConnection(conn net.Conn) *TcpServerConnection {
	c := TcpServerConnection{
		_conn:     conn,
		_endPoint: inet.NeoIPV4EndPointByAddr(conn.RemoteAddr()),
		_buffer:   make([]byte, 32768+4),
		_length:   0,
	}
	return &c
}
