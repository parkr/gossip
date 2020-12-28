REV:=$(shell git rev-parse HEAD)
GOSSIP_DB_PATH:=$(shell pwd)/data/gossip_test.sqlite3

all: build test

mod-download:
	go mod download

tools: mod-download
	cat cmd/gossip/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

.PHONY: statik
statik: tools
	statik -f -src=$(shell pwd)/public

build: statik
	go install ./...

pretest:
	GOSSIP_DB_PATH=$(GOSSIP_DB_PATH) gossip-db-init

test: pretest statik
	TZ=UTC GOSSIP_DB_PATH=$(GOSSIP_DB_PATH) go test ./...

clean:
	rm -rf gossip

docker-build: clean
	docker build -t parkr/gossip:$(REV) .

docker-test: docker-build
	docker run --name gossip-test --rm -it parkr/gossip:$(REV)

docker-release: docker-build
	docker push parkr/gossip:$(REV)
