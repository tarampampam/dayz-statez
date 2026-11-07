#!/usr/bin/make

tools/depbo:
	@mkdir -p ./tools
	curl -L -o ./tools/depbo.tgz \
		"https://mikero.bytex.digital/api/download?filename=depbo-tools-0.10.18-linux-amd64.tgz"
	tar xzfv ./tools/depbo.tgz -C ./tools
	mv -v ./tools/depbo-* ./tools/depbo
	rm -v ./tools/depbo.tgz
