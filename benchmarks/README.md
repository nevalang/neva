# Runtime Benchmarks

This directory holds the runtime benchmark baseline for CI and release artifacts.

- `benchmarks/bench_test.go`: shared Go harness that builds `cmd/neva` once, compiles each fixture once, and times only execution of the produced `output` binary.
- `benchmarks/simple`: focused single-concern runtime paths.
- `benchmarks/complex`: mixed-feature runtime paths that approximate more realistic programs.

The e2e suite intentionally does not report `allocs/op`: the measured process is an external compiled binary, so harness allocations would be misleading as a runtime signal.

CI uploads the raw benchmark output as an artifact. Compare two runs locally with:

```bash
benchstat old.txt new.txt
```

To run the suite locally:

```bash
go test ./benchmarks -run=^$ -bench BenchmarkRuntimeE2E -benchtime=1x -count=1
```
