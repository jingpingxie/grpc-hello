windows 下安装go protoc


1、下载最新版的golang安装程序

https://github.com/protocolbuffers/protobuf/releases/download/v21.8/protoc-21.8-win64.zip

加压后把protoc.exe放入gopath\bin

2、安装protoc-gen-go

go get google.golang.org/grpc

go install google.golang.org/protobuf/cmd/protoc-gen-go

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

3、将proto转换成go

protoc --go_out=. hello.proto

protoc --go-grpc_out=. hello.proto
