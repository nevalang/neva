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