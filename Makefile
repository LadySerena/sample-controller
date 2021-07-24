BINARY_NAME=pod-timestamp-controller

build:
	go build -o ${BINARY_NAME} main.go

lint:
	 golangci-lint -v run ./...

test:
	go test -v ./...

create-local-env:
	./scripts/create-kind-with-registry.sh

destroy-local-env:
	./scripts/destroy-kind-and-registry.sh
