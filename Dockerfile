FROM golang

WORKDIR /go/src/github.com/parkr/gossip

EXPOSE 3306

ADD . .

RUN go version

# Compile a standalone executable
RUN CGO_ENABLED=0 go install github.com/parkr/gossip

# Make `gossip` available to the `Dockerfile.release` build
CMD [ "./gossip" ]
