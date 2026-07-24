# Runtime Functions

Runtime functions implement standard-library components declared with
[`#extern`](../user/book/directives.md#extern). They are an implementation
boundary, not the default way to add behavior: prefer a Neva graph when the
behavior can be expressed clearly with existing public components.

`internal/runtime/` imports only the Go standard library. Runtime functions
under `internal/runtime/funcs/` may also import `internal/runtime`.

## Before Adding a Runtime Function

1. Inspect equivalent components in `std/`, `internal/runtime/funcs/`, and
   their tests.
2. State why a Neva composition is insufficient: required primitive semantics,
   unavailable state, or a measured hot-path requirement are valid reasons.
3. Keep the public component signature in `std/`, its `#extern` name, and the
   entry in `internal/runtime/funcs/registry.go` synchronized.

## Execution Contract

A creator resolves its ports once and returns a function that processes
messages until its context is cancelled or a port operation cannot continue.
`Receive` and `Send` return `false` when the operation stops because the
context is done; the runtime function must then return rather than continue
with a zero message.

Every message derived from received input must pass the received `OrderedMsg`
values to `Send` as causes. This preserves runtime ordering and dataflow
tracing.

## Typed Containers

The public Neva values remain `list<T>` and `dict<T>`, but scalar containers
can retain unboxed Go storage such as `[]int64` or `map[string]string`.
Runtime functions should use the matching `runtime.AsList...` or
`runtime.AsDict...` accessor on scalar-preserving hot paths. Keep the generic
`listToMsgs` and `dictToMsgs` helpers for boundaries that genuinely require one
`runtime.Msg` per element, such as conversion to a stream; they box every
scalar element by design.

## Concurrent Inputs

Inputs that belong to one logical operation must be received concurrently.
Use `receive2`, `receive3`, or `receive4`, or add a narrowly scoped equivalent
when the arity requires it. Receiving independent ports sequentially can block
the graph even when each sender is correct.

Sequential reception is permitted only when it is the component's deliberate
public protocol. Document that protocol near the code and cover it with a test.
Do not use sequential receives merely to make local control flow simpler.

## Stateful Functions

State is local to one runtime-function instance unless the component contract
explicitly requires a shared runtime service. Define the ownership, lifecycle,
and contention boundary before adding locks or shared tables. A lock around a
single shared instance can become a program-wide bottleneck even when the type
permits many instances.

Prefer an explicit public runtime API or a dependency passed through an
existing runtime boundary over package-level mutable singletons.

## Tests and Comments

Add focused unit tests for every runtime behavior changed, including normal,
termination, and meaningful corner cases. Add e2e coverage for the exposed
Neva component when a graph-level contract is affected. Benchmarks measure a
performance claim; they do not replace behavior tests.

New Go functions and types need doc comments. For non-obvious concurrency,
ordering, or state blocks, explain the invariant and why the chosen protocol is
safe.
