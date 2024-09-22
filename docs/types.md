# Types

Type entity (type definition) consist of an optional list of _type parameters_ followed by optional _type expression_ that is called _body_.

## Base Type

Type definition without body means _base_ type. Compiler is aware of base types and will throw error if non-base type has no body. Base types are only allowed inside `std/builtin` package. Some base types can be used inside _recursive type definitions_.

- _any_
- _maybe_
- _bool_
- _int_
- _float_
- _string_
- _dict_
- _list_
- _enum_
- _union_
- _struct_

### Any (Top-Type)

`any` is a _top-type_. It means that any other type is a _sub-type_ of any. That means you can pass any type everywhere `any` is expected. However, since `any` doesn't tell anything about the type, you cannot pass message of type `any` anywhere where more concrete type is expected. You need to either rewrite your code without using any or explicitly cast `any` to concrete type.

### Maybe

Maybe represents value that maybe do not exist. One must unwrap maybe before using the actual value.

### Boolean

Boolean type has only two possible values `true` and `false`

### Integer

Integer is 64 bit integer number

### Float

Integer is 64 bit floating point number

### Strings

Strings are immutable utf encoded byte arrays

### Maps

Maps are unordered key-value pairs with dynamic set of keys. All values must have the same type

### List

List is a dynamic array that grows as needed. All values in the list must have the same type.

### Enums

Enums are set fixed set of values (members) each with its own name. They are represented in memory like integer numbers.

### Union

Union is a _sum type_. It defines set of possible types.

### Struct

Structures are product types (records) - compile-time known set of fields with possibly different types.

## Custom Type

User is allowed to create custom types based on base-types.

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

<!-- TODO add examples -->
