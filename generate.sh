#!/usr/bin/env bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/service.proto
go clean -modcache
go get github.com/buhtigexa/naive-bayes@latest
go mod tidy

