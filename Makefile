.DEFAULT_GOAL := build

GO_BIN=${GOROOT}/bin/go
EXECUTABLE=wasabi-cleanup
VERSION=v1.3.2

GO_MAJOR_VERSION = $(shell $(GO_BIN) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell $(GO_BIN) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
MINIMUM_SUPPORTED_GO_MAJOR_VERSION = 1
MINIMUM_SUPPORTED_GO_MINOR_VERSION = 19
GO_VERSION_VALIDATION_ERR_MSG = Golang version is not supported, please update to at least $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION).$(MINIMUM_SUPPORTED_GO_MINOR_VERSION)

validate-go-version: ## Validates the installed version of go against Mattermost's minimum requirement.
	@if [ $(GO_MAJOR_VERSION) -gt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		exit 0 ;\
	elif [ $(GO_MAJOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MAJOR_VERSION) ]; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	elif [ $(GO_MINOR_VERSION) -lt $(MINIMUM_SUPPORTED_GO_MINOR_VERSION) ] ; then \
		echo '$(GO_VERSION_VALIDATION_ERR_MSG)';\
		exit 1; \
	fi

$(info Compiling with $(GOROOT))

fmt: validate-go-version
	${GO_BIN} fmt ./...
.PHONY:fmt

lint: fmt
	golangci-lint run ./...
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