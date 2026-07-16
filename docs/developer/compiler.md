# Compiler

The compiler pipeline is:

```text
parser -> analyzer -> desugarer -> IR generator -> Go backend
```

The [architecture page](./architecture.md) shows the surrounding build flow.
The user documentation owns public language semantics; this page records
implementation boundaries that preserve those semantics.

## Parser

The parser produces structured AST data. Do not leave a syntax-level
mini-language encoded in strings for later stages to parse again. When syntax
changes, update the grammar and generated parser artifacts together and run the
parser smoke tests.

## Analyzer and Standard-Library Contracts

Some standard-library directives are compiler contracts rather than ordinary
library calls. For example, struct selectors and directives such as `#extern`
and `#bind` require analysis that cannot be expressed solely in Neva source.
Keep that knowledge explicit and narrowly scoped; do not disguise it as a
general-purpose language exception.

Analyze source before desugaring. This preserves diagnostics in the user's
terms and prevents desugaring from operating on unchecked constructs. The IR
generator should lower an already analyzed program, not recover source syntax.

## Type-System Changes

Before changing grammar, AST, analyzer, or type-system code, create a minimal
reproducer and establish which existing language invariant it violates. Check
whether the requested behavior is already expressible by an existing type,
union, component, or standard-library contract. Do not introduce a special
case for a component when a general language rule is the actual model.
