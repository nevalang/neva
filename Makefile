# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate go sdk from ir proto
.PHONY: irproto
irproto:
	@protoc --go_out=. ./api/proto/ir.proto

# generate go sdk from src proto
.PHONY: srcproto
srcproto:
	@protoc --go_out=. ./api/proto/src.proto

# generate go gql sdk
.PHONY: gqlgo
gqlgo:
	@go run -mod=mod github.com/99designs/gqlgen --config ./api/graphql/gqlgen.yml

# generate ts gql sdk
.PHONY: gqlts
gqlts:
	@cd web && npm run gqlgen