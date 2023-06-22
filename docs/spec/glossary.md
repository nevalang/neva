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

Two type expressions can be checked for [compatibility](#compatibility) if they are both [resolved](#resolved-expression). TODO

## Resolved expression

Expression is resolved if it doesn't contain references to unresolved types anywhere in it's structure. This is recursive definition where base cases are 1) type expression refers to only [native types](#base-type)
Invalid expression cannot be resolved and thus must be threated as unresolved one.

## Resolving

Type expression considered resolved when:

- It's a valid expression itself
- It doesn't contain references to invalid expressions anywhere in it
- There's enough information in the scope to resolve references
- There's no arguments incompatible with corresponding parameters anywhere in it

## Valid expression

Valid type expression is an expression where invariant is not broken.

## Component

Structure defining computational unit that can be used as a node in network. All components have interfaces.

## Interface

Type information about component's interface. Includes optional parameters with optional constants and required IO interface with at least one inport.

## Native component

Component that isn't implemented in source code. Compiler is aware of that componennt and knows where to get its implementation. Some native component are [builtins](#builtins) (they're called [operators](#operators) (but not because they're native)) and some part of standard library being represented as just interfaces. It's important to note that the second group is something that must be replaced with normal components as much as possible in the future. User must not be able to create its own native component as it makes programs less portable and harder to maintain.

## Operator

Just another name for [builtin](#builtin-component). As follows from definition of a builtin components, operator doesn't have to be [native component](#native-component), even though it usually is.

## Normal component

Component that is implemented in source code. It have its own network and can depend on other components.

## Builtin component

Component that doesn't have implentation in source code and copiler is aware of it. The goal is to have as less such components as possible in order to make standard library (and the language itself) as cross-platform as possible. On the other hand all the necessary functionality must be present in a language so many builtin components are needed for the first time.

## Safety

...

## Sub-stream

Stream with nesting

## Nested sub-stream

Stream with nesting level more than 1

## Base type

Type that doesn't have underlying type. It can have parameters with or without constraints though.

## Recursive type

Not-base type that refers to itself in its body

## Type that can be used for recursive definitions

TODO create term for that

## Builtins

Set of [entities](#entitiy) that are not part of standart library. Compiler makes it possible to refer such entities from any package as if they were local. One may say that "builtins are available in global scope". It's important to note that builtin entity not necessary a [native entity](#native-entity), even though it usually is.

## Native Entity

[Entity](#entity) that isn't implemented in source code. Compiler is aware of such entity. There's 2 types of such entities: [native types](#native-types) and [native components](#native-component).

## Native type

A type that isn't implemented (doesn't have definition) in source code. All native types are [builtin](#builtins) (as opposed to [native components](#native-component) that can also live in standard library). These types are essential for [resolving](#resolving).