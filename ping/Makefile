.PHONY: build docker

docker:
	@make build
	docker build . -t ping-service:latest

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w' -i -o ./build/server ./cmd/server/main.go
