FROM golang:1.16 AS builder
WORKDIR /workspace
COPY --from=golangci/golangci-lint:v1.41.1 /usr/bin/golangci-lint /usr/bin/golangci-lint
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile ./
COPY main.go ./
COPY reconcile/ ./reconcile
COPY prestart/ ./prestart
RUN make build lint test

FROM gcr.io/distroless/base
COPY --from=builder /workspace/pod-timestamp-controller /pod-timestamp-controller
ENTRYPOINT ["/pod-timestamp-controller"]