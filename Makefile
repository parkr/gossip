all: build test

statik:
	statik || go get github.com/rakyll/statik
	statik -src=$(shell pwd)/public

build: statik
	godep go build

test:
	TZ=UTC godep go test ./...

clean:
	rm -rf gossip

docker: clean statik
	docker build -t parkr/gossip:$(shell git rev-parse HEAD) .

publish-to-docker:
	docker push parkr/gossip:$(shell git rev-parse HEAD)
