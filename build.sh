#! /bin/sh

export GO111MODULE=on
export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=1

export GOPROXY="https://proxy.golang.org"
export GOSUMDB=sum.golang.org

go mod tidy -v
go mod download -json
go build -v -o printbin-server
