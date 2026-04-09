.PHONY: build test run clean install uninstall

GOENV := source $$HOME/.gvm/scripts/gvm && gvm use $$(grep '^go ' go.mod | awk '{print "go"$$2}') >/dev/null 2>&1 &&

build:
	@bash -lc '$(GOENV) go build -o vocabmaster .'

test:
	@bash -lc '$(GOENV) go test ./...'

run: build
	./vocabmaster

install: build
	sudo cp vocabmaster /usr/local/bin/vocabmaster
	@echo "已安装到 /usr/local/bin/vocabmaster"

uninstall:
	sudo rm -f /usr/local/bin/vocabmaster
	@echo "已卸载"

clean:
	rm -f vocabmaster
