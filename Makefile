.PHONY: install-cfg build install

build:
	go build -o build

run:
	make build
	build/gofetch

install:
	install -Dm755 build/gofetch $$HOME/.local/bin/gofetch

install-cfg:
	mkdir $$HOME/.config/gofetch/
	install -Dm755 cfg.yml $$HOME/.config/gofetch/gofetch.yml
