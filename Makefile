# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate go sdk from ir proto
.PHONY: ir-proto
ir-proto:
	@protoc --go_out=. ./api/proto/ir.proto

# build language server and put executable to web/out
.PHONY: lsp
lsp:
	@go build -ldflags="-s -w" -o web/lsp ./cmd/lsp/*

# generate typescript types from golang src package to use in vscode extension
# https://github.com/gzuidhof/tygo
.PHONY: tygo
tygo:
	@tygo generate

