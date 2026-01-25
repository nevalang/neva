# Fan-in at chained connection head — semantics and motivation

## Context
The analyzer currently forbids multiple senders inside a chained connection. In
`analyzeChainedConnectionReceiver`:

```
if len(chainedConn.Normal.Senders) != 1 {
    return error("multiple senders are only allowed at the start of a connection")
}
```

This means `a -> b -> c` is allowed (single sender at chain head), but
`[a, b] -> x -> y` is rejected.

## Why this matters
The language already allows multiple senders at the *start* of a connection:

```
[a, b] -> receiver
```

But disallows multiple senders at a **chained head**:

```
[a, b] -> x -> y   // currently forbidden
```

We should decide whether this restriction is purely historical/implementation
or if it encodes a semantic rule that keeps chain connections linear.

## Code references
- `internal/compiler/analyzer/receivers.go`
  - `analyzeChainedConnectionReceiver` (single sender enforcement)
- `internal/compiler/desugarer/network.go`
  - `desugarChainedConnection` (how chains are lowered)

## Semantics to clarify
- Is a chain meant to be a **linear transform** of a single stream?
- If we allow `[a, b] -> x -> y`, should `x`:
  - accept either `a` or `b` independently (fan-in), or
  - wait for both (implicit sync), or
  - something else?
- Do we want the chain to be “just sugar” over explicit fan-in nodes (Lock/Join/etc)?

## Examples to evaluate
### Allowed today
```
:a -> x -> y
[a, b] -> receiver
```

### Forbidden today
```
[a, b] -> x -> y
```

### Possible interpretation if allowed
- fan-in at `x` (either `a` or `b` triggers `x`, then `y`)
- or implicit synchronization (wait for both) — but that is not currently encoded

## Decision points
- Keep restriction as a semantic rule (chains are linear)
- Allow fan-in at chain head and define explicit behavior
- If allowed, update analyzer check + add tests

