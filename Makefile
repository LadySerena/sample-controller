BINARY_NAME=pod-timestamp-controller
LOCAL_DOCKER_REPO=localhost:5000/sample-controller:latest

build:
	go build -o ${BINARY_NAME} main.go

docker-build:
	docker build -t ${LOCAL_DOCKER_REPO} .

docker-push-local: docker-build
	docker push ${LOCAL_DOCKER_REPO}

lint:
	 golangci-lint -v run ./...

test:
	go test -v ./...

create-local-env:
	./scripts/create-kind-with-registry.sh

destroy-local-env:
	./scripts/destroy-kind-and-registry.sh
