package main

import (
	"fmt"
	"net"
	"time"
)

func accr(lis net.Listener) {
	for {
		fmt.Println("等待客户端来连接...")
		conn, err := lis.Accept()
		//错误处理和输出当前连接客户端信息
		if err != nil {
			fmt.Println("Accept() err=", err)
		} else {
			fmt.Printf("Accept() suc con: %s -> %s\n", conn.RemoteAddr().String(), conn.LocalAddr().String())
		}

	}

}

func main() {
	fmt.Printf("服务器开始监听...")
	//监听一个端口
	listen, err := net.Listen("tcp", ":9998")
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	//延迟关闭连接
	defer listen.Close()
	listen2, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("listen err=", err)
		return
	}
	//延迟关闭连接
	defer listen2.Close()

	go accr(listen2)
	go accr(listen)

	//循环等待客户端来连接
	for {
		time.Sleep(10 * time.Second)
	}

}
