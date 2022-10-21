package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	gw "grpc-hello/proto/hello_http"
	"net/http"
	"time"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
	// OpenTLS 是否开启TLS认证
	OpenTLS = true
)

// customCredential 自定义认证
//定义了一个customCredential结构，并实现了两个方法GetRequestMetadata和RequireTransportSecurity。
//这是gRPC提供的自定义认证方式，每次RPC调用都会传输认证信息。
//customCredential其实是实现了grpc/credential包内的PerRPCCredentials接口
type customCredential struct{}

// GetRequestMetadata 实现自定义认证接口
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
		//"appkey": "i am not key",
	}, nil
}

// RequireTransportSecurity 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}

// interceptor 客户端拦截器
func interceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	grpclog.Infof("method=%s req=%v rep=%v duration=%s error=%v\n", method, req, reply, time.Since(start), err)
	return err
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// grpc服务地址
	endpoint := Address
	mux := runtime.NewServeMux()

	var opts []grpc.DialOption
	if OpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile("../cert/server/server.pem", "*.example.com")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// HTTP转grpc
	err := gw.RegisterHelloHTTPHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		grpclog.Fatalf("Register handler err:%v\n", err)
	}

	fmt.Println("HTTP Listen on 5555")
	http.ListenAndServe(":5555", mux)
}
