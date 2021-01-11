proto:
	mkdir -p ./proto/src
	protoc -I=. --go_out=./proto/src --go-grpc_out=./proto/src ./proto/api/port/port.proto

build: proto
	mkdir -p ./bin
	CGO_ENABLED=0 go build -o ./bin/ ./services/...

.PHONY:
	proto