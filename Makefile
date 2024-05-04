.DEFAULT_GOAL := all

BIN_NAME := chirpy

DEBUG ?= 0
ifeq ($(DEBUG), 1)
	FLAGS += --debug
endif

.PHONY: all build run test is-pretty

all: build run

build:
	go build -o $(BIN_NAME)

run:
	./$(BIN_NAME) $(FLAGS)

test:
	go test -v -cover ./...

is-pretty:
	-golangci-lint run
	-gofmt -d $(CURDIR) | colordiff
