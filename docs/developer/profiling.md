# Go runtime profiling

Profiling answers a performance or scheduling question about one reproducible
workload. It is not a CI quality gate: profiles add overhead and a profile
without a question is usually noise.

## Runtime baseline

Run the concurrent select workload:

```sh
make profile-runtime
```

The target executes `BenchmarkSelectHotpath` once for each artifact and writes
them, plus its generated test binary, to `/tmp/neva-runtime-profile`. Capturing
CPU, memory, blocking, mutex, and trace data in separate runs prevents one
profiler's overhead from distorting another's result. Set `PROFILE_DIR` to
retain them in a different location:

```sh
make profile-runtime PROFILE_DIR=/tmp/neva-profile-2026-08-22
```

The workload starts a runtime handler in one goroutine and exchanges messages
through two input channels. It is deliberately small: it provides a stable
starting point for investigating runtime scheduling without compiling every
benchmark fixture.

## Reading the artifacts

| File | Question it answers | First command |
| --- | --- | --- |
| `cpu.pprof` | Where did CPU time go? | `go tool pprof -top cpu.pprof` |
| `mem.pprof` | Which paths allocate, and which objects remain live? | `go tool pprof -top -sample_index=alloc_space mem.pprof` |
| `block.pprof` | Where do goroutines spend time blocked? | `go tool pprof -top block.pprof` |
| `mutex.pprof` | Where is contended mutex time attributed? | `go tool pprof -top mutex.pprof` |
| `trace.out` | How do goroutines run, block, wake, and schedule over time? | `go tool trace -http=localhost:0 trace.out` |

Run the commands from the profile directory, or provide an absolute path. For
memory retention rather than total allocation volume, use
`-sample_index=inuse_space`. `go tool pprof -http=localhost:0 <profile>` opens
an interactive local view for any pprof artifact.

`trace.out` can be much larger than the pprof files; treat it as a temporary
diagnostic artifact and use a dedicated `PROFILE_DIR` when retaining it.

An empty mutex profile means this workload did not observe mutex contention; it
does not prove that the runtime is mutex-free. Likewise, a block profile shows
where waiting happened, not whether that waiting is a deadlock. Use a trace to
investigate a suspected scheduling or blocking problem, then add a focused
regression test when the expected behavior is known.

## Choosing a different workload

Do not profile all e2e packages by default. Their Go test process is primarily
a harness that starts separately compiled CLI processes, so its profile does
not describe the generated program. For a new investigation, select one
benchmark or test that reproduces the behavior, use a dedicated output
directory, and record the command with the result.
