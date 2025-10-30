# Routers and Selectors Analysis Result

# Routers & Selectors — Design Evaluation & Verdict

This document provides the final design evaluation for routing/selection constructs in Neva/Nevalang, based on proposals #802–#807 and comprehensive analysis.

## Executive Summary
- **Verdict**: Adopt the full 2×2 matrix of constructs (if/switch/race × match/select) with all proposed features including selector-mode and final receivers.
- **Core insight**: The proposals correctly identify that dataflow needs both value-based and order-based decision making, and both routing and selection patterns.
- **Key strength**: Clean separation of concerns with predictable blocking semantics that map directly to the runtime's goroutine/channel model.
- **Main risk**: Surface complexity from multiple forms per construct, but this is justified by the expressiveness gains and clear composition rules.

## Design Evaluation Scores (0-5)
- **Expressiveness**: 4.5/5 — Covers all essential dataflow patterns; only missing built-in timeouts (composable)
- **Simplicity/Orthogonality**: 4/5 — Five constructs but clean 2×2 matrix; selector-mode adds complexity but is justified
- **Runtime Fit**: 4.5/5 — Maps perfectly to goroutines/channels; blocking semantics are implementable and predictable
- **Developer Experience**: 4/5 — Familiar control structures; clear error messages needed for common mistakes
- **Performance Predictability**: 4/5 — Backpressure model is intuitive; selector-mode latency scales with data senders

## Core Constructs (Final Design)
- **Routers (choose direction)**
  - `if`: Boolean routing with then/else branches
  - `switch`: N-way value routing with default `_` 
  - `race`: Order-based routing (first condition wins)
- **Selectors (choose message)**
  - `match`: Value-based message selection
  - `select`: Order-based message selection (no default)

## Key Design Decisions

### 1. Adopt All Proposed Constructs
**Decision**: Include all five constructs (if, switch, race, match, select) with their full feature sets.
**Rationale**: The 2×2 matrix (router/selector × value/order) covers essential dataflow patterns that cannot be cleanly expressed with simpler primitives.

### 2. Support Selector-Mode
**Decision**: Allow routers to also select messages, not just routes.
**Rationale**: Common pattern where you need both message selection and routing; atomic operation prevents race conditions.

### 3. Support Final Receivers
**Decision**: Allow constructs to have a continuation that receives result metadata.
**Rationale**: Enables observability and staging patterns; the "final" is a handoff barrier, not completion semantics.

### 4. Order-Based Constructs Are Intentionally Nondeterministic
**Decision**: `select` and `race` are nondeterministic by design; ties break by stable branch order.
**Rationale**: Order-based decisions are inherently timing-dependent; forcing determinism would be artificial and limiting.

### 5. Multiple Conditions Per Branch
**Decision**: Support union semantics for both value-based (deterministic) and order-based (first-wins) constructs.
**Rationale**: Common pattern for grouping related conditions; semantics are clear and composable.

## Coverage Analysis

### Value-Based Routing ✅
```neva
// Boolean routing
:is_admin -> if {
    then -> grant_access
    else -> deny_access
}

// Multi-way routing
:user_role -> switch {
    'admin' -> admin_dashboard
    'user'  -> user_dashboard
    _       -> error_page
}
```

### Order-Based Routing ✅
```neva
// First-wins routing
:request -> race {
    cache_hit -> serve_from_cache
    db_ready  -> serve_from_db
    _         -> serve_error
}
```

### Value-Based Selection ✅
```neva
// Message selection
:status_code -> match {
    200: 'success'
    404: 'not_found'
    _:   'error'
} -> log_message
```

### Order-Based Selection ✅
```neva
// Event-based selection
select {
    user_input: 'user_responded'
    timeout:    'timeout_occurred'
} -> handle_response
```

### Hybrid Selection + Routing ✅
```neva
// Selector-mode: choose both message and route
:condition -> if {
    then: success_data -> success_handler
    else: error_data   -> error_handler
} -> audit_log
```

## Static Semantics

### Determinism
- **Value-based** (if, switch, match): Deterministic given deterministic inputs
- **Order-based** (select, race): Intentionally nondeterministic; first condition wins

### Type System
- Branch payload types must unify within each construct
- Condition types must be comparable to data types for value-based constructs
- Result types: `IfResult<T>`, `SwitchResult<T,Y>`, `RaceResult<T>`

### Blocking Semantics
- Each construct waits for all required inputs before proceeding
- Sends to selected branch receivers, waits for handoff
- Then sends to final receiver (if present), waits for handoff
- Only then ready for next iteration

## Concurrency Model

### Runtime Mapping
- Each construct becomes a goroutine with unbuffered channels
- Blocking semantics map directly to channel send/receive
- Fan-out creates multiple sends; backpressure from slowest receiver

### Deadlock Prevention
- No cycles in the construct itself (enforced by type system)
- Order-based constructs can deadlock if conditions never fire; use timeouts or defaults

### Backpressure
- Fan-out: slowest receiver governs progression
- Selector-mode: waits for all data senders, even unchosen ones
- Final receiver: adds extra handoff barrier

## Design Alternatives Considered

### 1. Desugar to Smaller Core
**Alternative**: Core primitives (Join, Route, SelectFirst, FanOut) with surface sugar.
**Verdict**: Rejected. Surface constructs are more readable and the core would be too low-level.

### 2. Remove Selector-Mode
**Alternative**: Keep routers and selectors separate; compose them explicitly.
**Verdict**: Rejected. Atomic selection+routing is a common pattern; composition is more complex.

### 3. Make Some Constructs Library Components
**Alternative**: Only if/switch as syntax; match/select/race as stdlib.
**Verdict**: Rejected. All constructs are fundamental enough to warrant first-class syntax.

### 4. Different Condition Positioning
**Alternative**: Put conditions on the left for better chainability.
**Verdict**: Rejected. Right-positioned conditions match familiar if-statement syntax.

## Implementation Notes

### Error Messages
- "Branch payload types do not unify: expected T, found U on branch 'else'"
- "Duplicate receiver 'println' in construct body"
- "Condition type X not comparable to data type T"
- "Select cannot have default branch"

### Performance Characteristics
- Selector-mode latency scales with number of data senders
- Fan-out backpressure governed by slowest receiver
- Order-based constructs have O(1) decision time once condition fires

## Final Verdict

**Adopt the complete design as proposed.** The constructs provide essential dataflow primitives with clear semantics, good runtime fit, and sufficient expressiveness. The surface complexity is justified by the patterns they enable and the race-condition-free guarantees they provide.

The design successfully balances expressiveness with implementability, providing a solid foundation for dataflow programming in Neva.




---

*Notepad ID: 9a5caa7f-4870-492d-9bf1-f3ad2f105596*

*Created: 10/15/2025, 5:41:31 PM*

