REV:=$(shell git rev-parse HEAD)

all: build test

mod:
	go mod download

statik: mod
	statik || go install github.com/rakyll/statik
	statik -src=$(shell pwd)/public

build: statik
	go install ./...

pretest:
	gossip-db-init

test: pretest
	TZ=UTC go test ./...

clean:
	rm -rf gossip

docker-build: clean statik
	docker build -t parkr/gossip:$(REV) .

docker-test: docker-build
	docker run --name gosssip-test --rm -it parkr/gossip:$(REV)

docker-release: docker-build
	docker push parkr/gossip:$(REV)
