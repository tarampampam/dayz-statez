#!/usr/bin/make

.DEFAULT_GOAL: all
all: dist/mod dist/statez-server

tools/depbo:
	@mkdir -p ./tools
	curl -L -o ./tools/depbo.tgz \
		"https://mikero.bytex.digital/api/download?filename=depbo-tools-0.10.18-linux-amd64.tgz"
	tar xzfv ./tools/depbo.tgz -C ./tools
	mv -v ./tools/depbo-* ./tools/depbo
	rm -v ./tools/depbo.tgz

## Build the server-ready mod into ./dist/mod/
.PHONY: dist/mod
dist/mod: tools/depbo
	@rm -rf ./mod/.build ./dist/mod
	mkdir -pv ./mod/.build ./dist/mod/@statez/addons
	cp -rv ./mod/config.cpp ./mod/scripts ./mod/.build/
	LD_LIBRARY_PATH=./tools/depbo/lib ./tools/depbo/bin/makepbo \
		-P -D -@=statez ./mod/.build ./dist/mod/@statez/addons/statez.pbo
	@rm -rf ./mod/.build

## Build the server binary into ./dist/statez-server
.PHONY: dist/statez-server
dist/statez-server:
	CGO_ENABLED=0 go build -C ./server -ldflags "-s -w" -trimpath -buildvcs=false -o ./../dist/statez-server ./cmd/statez-server/

