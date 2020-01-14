REV:=$(shell git rev-parse HEAD)

all: build test

mod-download:
	go mod download

tools: mod-download
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: statik
statik: tools
	statik -src=$(shell pwd)/public

build: statik
	go install ./...

pretest:
	gossip-db-init

test: pretest statik
	TZ=UTC go test ./...

clean:
	rm -rf gossip

docker-build: clean
	docker build -t parkr/gossip:$(REV) .

docker-test: docker-build
	docker run --name gosssip-test --rm -it parkr/gossip:$(REV)

docker-release: docker-build
	docker push parkr/gossip:$(REV)
