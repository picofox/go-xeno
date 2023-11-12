package main

import (
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
	"strconv"
	"time"
	"xeno/RPCTest/generated/strmtest"
)

type Server struct {
	strmtest.UnimplementedStreamTestServer
}

func (ego *Server) ServerStream(request *strmtest.StreamRequest, server strmtest.StreamTest_ServerStreamServer) error {

	for i := 0; i < 10; i++ {
		server.Send(&strmtest.StreamReply{
			Data: "reply no" + strconv.Itoa(i),
		})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}

	fmt.Println("Server Stream Done")
	return nil

}

func (ego *Server) ClientStream(cli strmtest.StreamTest_ClientStreamServer) error {
	for {
		a, err := cli.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("done")
			} else {
				fmt.Println("recive error" + err.Error())
			}

			break
		} else {
			fmt.Println("got data" + strconv.FormatInt(a.Id, 10))

		}
	}

	return nil
}

func (ego *Server) DualStream(server strmtest.StreamTest_DualStreamServer) error {
	return nil
}

func main() {
	gs := grpc.NewServer()
	strmtest.RegisterStreamTestServer(gs, &Server{})
	listen, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		panic("Failed listen" + err.Error())
	}
	err = gs.Serve(listen)
	if err != nil {
		panic("Failed toStart " + err.Error())
	}

	fmt.Println("RPC Started")

}
