
PKG = $(shell go list ./... | grep -v '/vendor')

release:
	@echo "[+] re-generating"
	@go generate ./...
	@echo "[+] building"
	@$(MAKE) build
	@echo "[+] complete"
.PHONY: release

test:
	@go test -cover $(PKG)
.PHONY: test

build:
	@gox -os="linux darwin windows" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean
