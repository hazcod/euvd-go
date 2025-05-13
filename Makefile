
all: run

run:
	go run ./cmd/... -log=debug

lookup:
	go run ./cmd/... -log=debug lookup EUVD-2025-14349

search:
	go run ./cmd/... -log=debug search

build:
	$$GOPATH/bin/goreleaser build --config=.github/goreleaser.yml --clean --snapshot

clean:
	rm -r dist/ euvd || true
