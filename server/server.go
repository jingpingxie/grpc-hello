package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"

	hello "grpc-hello/proto"
	"net"
)

type server struct {
	hello.UnimplementedHelloServer
}

func (*server) SayHello(cn context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	fmt.Println("request:", req.Name)
	return &hello.HelloResponse{Message: "hello," + req.Name}, nil
}
func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("failed to listen:%v\n", err)
		return
	}
	srv := grpc.NewServer()
	hello.RegisterHelloServer(srv, &server{})
	defer func() {
		srv.Stop()
		listen.Close()
	}()
	fmt.Println("Serving 8080")
	err = srv.Serve(listen)
}
