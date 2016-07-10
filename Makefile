all: build test

build:
	godep go build

test:
	godep go test
