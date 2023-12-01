.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY: fmt

vet: fmt
	go vet ./...
.PHONY: vet

lint: fmt vet
	staticcheck ./...
.PHONY: lint

clean:
	rm -rf bin/
.PHONY: clean

build-linux: lint clean
	GOOS=linux GOARCH=amd64 go build -o bin/linux/ ./...
.PHONY: build-linux

build-mac: lint clean
	GOOS=darwin GOARCH=amd64 go build -o bin/mac/ ./...
.PHONY: build-mac

build-windows: lint clean
	GOOS=windows GOARCH=amd64 go build -o bin/windows/ ./...
.PHONY: build-windows

build: build-windows build-mac build-linux
.PHONY: build