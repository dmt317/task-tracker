.PHONY: build, run, test, lint

build:
	go build -o bin/main main.go

run: build
	./bin/main

test:
	go test ./storage/ -v -race -count=1

lint:
	golangci-lint run --config=.golangci.yml
