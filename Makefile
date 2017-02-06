all: build test

build:
	godep go build

test:
	godep go test

docker:
	docker build -t gossip:$(shell git rev-parse HEAD) .
	docker build -t gossip:latest .

publish-to-docker:
	docker push parkr/gossip
