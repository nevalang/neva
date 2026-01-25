# Context: deriveNodeConstraintsFromNetwork receiver pairing fix follow-up

## Current focus
We recently fixed incorrect sender→receiver pairing for incoming constraints in `deriveNodeConstraintsFromNetwork` by replacing:

```
recvPortAddrs := a.flattenReceiversPortAddrs(conn.Normal.Receivers)
```

with:

```
recvPairs := a.collectReceiverSenderPairs(conn.Normal.Receivers, conn.Normal.Senders)
```

This ensures chained connections use the correct sender for each receiver (e.g., in `a -> b -> c`, receiver `c` is paired with sender `b`, not `a`). This fix was motivated by union/switch regression where overload resolution saw union/any instead of concrete types.

## Remaining question
`flattenReceiversPortAddrs` is still used in **three places** inside `deriveNodeConstraintsFromNetwork` for **outgoing constraints** (when the node is a sender):

1) Deriving outgoing constraints from direct receivers.
2) Deriving outgoing constraints from chained receivers.
3) Similar logic in nested chained scans.

These usages may still over-constrain a sender by propagating constraints from deep receivers in chains (e.g., `a -> b -> c` makes `a` constrained by `c`), rather than only by the immediate chain head (`b`).

## Task
Review whether outgoing constraint collection should also use sender/receiver pairing like `collectReceiverSenderPairs` to avoid propagating constraints across chain hops.

If yes:
- Replace `flattenReceiversPortAddrs(...)` usage(s) in `deriveNodeConstraintsFromNetwork` with a pairing-aware approach.
- Ensure only the *immediate* receiver (chain head) influences the sender’s outgoing constraints.
- Re-run targeted tests around unions/switch + any relevant analyzer tests.

If no:
- Document why deeper constraints are desired/acceptable.

## Key files
- `internal/compiler/analyzer/nodes.go` (`deriveNodeConstraintsFromNetwork`, `collectReceiverSenderPairs`)
- `internal/compiler/analyzer/switch_logic.go` (context for switch-case typing)
- `internal/compiler/analyzer/union_logic.go` (context for union tag/data typing)

## Tests previously used
- `go test ./internal/compiler/analyzer`
- `go test ./e2e/switch_union_matching_no_data ./e2e/switch_union_matching_with_data ./e2e/union_wrapping_test` (when relevant)

