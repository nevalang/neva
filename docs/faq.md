# FAQ

This FAQ addresses questions not covered in other documentation pages. It serves as a reference for the reasoning behind certain language design decisions.

## Is Neva "classical FBP"?

No, but it many ideas from Flow-Based Programming (FBP). You can read more on [FBP page](./fbp.md)

## Why are array-ports needed?

Array-ports are necessary for combining data from multiple sources.

## Why can't components read from their own array-inports by index?

Allowing this would lead to potential blocking or crashes if the parent node doesn't use all slots. Enforcing exact slot usage would negate the flexibility of array-ports.

## Why can components read from sub-node's array-outports by index?

This is necessary for critical cases like "routing" where data needs to be sent to specific handlers based on certain conditions. While it could be implemented with if-else chaining, that approach would be tedious and less efficient.

## Why is outport usage optional while inport usage is required?

In general it's impossible to tell if component is able to produce result without all inports being used. Inports are requirements to trigger computation. Outports, on the other hand, are results you can use anyway you want. This design allows for discarding unwanted data.

## Why are there no int32, float32, etc.?

Neva opts for simplicity, providing only int and float types as a compromise between flexibility and ease of use. This way you may have much less type-convertsion in your code.

## Why have integers and floats instead of just numbers?

Separate int and float types provide:

1. Better handling of large numbers and overflow issues
2. Improved performance for integer operations
3. More predictable comparisons
4. Enhanced type safety

## What determines which entities are in the builtin package?

Entities in the builtin package are:

1. Frequently used
2. Used internally by the compiler

## Why Emitter implemented like an infinite loop?

Const nodes are implemented like infinite loops that constantly sends messags across their receivers. This covers all the usecases but also requires locks because we usually want control when we send constant messages.

Alternative to that design would be "trigger" semantics where we have some kind of `sig` inport for const nodes. Instead of starting at the program startup such trigger component would wait for external signal and then do the usual const stuff (infinite loop with sending messages).

**The problem #1 with this approach - we still needs locks**. It's not enough to trigger infinite loop. E.g. in "hello world" example nothing would stop `msg` const node to send more than 1 `hello world` message to `print`.

**Possible solution for that would be to change semantics and remove infinite loop logic**. Make const node send signal only after we trigger it via sig port. The problem with this approach is that there is many usecases where we want infinite loop behavior. Think of initial inports - e.g. `requestSender` component with `data` and `url` inports where `data` is dynamic and `url` is static. It's not enough to send static url value just once (`requestSender` must remember it, we don't go that way because that leads to bad design where components know where they get data from - this is huge violation of transport vs logic separation).

This problem by itself is fixable with using external sources like signals. When we have some static inport we usually have some kind of dynamic data that must be used in combination with it. Even though it would lead to making networks more complicated (locks do this too though), it's possible solution. But we have another problem:

**Problem #2** - `$` syntax sugar.

Another problem with previous solution (const nodes have sig inport and they send one message per one signal) is how use `$` syntax sugar.

Currently it's possible to _refer constants_ in network like this:`$msg -> ...`

This won't be the thing because we have to have not just entity reference but regular ports like `$msg.sig` and `$msg.v`. This is not a disaster but rather strange API.

Where this `$msg` node come from? Is it clear to me as a user that there are implicit nodes $ prefix for every constant that I can refer? Why these `sig` and `v` inports? Because this is how `std/builtin.Const` works? Why do I have to know this? Why do I have to know how syntax sugar is implemented under the hood in compiler?

Finally another possible solution to that could be `Trigger` components in combination with `Const` components. The difference would be that const behaves like infinite loops that requires locks and triggers behaves like single sending triggers (no lock required).

Problems with this solution are:

1. Now we have 2 ways to do the same thing. Do I need to use const? Or trigger? When to use what?
2. Trigger either must be used in combination with `#runtime_func_msg` directive (this violates principle that user must be able to program without directvies) or there must be sugar for triggers.

It's possible in theory to create sugar for triggers but the language could be too complicated with so many syntax features. This conclusion by itself is questionable but in combination with the first problem - having 2 ways to use static values. Looks like it's better not to have triggers.

All this leads to a conclusion that the only good semantic for constants is the current ones - infinite loops that requires locks. The need to have locks is not a fancy thing by itself, but "then connections" sugar made this pretty simple.

## Why have streams builtin?

<!-- TODO add about streaming meaning in language -->

Sub-streams solve the problem of iterating over collections in FBP, where we lack mutable state and code-as-data. They provide a way to know when a list ends, which is crucial for implementing patterns like `map`.

## Why are Neva's streams different from classical FBP?

Neva supports infinite nesting for streams because they are implemented with just structs. However, nested streams are not used to represent structured data because there are structs, lists and dictionaries.

## How to work with components that expect `T` when you have `stream<T>`?

It depends on what you want to do with it. Generally you're ok with `Map/Filter/Reduce` for data transformations and `For` for side-effects. If you have more complex case you can access `.data` on stream item directly.

## Why isn't `Main:stop` has `int` type?

Using `any` allows for more flexible exit conditions. Interpreting any `int` as an exit code could lead to unintended behavior, especially with union types or any outports.

## Why use structural subtyping?

1. Reduces code, especially for mappings between records, vectors, and dictionaries
2. Nominal subtyping doesn't prevent mistakes like passing wrong values to type-casts anyway

## Why have `any`?

Neva's `any` similar to Go's `any` or TypeScript's `unknown`. Is necessary for certain critical cases where the alternative would be an overly complicated type system.

## Why can only primitive messages be used as "literal network senders"?

1. Easier type inference for the compiler
2. Keeps networks readable (complex literals would be hard to read)

## Why is there no syntax sugar for `list<T>` or `maybe<T>`?

Consistency with other type syntax and avoiding confusion with different syntaxes in other languages (should it be `[]T` or `[]T`? Or `[T]`? Etc.)

## Why is there inconsistent naming in stdlib?

Some basic components (e.g., Len, Eq) follow naming conventions from other languages for familiarity.

## What's the reasoning behind Neva's naming conventions?

Names are chosen to be familiar to most programmers, easing the paradigm shift.

## Why `struct` and `dict` literals require `:` and `,` while struct declarations don't?

This makes literals similar to JSON and consistent with languages like Python, JavaScript, and Go. It also distinguishes between `type` and `const` expressions.
