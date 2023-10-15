# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/compiler/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate go sdk from ir proto
.PHONY: irproto
irproto:
	@protoc --go_out=. ./api/proto/ir.proto

# run frontend devserver
.PHONY: web
web:
	@cd web && npm start

# generate golang graphql SDK
.PHONY: gqlgo
gql:
	@go run -mod=mod github.com/99designs/gqlgen --config ./api/graphql/gqlgen.yml

# generate typescript graphql SDK
.PHONY: gqlts
gqlts:
	@cd web && npm i && npm run gqlgen