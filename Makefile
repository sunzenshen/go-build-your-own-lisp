all: dependencies test

test:
	go test -v ./...

dependencies:
	which go

.PHONY: all dependencies test
