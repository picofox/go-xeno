package main

import (
	"net"
	"time"
)

func main() {

	sa := "192.168.0.100:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp", sa)
	if err != nil {
		return
	}

	_, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}
	time.Sleep(100 * time.Second)

}
