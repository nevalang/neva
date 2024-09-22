# Directives

Compiler directives are special instructions for the compiler, not intended for daily use but important for understanding language features.

## `#extern`

Tells compiler a component lacks source code implementation and requires a runtime function call.

## `#bind`

Instructs compiler to insert a given message into a runtime function call for nodes with `extern` components.

## `#autoports`

Derives component inports from its type-argument structure fields, rather than defining them in source code.
