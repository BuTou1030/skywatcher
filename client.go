package main

import (
	"context"
	"flag"
	"log"
	"skywatcher/etcd_lb"
	"time"

	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	var projectName = "gRPC project1"
	var moduleName = "module1"

	// 注册中心
	var targetAddr = "http://127.0.0.1:12379"
	var userName = ""
	var passWord = ""

	r := etcd_lb.NewResolver(projectName, moduleName, userName, passWord)
	b := grpc.RoundRobin(r)

	if conn, err := grpc.DialContext(context.Background(), targetAddr, grpc.WithInsecure(), grpc.WithBalancer(b)); err != nil {
		panic(err)
	} else {
		ticker := time.NewTicker(1 * time.Second)

		for _ = range ticker.C {
			client := pb.NewDemoClient(conn)
			req := &pb.SayReq{Msg: "Hi"}
			log.Printf("client send:[%s]\n", req.Msg)
			if rsp, err := client.Say(context.Background(), req); err != nil {
				log.Printf("client send failed, err=%+v\n", err)
			} else {
				log.Printf("client receive:[%s]\n", rsp.Msg)
			}
		}
	}
}
