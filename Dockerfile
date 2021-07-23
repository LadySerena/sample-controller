FROM golang:1.16 AS builder

WORKDIR /workspace
COPY Makefile ./
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY reconcile ./


RUN make build

FROM gcr.io/distroless/base

COPY --from=builder ./pod-timestamp-controller /pod-timestamp-controller
ENTRYPOINT ["/pod-timestamp-controller"]