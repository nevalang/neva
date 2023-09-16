# generate parser from antlr grammar
.PHONY: antlr
antlr:
	@cd internal/parser && \
	antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated