# Switch Removal: Rationale for Analyzer, IR, and Backend Changes

This document explains why the recent changes were made, using concrete Neva and Go examples, and highlights open questions about necessity. The goal is to make the intent explicit and to capture the edge case that triggered the runtime deadlock.

## 1) Why Switch Coverage Is a Separate Analyzer Pass

### The problem being solved

Switch coverage rules depend on **global usage**, not on a single connection. For example:

```neva
switch Switch<string>
---
'Alice' -> switch:case[0]
// else not connected
```

Whether `switch:else` must be used depends on **what `T` is** and **what case ports are connected** across the entire network. This is not local to one `analyzeConnection` call.

### Why not enforce in `analyzeConnection`

`analyzeConnection` is connection-local: it sees one sender-to-receiver link at a time. It does not know:

- Which other `case[i]` slots were used
- Whether `switch:else` is used somewhere else
- Whether `T` resolves to `bool` vs `union` vs other

So if we tried to enforce coverage there, we’d need **cross-connection state** (a node-level aggregator) or risk false positives/negatives.

### Why place it after `analyzeNetPortsUsage`

Coverage checks depend on the **final port-usage map** (`nodesUsage`) built during connection analysis. This map already detects:

- which slots are used for array ports
- whether `switch:else` is connected

That makes it the natural data source for Switch coverage rules. In other words:

- `analyzeConnection`: validate individual links
- `analyzeNetPortsUsage`: validate global port usage
- `analyzeSwitchUsage` (new): validate global Switch coverage

It is conceptually a **post-pass over the usage graph**.

### Could it be inside `analyzeNetPortsUsage`?

Yes, and that would also be reasonable. The current placement is a separate call immediately after `analyzeNetPortsUsage`. This keeps the logic isolated and easier to reason about, but it could be merged into `analyzeNetPortsUsage` if you prefer a single global-pass function.

### Example: Switch<bool> coverage needs global information

```neva
switch Switch<bool>
---
true -> switch:case[0]
false -> switch:case[1]
```

Coverage can be proven only if both literals are used; this needs the usage map of all `case` slots, which is not visible in a single connection check.

## 2) Why Switch Example Programs Were Adjusted

### What changed

The `switch` examples now forward `switch:else` directly to `:err` and define `err` as `any`:

```neva
// We keep :err as any to forward the unmatched string to Panic.
def App(start any) (stop any, err any) {
    switch Switch<string>
    ---
    switch:else -> :err
}
```

This was done so the user-facing behavior matches test expectations such as:

```
Enter the name: panic: Bob
```

If `switch:else` goes through `errors.New`, the panic message becomes the error object (struct) instead of the original string. The tests expect the raw string. This change is semantically consistent and improves the examples.

## 3) The Deadlock: Root Cause and the IR/Backend Changes

### The runtime symptom

The `examples/switch_fan_out` test deadlocked on the else path (e.g., input `Bob`). Trace shows:

```
sent | app/switch:else | "Bob"
recv | app/err_new:data | "..."
sent | app/err_new:res | {"text": "..."}
```

But **panic never receives** the error, so the program stalls.

### The underlying issue (edge case)

This was triggered by **fan-in to the virtual outport `:err`**.

Neva network:

```neva
[switch:else, scanln:err, println:err] -> :err
```

This is legal and common. It means *multiple senders converge on one virtual port*.

When IR is reduced and Go channels are generated, the virtual port (`out:err`) was accidentally split into **multiple channels**, one per sender. The `panic` node only listened to one of them. If another sender fired (like `err_new` after `switch:else`), the message went to a different channel and got stuck with no receiver.

### Why this didn’t show up before

It appears the old `switch` syntax path avoided the specific fan-in shape (or the paths weren’t exercised), so the bug stayed hidden. The new switch-as-node form routes else as a normal outport, which surfaced the existing flaw.

So yes: **this is an edge case that existed but was not triggered** prior to the switch refactor.

### Fix 1: Graph reduction must not collapse fan-in targets

Change in `internal/compiler/ir/graph_reduction.go`:

- Graph reduction now stops at any receiver with multiple incoming connections.
- This prevents collapsing multiple senders into independent channels.

Example:

```neva
[print:err, scanln:err] -> :err
```

This must preserve a shared receiver node instead of “flattening” to:

```
print:err -> panic
scanln:err -> panic
```

Without fan-in preservation, each becomes its own channel and panic only listens to one.

### Fix 2: Backend must reuse channels for virtual ports

Change in `internal/compiler/backend/golang/backend.go`:

- Virtual ports (`in`/`out`) do not correspond to runtime funcs, so they are **special**.
- When multiple connections target the same virtual port, they must **reuse the same channel**.

Example:

```neva
[scanln:err, err_new:res] -> :err
:err -> panic
```

Now all fan-in senders share a single channel feeding `panic`.

### Are these changes necessary?

Given the deadlock trace, yes — at least one of these changes is required:

- **Option A**: Make graph reduction fan-in aware (current fix).
- **Option B**: Allow reduction but fix backend channel assignment to merge fan-in targets.

The fix currently implements **both**, which is arguably conservative. If you want a narrower change, we can drop one of them and rely on the other, but we would need to verify that all fan-in cases still behave correctly.

## 4) Summary of Rationale

- **Switch coverage pass** is global by nature and requires full port usage. It is not a local connection check.
- **Switch examples** were updated to match expected user-facing output and align with `panic` behavior.
- **Graph reduction and Go backend** changes address a fan-in-to-virtual-port edge case that surfaced with switch refactor; this deadlock existed but was previously untested.

## 5) Open Questions / Follow-ups

- Do we want `analyzeSwitchUsage` folded into `analyzeNetPortsUsage` for cohesion?
- Should we keep both graph reduction and backend channel fixes, or narrow to one?
- Do we want a focused test case for fan-in to `:err` to lock this behavior?

If you want, I can add a minimal e2e test specifically for the virtual `:err` fan-in to prove the necessity of the fix.
