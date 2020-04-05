GOCMD=go
GOBUILD=$(GOCMD) build

apt-dev:
	@mkdir -p ./build
	@$(GOBUILD) -o build/apt ./cmd/apt

setup-web: web-deps run-web

run-web:
	@cd web; BROWSER=none cnpm start

web-deps:
	@cd web; cnpm i

web-build: web-deps
	@cd web; cnpm build
	@go generate ./web