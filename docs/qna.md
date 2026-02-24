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

## Why outport usage is optional while inport usage is required?

Inports are requirements to trigger computation. Outports are optional results. It's possible to implicitly discard unwanted data, but it's not possible to implicitly provide default values for ports not knowing the specifics of the usecase.

## Why are there no int32, float32, etc.?

Neva opts for simplicity with only int and float types, reducing type-conversion needs.

## Why have integers and floats instead of just numbers?

Separate int and float types provide better handling of large numbers, improved integer operation performance, more predictable comparisons, and enhanced type safety.

## Why doesn't Neva have type-casting as a language feature?

Neva intentionally keeps conversions as explicit components instead of syntax
(`-> int ->`, `type(...)`, etc.). There are three reasons:

1. Preserve the 1:1 graph model. In Neva, every computation should be visible as
   a node and edge. Cast syntax would hide conversion nodes behind parser sugar.
2. Keep language core small. Conversion behavior (rounding, parse policy, error
   handling) is a library concern and can evolve in stdlib without growing the
   compiler surface.
3. Keep failures explicit. Many conversions are partial (for example string to
   number parsing). Component APIs can expose `err` outports and integrate with
   `?` propagation naturally.
4. Keep conversions composable. As normal components, converters can be passed
   as dependencies into HOCs/DI flows instead of being hardcoded syntax.

In short: conversions exist, but they are modeled as normal components in
stdlib, not as special language-level casts.

## What is the Go-like split for scalar conversions?

Use a simple split:

1. `builtin`: only total scalar casts that cannot fail at runtime.
2. `strconv`: text parsing/formatting (`string` <-> scalar), where input can be invalid.

In practice, this means:

- `Int(float) -> int`, `Float(int) -> float`, and `String(int) -> string` (code-point cast)
  belong to `builtin` and have no `err` outport.
- `string -> int/float/bool` belongs to `strconv` and should return `err` on invalid input.
- Human-readable scalar-to-string formatting should also live in `strconv` (Go style),
  for example `strconv.Itoa`/`FormatFloat`/`FormatBool`.

## What determines which entities are in the builtin package?

Builtin is Neva's implicit prelude. Every file can reference builtin entities
without imports, so this surface should stay small and stable.

In practice, builtin is for:

1. Primitive language-level types (`int`, `float`, `list<T>`, etc.)
2. Compiler-coupled contracts (for example directives/special analyzer behavior)
3. Very common low-level building blocks that the language model relies on

If Neva later splits builtin into `core` and `prelude`, the same principle
still applies: compiler contracts stay close to core, policy-heavy APIs stay
outside.

## Why are `Union` and `Struct` in builtin?

Because they are part of compiler contracts, not just convenience utilities.

- `Union` has analyzer-aware logic for tag/data compatibility and union-member
  checks. It is not treated as a regular arbitrary helper component.
- `Struct` is the canonical `#autoports` builder. Analyzer/desugarer flows
  assume this pattern for deriving inports from struct type arguments.

So keeping them in builtin makes their special role explicit and keeps them
always available without import noise.

Note: coupling strength differs. `Union` has stronger explicit analyzer coupling;
`Struct` is mostly coupled through `#autoports` conventions and desugaring flow.
If this changes in compiler architecture, this answer should be updated.

## What is builtin `Type` (`type Type any`) and why does it exist?

`Type` is a semantic marker alias over `any` used by compiler-aware builtin
components such as `Union` and `Switch`:

- `Union<T Type>(data Type, tag T) (res T)`
- `Switch<T>(data T, [case] T) ([case] Type, else T)`

It marks ports that may carry different concrete payload types depending on
tag/case analysis. Runtime representation is still `any`, but the alias keeps
signatures readable and communicates "this is intentionally heterogeneous".

In short, `Type` exists because current type-system expressiveness is not enough
to model these ports with fully precise static types while preserving today's
ergonomic component APIs.

## How is `strconv` different from `fmt`?

`strconv` and `fmt` solve different problems:

- `strconv`: pure value conversion/parsing contracts (`string` <-> numbers, etc.)
- `fmt`: presentation and I/O-oriented formatting

Even if both are deterministic, their compatibility promises differ. Conversion
APIs are expected to be canonical and stable for machine-to-machine flows.
Formatting APIs are user-facing and may prioritize readability or template
flexibility.

## Why is `string(42)` in Go not `"42"`?

In Go, integer-to-string conversion in the language is a Unicode code point
conversion, not decimal formatting:

- `string(42)` is `"*"` because 42 is `U+002A`.
- `string(1)` is a one-byte control character (`U+0001`), often shown as `"\x01"`
  in escaped form.

So this is different from "number to decimal text". Decimal formatting is done
via `strconv`/`fmt` APIs.

## Why not allow every `bool <-> number` conversion by default?

Because these conversions are policy-heavy and easy to misuse when global:

- `bool -> int`: should `true` always be `1` and `false` `0`? usually yes.
- `int -> bool`: should non-zero mean `true` or only `1` mean `true`?
- `float -> bool`: what about `0.0`, `-0.0`, `NaN`, `Inf`?

