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

