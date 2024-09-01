# Component Entity

Component always has signature (basically embedded _interface_) and optional _compiler directives_, _nodes_ and _network_. There are two kinds of components: _normal_ and _native_ ones.

## Main Component

Executable package must have _component_ called `Main`. This component must follow specific set of rules:

- Must be _normal_
- Must be _private_
- Must have exactly 1 inport `start`
- Must have exactly one outport `stop`
- Both ports must have type `any`
- Must have no _abstract nodes_

Main component doesn't have to have _network_ but it is usually the case because it's the _root_ component in the program.

## Native Components

Component without implementation (without nodes and network) must use `#extern` directive to refer to _runtime function_. Such component called _native component_. Native components normally only exist inside `std` module, but there should be no forced restriction for that.

## Normal Component

Normal component is implemented in source code i.e. it has not only interface but also nodes and network, or at least just network. Normal components must never use `#extern` directive.

## Nodes

Nodes are things that have inports and outports that can be connected in network. There's two kinds of nodes:

1. IO Nodes
2. Computational nodes

IO nodes are created implicitly. Every component have one `in` and one `out` node. Node `in` has outports corresponding to component's interface's inports. And vice versa - `out` node has inports corresponding to component interface's inports.

Computational nodes are nodes that are instantiated from entities - components or interfaces. There's 2 types of computational nodes: concrete and abstract. Nodes that are instantiated from components are _concrete nodes_ and those that instantiated from interfaces are _abstract nodes_.

Interfaces and component's interfaces can have type parameters. In this case node must specify type arguments in instantiation expression.

## Dependency Injection (DI)

Normal component can have _abstract node_ that is instantiated from an interface instead of a component. Such components with abstract nodes needs what's called dependency injection.

I.e. if a component has dependency node `n` instantiated with interface `I` one must provide concrete component that _implements_ this interface.

Dependency Injection can be infinitely nested. Component `Main` cannot use dependency injection.

## Component and Interface Compatability (Implementation)

Component _implements_ interface (is _compatible_ with it) if type paremeters, inports and outports are compatible.

Type parameters are compatible if their count, order and names are equal. Constraints of component's type parameters must be compatible with the constraints of the corresponding interface's type parameter's constraints.

Component's inports are compatible with the interface's if:

1. Amount is exactly equal
2. They have exactly the same names and _kind_ (array or single)
3. Their types are _compatible_ (are _subtypes_ of) with the corresponding interface's inports

Outports of a component are compatible with the interface's if:

1. Amount is equal or more (this is only difference with inports)
2. Exactly the same names and _kind_
3. Their types are _compatible_
