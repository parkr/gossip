FROM golang as builder
COPY . /go/src/github.com/parkr/gossip
RUN go version
RUN CGO_ENABLED=0 go install github.com/parkr/gossip/...

FROM scratch
HEALTHCHECK --start-period=1ms --interval=30s --timeout=5s --retries=1 \
  CMD [ "/bin/gossip-healthcheck" ]
COPY --from=builder /go/bin/gossip-healthcheck /bin/gossip-healthcheck
COPY --from=builder /go/bin/gossip /bin/gossip
ENTRYPOINT [ "/bin/gossip" ]
