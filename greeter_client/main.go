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

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"time"
	"flag"
	"skywatcher/etcd_lb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	flag.Parse()
	var projectName = "gRPC project1"
	var moduleName = "module1"

	// 注册中心
	var targetAddr = "http://127.0.0.1:12379"
	var userName = ""
	var passWord = ""
	r := etcd_lb.NewResolver(projectName, moduleName, userName, passWord, targetAddr)
	resolver.Register(r)
	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// Set up a connection to the server.
	conn, err := grpc.DialContext(context.Background(), r.Scheme()+"://authority/"+targetAddr, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancerName(roundrobin.Name))
	// cancel()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		c := pb.NewGreeterClient(conn)
		name := "client"
		req := &pb.HelloRequest{Name: name}
		log.Printf("client send:[%s]\n", req.Name)
		if rsp, err := c.SayHello(context.Background(), req); err != nil {
			log.Printf("client send failed, err=%+v\n", err)
		} else {
			log.Printf("client receive:[%s]\n", rsp.Message)
		}
	}
}
