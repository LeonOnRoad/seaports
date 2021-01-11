DOCKER_IMAGE_TAG=local-build

proto-gen:
	mkdir -p ./proto/src
	protoc -I=. --go_out=./proto/src --go-grpc_out=./proto/src ./proto/api/port/port.proto

build: proto-gen
	go mod vendor
	mkdir -p ./bin
	CGO_ENABLED=0 go build -o ./bin/ ./services/...

images: build
	docker build -f ./services/ports-api/Dockerfile . -t "ports-api:$(DOCKER_IMAGE_TAG)"
	docker build -f ./services/ports-service/Dockerfile . -t "ports-service:$(DOCKER_IMAGE_TAG)"

.PHONY:
	proto-gen build images