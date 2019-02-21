FROM golang
WORKDIR /go/src/github.com/parkr/gossip
COPY . .
RUN go version
RUN CGO_ENABLED=0 go install github.com/parkr/gossip

FROM scratch
COPY --from=0 /go/bin/gossip /bin/gossip
EXPOSE 3306
ENTRYPOINT [ "/bin/gossip" ]
