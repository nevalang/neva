# build neva cli
.PHONY: install
install:
	@go build -ldflags="-s -w" cmd/neva/*.go && \
	rm -rf /usr/local/bin/neva && \
	mv cli /usr/local/bin/neva

# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# build language server and put executable to web/out
.PHONY: lsp
lsp:
	@go build -ldflags="-s -w" -o web/lsp ./cmd/language-server/*

# generate typescript types from golang src package to use in vscode extension
# https://github.com/gzuidhof/tygo
.PHONY: tygo
tygo:
	@tygo generate

# static Analysis tool to detect potential nil panics in go code
.PHONY: nilaway
nilaway:
	@nilaway ./...
