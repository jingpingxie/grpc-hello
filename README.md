windows 下安装go protoc

# 下载最新版的protoc安装程序

https://github.com/protocolbuffers/protobuf/releases/download/v21.8/protoc-21.8-win64.zip

加压后把protoc.exe放入gopath\bin

# 安装protoc-gen-go

go get google.golang.org/grpc

go install google.golang.org/protobuf/cmd/protoc-gen-go

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

# 将proto转换成go

protoc --go_out=. hello.proto

protoc --go-grpc_out=. hello.proto

# 使用openssl生成SAN证书
下载安装
http://slproweb.com/download/Win64OpenSSL_Light-3_0_5.msi

参考：
https://blog.csdn.net/a145127/article/details/126311442

1.创建一个“cert”目录用于，保存证书和配置文件。
2.创建配置文件(openssl.cnf)，并保存到“cert”目录下。
3.生成根证书（rootCa）
使用命令行工具，进入到“cert”目录下，并执行如下命令：
生成私钥，密码可以输入123456
$ openssl genrsa -des3 -out ca.key 2048

用私钥来签名证书
$ openssl req -new -key ca.key -out ca.csr

使用私钥+证书来生成公钥
$ openssl x509 -req -days 365 -in ca.csr -signkey ca.key -out ca.crt
4.在“cert”目录下，创建“server”目录，它们用来保存服务器密钥。
5.生成服务器密钥。
使用命令行工具，进入到“cert”目录下，并执行如下命令：
生成服务器私钥，密码输入123456
$ openssl genpkey -algorithm RSA -out server/server.key

使用私钥来签名证书
$ openssl req -new -nodes -key server/server.key -out server/server.csr -config openssl.cnf -extensions 'v3_req'

生成SAN证书
$ openssl x509 -req -in server/server.csr -out server/server.pem -CA ca.crt -CAkey ca.key -CAcreateserial -extfile ./openssl.cnf -extensions 'v3_req'

# grpc-gateway
参考：https://blog.csdn.net/Mr_XiMu/article/details/125000670
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

cd E:\source\gopath\src\grpc-hello\proto
编译hello_http.proto
protoc -I . --go_out ./ --go_opt paths=source_relative hello_http/hello_http.proto
protoc -I . --go-grpc_out ./ --go-grpc_opt paths=source_relative hello_http/hello_http.proto

编译hello_http.proto gateway
protoc -I . --grpc-gateway_out ./ --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative hello_http/hello_http.proto