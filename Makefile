proto:
	mkdir -p ./proto/src
	protoc -I=. --go_out=./proto/src ./proto/api/port/port.proto

.PHONY:
	proto