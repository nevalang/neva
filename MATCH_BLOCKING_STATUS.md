### current problem (low-level)

- **symptom**: `examples/match` prints "one" then blocks.
- **blocking edge**: the termination path’s `switch_router` never receives its `data` input (`wait/__switch__7:in:data`), so it waits forever.
- **where it stalls**: `case[0]` is received (`true`), but the `data` flow that should carry the `.last` field never arrives.

### evidence

- **trace**: shows repeated sends on the third fan-out branch but no receives on the “wait” chain’s `field` or `switch` data input.
  - sends exist: `sent | for/__fan_out__22:data[2] | {..., "last": false}`
  - case received: `recv | wait/__switch__7:case[0] | true`
  - missing: any `recv | wait/__field__9:in:data` or `recv | wait/__switch__7:in:data`
- **expectation from IR**: “wait” path should be wired: `fan_out.data[2] -> wait/__field__/in:data -> wait/__switch__/in:data`.

### status & progress

- **root-cause narrowed**: likely a backend channel wiring issue for array ports (3rd slot) into the `wait` chain. Either:
  - a missing/misnamed connection in IR `Connections`, or
  - incorrect path/port handling when building function IO in the Go backend.
- **verified components**:
  - `switch_router` waits on both `data` and all `case`s; receiving `case` without `data` will block as observed.
  - `GraphReduction` preserves array indices; unlikely the culprit.
  - Go backend groups array ports and constructs `runtime.NewArray(In|Out)port`; the issue is likely that the specific edge is missing/misnamed upstream or mismatched during channel map assembly.
- **tooling improvement**: enabled IR YAML to include connections to aid inspection going forward. See `ir.yml` and `trace.log` at repo root when running with `--emit-ir --emit-trace`.

### next steps

1. **dump/check connections just before Go templating**:
   - assert presence of the following addresses in `prog.Connections`:
     - `for/__fan_out__22/out:data[2]`
     - `wait/__field__9/in:data`
     - `wait/__switch__7/in:data`
2. **emit and inspect IR (JSON)** for `examples/match` to confirm the mapping from `fan_out.data[2]` → `wait/__field__9` → `wait/__switch__7`.
3. **inspect generated main.go**:
   - verify the channel var for `for/__fan_out__22/out:data[2]` is shared with `wait/__field__9/in:data`, and that the downstream from field to switch is present.
4. **unit test (IR-level)**:
   - assert the existence of edges: `fan_out.data[2] -> wait/__field__/data` and `wait/__field__/res -> wait/__switch__/data` for the `examples/match` network.
5. **fix**:
   - if missing at IRGen: correct receiver path normalization (`.../in` vs `.../out`, array flags/indices).
   - if present in IR but broken in Go: correct Go backend `buildPortChanMap`/IO assembly or path sanitization so lookups match and channels link.
