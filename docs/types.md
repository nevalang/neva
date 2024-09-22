# Types

Nevalang is statically typed language which means types of constants or ports in interfaces must be known at compile time. That's why when you define constants or interface/component signature you write type expressions, that refer to type definitions.

## Expression

There are 2 types of type expression: instantiation and literal.

### Instantiaton

Consist of 2 parts: entity reference and (optional) type-arguments. Type-arguments can only be provided if type-definition we referring to contains type-parameters.

```neva
int // instantiation without type-arguments
list<int> // instantiation with 1 type-argument
streams.ZipResult<int, float> // instantiation with 2 type-arguments
```

Type expressions can be infinitely nested. Example:

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

Literal is form of expression that is used for structs, enums and unions. They are cannot be expressed in form of instantiation.

```neva
struct { a int, b float } // struct with 2 fields
enum { Foo, Bar, Baz } // enum with 3 members
int | string | float | struct{} // union with 4 elements
```

## Definition

Type Definition consist of identifier, body-expression (except base types) and (optional) type-parameter list.

Examples:

```neva
// identifier `foo` with expression `int` and no type-parameters
type foo int
// identifier `bar` with expression `list<T>` and 1 type-parameter `T`
type bar<T> list<T>
// identifier `baz` with expression `dict<T, Y>`
// and 2 type-parameter `T` and `Y`, Y has `int | float` constraint
type baz<T, Y int | float> dict<T, Y>
```

Having these type-definitions in the scope we can use `foo` as `int`, `bar` as `list` and `baz` as `dict` respectfully.

### Parameters (Generics)

Every type paremeter has name that must be unique across all other type parameters in this definition and constrant.

#### Parameter Constraint

Constraint is a type expression that is used as _supertype_ to ensure _type compatibility_ between _type argument_ and type corresponding parameter. If no constrained explicitly defined then `any` is implicitly used.

### Parameters and Arguments Compatibility

Argument `A` compatible with parameter `P` if there's subtyping relation of form `A <: C` where`C` is a constraint of `P`. If there's several parameters/arguments, every one of them must be compatible. Count of arguments must always be equal to the count of parameters.

### Recursive Definition

If type refers to itself inside its own definition, then it's recursive definition. Example: `type l list<l>`. In this case `list` must be base type that supports recursive definitions. Compiler knows which types supports recursion and which do not.

## Expression Resolving

> Please note that this section describes conceptual algorithm rather than actual implementation. In reality type resolution is more complex problem. If you interested in details better refer to source code in typesystem package.

In order to ensure type-safety and generate IR compiler must first resolve each type-expression it encounters. Resolving is reduction of the type-expression to a form where each entity reference refers to a base type.

Let's say we have these type definitions in our scope:

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

Resolving could be thought of as a lambda-reduction-like step-by-step process where type-expression is like a function-call. Entity-reference is name of the function to call and type-arguments are arguments to pass to that function. Then references to type-parameters inside function's body replaced with provided arguments. Then this process recursively starts again and goes on until whole expression is resolved:

**Step 1 - resolve type-arguments**

`maybe` and `float` are base-types so they are already resolved, but `foo` isn't, so `bax<maybe<float>, foo>` becomes `bax<maybe<float>, int>`:

```neva
bax<maybe<float>, foo>
// =>
bax<maybe<float>, int>
```

**Step 2 - replace type-ref with its def body and type-params with provided args**

We need to find type-definition of type we're referring to and take its body, in our case definition of `bax` type looks like this:

```neva
type bax<T, Y int> struct { x T, y baz<Y> }
```

Therefore, we need to replace `bax` in our type-expression with `struct { x T, y baz<T> }` and then replace `T` and `Y` type parameters with provided type arguments `maybe<float>` and `int` resplectfully:

```neva
bax<maybe<float>, int>
// =>
struct { x maybe<float>, y baz<int> }
```

