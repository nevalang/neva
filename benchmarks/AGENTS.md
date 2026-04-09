# AGENTS.md

This file defines benchmark authoring rules for both humans and machines.

## Scope

- Applies to all packages under `benchmarks/`.
- Current benchmark taxonomy has three tiers: `atomic`, `simple`, `complex`.
- Keep benchmark package paths shallow and stable.
- Prefer `benchmarks/<tier>/<pkg>/main.neva` for new work.
- Use one scoped prefix only when it carries real scenario meaning (for example `complex/control_flow/selectors_routers`).
- Do not reintroduce old taxonomy trees such as `benchmarks/simple/operators/<type>/...`.

## Tier Definitions

- `atomic`: one primary component under test.
- `simple`: a small composition validating one focused scenario.
- `complex`: multi-domain scenario closer to realistic pipelines.
- Startup baseline is a deliberate `atomic` exception (`startup_noop`) used to measure process/runtime overhead.

## Naming

- Benchmark package naming should be explicit and stable.
- For `atomic` builtins, prefer `<pkg>_<component>` with optional `_<context>` when needed.
- For `simple` and `complex`, prefer scenario names that describe the measured behavior, not implementation detail buckets.
- Do not add `_<type>` suffixes by default.
- Add type/context suffix only when a benchmark intentionally targets a distinct Neva runtime path or semantic mechanism.
- Do not split benchmarks by surface type when the difference would only re-measure native Go operators instead of a distinct Neva path.
- Avoid generic names like `basic`.

## Atomic Caveat

- Absolute atomic isolation is a goal, not a hard guarantee.
- Some components require technical support wiring to produce a valid standalone program.
- Keep support wiring minimal and document why it is required.
- Example: `streams.Range` requires a terminal sink such as `streams.Wait`.

## Iteration Policy

- By default, iteration belongs to the Go benchmark harness (`testing.B.Loop`).
- Prefer one-shot benchmark programs for builtins and small components.
- Use internal loops inside `.neva` only for deliberate throughput-oriented benchmarks.
- Treat those throughput cases as explicit special cases, not as the default `atomic` pattern.
- Example: `streams_range` measures stream throughput and may legitimately keep an internal range loop.

## Review Checklist

- Is the measured hot path explicit?
- Is support wiring minimal and justified?
- Is the benchmark name descriptive and non-ambiguous?
- Does the program compile and terminate deterministically?
