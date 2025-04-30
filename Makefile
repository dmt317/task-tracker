.PHONY: build, run, test, lint

build:
	go build -o bin/main cmd/server/main.go

run: build
	./bin/main

test:
	go test ./...

lint:
	golangci-lint run --config=.golangci.yml
