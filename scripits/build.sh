#!/usr/bin/env bash

OUTPUT=_output
ROOTPATH="/Users/donggang/Documents/Code/blog/go-grpc-http-rest-microservice-tutorial/"

# build gRPC client
echo -e "build grpc client..."
cd "${ROOTPATH}/cmd/client-grpc"
go build -o "${ROOTPATH}/${OUTPUT}/grpc-client" -v
echo -e "build grpc client sucessfully"

## build gRPC server
echo -e "build grpc server..."
cd "${ROOTPATH}/cmd/server"
go build -o "${ROOTPATH}/${OUTPUT}/grpc-server" -v
echo -e "build grpc server sucessfully"

