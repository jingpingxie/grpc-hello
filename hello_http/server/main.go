package main

import (
	"context"
	"fmt"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	//hello2 "grpc-hello/hello_http/proto"
	hello3 "grpc-hello/proto/hello_http"
	"net"
	"net/http"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:50052"
)

//type HelloService struct {
//	hello2.UnimplementedHelloServer
//}

type HelloService2 struct {
	hello3.UnimplementedHelloHTTPServer
}

//func (*HelloService) SayHello(ctx context.Context, in *hello2.HelloRequest) (*hello2.HelloResponse, error) {
//	resp := new(hello2.HelloResponse)
//	resp.Message = fmt.Sprintf("Hello %s.", in.Name)
//
//	return resp, nil
//}

func (*HelloService2) SayHello(ctx context.Context, in *hello3.HelloHTTPRequest) (*hello3.HelloHTTPResponse, error) {
	resp := new(hello3.HelloHTTPResponse)
	resp.Message = fmt.Sprintf("Hello %s.", in.Name)

	return resp, nil
}

// auth 验证Token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return status.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}

// interceptor 拦截器
func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	// 继续处理请求
	return handler(ctx, req)
}

func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}

	grpc.EnableTracing = true
	go http.ListenAndServe(":50051", nil)
	grpclog.Infoln("Trace listen on 50051")
}
func main() {
	// 开启trace
	go startTrace()

	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Infof("failed to listen:%v\n", err)
		return
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("../../cert/server/server.pem", "../../cert/server/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	var opts []grpc.ServerOption
	opts = append(opts, grpc.Creds(creds))
	// 注册interceptor
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	srv := grpc.NewServer(opts...)
	//hello2.RegisterHelloServer(srv, &HelloService{})
	hello3.RegisterHelloHTTPServer(srv, &HelloService2{})
	defer func() {
		srv.Stop()
		listen.Close()
	}()

	fmt.Printf("Listen on " + Address + " with TLS + Token")
	err = srv.Serve(listen)
}
