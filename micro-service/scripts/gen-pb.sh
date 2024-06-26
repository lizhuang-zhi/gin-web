#!/bin/bash

# 生成pb.go和grpc pb.go文件

# 确保gen-pb目录存在且可写
mkdir -p protobuf/gen-pb
chmod -R u+w protobuf/gen-pb

# 遍历proto目录中的所有.proto文件
for file in protobuf/proto/*.proto; do
    # 生成普通的pb.go文件
    protoc --go_out=protobuf/gen-pb "$file"

    # 如果.proto文件包含gRPC服务，则生成对应的grpc pb.go文件
    protoc --go-grpc_out=protobuf/gen-pb "$file"
done
