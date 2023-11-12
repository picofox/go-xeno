package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"xeno/RPCTest/generated/account"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		panic("conn Failed" + err.Error())
	}

	defer conn.Close()

	cli := account.NewAccountServiceClient(conn)
	r, err := cli.OnRegister(context.Background(), &account.AccountRegister{
		Name:   "fox",
		Email:  "picobsd@qq.com",
		Passwd: "dandan",
	})
	if err != nil {
		panic("call" + err.Error())
	}

	fmt.Println("uid", r.Uid)

}
