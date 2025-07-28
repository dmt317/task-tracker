.PHONY: build, run, test, lint

build:
	go build -o bin/app cmd/server/main.go

build-auth:
	go build -o bin/auth cmd/auth/main.go

run: build
	./bin/app

run-auth: build-auth
	./bin/auth

test:
	docker-compose up --build -d postgres migrate
	go test -v -race -cover ./...
	docker-compose down

lint:
	golangci-lint run --config=.golangci.yml
