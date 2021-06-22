.PHONY: build clean docker

default: build

all: build docker

build: build-darwin build-linux

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/tksinfo-darwin-amd64 ./cmd/server/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/tksinfo-appclient-darwin-amd64 ./examples/client.go

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/tksinfo-linux-amd64 ./cmd/server/
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/tksinfo-appclient-linux-amd64 ./examples/client.go

clean:
	rm -rf ./bin

docker:
	docker build --no-cache -t tks-info/tks-info -f Dockerfile .
