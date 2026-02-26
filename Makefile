VERSION ?= dev

.PHONY: build run dev docker-build docker-up clean lint release

build:
	go build -ldflags "-X altcha/pkg/handler.Version=$(VERSION)" -o bin/server ./cmd/server

run: build
	./bin/server

dev:
	air

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
