package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"
	"xeno/RPCTest/generated/strmtest"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		panic("conn Failed" + err.Error())
	}

	defer conn.Close()

	cli := strmtest.NewStreamTestClient(conn)
	r, err := cli.ServerStream(context.Background(), &strmtest.StreamRequest{
		Id: 10001,
	})
	if err != nil {
		panic("call" + err.Error())
	}

	for {
		data, err := r.Recv()
		if err != nil {
			fmt.Printf("recv done")
			break
		}
		fmt.Println("Got " + data.Data)
	}

	fmt.Println("Done ")

	s, _ := cli.ClientStream(context.Background())
	for i := 0; i < 10; i++ {
		s.Send(&strmtest.StreamRequest{Id: int64(i)})
		time.Sleep(time.Second)
	}

}
