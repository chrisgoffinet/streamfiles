grpc:
	protoc -I api -I${GOPATH}/src api/api.proto --go_out=plugins=grpc:api

build:
	go build