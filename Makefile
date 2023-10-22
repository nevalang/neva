# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate go sdk from ir proto
.PHONY: ir-proto
ir-proto:
	@protoc --go_out=. ./api/proto/ir.proto

# generate go sdk from src proto
.PHONY: src-proto
src-proto:
	@protoc --go_out=. ./api/proto/src.proto

# generate ts sdk from src proto
.PHONY: src-proto-ts
src-proto-ts:
	@mkdir -p web/src/generated
	@protoc \
		--plugin=./node_modules/.bin/protoc-gen-ts_proto \
		--ts_proto_opt=esModuleInterop=true \
		--ts_proto_out=./web/src \
		./api/proto/src.proto

# build language server for neva written in go and put to vscode extension out
.PHONY: lsp
lsp:
	@go build -ldflags="-s -w" -o web/out/lsp ./cmd/lsp/main.go

# # generate go gql sdk
# .PHONY: gqlgo
# gqlgo:
# 	@go run -mod=mod github.com/99designs/gqlgen --config ./api/graphql/gqlgen.yml

# # generate ts gql sdk
# .PHONY: gqlts
# gqlts:
# 	@cd web && npm run gqlgen