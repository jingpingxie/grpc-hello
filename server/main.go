package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	hello "grpc-hello/proto"
	"net"
	"os"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

type server struct {
	hello.UnimplementedHelloServer
}

func (*server) SayHello(cn context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	grpclog.Infoln("request:", req.Name)
	return &hello.HelloResponse{Message: "hello," + req.Name}, nil
}

// 判断所给路径文件/文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}
func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Infof("failed to listen:%v\n", err)
		return
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("../cert/server/server.pem", "../cert/server/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	srv := grpc.NewServer(grpc.Creds(creds))
	hello.RegisterHelloServer(srv, &server{})
	defer func() {
		srv.Stop()
		listen.Close()
	}()
	fmt.Println("Listen on " + Address + " with TLS")
	err = srv.Serve(listen)
}
