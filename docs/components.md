# Components

Components have a signature (interface), optional compiler directives, nodes, and network. They can be normal or native.

## Main Component

Executable packages must have a private `Main` component with:

- One `start` inport and one `stop` outport, both of type `any`
- No interface nodes
- Both nodes and network
- Not public

## Native Components

Components without implementation use `#extern` directive to refer to runtime functions. Typically found in `std` module.

## Runtime Function Overloading

Native components can use overloading with `#extern(t1 f1, t2 f2, ...)`. Requires one type parameter of type `union`.

## Normal Component

Implemented in source code with network and maybe nodes. Must not use `#extern` directive.

## Nodes

Nodes are instances of other components that can be used in component's network to perform computation.

If entity that node refers to (component or interface) have type-parameters, then type-arguments must be provided with node instantiation. If node is instantiated from component that requires dependencies, then other node instantiations must be provided as those dependencies.

There are 2 types of nodes:

1. IO Nodes - `in` and `out`, created implicitly, you omit node name when you refer to them (e.g. `:start`, `:stop` instead of `in:start` and `in:stop`)
2. Computational nodes - explicitly created by user, instances of entities:
   - Components (concrete/component nodes)
   - Interfaces (abstract/interface nodes)

## Dependency Injection (DI)

Normal components can have interface nodes, requiring DI. `Main` component cannot use DI.

## Component and Interface Compatibility

A component implements an interface if:

- Type parameters are compatible (count, order, names, constraints)
- Inports are compatible (equal amount, same names/kind, compatible types)
- Outports are compatible (equal or more amount, same names/kind, compatible types)
