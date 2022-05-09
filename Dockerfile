FROM golang:1.18-buster as builder
WORKDIR /srv/app
COPY . /srv/app
RUN go version
RUN make statik
RUN go install ./...

FROM debian:buster-slim
HEALTHCHECK --start-period=1ms --interval=30s --timeout=5s --retries=1 \
  CMD [ "/bin/gossip-healthcheck" ]
COPY --from=builder /go/bin/gossip-healthcheck /bin/gossip-healthcheck
COPY --from=builder /go/bin/gossip /bin/gossip
ENTRYPOINT [ "/bin/gossip" ]
