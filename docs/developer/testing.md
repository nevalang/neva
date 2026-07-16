# Testing

- `internal/**` unit tests validate implementation-level behavior.
- `e2e/**` contains isolated Neva modules with Go harnesses for language and
  standard-library regressions.
- `examples/**` are executable user-facing documentation; all examples belong
  to one module and must compile together.
- `internal/runtime/**/*_bench_test.go` measures native runtime functions and
  runtime primitives in isolation.
- `benchmarks/**` measure explicit language-level and composed runtime
  performance questions.

Run the smallest meaningful scope while iterating, then widen validation when a
change crosses compiler, runtime, or public standard-library boundaries.
Generated tests should carry a short intent comment. Runtime benchmarks never
substitute for unit tests.

## E2E and Examples

Each `e2e/` package is an independent Neva module. Run focused e2e packages
while iterating; use a broader e2e run when the change crosses module or
standard-library boundaries.

`examples/` is one module. A single example run still requires all examples to
compile, so examples must remain readable, deterministic where their topology
requires it, and suitable as user documentation.

## Benchmarks

Benchmarks have three tiers:

- `atomic`: one primary component under test;
- `simple`: a small composition for one focused scenario;
- `complex`: a multi-domain, more realistic pipeline.

Keep the layout flat by tier:
`benchmarks/<tier>/<pkg>_<component>[_<context>]/main.neva`. Support wiring is
allowed only when necessary to make a standalone program valid and should be
documented. Use the Go benchmark harness for iteration by default; internal
Neva loops are reserved for deliberate throughput measurements.
