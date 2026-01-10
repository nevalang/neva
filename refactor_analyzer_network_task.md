# Refactor Task: Analyzer Network Senders

## Problem Statement

The current implementation of sender analysis in the Neva compiler's analyzer is fragmented and contains significant logic duplication, specifically regarding `Union` senders.

## Identified Issues

### 1. Duplicated Logic for Union Senders

The logic for resolving union types, validating tags, and handling wrapped data exists in two places:

- `analyzeSender` in `internal/compiler/analyzer/senders.go`
- `getResolvedSenderType` in `internal/compiler/analyzer/network.go`

A comment in `network.go` (line 612) acknowledges this duplication:
_"logic of getting type for union sender partially duplicates logic of validating it... but it should be possible to refactor"_

### 2. Ambiguous Reachability and Ownership

In the primary connection analysis flow, the `Union` block in `getResolvedSenderType` is mathematically unreachable because `analyzeSender` handles the union and returns early. It is only called recursively by `getResolvedSenderType` itself when resolving predecessors of struct selectors.

### 3. Indirect Recursion loop

There is a complex loop between `analyzeSender` and `getResolvedSenderType`:

- `analyzeSender` calls `getResolvedSenderType` (for non-union senders).
- `getResolvedSenderType` calls `analyzeSender` (to resolve wrapped union data).
- `analyzeSender` calls itself (recursive union data analysis).

### 4. Computational Inefficiency in Chains

When analyzing struct selector chains (e.g., `node:port -> .field1 -> .field2`), the `getResolvedSenderType` function re-resolves the entire preceding chain for every link. For a chain of length $N$, this leads to $O(N^2)$ type resolutions because previous results are not cached or passed through.

## Proposed Resolution

1. **Unify Type Resolution**: Consolidate all "get type" logic into `getResolvedSenderType`. `analyzeSender` should rely on this function for types and focus on validation.
2. **Flatten Recursion**: Refactor the relationship between `analyzeSender` and `getResolvedSenderType` to remove indirect recursion.
3. **Optimized Chain Analysis**: Modify the connection analysis to pass the resolved type of the previous link forward, eliminating the $O(N^2)$ re-resolution of chain predecessors.
