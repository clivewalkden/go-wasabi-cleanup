.DEFAULT_GOAL := build

GO_BIN=${HOME}/go/go1.16.15/bin/go
EXECUTABLE=wasabi-cleanup
VERSION=v1.2.0

fmt:
	${GO_BIN} fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	${GO_BIN} vet ./...
.PHONY:vet

build: vet
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 ${GO_BIN} build -o bin/${VERSION}/${EXECUTABLE}-freebsd-amd64 main.go
	GOOS=darwin GOARCH=arm64 ${GO_BIN} build -o bin/${VERSION}/${EXECUTABLE}-macos-arm64 main.go
	GOOS=linux GOARCH=amd64 ${GO_BIN} build -o bin/${VERSION}/${EXECUTABLE}-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 ${GO_BIN} build -o bin/${VERSION}/${EXECUTABLE}-windows-amd64.exe main.go
.PHONY:build