Note that `Y` type-parameter of the `bax` type has type-constraint `int` which means type provided argument must be compatible with it. However, since we passed `int` as a type-argument, everything is ok because `int` is, of course` compatible with itself.

**Step 3 - recursively apply algorithm to def's body**

In our case def's body is body of `bax` which is `struct { x maybe<float>, y baz<int> }`. As we can see it contains only one unresolved reference at the moment - `baz` with `dict<int, T>`.

```
struct { x maybe<float>, y baz<int> }
// =>
struct { x maybe<float>, y dict<int, int> }
```

**All Together**

```neva
bax<maybe<float>, foo>
// =>
bax<maybe<float>, int>
// =>
struct { x maybe<float>, y baz<int> }
// =>
struct { x maybe<float>, y dict<int, int> }
```

## Base Types

Type definition without body called base. There's limited set of base types located in `std/builtin` and compiler is aware of them. If user'll try to define type without body, compiler with throw an error. Some of base types can be used inside recursive type definitions. Here's list of base types:

```neva
pub type any
pub type bool
pub type int
pub type float
pub type string
pub type dict<T>
pub type list<T>
pub type maybe<T>
```

### `any`

Any is a [top-type](https://en.wikipedia.org/wiki/Top_type) - any other type is a sub-type (compatable with) of any. You can pass anything, if `any` is expected. However, since `any` doesn't tell anything about the shape of the data, you cannot pass message of type `any` anywhere where more concrete type is expected. You need to either rewrite your code without using any or explicitly cast `any` to concrete type, while checking the error.

Example 1: `Foo:foo` expects `any` and we send `int`. It's ok because `int` is compatable with `any`:

```neva
flow Main(start any) (stop any) {
    Foo
    ---
    :start -> (42 -> foo -> :stop)
}

flow Foo(foo any) (bar any) {
    // ...
}
```

Example 2: `Main:start` of type `any` is sent to `Foo:foo` which lead to compiler error because `any` is not compatible with `int`:

```neva
flow Main(start any) (stop any) {
    Foo
    ---
    :start -> foo -> :stop
}
```

### `maybe<T>`

Maybe is an [option-type](https://en.wikipedia.org/wiki/Option_type) - it represents message that might not have meaningful value. It's an alternative to have `nil` value in a language. If `T` is expected and you have `Maybe<T>` you have to explicitly check if it has meaningful message by using special unwrap component, before using the value. Ability to work with potentially absent values this way also known as [null/void-safety](https://en.wikipedia.org/wiki/Void_safety), it allows to avoid [billion dollar mistake](https://en.wikipedia.org/wiki/Null_pointer). Maybe expression is compatible with another if type-argument are compatible.

Example 1: `int` is sent to `Foo:foo` which is ok because `int` is compatible with `int|float`:

```neva
flow Main(start any) (stop any) {
    Foo
    ---
    :start -> (42 -> foo -> :stop)
}

flow Foo(foo int|float) (bar any) {
    // ...
}
```

Example 2: `int|string` is sent to `Foo:foo` which lead to compiler error because `int|float` is not compatible with `int`:

```neva
const intOrStr int | string = 42

flow Main(start any) (stop any) {
    Foo
    ---
    :start -> ($intOrStr -> foo -> :stop)
}

flow Foo(foo int) (bar any) {
    // ...
}
```

### `bool`

Boolean type has only two possible values `true` and `false`. It could be though of as special case of enum with 2 members. Booleans are mostly used for `if/then/else` logic and routing. Boolean is only compatible with type that resolves to boolean.

### `int`

Integer is 64-bit signed integer number. Integer is only compatible with type that resolves to integer.

### `float`

Float is 64-bit floating point number. Float is only compatible with type that resolves to float.

### `string`

Strings are UTF-8 encoded byte arrays. It's possible to access string character by index, but you'll have to handle possible absortion of value. Strings also could be terned into streams to iterate on. String is only compatible with type that resolves to string.

### `list`

List is a dynamic array that grows as needed. All messages in the list must have the same type. Static arrays (of fixed size) are not supported. Just like strings, it's possible to access list element by index and turn list into stream to iterate over it. You can access list element by index by O(1) time, but you'll have to handle possible absortion of value. Lists are compatible if their type-argument (types of their elements) are compatible.

### `dict`

Dictionary is an [associative array](https://en.wikipedia.org/wiki/Associative_array) - unordered set of key-value pairs with dynamic keys. All values must have the same type. Just like strings and lists, dictionaries could be converted to streams (each stream item would be structure with key and value in that case) to be iterated over. You can access dictionary key by O(1) time, but you'll have to handle possible absortion of value. Dictionary is compatible with another dictionary if their type-arguments (types of keys and values) are compatible.

### `enum`

Enums are set fixed set of members each with its own unique name. They are represented in memory like integer numbers. When you handle enum member you have to either check all possible cases or handle a default case. Enum is compatible with other enum if its a sub-set of it. Example: `enum {A, B, C}` is compatible with `enum {A, B, C, D}` because everywhere 4 possible values are expected, 3 of those values are also expected.

### `union`

Union is a [sum type](https://en.wikipedia.org/wiki/Tagged_union). It defines set of possible types that message can have. Union is compatible with other union if it a sub-set of it. Example: `int|float` is compatible with `int|float|string` because if message that is either integer, float or string is expected, then a message that is only integer or string is also expected.

### `struct`

Structures are [product types](https://en.wikipedia.org/wiki/Product_type) - compile-time known set of fields with possibly different types.
