# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated

# generate go sdk from ir proto
.PHONY: irproto
irproto:
	@protoc --go_out=pkg/ir ./api/ir.proto