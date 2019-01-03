@echo off
call ./scripts/win/install-protobuf.bat

go get github.com/golang/protobuf/protoc-gen-go
set GO111MODULE=off
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
set GO111MODULE=on
