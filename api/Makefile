.PHONY: build docker swagger

docker:
	@make build
	docker build . -t api-service:latest

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o ./build/server ./cmd/server/main.go


swagger:
	@echo "Build swagger openapis binary data..."
	@go-bindata --nocompress -pkg swagger -o web/ui/swagger/data/datafile.go third_party/swagger-ui/...
	@echo ✓ compiled;