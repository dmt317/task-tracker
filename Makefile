.PHONY: build, run, test, lint

build:
	go build -o bin/main cmd/server/main.go

run: build
	./bin/main

test:
	go test ./internal/storage/ -v -race -count=1
	go test ./internal/config -v -count=1

lint:
	golangci-lint run --config=.golangci.yml
