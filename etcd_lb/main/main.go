package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"skywatcher/etcd_lb"
	"strconv"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	var port, _ = strconv.Atoi(flag.Arg(0))
	log.Printf("listen port=%d\n", port)

	var projectName = "gRPC project1"
	var moduleName = "module1"
	var addr = fmt.Sprintf("%s:%d", "127.0.0.1", port)

	// 注册中心
	var targetAddr = "http://127.0.0.1:12379"
	var userName = ""
	var passWord = ""
	var interval uint64 = 20 // 心跳周期(s)
	var lyfeCycle uint64 = 5 // 生命周期(interval

	// 注册服务
	if err := etcd_lb.Register(projectName, moduleName, addr, targetAddr, userName, passWord, interval, lyfeCycle); err != nil {
		log.Printf("Register failed, %+v\n", err)
	} else {
		log.Printf("Register success\n")
	}

	if lis, err := net.Listen("tcp", addr); err != nil {
		panic(err)
	} else {
		s := grpc.NewServer()
		pb.RegisterDemoServer(s, new(service.MyServer))

		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}
}
