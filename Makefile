.DEFAULT_GOAL := all

BIN_NAME := chirpy

.PHONY:all
all: build run

.PHONY:build
build:
	go build -o $(BIN_NAME)

.PHONY:run
run:
	./$(BIN_NAME)

.PHONY:test
test:
	go test -v -cover ./...

.PHONY:is-pretty
is-pretty:
	-golangci-lint run
	-gofmt -d $(CURDIR) | colordiff
