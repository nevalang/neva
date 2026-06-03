# AGENTS.md (parser)

Scoped guidance for `internal/compiler/parser`.

## Parser Contract

1. Parser output must be structurally parsed AST data.
2. After parser stage, there should be no remaining syntax-level mini-languages
   that still require parsing to recover intended structure.
3. If current implementation violates this rule for compatibility or incremental
   delivery reasons, document the exception near the code and keep it temporary.
4. Prefer explicit AST fields over deferred string parsing in later compiler stages.
