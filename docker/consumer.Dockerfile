FROM golang:1.23-alpine AS builder
WORKDIR /build
COPY . .
RUN go get -d -v ./...
RUN go build -ldflags="-w -s" -o consumer-cli cmd/consumer-cli/main.go
FROM scratch
COPY --from=builder /build/consumer-cli /go/bin/consumer-cli
ENTRYPOINT ["/go/bin/consumer-cli"]
