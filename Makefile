SHELL=bash

.PHONY: all

all: clean fmt test

clean:
	go clean ./...

fmt:
	go fmt ./...

test:
	go test -v -count=1 ./...

