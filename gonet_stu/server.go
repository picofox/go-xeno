package main

import (
	"fmt"
	"net"
	"github.com/cloudwego/netpoll"
)


func main() {
	listener, err := netpoll.CreateListener("tcp", "0.0.0.0")
	if err != nil {
		panic("create netpoll listener failed")
	}
	.
}