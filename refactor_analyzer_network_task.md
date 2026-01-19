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

## Additional Idea: Decouple Overload Resolution From Network Analysis

### Background: Current Cycle

Overload resolution for nodes lives in `getInterfaceAndOverloadingIndexForNode` and `getNodeOverloadVersionAndIndex` in `internal/compiler/analyzer/nodes.go`. The current flow looks like this:

- **Node analysis** needs a resolved component interface to validate node instantiation (ports, type params, `#autoports`, `#extern`, etc.).
- **Overload selection** depends on **network usage**: `deriveNodeConstraintsFromNetwork` inspects how the node is wired (senders/receivers) to infer constraints and select the best overload.
- **Network analysis** in turn needs resolved node interfaces to type-check connections.

This introduces a circular dependency: nodes -> network -> nodes. The current implementation sidesteps this by doing a lighter network scan during node analysis (derive usage constraints early), then doing the full network validation later.

### Proposed Split

Treat node analysis as **two separate responsibilities**:

1. **Instantiation validation** (lightweight):
   - Resolve entity reference.
   - Validate number and shape of type arguments (`TypeArgs`, `#autoports` preconditions).
   - Validate DI arguments and syntactic correctness (but do not select a specific overload).
   - Return *candidate* overloads for later filtering.

2. **Overload selection** (network-driven):
   - During network analysis, use usage constraints (senders/receivers, chaining, inferred port types) to choose the best overload.
   - The "node interface" used by network analysis becomes a *set* of possible interfaces rather than a single resolved one.

### Concrete Data Shape Change

Instead of `nodesIfaces map[string]foundInterface`, return something like:

- `map[string][]foundInterface` or
- `map[string]nodeCandidates` where `nodeCandidates` contains:
  - list of candidate interfaces (one per overload)
  - precomputed metadata (e.g., overload index, `#extern` flag, `#autoports` eligibility)
  - resolved type args (if any)

This would keep the node analysis deterministic while deferring overload selection until network constraints are available.

### How This Helps

- Eliminates the "double check" on the network during node analysis.
- Removes the implicit coupling between `analyzeNodes` and `analyzeNetwork`.
- Makes overload resolution a first-class constraint solving step rather than a side-effect of node analysis.

### Implementation Notes and Risks

- **Constraint propagation**: `deriveNodeConstraintsFromNetwork` currently depends on resolved interfaces. It would need to accept candidate interfaces and filter them iteratively as constraints are gathered.
- **Error reporting**: "no compatible overload" vs "ambiguous overload" should be reported after network analysis (or when constraints become unsatisfiable).
- **Autoports (#autoports)**: This logic currently needs a concrete overload and type args to generate a synthesized interface. If overloads are deferred, the analyzer must either:
  - generate interfaces for each candidate up front (if type args are fixed), or
  - keep a lazy generator that can be resolved once constraints narrow down.
- **Native component overloads**: `isNativeComponentWithMultipleExterns` currently shortcuts via type args or inferred constraints. This logic should move to the network-driven phase, so it can leverage full usage info and avoid premature pruning.
- **Complexity risk**: If implemented naively, network analysis might have to reason across many candidates. A practical approach would be:
  - first filter by static checks (type arg count, `#extern`, interface compatibility with explicit ports),
  - then apply dynamic constraints from network senders/receivers.

### Suggested Direction

A good end state is a two-pass approach:

1. **Node candidate pass**: validate instantiation and build candidate overload set.
2. **Network constraint pass**: resolve overloads using full wiring constraints, then validate network strictly using the chosen overloads.

This preserves existing behavior but makes the dependency graph explicit and removes the need for a partial network scan during node analysis. It should also make the overload subsystem easier to test in isolation.*** End Patch"}}
