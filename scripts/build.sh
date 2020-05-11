#!/bin/bash

cd $(dirname $0)/..

echo "Fetching dependecies..."
go get
mkdir -p build/linux build/darwin
echo "Build for linux/amd64"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/terraform-provider-mongodb_$1
echo "Build for darwin/amd64"
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/darwin/terraform-provider-mongodb_$1

echo "Compress..."
mkdir -p release
cd build/linux
GZIP=-9 tar -czf ../../release/terraform-provider-mongodb_$1_linux_amd64.tar.gz .
cd ../../build/darwin
GZIP=-9 tar -czf ../../release/terraform-provider-mongodb_$1_darwin_amd64.tar.gz .

rm -rf ../../build
