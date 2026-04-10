# Runtime Benchmarks

This directory holds the runtime benchmark baseline for CI and release artifacts.

- `benchmarks/bench_test.go`: shared Go harness that builds `cmd/neva` once, compiles each fixture once, and times only execution of the produced `output` binary.
- `benchmarks/atomic`: one-shot hot-path baselines for individual components or explicit runtime overhead baselines.
- `benchmarks/simple`: focused single-concern runtime paths.
- `benchmarks/complex`: mixed-feature runtime paths that approximate more realistic programs.

The e2e suite intentionally does not report `allocs/op`: the measured process is an external compiled binary, so harness allocations would be misleading as a runtime signal.

CI uploads the raw benchmark output as an artifact. Compare two runs locally with:

```bash
benchstat old.txt new.txt
```

To run the suite locally:

```bash
make bench-runtime
```

Use this suite as the baseline for future runtime work, including comparisons against `#1004`.
