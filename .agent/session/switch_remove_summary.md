# Switch Syntax Removal Task Summary

## Original Task (from first message)

Remove all `switch` syntax construct.

**Before**: special syntax

```neva
data -> switch {
  case1 -> receiver1
  case2 -> receiver2
  _ -> default
}
```

**After**: normal node (with a bit of compiler magic (type inference for `case` outport when working with tagged unions)).

```neva
switch Switch<T>
---
data -> switch:data
case1 -> switch:case[0] -> receiver1
case2 -> switch:case[1] -> receiver2
sw:else -> default
```

## The Problem: Union Elimination

With the old language feature, `switch` wasn't "just wiring sugar". It acted like a _typed eliminator_ for tagged unions:

- Input: `U = union{ Int: int, Str: string, ... }`
- In each branch, payload type is **refined** (`int` in `Int` branch, `string` in `Str` branch)
- Branches can then continue with different component graphs that expect different payload types.

A normal generic component `Switch<T>` can only have a single `T` flowing through its ports. If `T = U`, then the output(s) are also `U`. But for pattern matching you need **N different payload types** (one per union member). In other words, a single switch wants outputs like `[int, string, ...]` (heterogeneous), which conflicts with Neva's typing rules (e.g., array ports share the same type). So "pattern matching + branch-specific graph execution" can't be modeled as an ordinary `Switch<T>` node.

Also, since we cannot pass "code/continuations" as messages, `Switch<T>` cannot take "branch handlers" as inputs the way functional languages do.

## Introduce `Type` Type

We could get to the middle ground in Go way similar to how it handled generic functions before generics.

For example:

```go
func len(v Type) int
```

Where `Type` is defined as

```go
// Type is here for the purposes of documentation only. It is a stand-in for any Go type, but represents the same type for any given function invocation.
type Type int
```

In our case that could be:

```neva
def Switch<T>(data T, [case] T) ([case] Type, else Type)
```

Compiler is aware of `Switch` and infers the type for its every usage.

## Invariant

This is expected behaviour of `Switch` (and tagged union network senders) as the outcome of this PR. We refer to `T` meaning type argument passed to `Switch<T>`.

- When `T` is not (resolved to) `bool` or `union`
  - Type of the `case[i]` array output port slot is the `T` itself, same for `else` outport
  - It is required to use `else` output port regardless how many `case` slots are used (compiler requires to use all the inports already, but must enforse using all the outports too, in this specific case)
- When `T` is (resolved to) `bool`
  - Type of the `case[i]` array outport slot is `T` itself (`bool` in this case), same for `else` outport
  - `else` is required if compiler can't prove that both `true` and `false` are covered. It's possible to prove only if both `case[0]` and `case[1]` inport slots are connected to constant references or constant literals. This is NOT top priority and should be implemented as separate step or at least commented as clear TODO.
