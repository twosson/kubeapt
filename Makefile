SHELL=/bin/bash
BUILD_TIME=$(shell date --utc --rfc-3339 ns | sed -e 's/ /T/')
GIT_COMMIT=$(shell git rev-parse --short HEAD)

LD_FLAGS= '-X "main.buildTime=$(BUILD_TIME)" -X main.gitCommit=$(GIT_COMMIT)'
GO_FLAGS= -ldflags=$(LD_FLAGS)
GOCMD=go
GOBUILD=$(GOCMD) build

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