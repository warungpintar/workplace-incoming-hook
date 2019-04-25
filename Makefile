.PHONY: build run install lint

export GO111MODULE=on

build:
	go build workplace-incoming-hook.go

run:
	go run workplace-incoming-hook.go

install:
	go mod download

lint:
	golangci-lint run
