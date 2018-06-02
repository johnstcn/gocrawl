SHELL=bash

.PHONY: all

all: clean generate fmt test install

clean:
	go clean ./...

generate:
	go generate ./...

fmt:
	go fmt ./...

test:
	go test ./...

install:
	go install ./...

