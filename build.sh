#!/bin/bash
rm -rf bin
GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -mod=vendor -ldflags="-s -w -X main.version=${VERSION}" -a -o ./bin/mkdocs-generator

