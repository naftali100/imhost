# -----------------------------------------------------------------------------
#  Build Stage
# -----------------------------------------------------------------------------
FROM golang:latest AS builder

COPY . /go/src/app
ENV CGO_ENABLED 0

WORKDIR /go/src/app

RUN go mod download && \
    go mod verify && \
    go test -v ./...

RUN go build \
    -ldflags="-s -w -extldflags \"-static\"" \
    -o /tmp/imhost \
    /go/src/app/imhost/main.go

# -----------------------------------------------------------------------------
#  Main Stage
# -----------------------------------------------------------------------------
FROM scratch

COPY --from=builder /tmp/imhost /usr/bin/imhost
EXPOSE 80/tcp
WORKDIR /tmp
ENTRYPOINT [ "/usr/bin/imhost" ]

HEALTHCHECK \
    --start-period=15s \
    --interval=5m \
    --timeout=15s \
    --retries=3 \
    CMD [ "/usr/bin/imhost", "--ping"]
