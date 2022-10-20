package main

import (
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	hello "grpc-hello/proto"
	"os"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("../cert/server/server.pem", "*.example.com")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	client := hello.NewHelloClient(conn)

	rsp, err := client.SayHello(context.Background(), &hello.HelloRequest{Name: "JingpingXie"})
	if err != nil {
		grpclog.Fatalln(err)
	}
	fmt.Println(rsp)
	fmt.Println("按回车键退出程序...")
	in := bufio.NewReader(os.Stdin)
	_, _, _ = in.ReadLine()
}
