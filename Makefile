SHELL=/bin/bash
BUILD_TIME=$(shell date -u +%Y-%m-%dT%T%z)
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build

VERSION ?= v0.1.0

.PHONY: version
version:
	@echo $(VERSION)

# Run all tests
.PHONY: test
test:
	@echo "-> $@"
	@go test -v ./{cmd,internal}/...

# Run govet
.PHONY: vet
vet:
	@echo "-> $@"
	@go vet ./...

apt-dev:
	@mkdir -p ./build
	@$(GOBUILD) -o build/apt $(GO_FLAGS) ./cmd/apt

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none cnpm start

web-deps:
	@cd web; cnpm i

web-build: web-deps
	@cd web; cnpm build
	@go generate ./web

ui-server:
	DASH_DISABLE_OPEN_BROWSER=false DASH_LISTENER_ADDR=localhost:3001 $(GOCMD) run ./cmd/apt/main.go dash $(DASH_FLAGS)

ui-client:
	cd web; API_BASE=http://localhost:3001 npm run start

.PHONY: release
release:
	git tag -a $(VERSION) -m "Release $(VERSION)"
	git push --follow-tags


.PHONY: ci
ci: test vet web-build apt-dev
