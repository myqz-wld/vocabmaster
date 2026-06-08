.PHONY: build test run clean install uninstall

GO_VERSION := $(shell awk '/^go / {print "go" $$2}' go.mod)
GOENV := if [ -f "$$HOME/.gvm/scripts/gvm" ]; then source "$$HOME/.gvm/scripts/gvm"; if gvm list | grep -q "$(GO_VERSION)"; then gvm use "$(GO_VERSION)" >/dev/null; fi; fi;
BINDIR ?=

build:
	@bash -lc '$(GOENV) go build -o build/vocabmaster .'

test:
	@bash -lc '$(GOENV) go test ./...'

run: build
	./build/vocabmaster

install: build
	@bash -lc '$(GOENV) \
		bindir="$(BINDIR)"; \
		if [ -z "$$bindir" ]; then \
			gopath="$$(go env GOPATH)"; \
			if [ -z "$$gopath" ]; then echo "无法解析 GOPATH；请使用 make install BINDIR=/path/to/bin" >&2; exit 1; fi; \
			gopath="$${gopath%%:*}"; \
			bindir="$$gopath/bin"; \
		fi; \
		mkdir -p "$$bindir"; \
		if [ -e "$$bindir/vm" ] || [ -L "$$bindir/vm" ]; then \
			target="$$(readlink "$$bindir/vm" 2>/dev/null || true)"; \
			if [ "$$target" != "vocabmaster" ] && [ "$$target" != "$$bindir/vocabmaster" ]; then \
				echo "安装失败: $$bindir/vm 已存在，请先移除或使用 BINDIR 指定其他目录" >&2; \
				exit 1; \
			fi; \
		fi; \
		install -m 0755 build/vocabmaster "$$bindir/vocabmaster"; \
		ln -sfn "vocabmaster" "$$bindir/vm"; \
		echo "已安装到 $$bindir/vocabmaster"; \
		echo "已创建别名命令 $$bindir/vm"; \
		case ":$$PATH:" in *":$$bindir:"*) ;; *) echo "提示: $$bindir 不在 PATH 中，请加入 PATH 后使用 vocabmaster 或 vm";; esac'

uninstall:
	@bash -lc '$(GOENV) \
		bindir="$(BINDIR)"; \
		if [ -z "$$bindir" ]; then \
			gopath="$$(go env GOPATH)"; \
			if [ -z "$$gopath" ]; then echo "无法解析 GOPATH；请使用 make uninstall BINDIR=/path/to/bin" >&2; exit 1; fi; \
			gopath="$${gopath%%:*}"; \
			bindir="$$gopath/bin"; \
		fi; \
		rm -f "$$bindir/vocabmaster" "$$bindir/vm"; \
		echo "已卸载 $$bindir/vocabmaster 和 $$bindir/vm"'

clean:
	rm -rf build/ vocabmaster
