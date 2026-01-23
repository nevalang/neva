# Union Task Continuation (Fresh Context)

## Context Summary

We removed union sender syntax and now use union literal constants (e.g., `MyUnion::Tag` or `MyUnion::Tag(value)`) as **const senders** in networks. This required analyzer changes for `Union<T Type>` and exposed regressions in Switch pattern matching.

There is also an open question about **Go test caching in VSCode**: e2e tests are invoked via `go test` (which can cache test results), while the CLI in e2e uses `go run` on `cmd/neva`. This may mean changed `*.neva` files don’t invalidate Go test cache, leading to stale “pass” results in VSCode UI. Need a mitigation (e.g., `-count=1`, explicit cache-busting in e2e harness, or test wrapper).

## Current Failing Tests (from `go test ./...`)

```
--- FAIL: Test (2.15s)
    e2e_test.go:11:
                Error Trace:    /Users/emil/projects/neva/pkg/e2e/e2e.go:97
                                                        /Users/emil/projects/neva/e2e/switch_union_matching_no_data/e2e_test.go:11
                Error:          Not equal:
                                expected: 0
                                actual  : 1
                Test:           Test
                Messages:       neva execution exit code mismatch. stdout: ""
                                stderr: "main/main.neva:18:18: sender for switch case inport not found\n"
FAIL
FAIL    github.com/nevalang/neva/e2e/switch_union_matching_no_data      3.660s
--- FAIL: Test (2.52s)
    e2e_test.go:11:
                Error Trace:    /Users/emil/projects/neva/pkg/e2e/e2e.go:97
                                                        /Users/emil/projects/neva/e2e/switch_union_matching_with_data/e2e_test.go:11
                Error:          Not equal:
                                expected: 0
                                actual  : 1
                Test:           Test
                Messages:       neva execution exit code mismatch. stdout: ""
                                stderr: "main/main.neva:13:1: std@0.34.0/builtin/operators.neva:31:8: no compatible overload found for node add (total components: 3, remaining: 0)\n"
FAIL
FAIL    github.com/nevalang/neva/e2e/switch_union_matching_with_data    4.123s
ok      github.com/nevalang/neva/e2e/type_expr_with_imported_type_arg   7.223s
ok      github.com/nevalang/neva/e2e/union_data_type_mismatch   4.495s
ok      github.com/nevalang/neva/e2e/union_invalid_type_arg     4.706s
ok      github.com/nevalang/neva/e2e/union_tag_not_const        4.821s
ok      github.com/nevalang/neva/e2e/union_tag_requires_value   4.476s
ok      github.com/nevalang/neva/e2e/union_tag_type_mismatch    4.632s
--- FAIL: Test (2.39s)
    e2e_test.go:11:
                Error Trace:    /Users/emil/projects/neva/pkg/e2e/e2e.go:97
                                                        /Users/emil/projects/neva/e2e/union_wrapping_test/e2e_test.go:11
                Error:          Not equal:
                                expected: 0
                                actual  : 1
                Test:           Test
                Messages:       neva execution exit code mismatch. stdout: ""
                                stderr: "main/main.neva:31:25: sender for switch case inport not found\n"
FAIL
FAIL    github.com/nevalang/neva/e2e/union_wrapping_test        4.651s
```

## Primary Regression to Fix: Switch + Union Pattern Matching

### Symptoms
- `switch_union_matching_no_data` fails with `sender for switch case inport not found`.
- `union_wrapping_test` fails with the same error.
- `switch_union_matching_with_data` fails with overload resolution for `Add`, likely because the `Switch` case output type is not inferred as `int` (tag type inference failed).

### Likely Root Cause
Old union sender syntax was removed. Switch analyzer logic still assumes case patterns are union senders or port‑addr senders. Now they are **union literal const senders** (with or without payload), which are not being recognized as valid pattern senders, especially in chained connections. That breaks:
- locating the sender for `switch:case[i]`
- inferring correct case output type
- overload resolution downstream (e.g. `Add`)

### Expected Fix Direction
- Update analyzer switch logic to treat **union literal const senders** (including refs) as pattern senders.
- Ensure `findSenderForSwitchCaseInput` can find the sender for `switch:case[i]` even when the chain head is a union literal const sender (and not a port addr).
- Update case type inference to derive tag from union literal const sender (similar to new Union tag logic).

## e2e Union Invariant Tests (now valid programs)

Four new e2e cases were added and recently fixed to avoid double‑use of `:start` by using syntax‑level fanout (`:start -> [ ... ]`). These are expected to fail only because of union invariant violations:
- `e2e/union_invalid_type_arg`
- `e2e/union_data_type_mismatch`
- `e2e/union_tag_requires_value`
- `e2e/union_tag_type_mismatch`

`e2e/union_tag_not_const` already committed and is valid.

## Open Question: VSCode / go test cache

We observed tests appearing to pass in VSCode UI due to Go test caching, even when `*.neva` changed. Since e2e uses `go run` to execute the CLI, Go’s cache may not be invalidated by changes in `*.neva` files. Investigate and propose a fix, likely one of:
- run `go test -count=1` for e2e
- modify `pkg/e2e` to include input hashes (or file mtimes) in test names to bust cache
- add a flag/env to disable cache for e2e runs

## Immediate Next Tasks (ordered)

1) Fix Switch regression so union literal const senders work as pattern senders (tag‑only and with payload).
2) Re‑run the three failing e2e tests:
   - `e2e/switch_union_matching_no_data`
   - `e2e/switch_union_matching_with_data`
   - `e2e/union_wrapping_test`
3) Only after Switch is fixed, investigate Go test caching / VSCode behavior.

## Notes / Reminders

- Union sender syntax must remain fully removed, including analyzer stages.
- Pattern matching should now assume union literals (const sender literal or const ref) as the tag pattern.
- Avoid reusing the same port as sender twice; use `-> [...]` fanout syntax.
- `:start` should trigger computation; do not connect `:start -> :stop` as a no‑op.
