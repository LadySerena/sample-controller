BINARY_NAME=pod-timestamp-controller
TAG=foo
LOCAL_DOCKER_REPO=localhost:5000/sample-controller:${TAG}

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

install: docker-push-local
	helm install serena-test ./config/sample-controller --set image.tag=${TAG}

uninstall:
	helm uninstall serena-test ./config/sample-controller

upgrade: docker-push-local
	helm upgrade serena-test ./config/sample-controller --set image.tag=${TAG}

destroy-local-env:
	./scripts/destroy-kind-and-registry.sh
