# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate go sdk from ir proto
.PHONY: ir-proto
ir-proto:
	@protoc --go_out=. ./api/proto/ir.proto

# build language server for neva written in go and put to vscode extension out
.PHONY: lsp
lsp:
	@go build -ldflags="-s -w" -o web/out/lsp ./cmd/lsp/main.go

# generate typescript types from golang src package to use in vscode extension
# https://github.com/gzuidhof/tygo
.PHONY: tygo
tygo:
	@tygo generate

