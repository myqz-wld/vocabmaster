.PHONY: build test run clean pack install uninstall

GO_VERSION := $(shell awk '/^go / {print "go" $$2}' go.mod)
GOENV := if [ -f "$$HOME/.gvm/scripts/gvm" ]; then source "$$HOME/.gvm/scripts/gvm"; if gvm list | grep -q "$(GO_VERSION)"; then gvm use "$(GO_VERSION)" >/dev/null; fi; fi;
BINDIR ?=
UPDATE_SHELL_RC ?= 1
SHELL_RC ?=
export BINDIR
export UPDATE_SHELL_RC
export SHELL_RC

build:
	@bash -lc '$(GOENV) go run scripts/write-build-info.go && go build -o build/vocabmaster ./src'

test:
	@bash -lc '$(GOENV) go test ./...'
	@scripts/install_test.sh

run: build
	./build/vocabmaster

pack: build
	@scripts/package.sh

install: build
	@bash -lc '$(GOENV) scripts/install.sh install'

uninstall:
	@bash -lc '$(GOENV) scripts/install.sh uninstall'

clean:
	rm -rf build/ vocabmaster
