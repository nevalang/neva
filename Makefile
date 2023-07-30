.PHONY: antlr
antlr:
	@cd internal/compiler/frontend && \
	@antlr4 -Dlanguage=Go -no-visitor -package parsing ./neva.g4 -o generated