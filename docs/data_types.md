# Data Types

Nevalang has strong static structural type-system

- **Strong** means that there's no implicit type convertions
- **Static** means that semantic analysis for type-safety happens at compile time, not run-time
- **Structural** means that types are compatible if their structure is compatible. Unlike nominative sub-typing, names of the types does not matter. This also means that type can carry more information than needed and still be compatible.

## Base Types

These are types that doesn't have definition in neva source code. Instead compiler is aware of them and knows how to handle them:

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
- _structure_

## Any (Top-Type)

`any` is a _top-type_. It means that any other type is a _sub-type_ of any. That means you can pass any type everywhere `any` is expected. However, since `any` doesn't tell anything about the type, you cannot pass message of type `any` anywhere where more concrete type is expected. You need to either rewrite your code without using any or explicitly cast `any` to concrete type.

## Maybe

Maybe represents value that maybe do not exist. One must unwrap maybe before using the actual value.

## Boolean

Boolean type has only two possible values `true` and `false`

## Integer

Integer is 64 bit integer number

## Float

Integer is 64 bit floating point number

## Strings

Strings are immutable utf encoded byte arrays

## Maps

Maps are unordered key-value pairs with dynamic set of keys. All values must have the same type

## List

List is a dynamic array that grows as needed. All values in the list must have the same type.

## Enums

Enums are set fixed set of values (members) each with its own name. They are represented in memory like integer numbers.

## Union

Union is a _sum type_. It defines set of possible types.

## Struct

Structures are product types (records) - compile-time known set of fields with possibly different types.

## Custom Types

User is allowed to create custom types based on base-types.

---

Further section needs some work. This is WIP document.
