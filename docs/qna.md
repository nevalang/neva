# Questions And Answers

This Q&A addresses questions not covered in other documentation pages and explains certain language design decisions.

## Is Neva "classical FBP"?

No, but it incorporates many ideas from Flow-Based Programming (FBP). [Read more](./book/about.md#flow-based-programming).

## Why are array-ports needed?

Array-ports are necessary for combining data from multiple sources. [Read more](./book/interfaces.md#array-ports).

## Why can't components read from their own array-inports by index?

This could lead to deadlock/panic if the parent component doesn't use specific amount slots, which would defeat the purpose of array-ports. [Read more](./book/networks.md#array-ports-constraints).

## Why components _can_ read from sub-node's array-outports by index?

This is crucial for "routing" cases where data needs to be sent to specific handlers based on conditions.

## Why is outport usage optional while inport usage is required?

Inports are requirements to trigger computation. Outports are optional results. It's possible to implicitly discard unwanted data, but it's not possible to implicitly provide default values for ports not knowing the specifics of the usecase.

## Why are there no int32, float32, etc.?

Neva opts for simplicity with only int and float types, reducing type-conversion needs.

## Why have integers and floats instead of just numbers?

Separate int and float types provide better handling of large numbers, improved integer operation performance, more predictable comparisons, and enhanced type safety.

## What determines which entities are in the builtin package?

Builtin package entities are frequently used or used internally by the compiler.

## Why is `Const` implemented like an infinite loop?

Const nodes are implemented as infinite loops that constantly send messages. This covers all use cases but requires locks for controlled constant message sending. Alternative designs like "trigger" semantics or changing to single-send behavior have been considered, but they introduce other complexities or limitations.

## Why have `streams.Item`?

`streams.Item` values solve the problem of iterating over collections in dataflow, providing a way to know when a collection ends—crucial for implementing patterns like `streams.Map`/`streams.Filter`/`streams.Reduce`.

## Why are Neva's streams different from classical FBP?

Neva supports infinite nesting for streams, but nested streams aren't used to represent structured data due to the existence of `struct`, `list`, and `dict`.

## How to work with components that expect `T` when you have `streams.Item<T>`?

Use `streams.Map`/`streams.Filter`/`streams.Reduce` for data transformations and `streams.ForEach` for side-effects. For complex cases, access `.data` on the stream item directly.

## Why isn't `Main:stop` of `int` type?

Using `any` allows for more flexible exit conditions, preventing unintended behavior with union types or any outports.

## Why use structural subtyping?

It reduces code, especially for mappings between records, vectors, and dictionaries. Nominal subtyping doesn't prevent mistakes in type-casts anyway. [Read more](./book/about.md#structural-subtyping).

## Why have `any`?

Neva's `any` is similar to Go's `any` or TypeScript's `unknown`. It's necessary for certain critical cases where the alternative would be an overly complicated type system.

## Why can only primitive messages be used as "literal network senders"?

It enables easier type inference and keeps networks readable.

## Why is there no syntax sugar for `list<T>` or `maybe<T>`?

For consistency with other type syntax and to avoid confusion with different syntaxes in other languages.

## Why is there inconsistent naming in stdlib?

Some basic components follow naming conventions from other languages for familiarity.

## What's the reasoning behind Neva's naming conventions?

Names are chosen to be familiar to most programmers, easing the paradigm shift.

## Why do `struct` and `dict` literals require `:` and `,` while struct declarations don't?

This makes literals similar to JSON and consistent with languages like Python, JavaScript, and Go. It also distinguishes between `type` and `const` expressions.

## Why keyword `def` has been chosen to define a component?

Previously used `flow` keyword had drawbacks: it's not used in other languages and may confuse since the abstraction is called "component" not "flow".

Common keywords like `fn`/`fun`/`func` could be misleading since Neva components are not functions but rather coroutines. Keywords like `coro`/`routine` focus too much on implementation details rather than semantics.

While `class` is familiar and components are similar to classes as blueprints, it risks confusion with OOP which Neva does not follow.

`def` was chosen because it's familiar from Python, short, and generic enough to define components without implying specific semantics.

## Why syntax for ternary operator is like this (without incoming `->`)?

It would be in the spirit of Neva to have incoming `->` that would visually "trigger" ternary:

```neva
condition -> (? thenValue : elseValue) -> ...
```

This way we have clear incoming and outgoing parts of the connection. However, despite that looks weird (not a single language have this ternary syntax), it's just inconsistent - all 3 (condition, then and else) parts of the ternary are incoming messages ("arguments") for ternary component. Consistent syntax could be like this:

```neva
condition -> ternary:cond
thenValue -> ternary:then
elseValue -> ternary:else
ternary:res -> println
```

And guess what? This is exactly how desugared `Ternary` component works! But these are 4 (!) connections! Compare it with `(condition ? thenValue : elseValue) -> println`.

## Why operators and reducers have `left` and `right` naming for ports?

Operators should follow same pattern for simplicity of desugarer and usage by user and they also also should be able to be used as reducers by `Reduce`. It means we need to choose between `left/right` which is convinient for operators and `acc/el` for reduce. Binary expressions (infix form) are more common than reduce operations so desicion was made to sacrifice reduce clarity a little bit.
