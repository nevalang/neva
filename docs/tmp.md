  // TODO handle interface case: extend nodeCtx with DIArgs and use it if current node refers to interface
  // solution2 - make this func always called with component
  // (probably not because you would have to lookup for every node)

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
