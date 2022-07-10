/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"flag"
	"strconv"
	"net"
	"skywatcher/etcd_lb"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

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
	var interval uint64 = 20 // 心跳周期 (s)
	var lifeCycle uint64 = 5 // 生命周期 (interval)

	// 注册服务
	if err := etcd_lb.Register(projectName, moduleName, addr, targetAddr, userName, passWord, interval, lifeCycle); err != nil {
		log.Printf("Register failed, %+v\n", err)
	} else {
		log.Printf("Register success\n")
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
