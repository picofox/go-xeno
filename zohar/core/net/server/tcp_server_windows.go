package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

func inet_addr(ipaddr string) [4]byte {
	var (
		ips = strings.Split(ipaddr, ".")
		ip  [4]uint64
		ret [4]byte
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(ips[i], 10, 8)
	}
	for i := 0; i < 4; i++ {
		ret[i] = byte(ip[i])
	}
	return ret
}

func TestSock() {

	listen, er := net.Listen("tcp", "127.0.0.1:9090")
	if er != nil {
		fmt.Printf("listen failed, err:%v\n", er)
		return
	}
	for {
		// 等待客户端建立连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("accept failed, err:%v\n", err)
			continue
		}

		// 启动一个单独的 goroutine 去处理连接

		fmt.Printf(conn.LocalAddr().String())
	}

	fmt.Printf(listen.Addr().String())

	wsa := syscall.WSAData{}
	err := syscall.WSAStartup(uint32(514), &wsa)
	if err != nil {
		fmt.Printf("WsaStartup FAiled %s", err.Error())
	}

	sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket FAiled %d", sock)
	}

	k32dll := syscall.NewLazyDLL("kernel32.dll")
	sp := k32dll.NewProc("select")
	if sp == nil {
		fmt.Printf("error")
	}

	addr := syscall.SockaddrInet4{
		Port: 10000,
		Addr: inet_addr("0.0.0.0"),
	}
	err = syscall.Bind(sock, &addr)
	if err != nil {
		fmt.Printf("bind FAiled %s", err.Error())
	}

	err = syscall.Listen(sock, 1024)
	if err != nil {
		fmt.Printf("listen FAiled %s", err.Error())
	}

	var rsan int32
	var rawsa [2]syscall.RawSockaddrAny
	rsan = 116
	var afd syscall.Handle
	var qty uint32
	var o syscall.Overlapped = syscall.Overlapped{}

	afd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Printf("socket FAiled %d", sock)
	}

	err = syscall.AcceptEx(sock, afd, (*byte)(unsafe.Pointer(&rawsa[0])), 116*2, uint32(rsan), uint32(rsan), &qty, &o)
	if err != nil {
		fmt.Printf("acc failed %s", err.Error())
	}

	for {

	}

}