- When `T` is (resolved to) `union`:
  - It is required to cover all union members and compiler must either prove they all covered or throw compilation error with explicit message
    - Compiler can prove that all members are covered if either 1) `else` outport is used or 2) there's exactly one `case` inport slot covering each union's member; It's ok have both all `case` covered manually and using `else` (in case union will be extended so code still compiles). What must NOT be possible is to have at least one union member not covered. This might lead to program deadlock (no branch executed given uncovered union member).
  - Type of the `else` outport is always the `T` itself because there's no unboxing possible in that case since else is triggered when no pattern matched.
  - Type of `case[i]` outport slot depends on the correspinding union member (matched pattern):
    - If there's a type-express (it's not tag-only union member), then that type-expression is the output (resolved sender) type for `case[i]`
    - Otherwise (if the correspinding union member has no type expr, i.e. it's a "tag-only" member) the type is just `T` (the union in this case) itself.
- Union network sender rules.
  - There are 2 types of union senders, the ones that are wrapping data and the ones that aren't. Wrapping: `Type::Member(wrapped_sender) -> ...`; Non-wrapping: `Type::Member -> ...`.
  - If corresponding union member has NO type expr (it's tag-only), it's required that it's NOT wrapping anything and used just like `U::Foo -> ...`. Example of such member: `type U union { Foo }`
  - If corresponding union member HAS type expr (body) it IS required for the network sender to wrap the data. For example in the `type U union {Foo int}` member `Foo` has type (`int`) i.e. it's not "tag-only" member. It means when referred from the network as sender `U::Foo` it must wrap the message (of type that resolves to integer) like this: `U::Foo(some_int_sender) -> ...`. This is already enforsed (but check and read the code for perfect understanding). **This rule has an EXCEPTION: if the network sender is used as sender for `Switch:case[i]` inport slot.** This is what's called "pattern matching context" in the `analyzer`. This is because we must have a way to send "just tag" to the `Switch` to compare the data message with the pattern message. This is the most complecated case where "unboxing" happens.

## Current State of the Code

We've made a LOT of changes while going to this but it's still not finished.

Use `git` to find out what changes were made to this proejct inside `remove_switch` branch. Ignore `*.md` files. If checking `*.neva` files only focus on failing one or stdlib (loads for every run).

To debug deadlock use neva CLI flags that allow you to produce IR and `trace.log` (see AGENTS.md and internal/cli for details).

Be attentive to changes made to go files in runtime and runtime/funcs. Be aware of Go 1.25 features such as `wg.Go(...)`.

## Source Code References

Important references (except mentioned ones plus the ones you'll find ourself (and you should lookup the repo for better understanding)).

Analyzer:

- network.go
- nodes.go
- senders.go
- receivers.go
- switch_logic.go

Stdlib (neva):

- routers.neva
- types.neva
- maybe something else in `builtin` pkg or `std` mod in general?

Runtime:

- message.go
- race.go
- switch.go

In general check `git` diff with `main` branch (but be VERY careful not to overload the context).

---

## Previous Agent Analysis

The previous agent analyzed the repo and identified:

**Task Summary**
You want to remove the `switch` syntax construct from the Neva language, replacing it with a normal node (Switch<T>) with special compiler magic for type inference. The work is already significantly in progress on the remove_switch branch, but there's a deadlock issue when running the examples/switch_fan_out test.

**My Understanding**

- Old Syntax: data -> switch { case1 -> receiver1; case2 -> receiver2; \_ -> default }
- New Syntax: Use Switch<T> component with explicit port connections
- The Problem: A deadlock occurs at runtime when executing switch_fan_out example
- Key Complexity: The Switch component needs special type inference for union types (unboxing)

**Clarifying Questions & Answers**

1. **Regarding the deadlock error**: The stack trace shows ArrayInport.ReceiveAll waiting forever on switch:case. Looking at main.neva:
   `'Alice' -> switch:case[0] -> [upper, lower]`
   This connects a constant 'Alice' to switch:case[0], but only uses one case slot. The runtime ReceiveAll waits for all slots to receive. Is this the intended behavior? Or should Switch not require all case slots to receive messages before matching?

   **Answer**: ReceiveAll implemented correctly as far as I remember. The thing is there should be generated only one slot (go channel if we talk about runtime) because only one slot is used. This can be checked with IR I believe.

2. **Are there supposed to be multiple case patterns?** The example seems to only send one pattern 'Alice', but ReceiveAll will block until it receives on all slots of the array inport. This seems like a semantic mismatch—could the runtime logic need to change to something like Select (race to first available) rather than ReceiveAll (wait for all)?

   **Answer**: If I'm not mistaken single `case[0]` is enough. We only care about single case, everything else is routed (should be) to `:else -> ...`.

3. **Pattern matching with Switch**: From your invariant description, it seems like case patterns should be sent once per cycle, not continuously. Currently, 'Alice' -> switch:case[0] sends a constant pattern on each cycle. Is this correct, or should patterns be statically defined (like part of component configuration)?

   **Answer**: In case of this example `Alice` is sent to `case[0]` exactly once (this is how `:start ->` works). In general I don't see an issue here, because everything (should be) synchronized by `switch`.

4. **What is the expected flow?**
   - User types "Alice" → scanln outputs "Alice" → goes to switch:data
   - 'Alice' constant is also sent to switch:case[0]
   - Switch compares data with case[0], matches, sends unboxed data to upper and lower
   - But do you expect the constant 'Alice' to be re-sent for each incoming data message? Or is the pattern static?

   **Answer**: Expected behaviour is somewhat literally described in the main.neva. Start arrives, msg sent to print, print finishes and triggers scanln, scanln reads and finishes, sends data scanned to switch:data, switch waits for data and all cases (single in this case), if mathing happens (if Alice was scanned) it sends to upper and lower (fan out), otherwise it sends a signal to `-> '...' -> err_new -> :err` which means "send ... literal to create err and router that err to err outport of us"

5. **Is the else outport required in this example?** The code uses switch:else -> '...' -> err_new -> :err, so it is connected. But the deadlock may be because only case[0] has a sender—meaning ReceiveAll waits for case slots that never receive.

   **Answer**: Else outport at the moment is not required by compiler, because compiler rules for outports is "at least use one". However invariant we strive to should require it, because string is not bool and not union (see invariant below).

## Previous Agent Implementation Plan

**Fix Switch Case Type Inference for Overload Resolution**

**Problem Description**
When a switch:case[i] output is used as a sender to an overloaded component (like `Add`), the overload resolution fails because it doesn't know the correct unboxed type.

**Example failing code**:

```neva
type IntOrString union { Int int; String string }
def Main(start any) (stop any) {
    switch Switch<IntOrString>
    add Add  // Add is overloaded: Add(int,int)->int, Add(float,float)->float, Add(string,string)->string
    ---
    IntOrString::Int -> switch:case[0] -> add:right  // Should infer `int`
    10 -> add:left
}
```

The switch:case[0] output should have type `int` (unboxed from IntOrString::Int), but overload resolution sees `Type` (the generic output type from Switch's signature).

**Root Cause**

- Overload resolution happens BEFORE network analysis in `analyzeNodes`
- `collectUsageDerivedTypeConstraintsForNode` calls `getPossibleSenderTypes` to gather type constraints
- `getPossibleSenderTypes` (line 1097-1115) doesn't handle Switch case outputs specially - it just returns the port type from the component's interface
- The Switch interface declares `([case] Type, else T)` - so it returns `Type` or any instead of the actual unboxed type
- The `getSwitchCaseOutportType` fix runs in `analyzeNormalConnection`, which is AFTER nodes are analyzed

**Proposed Changes**

[Analyzer - Overload Resolution]
[MODIFY] `nodes.go`
Add special handling in `getPossibleSenderTypes` for Switch:case outputs:

```go
// port-addr
if sender.PortAddr != nil {
-    // We don't care if it's Switch
-    // Because Switch is not overloaded.
+    // Special handling for Switch:case[i] outputs.
+    // When the input to case[i] is a union tag pattern, the output
+    // should be the unboxed member type, not the generic Type.
+    if isSwitchCasePort(*sender.PortAddr, nodes) {
+        if t := a.getSwitchCaseOutputTypeForOverloading(
+            *sender.PortAddr, nodes, scope, net,
+        ); t != nil {
+            return []typesystem.Expr{*t}
+        }
+    }
     pa := *sender.PortAddr
```

[MODIFY] `switch_logic.go`
Refactor `getSwitchCaseOutportType` to be usable from `nodes.go`:

- Add a new function getSwitchCaseOutputTypeForOverloading that takes similar parameters but works during overload resolution (without requiring fully resolved nodes)
- The function should find the sender for the corresponding switch:case[i] input and determine the union member type

**Alternative Simpler Approach**
Since the net (network) is available at overload resolution time, we can extract the sender type directly:

[MODIFY] `nodes.go`

- Add a helper function that checks if a port is Switch:case[i]
- Find the sender to that case input in the network
- If it's a union tag sender, return the unboxed type

The key insight is that `collectUsageDerivedTypeConstraintsForNode` already has access to net []src.Connection, so we can scan it to find what feeds into switch:case[i].

## Current Status & Issues Found

**Task List Status**

- Initial investigation - deadlock reproduction
- Analyze test failures
- Trace type inference flow for Switch
- Write implementation plan
- Fix Switch case type inference for chained connections
- Fix overload resolution with Switch unboxed types
- Run tests and verify
- Run linting
- Final verification

**Issues Found**

- `switch_union_matching_with_data` test fails: "no compatible overload found for node add"
- Type inference for switch:case[0] works in `getSwitchCaseOutportType` but doesn't propagate correctly to overload resolution for `Add`
- The `switch_fan_out` test sometimes passes, sometimes deadlocks (may be timing-related)

**Current Test Failures**

- `examples/switch_fan_out`: Deadlock with ArrayInport.ReceiveAll waiting forever
- `examples/switch`: Test timeout after 10 minutes (likely similar deadlock)

## Key Files Modified in Branch

**Analyzer Changes:**

- `internal/compiler/analyzer/switch_logic.go` - New file with Switch type inference logic
- `internal/compiler/analyzer/network.go` - Modified for Switch handling
- `internal/compiler/analyzer/nodes.go` - Modified for overload resolution
- `internal/compiler/analyzer/senders.go` - Modified for union sender rules
- `internal/compiler/analyzer/receivers.go` - Modified

**Runtime Changes:**

- `internal/runtime/funcs/switch.go` - Switch router implementation
- `internal/runtime/message.go` - Message matching logic
- `internal/runtime/funcs/registry.go` - Added switch_router

**Stdlib Changes:**

- `std/builtin/routers.neva` - Switch component definition
- `std/builtin/types.neva` - Type definitions
- Various .neva files updated from old switch syntax to new Switch<T> component

**New E2E Tests Added:**

- `e2e/switch_union_matching_no_data/` - Tests Switch with tag-only unions
- `e2e/switch_union_matching_with_data/` - Tests Switch with data-carrying unions

## Next Steps

1. **Diagnose switch_fan_out deadlock** - Generate IR to see if multiple case slots are being created when only one is used
2. **Fix overload resolution** - Implement getPossibleSenderTypes special handling for Switch:case outputs
3. **Add missing invariant tests** - Create e2e tests for all Switch invariant cases
4. **Run full test suite** - Ensure no regressions
5. **Lint and document** - Update AGENTS.md with insights learned
