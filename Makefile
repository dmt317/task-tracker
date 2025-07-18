.PHONY: build, run, test, lint

build:
	go build -o bin/main cmd/server/main.go

run: build
	./bin/main

test:
	docker-compose up --build -d postgres migrate
	go test -cover ./...
	docker-compose down

lint:
	golangci-lint run --config=.golangci.yml
