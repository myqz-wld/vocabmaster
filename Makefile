.PHONY: build test run clean install uninstall

GO_VERSION := $(shell awk '/^go / {print "go" $$2}' go.mod)
GOENV := if [ -f "$$HOME/.gvm/scripts/gvm" ]; then source "$$HOME/.gvm/scripts/gvm"; if gvm list | grep -q "$(GO_VERSION)"; then gvm use "$(GO_VERSION)" >/dev/null; fi; fi;
BINDIR ?=
UPDATE_SHELL_RC ?= 1
SHELL_RC ?=
export BINDIR
export UPDATE_SHELL_RC
export SHELL_RC

build:
	@bash -lc '$(GOENV) go build -o build/vocabmaster ./src'

test:
	@bash -lc '$(GOENV) go test ./...'
	@src/tools/install_test.sh

run: build
	./build/vocabmaster

install: build
	@bash -lc '$(GOENV) src/tools/install.sh install'

uninstall:
	@bash -lc '$(GOENV) src/tools/install.sh uninstall'

clean:
	rm -rf build/ vocabmaster
