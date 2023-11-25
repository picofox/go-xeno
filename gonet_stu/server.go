package main

import (
	"github.com/cloudwego/netpoll"
	"time"
)

func main() {

	//listener, err := netpoll.CreateListener("tcp", "0.0.0.0")
	//if err != nil {
	//	panic("create netpoll listener failed")
	//}
	//.

	eventLoop, _ := netpoll.NewEventLoop(
		handle,
		netpoll.WithOnPrepare(prepare),
		netpoll.WithReadTimeout(time.Second),
	)
}