Neva prefers explicitness for these cases. If such conversions are added, they
should have narrowly named components with documented semantics and error rules,
instead of one broad "magic cast" behavior.

## Why `New` is implemented like an infinite loop?

Const nodes are implemented as infinite loops that constantly send messages. This covers all use cases but requires locks for controlled constant message sending. Alternative designs like "trigger" semantics or changing to single-send behavior have been considered, but they introduce other complexities or limitations.

<!-- ## Why have `stream` builtin?

TODO upd

Streams solve the problem of iterating over collections in dataflow, providing a way to know when a collection ends, crucial for implementing patterns like `map/filter/reduce`. -->

## Why Neva's streams are different from classical FBP?

Neva supports infinite nesting for streams, but nested streams aren't used to represent structured data due to the existence of `struct`, `list`, and `dict`.

## How to work with components that expect `T` when you have `stream<T>`?

Use `Map/Filter/Reduce` for data transformations and `ForEach` for side-effects. For complex cases, access `.data` on stream item directly.

## Why isn't `Main:stop` of `int` type?

Using `any` allows for more flexible exit conditions, preventing unintended behavior with union types or any outports.

## Why use structural subtyping?

It reduces code, especially for mappings between records, vectors, and dictionaries. Nominal subtyping doesn't prevent mistakes in type-casts anyway. [Read more](./book/about.md#structural-subtyping).

## Why have `any`?

Neva's `any` is similar to Go's `any` or TypeScript's `unknown`. It's necessary for certain critical cases where the alternative would be an overly complicated type system.

## Why can only primitive messages be used as "literal network senders"?

It enables easier type inference and keeps networks readable.

## Why is there no syntax for `switch`, `for`, `if`, binary and ternary expressions etc.?

Neva is a dataflow language first. Control-flow syntax assumes a tree of expressions that "pulls" data, but in dataflow every computation is a node and every edge is explicit. When we added composite syntax like `switch`, `for`, and expression-style senders, we had to smuggle hidden nodes and edges into the graph. That broke the 1:1 mapping between the visual editor and the text, created context-dependent rules (what is allowed in sender position vs inside a composite), and forced the compiler to know stdlib internals. We moved in a more conventional direction before, and it caused real problems across the parser, analyzer, desugarer, and visualization. Turning back and removing that sugar was hard, but it restored a single, consistent model.

Neva is intentionally hybrid. That means we optimize for a small, pure dataflow core rather than for compact control-flow syntax. The language is also static and code is not data, so we use higher-order components (HOCs) for many patterns that other languages express with syntax features like `for {}`.

In short, we pay a price for being cool.

## Why is there no syntax sugar for `list<T>` or `maybe<T>`?

For consistency with other type syntax and to avoid confusion with different syntaxes in other languages.

## Why is there inconsistent naming in stdlib?

Historically, stdlib mixed multiple naming styles while the language and APIs were evolving. The current direction is to standardize on clear port roles: `res` for primary output, `err` for failures, and `data` only as a generic input placeholder.

## What's the reasoning behind Neva's naming conventions?

Naming follows three heuristics:

1. Semantic clarity in local context (`url`, `filename`, `left`, `right`).
2. Stable control/dataflow conventions (`sig`, `res`, `err`, `then`, `else`).
3. Familiarity with common ecosystems (for example Go stdlib) when it does not reduce clarity.

When heuristics conflict, semantic clarity wins.

## Why do `struct` and `dict` literals require `:` and `,` while struct declarations don't?

This makes literals similar to JSON and consistent with languages like Python, JavaScript, and Go. It also distinguishes between `type` and `const` expressions.

## Why keyword `def` has been chosen to define a component?

Previously used `flow` keyword had drawbacks: it's not used in other languages and may confuse since the abstraction is called "component" not "flow".

Common keywords like `fn`/`fun`/`func` could be misleading since Neva components are not functions but rather coroutines. Keywords like `coro`/`routine` focus too much on implementation details rather than semantics.

While `class` is familiar and components are similar to classes as blueprints, it risks confusion with OOP which Neva does not follow.

`def` was chosen because it's familiar from Python, short, and generic enough to define components without implying specific semantics.

## Why operators and reducers have `left` and `right` naming for ports?

Operators should follow same pattern for simplicity of desugarer and usage by user and they also should be able to be used as reducers by `Reduce`. It means we need to choose between `left/right` which is convenient for operators and `acc/el` for reduce. To keep syntax minimal, we align on `left/right`.

## Why does Neva have overloading, and why can’t users define it?

Neva uses overloading internally to keep the standard library conceptually small without forcing users to deal with tagged unions, explicit type arguments, or name explosions for every concrete type. This improves day-to-day DX in an operator-less, component-only language. However, exposing overloading to users would introduce mental-model complexity: it makes APIs harder to reason about as systems grow. By keeping overloading limited to the `builtin` package, Neva captures the ergonomic benefits where they matter most while preserving a simple, predictable model for user code.
