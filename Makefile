# build neva cli for host OS and put to the PATH
.PHONY: install
install:
	@go build -ldflags="-s -w" cmd/cli/*.go && \
	rm -rf /usr/local/bin/neva && \
	mv cli /usr/local/bin/neva

# generate go parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate ts types from go src pkg
.PHONY: tygo
tygo:
	@tygo generate

# === Release Artifacts ===

# build neva cli for amd64 linux
.PHONY: build-linux-amd64
build-linux-amd64:
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" cmd/cli/main.go

# build neva cli for amd64 mac
.PHONY: build-mac-amd64
build-mac-amd64:
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" cmd/cli -o neva-mac-amd64

# build neva cli for amd64 windows
.PHONY: build-windows-amd64
build-windows-amd64:
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" cmd/cli -o neva-windows-amd64.exe

# build neva cli for arm64 linux
.PHONY: build-linux-arm64
build-linux-arm64:
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" cmd/cli -o neva-linux-arm64
