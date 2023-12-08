package main

import (
	_ "github.com/go-swagger/go-swagger/cmd/swagger"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/vektra/mockery/v2"
	_ "gopkg.in/reform.v1/reform"
)
