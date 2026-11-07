#!/usr/bin/make

tools/depbo:
	@mkdir -p ./tools
	curl -L -o ./tools/depbo.tgz \
		"https://mikero.bytex.digital/api/download?filename=depbo-tools-0.10.18-linux-amd64.tgz"
	tar xzfv ./tools/depbo.tgz -C ./tools
	mv -v ./tools/depbo-* ./tools/depbo
	rm -v ./tools/depbo.tgz

## Build the server-ready mod into ./dist/mod/
dist/mod: tools/depbo
	rm -rf ./.build-mod ./dist/mod
	mkdir -p ./.build-mod ./dist/mod/@metricz/addons
	cp config.cpp ./.build-mod/
	cp -r scripts ./.build-mod/
	LD_LIBRARY_PATH=./tools/depbo/lib ./tools/depbo/bin/makepbo -P -D -@=metricz ./.build-mod ./dist/mod/@metricz/addons/metricz.pbo
	rm -rf ./.build-mod
