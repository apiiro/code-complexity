GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=complexity

.PHONY: all test build vendor

all: vet lint test clean build verify-binaries compress-bin

build: build-osx build-linux

build-osx:
ifneq ($(shell uname -s),Darwin)
	$(error this makefile assumes you're building from mac env)
endif
	GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -o bin/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-osx-x86_64 .
	GOARCH=arm64 GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -o bin/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-osx-arm64 .

build-linux:
	docker run --rm -v $(shell pwd):/app -w /app golang:1.22-alpine /bin/sh -c "GOARCH=amd64 GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -o bin/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-linux-x86_64 ."
	docker run --rm -v $(shell pwd):/app -w /app golang:1.22-alpine /bin/sh -c "GOARCH=arm64 GO111MODULE=on CGO_ENABLED=0 $(GOCMD) build -o bin/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-linux-aarch64 ."

clean:
	rm -rf ./bin

vet:
	go vet

lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:latest-alpine golangci-lint run

test:
	$(GOTEST) -v ./...

verify-binaries:
	$(info running binaries verification on osx, alpine, debiand, centos)
	./bin/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-osx-x86_64 --version
	./bin/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-osx-arm64 --version
	docker run --rm -v $(shell pwd)/bin:/app -w /app alpine /app/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-linux-aarch64 --version
	docker run --rm -v $(shell pwd)/bin:/app -w /app alpine /app/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-linux-x86_64 --version
	docker run --rm -v $(shell pwd)/bin:/app -w /app debian:buster /app/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-linux-aarch64 --version
	docker run --rm -v $(shell pwd)/bin:/app -w /app debian:buster /app/$(BINARY_NAME)-$(shell $(GOCMD) run . --version | cut -d" " -f 3)-linux-x86_64 --version

compress-bin:
	find bin -type f -print -exec zip -j '{}'.zip '{}' \;

benchmark:
	go test -tags bench ./...
