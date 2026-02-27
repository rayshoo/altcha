VERSION ?= dev

.PHONY: build build-dashboard build-all run dev dev-server dev-dashboard psql docker-build docker-up clean lint release

build:
	go build -ldflags "-X altcha/pkg/handler.Version=$(VERSION)" -o bin/server ./cmd/server

build-dashboard:
	go build -ldflags "-X altcha/pkg/handler.Version=$(VERSION)" -o bin/dashboard ./cmd/dashboard

build-all: build build-dashboard

run: build
	./bin/server

dev-server:
	air -c .air.server.toml

dev-dashboard:
	air -c .air.dashboard.toml

dev: dev-server

psql:
	docker compose up -d postgres

docker-build:
	docker compose build

docker-up:
	docker compose up --build

clean:
	rm -rf bin/

lint:
	golangci-lint run ./...

release:
	@if [ -z "$(version)" ]; then echo "Usage: make release version=x.x.x"; exit 1; fi
	sed -i 's/newTag: .*/newTag: $(version)/' examples/k8s/kustomization.yaml
	git add examples/k8s/kustomization.yaml
	git commit -m "release: $(version)"
	git tag -f $(version)
	git push origin HEAD --tags
