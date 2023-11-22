package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"xeno/RPCTest/generated/account"
)

type ACCServer struct {
	account.UnimplementedAccountServiceServer
}

func (ego *ACCServer) OnRegister(ctx context.Context, register *account.AccountRegister) (*account.AccountRegisterResult, error) {
	fmt.Println("reg get")
	return &account.AccountRegisterResult{
		Ok:  false,
		Uid: 0,
	}, nil
}

func MyUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	fmt.Println("lanjie le ")
	return handler(ctx, req)
}

func main() {
	//req := pbs.AccountRegister{
	//	Name:   "fox",
	//	Email:  "picobsd@sina.com",
	//	Passwd: "yaoyaolingxian",
	//}
	//
	//rsp, _ := proto.Marshal(&req)
	//
	//newReq := pbs.AccountRegister{}
	//proto.Unmarshal(rsp, &newReq)
	opt := grpc.UnaryInterceptor(MyUnaryServerInterceptor)
	gs := grpc.NewServer(opt)
	account.RegisterAccountServiceServer(gs, &ACCServer{})
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
