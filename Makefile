VERSION ?= dev

.PHONY: build run dev docker-build docker-up clean lint

build:
	go build -ldflags "-X altcha/pkg/handler.Version=$(VERSION)" -o bin/server ./cmd/server

run: build
	./bin/server

dev:
	go run ./cmd/server

docker-build:
	docker compose build

docker-up:
	docker compose up --build

clean:
	rm -rf bin/

lint:
	golangci-lint run ./...
