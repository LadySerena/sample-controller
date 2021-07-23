BINARY_NAME=pod-timestamp-controller

build:
	go build -o ${BINARY_NAME} main.go

test:
	go test -v ./...