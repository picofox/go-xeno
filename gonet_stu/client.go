package main

import (
	"fmt"
	"xeno/zohar/core/inet"
	"xeno/zohar/core/inet/client"
)

func main() {
	ep := inet.NeoIPV4EndPointByEPStr(inet.EP_PROTO_TCP, 0, 0, "www.sina.com.cn:8080")
	fmt.Println(ep.String())

	c := client.NeoTcpClientConnection("192.168.0.20:9998")
	if c == nil {
		fmt.Println("connect Failed")
	}

	c.Connect()

}
