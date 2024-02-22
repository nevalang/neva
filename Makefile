# build neva cli binary and insstall to PATH
# example usage: sudo make install
.PHONY: install
install:
	@go build -ldflags="-s -w" cmd/neva/*.go && \
	rm -rf /usr/local/bin/neva && \
	mv cli /usr/local/bin/neva

# generate go parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# build language-server
.PHONY: lsp
lsp:
	@go build -ldflags="-s -w" ./cmd/lsp/*

# generate ts types from go src pkg
.PHONY: tygo
tygo:
	@tygo generate

# lint go code
.PHONY: lint
lint:
	@golangci-lint run ./... --new-from-rev=HEAD~1
