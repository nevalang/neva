<!-- nodeCtx {
    path string
    ports [portAddr(name,idx)]: {
        ir.Msg|nil
        incomingConnectionsCount uint8
    }
    ??? di? | component? or component as a second arg?
} -->

- find main pkg and main component to use it as a root node
- get its io and create nodeCtx for root node with path ""
<!-- - call "generate" func with root node ??? -->
- for every connection generate IR connection add IR connection with the path from nodeCtx ---DONE---
  - while doing so, count incoming connections for every node inport in the network (to create new nodeCtx) and outgoing connections for every outport (to create void connections later when iterating nodes)
  - don't forget about giver and void connections! (void connections could be added later)
- iterate over nodes
  - if its refs to native then add IR func ref
  - if it does have static ports then create IR giver for every inport and corresponding message and add giver connection

generate(nodeCtx)

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

# Naming/Refactoring

- Rename trigers to lockers?

# Examples why FBP better

1. implicit state mutations are impossible
2. race conditions are impossible due to messages immutability (and lack of shared state?)
3. deadlocks are impossible because of no flow control (it's runtime who run goroutines and not end-user)
4. concurrent code isn't harder to maintain (real world is asynchronous)
5. type system: nominal typing and freedom from unnecessary type casts and mappings
6. type system: no nil pointer dereference
7. things goes exactly where we want (e.g. impossible to handle invalid response because of unhandled error. it's possible not to handle error, but that only means nothing will happen)
8. performance - concurrency by default everywhere (things to parallel are everywhere)
9. ready for future with lots of cores
10. always be as fast as go due to usage of go source code as a perfect and very high-level compile-target
11. perfect for visualizations and perfectly shows execution flow - data-charts as programs done right
12.
