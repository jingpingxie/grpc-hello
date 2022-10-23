package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	hello3 "grpc-hello/proto/hello_http"
	"net/http"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// grpc服务地址
	endpoint := Address
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// HTTP转grpc
	err := hello3.RegisterHelloHTTPHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		grpclog.Fatalf("Register handler err:%v\n", err)
	}

	fmt.Println("HTTP Listen on 8080")
	http.ListenAndServe(":8080", mux)
}
