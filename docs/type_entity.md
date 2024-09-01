# Type Entity

Type entity (type definition) consist of an optional list of _type parameters_ followed by optional _type expression_ that is called _body_.

## Base Type

Type definition without body means _base_ type. Compiler is aware of base types and will throw error if non-base type has no body. Base types are only allowed inside `std/builtin` package. Some base types can be used inside _recursive type definitions_.

## Recursive Type Definition

If type refers to itself inside its own definition, then it's recursive definition. Example: `type l list<l>`. In this case `list` must be base type that supports recursive definitions. Compiler knows which types supports recursion and which do not.

## Type Parameters (Generics)

Every type paremeter has name that must be unique across all other type parameters in this definition and constrant.

## Type Parameter Constraint

Constraint is a type expression that is used as _supertype_ to ensure _type compatibility_ between _type argument_ and type corresponding parameter. If no constrained explicitly defined then `any` is implicitly used.

## Type Parameters and Arguments Compatibility

Argument `A` compatible with parameter `P` if there's subtyping relation of form `A <: C` where`C` is a constraint of `P`. If there's several parameters/arguments, every one of them must be compatible. Count of arguments must always be equal to the count of parameters.

## Type Expression

There is 2 _kinds_ of type expressions:

1. Instantiation expressions
2. Literals expressions

Type expressions can be infinitely nested. Process of reducing the type expression called _type resolving_. Resolved type expression is an expression that cannot be reduced to a more simple form.

## Type Instantiation Expression

Such expression consist of _entity reference_ (that must refer to existing type definition or type parameter) and optional list of _type arguments_. Type arguments themselves are arbitrary type expressions (they follows the same rules described above).

## Literal Type Expression

Type expressions that cannot be described in a instantiation form.

