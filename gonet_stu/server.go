package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/netpoll"
	"runtime"
	"time"
)

func OnReq(ctx context.Context, connection netpoll.Connection) error {
	return nil
}

func OnPrepare(connection netpoll.Connection) context.Context {
	return context.Background()
}

func main() {
	fmt.Println(runtime.GOMAXPROCS(0))
	listener, err := netpoll.CreateListener("tcp", "0.0.0.0:9999")
	if err != nil {
		panic("create netpoll listener failed")
	}

	if listener == nil {
		return
	}

	el, _ := netpoll.NewEventLoop(OnReq, netpoll.WithOnPrepare(OnPrepare), netpoll.WithReadTimeout(time.Second))
	if el == nil {
		fmt.Print("fail")
	}

	el.Serve(listener)

}
