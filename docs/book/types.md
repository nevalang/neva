# Types

Nevalang is statically typed, requiring compile-time known types for constants and ports. Type expressions refer to type definitions. There is no [type inference](https://en.m.wikipedia.org/wiki/Type_inference) so all types must be explicitly present in the source code.

## Expression

Two types of type expressions: instantiation and literal.

### Instantiation

Consists of an entity reference and optional [type-arguments](https://en.wikipedia.org/wiki/Parametric_polymorphism) (if the type definition has type-parameters).

```neva
int // instantiation without type-arguments
list<int> // instantiation with 1 type-argument
streams.ZipResult<int, float> // instantiation with 2 type-arguments
```

Type expressions can be infinitely nested:

```neva
dict<
    string,
    list<struct{
        foo list<dict<int, float>>
        bar dict<float, list<int | string>>
    }>
>
```

### Literal

Literal expressions are used for structs and unions, which cannot be expressed as instantiations.

```neva
struct { a int, b float } // struct with 2 fields
union { Foo, Bar, Baz } // tagged union with 3 variants
int | string | float | struct{} // union with 4 elements
```

## Definition

Type Definition starts with `type` keyword and consists of an id, optional type parameters, and a body expression (except for base types).

Examples:

```neva
// id `foo` with expr `int` and no type-params
type foo int
// id `bar` with expr `list<T>` and 1 type-param `T`
type bar<T> list<T>
// id `baz` with expr `dict<T, Y>`
// and 2 type-param `T` and `Y`, Y has `int | float` constraint
type baz<T, Y int | float> dict<T, Y>
```

With these type definitions, `foo`, `bar`, and `baz` can be used as `int`, `list`, and `dict` respectively.

### Parameters (Generics)

Every type parameter has a unique name and an optional constraint.

#### Parameter Constraint

Constraint is a type expression used as a supertype to ensure compatibility between type argument and parameter. If not explicitly defined, `any` is implicitly used.

### Parameters and Arguments Compatibility

> Word "compatible" has the same meaning as "is subtype of". Example: "T1 compatible with T2" means "T1 is a sub-type of T2"

Argument `A` is [compatible](https://en.wikipedia.org/wiki/Subtyping) with parameter `P` if `A' <: C`, where `A'` is resolved form of `A` and `C` is resolved `P`'s constraint. For multiple parameters/arguments, each must be compatible. The number of arguments must match the number of parameters. This is a special case of type expressions compatibility.

### Recursive Definition

A type is recursive if it refers to itself in its own definition. For example: `type l list<l>`. The compiler determines which type support recursion.

## Expression Resolving

> Note: This section describes a simplified algorithm. For actual implementation, refer to the typesystem package source code.

Type resolution is the process of reducing type expressions to their base types, ensuring type-safety and enabling IR generation. Consider the following type definitions:

```neva
type foo int
type bar list<int>
type baz<T> dict<int, T>
type bax<T, Y foo> struct { x T, y baz<Y> }
```

And this type-expression:

```neva
bax<maybe<float>, foo>
```

Resolving can be thought of as a [lambda-reduction-like process](https://en.wikipedia.org/wiki/Lambda_calculus#Reduction) where type-expression is like a function-call. Entity-reference is the function name, and type-arguments are the arguments. References to type-parameters inside the function's body are replaced with provided arguments. This process recursively continues until the whole expression is resolved:

**Step 1 - Resolve type arguments**

`maybe` and `float` are base types, so they are resolved. `foo` isn't, so `bax<maybe<float>, foo>` becomes `bax<maybe<float>, int>`:

```neva
bax<maybe<float>, foo>
// =>
bax<maybe<float>, int>
```

**Step 2 - Replace type reference and parameters**

Find the type definition of `bax`:

```neva
type bax<T, Y int> struct { x T, y baz<Y> }
```

Replace `bax` with `struct { x T, y baz<Y> }` and substitute `T` with `maybe<float>` and `Y` with `int`:

```neva
bax<maybe<float>, int>
// =>
struct { x maybe<float>, y baz<int> }
```

**Step 3 - Apply the algorithm recursively**

Now apply the algorithm to the body:

```
struct { x maybe<float>, y baz<int> }
// =>
struct { x maybe<float>, y dict<int, int> }
```

**Final Result**

```neva
bax<maybe<float>, foo>
// =>
bax<maybe<float>, int>
// =>
struct { x maybe<float>, y baz<int> }
// =>
struct { x maybe<float>, y dict<int, int> }
```

### Type Compatibility

Expression `E1` is compatible with `E2`, if `E1` can be used wherever `E2` is expected. For instantiation-expressions:

1. Same entity reference (after both resolved)
2. Equal number of type-arguments
3. Each (resolved) argument of `E1` is compatible with the corresponding (resolved) argument from `E2`

Exception: `any` is a super-type, every type-expression is compatible with it.

For literals rules are different:

#### Struct Literals

Struct literal `S1` is compatible with `S2` if `S1` is a superset of `S2` and each field type in `S1` is compatible with its counterpart in `S2`.

#### Tagged Union Literals

Tagged union literal `U1` is compatible with `U2` if `U1` is a subset of `U2` (all variants in `U1` exist in `U2`).

#### Union Literals

Union literal `U1` is compatible with `U2` if:

1. `U1` has fewer or equal elements than `U2`
2. Each element of `U1` has a compatible element in `U2`

## Base Types

Base types are type definitions without bodies, located in `std/builtin`. The compiler recognizes these types and prevents users from defining bodyless types. Some base types can be used in recursive type definitions. Here's the list:

```neva
pub type any
pub type bool
pub type int
pub type float
pub type string
pub type bytes
pub type dict<T>
pub type list<T>
pub type maybe<T>
```

### `any`

Any is a [top-type](https://en.wikipedia.org/wiki/Top_type). All types are subtypes of `any`. You can pass anything where `any` is expected, but not vice versa. To use `any` where a specific type is needed, you must explicitly cast it, checking for errors.

### `maybe<T>`

Maybe is an [option-type](https://en.wikipedia.org/wiki/Option_type) representing a potentially absent value. It's an alternative to `nil`. Using `Maybe<T>` requires explicit unwrapping before use, ensuring [null-safety](https://en.wikipedia.org/wiki/Void_safety) and avoiding the [billion dollar mistake](https://en.wikipedia.org/wiki/Null_pointer).

### `bool`

Boolean type has two possible values: `true` and `false`. It's similar to a tagged union with 2 variants. Booleans are used for conditional logic and routing.

### `int`

Integer is a 64-bit signed number.

### `float`

Float is 64-bit floating point number.

### `string`

String is UTF-8 text.

For raw binary payloads use `bytes`.

### `bytes`

Bytes is an immutable binary blob type for payloads such as file contents, HTTP bodies, and encoded images.

### `list<T>`

List is a dynamic array of elements with the same type. It can be accessed by index (O(1) time, handling possible absence) or converted to a stream for iteration.

### `dict<T>`

Dictionary is an [associative array](https://en.wikipedia.org/wiki/Associative_array) of key-value pairs. All values have the same type, keys are always strings. Dictionaries can be converted to streams for iteration. Key access is O(1), but require handling absent values.

### Tagged Unions

Tagged unions are fixed sets of named variants, each potentially carrying data. They enable pattern matching and exhaustive case handling. Tagged unions are compatible if one is a subset of another.

### `union`

Union is a [sum type](https://en.wikipedia.org/wiki/Tagged_union) defining possible message types.

### `struct`

Structures are [product types](https://en.wikipedia.org/wiki/Product_type) - compile-time known set of fields with possibly different types.

## Non-Base Builtin Types

There are 2 more types in `std/builtin` that are expressed in terms of language itself - they have bodies and therefore they are not base. Yet they are embedded into builtin package for simplicity, because they are used heavily in the language. These types are `error` and `stream<T>`.

### `error`

Error type for components that can send errors. Similar to [Go's error interface](https://cs.opensource.google/go/go/+/refs/tags/go1.23.1:src/builtin/builtin.go;l=308) but as a structure, since interfaces in Nevalang are for components, not messages.

```neva
pub type error struct {
    text string
    child maybe<error>
}
```

Component/interface that sends error should name outport as `:err`.

### `stream<T>`

Stream structure represents an element of type `T` in a sequence. Streams handle sequences of data in a flow-based programming fashion, allowing operations like mapping, filtering, reducing, and more.

```neva
pub type stream<T> struct {
    data T // current element of the stream
    idx int // index of the current element
    last bool // whether this is the last element in the stream
}
```

Streams can be infinitely nested:

```neva
stream<stream<int>> // 2 levels
stream<stream<stream<int>>> 3 levels
// etc.
```

Stream processing is first-class citizen in Nevalang, so there's dedicated page about that. From data-type perspective, streams are just structures.
