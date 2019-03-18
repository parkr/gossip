FROM golang as builder
COPY . /go/src/github.com/parkr/gossip
RUN go version
RUN CGO_ENABLED=0 go install github.com/parkr/gossip

FROM golang as curler
RUN CGO_ENABLED=0 go get github.com/parkr/go-curl/cmd/go-curl

FROM scratch
HEALTHCHECK --interval=30s --timeout=3s \
  CMD [ "/bin/go-curl", "-f", "http://127.0.0.1:8080/_health" ]
COPY --from=curler /go/bin/go-curl /bin/go-curl
COPY --from=builder /go/bin/gossip /bin/gossip
ENTRYPOINT [ "/bin/gossip" ]
