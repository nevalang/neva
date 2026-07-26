# Runtime Functions

Runtime functions implement standard-library components declared with
[`#extern`](../user/book/directives.md#extern). They are an implementation
boundary, not the default way to add behavior: prefer a Neva graph when the
behavior can be expressed clearly with existing public components.

`internal/runtime/messages/` imports only the Go standard library and owns
immutable language values plus pure value operations. `internal/runtime/` may
also import `internal/runtime/messages` for ports, ordering, tracing, and
program execution. Runtime functions under `internal/runtime/funcs/` may import
both packages: use `messages` for pure value work and `runtime` for transport.
Import `messages` directly. Prefer the change that minimizes total system
complexity, even when it requires more work now; do not add transitional
adapters, aliases, or duplicate paths solely to postpone that work.

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
Runtime functions should use the matching `messages.AsList...` or
`messages.AsDict...` accessor on scalar-preserving hot paths. Use
`messages.ListToMsgs` and `messages.DictToMsgs` only at boundaries that
genuinely require one `messages.Msg` per element, such as conversion to a
stream. They return existing boxed storage unchanged, but typed scalar storage
is deliberately boxed into a newly allocated slice or map.

Use `messages.Equal(left, right)` for message equality. Equality is a pure
runtime operation that compares equivalent typed and untyped container storage;
runtime functions must not reimplement it or depend on representation-specific
`Equal` methods.

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
