FROM golang:1.16 AS builder

WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY Makefile ./
COPY main.go ./
COPY reconcile/ ./reconcile


RUN make build

FROM gcr.io/distroless/base

COPY --from=builder /workspace/pod-timestamp-controller /pod-timestamp-controller
ENTRYPOINT ["/pod-timestamp-controller"]