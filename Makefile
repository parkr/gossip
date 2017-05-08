all: build test

build:
	godep go build

test:
	TZ=UTC godep go test ./...

clean:
	rm -rf gossip

docker: clean
	docker build -t parkr/gossip:$(shell git rev-parse HEAD) .

publish-to-docker:
	docker push parkr/gossip:$(shell git rev-parse HEAD)
