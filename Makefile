.PHONY: all build test fmt lint install-hooks clean

all: build

build:
	go build -v -o bin/app ./cmd/app

test:
	go test ./...

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

install-hooks:
	sh scripts/install-hooks.sh

clean:
	rm -rf bin
