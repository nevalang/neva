# AGENTS.md

This file defines benchmark authoring rules for both humans and machines.

## Scope

- Applies to all packages under `benchmarks/`.
- Current benchmark taxonomy has three tiers: `atomic`, `simple`, `complex`.

## Tier Definitions

- `atomic`: one primary component or operation under test.
- `simple`: a small composition validating one focused scenario.
- `complex`: multi-domain scenario closer to realistic pipelines.

## Naming

- Benchmark package naming should be explicit and stable.
- For `atomic`, prefer `domain_component_case` format.
- Avoid generic names like `basic`.

## Atomic Caveat

- Absolute atomic isolation is a goal, not a hard guarantee.
- Some components require technical support wiring to produce a valid standalone program.
- Keep support wiring minimal and document why it is required.
- Example: `streams.Range` requires a terminal sink such as `streams.Wait`.

## Review Checklist

- Is the measured hot path explicit?
- Is support wiring minimal and justified?
- Is the benchmark name descriptive and non-ambiguous?
- Does the program compile and terminate deterministically?
