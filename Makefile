GOPATH:=$(shell go env GOPATH)

proto:
	protoc -I . \
	-I${GOPATH}/src \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
	--grpc-gateway_out=logtostderr=true,repeated_path_param_separator=csv,paths=source_relative:. \
	--swagger_out=allow_delete_body=true,allow_merge=true,repeated_path_param_separator=csv,logtostderr=true:. \
	--go_out=plugins=grpc,paths=source_relative:. protos/*.proto; \

	@mv -f apidocs.swagger.json api/apis/v1/apidocs.swagger.json

	@echo ✓ protobuf compiled; \

