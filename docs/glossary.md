# Network

Network is graph of `nodes` connected through `ports`.

# Node

Node is member of `network` graph.
Like components nodes have their `inports` and `outports`.
There are always nodes created from `workers` and two special `in` and `out` nodes for io.

# Component

Component represents unit of computation.
Every component has interface - it's `inports` and `outports`.
There are two types of components: `modules` and `operators`.

# Module

Module is a component that depends on other components.
It's a composition of other modules and/or operators.
Every component created by user is a module.

# Operator

Operator is a built-in component that do not depend on any other component.
Operators could be threated like atomics modules without internal network.
Nodes created from operators are always leafs.

# Ports

Ports are interface of a node.
At runtime they are actually channels.
There are two types of ports - `normal port` and `array port`.

## Normal Port

## Array Port (PortGroup?)

Array port is a group of ports.
