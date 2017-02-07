all: build test

build:
	godep go build

test:
	godep go test ./...

clean:
	rm -rf gossip

docker: clean
	docker build -t parkr/gossip:$(shell git rev-parse HEAD) .
	docker build -t parkr/gossip:latest .

publish-to-docker:
	docker push parkr/gossip
