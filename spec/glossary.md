# Glossary

## Mutual compatibility

Types that are [compatible](#compatibility) with eachother. E.g. type `a` is mutually compatible with type `b` if `a` is compatible with `b` and `b` is compatible `a`. So `a <: b && b <: c`

## Compatibility

Type `a` is compatible with the type `b` if it's a [subtype](#subtyping) of `b`: `a <: b`

## Subtyping

Type expr `a` is `<:` (subtype of) type expr `b` if it's [safe](#safety) to use it anywhere where type `b` is expected. There's special [algorithm](#subtype-checking) for checking this.

## Type expression

Type expression ...

## Subtype checking

Two type expressions can be checked for [compatibility](#compatibility) if they are both [resolved](#resolved-expression).

## Resolved expression

Expression is resolved if it doesn't contain references to unresolved types anywhere in it's structure. Invalid expression cannot be resolved and thus must be threated as unresolved one.

## Resolving

Type expression can be resolved if:

- It's a valid expression itself
- It doesn't contain references to invalid expressions anywhere in it
- There's enough information in the scope to resolve references
- There's no arguments incompatible with corresponding parameters anywhere in it

## Valid expression

Valid type expression is an expression where invariant is not broken.

## Component

Everything that can be used as a node in program's network. I.e has component header.

## Component header

Type information about component's interface. Includes optional parameters with optional constants and required IO interface with at least one inport.

## Native component

## Custom component

Component that depends on other components.

## Builtin component

Component that platform must implement to satisfy the specification. Such component must be part of standart library. It could be native or custom component.

## Safety

...

## Sub-stream

Stream with nesting

## Nested sub-stream

Stream with nesting level more than 1

## Base type

Type that doesn't have underlying type. It can have parameters with or without constraints though.
