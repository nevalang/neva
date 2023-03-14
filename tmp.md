# Generated go program structure

main.go
go.mod
runtime/runtime.go
std/...

# Analyzing

...

# Src to IR

- generating IR should happen in a tree manner from root to leafs
- every analyzed node on that path means adding ports to program
- every node means component. if it's interface then find out how it's instantiated
- every component means adding new connections (network) or operator (effect)

# IR to Go

- Simply include imports (for runtime, operators, code) or exactly copy/generate that code?
- Do we need `core` as a separate package?

# Issues

- deadlock detection?
- sure races are impossible?
- do we need transactions?
- ensure that any backwards compatible package could be used everywhere we assume lower forward compatible package (make semver incompatibility impossible)
- Blockchain package and docs management?


# Ideas
## GPT-like model for docs search

## Strict mode?

Probably bad idea

# Naming/Refactoring

- Rename trigers to lockers?
- Move givers and 


# Examples why FBP better

1) implicit state mutations are impossible
2) race conditions are impossible due to messages immutability (and lack of shared state?)
8) deadlocks are impossible because of no flow control (it's runtime who run goroutines and not end-user) 
3) concurrent code isn't harder to maintain (real world is asynchronous)
4) type system: nominal typing and freedom from unnecessary type casts and mappings
5) type system: no nil pointer dereference
6) things goes exactly where we want (e.g. impossible to handle invalid response because of unhandled error. it's possible not to handle error, but that only means nothing will happen)
7) performance - concurrency by default everywhere (things to parallel are everywhere)
8) ready for future with lots of cores
9) always be as fast as go due to usage of go source code as a perfect and very high-level compile-target
10) perfect for visualizations and perfectly shows execution flow - data-charts as programs done right
11) 