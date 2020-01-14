FROM golang as builder
WORKDIR /srv/app
COPY . /srv/app
RUN go version
RUN CGO_ENABLED=0 go install ./...

FROM scratch
HEALTHCHECK --start-period=1ms --interval=30s --timeout=5s --retries=1 \
  CMD [ "/bin/gossip-healthcheck" ]
COPY --from=builder /go/bin/gossip-healthcheck /bin/gossip-healthcheck
COPY --from=builder /go/bin/gossip /bin/gossip
ENTRYPOINT [ "/bin/gossip" ]
