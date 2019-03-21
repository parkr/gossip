REV=$(shell git rev-parse HEAD)

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

docker-build: clean statik
	docker build -t parkr/gossip:$(REV) .

docker-test: docker-build
	docker run --name gossip-test --rm -it parkr/gossip:$(REV)

docker-release: docker-build
	docker push parkr/gossip:$(REV)
