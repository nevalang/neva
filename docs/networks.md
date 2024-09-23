# Networks

Component network is a directed graph connecting ports of component nodes for message passing. Vertices represent port-addresses, and edges are connections between them. Messages flow through these connections, which always have a sender and receiver side. The graph may contain cycles.

Simplest connection with port-address `in:start` on sender-side and `out:stop` on receiver-side - each time message is send from `in:start` it goes to `out:stop`.

```
:start -> :stop
```

Program execution occurs through data transformations and side-effects as messages move from one node to another. Starting from Main, execution flows through sub-components recursively, with native components as leaf nodes. Unlike controlflow, execution happens concurrently across multiple connections and components simultaneously.

## Port Address

The basic form of sender and receiver sides. Port-address consists of node-name, port-name, and array-port-slot-index. At least one part is always present. Node and port are separated by `:`, followed by slot index in `[]` if present.

### Only Port

```
:start -> :stop
```

When port-address starts with `:`, it means node names are omitted. This is only possible when a component refers to its own ports. As explained in the nodes section, this is syntax sugar for `in` and `out` IO nodes, e.g. `:start` and `:stop` are actually `in:start` and `out:stop`. Notably, `in` only has outports (as it's a sender for the internal network), while `out` only has inports (as it's a receiver).

### Node and Port

```neva
Foo, Baz
---
foo:bar -> baz:bax
```

Another basic form of port-address includes both node and port names. This is necessary when referring to nodes other than IO nodes.

### Only Node

It's recommended to omit the port name if a node has only one port. For example, if `Foo` has only one outport `bar`, but `Baz` has multiple inports, we can write `foo ->` instead of `foo:bar`. However, we can't omit the inport for `Baz` as the compiler needs to know which specific inport to use.

```neva
foo -> baz:bax
```

If `Baz` has only one inport, we can omit its port name too:

```neva
foo -> baz
```

### Node, Port and Index

When using array-ports (except for array-bypass connections), it's crucial to specify the port-slot index:

```neva
foo:bar[0] -> baz:bax
foo:bar[1] -> baz:bax2
foo:bar[2] -> baz:bax3
```

When working with array-ports, there must be no "holes". Always start with `[0]` and use all slots up to the maximum used. For example, this is incorrect:

```neva
foo:bar[0] -> baz:bax
foo:bar[1] -> baz:bax2
foo:bar[3] -> baz:bax3
```

`foo:bar[2]` is missing, creating a gap between `[1]` and `[3]`.

### Node and Index

The rule for omitting port names also applies to array-ports. If `Foo` has only one outport (even if it's an array-outport), we can omit its name. For simplicity, let's assume `baz` also has only one inport:

```neva
foo[0] -> baz
foo[1] -> baz
foo[3] -> baz
```

## Sender Side and Receiver Side

Each connection always has a sender and receiver side. There are 3 types of each, leading to 9 possible forms of a single connection. Connections can be infinitely nested, resulting in countless options. This is similar to how control flow languages combine expressions, operators, and statements. Don't be intimidated; we don't need to learn every possible combination.

### Sender Side Forms

There are 3 sender-side forms:

1. Port address
2. Constant reference
3. Message literal

#### Port Address Sender

We've seen these port-address senders:

```neva
:foo ->
foo ->
foo:bar
foo:bar[0]
foo[0] ->
```

#### Constant Reference Sender

Constants can be referenced in the network.

```neva
$foo -> // local constant
$foo.bar -> // important constant
```

It acts as an infinite loop, repeatedly sending the same message at the receiver's speed. To prevent potential resource leaks from message spamming, constant senders are always used in conjunction with special technics.

One way to work with them is to have a node with multiple inports, where at least one is connected to a port, not a constant. This limits the constants' speed to that of the port. Here's a simple example:

We'll create a component that increments a number using an addition component with 2 inports: `:acc` and `:el`. We'll use a constant `$one` for `:acc`, while `:el` receives dynamic values:

```neva
const one int = 1

flow Inc(data int) (res int) {
    Add
    ---
    $one -> add:acc
    :data -> add:el
    add -> :res
}
```

In this example, `add:acc` and `add:el` are synchronized. When `:data -> add:el` has a message, `add:acc` can receive. If the parent of `Inc` sends `1, 2, 3`, `add` will receive `acc=1 el=1; acc=1 el=2; acc=1 el=3` and produce `2, 3, 4` respectively.

Another way to synchronize constants with real data is to use deferred connections. We'll explore this in the receiver-side forms section.

#### Message Literal Sender

Sometimes it's convenient to refer to message values directly in the network without creating a dedicated constant. This works the same as using constants, as both are syntax sugar for creating an emitter-node with a bound message.

```neva
flow Inc(data int) (res int) {
    Add
    ---
    1 -> add:acc
    :data -> add:el
    add -> :res
}
```

Only primitive data-types (`bool`, `int`, `float`, `string`) can be used like this. `struct`, `list`, and `dict` literals are not allowed in the network.

<!-- Connections forms component's network. There are array-bypass connections and normal connections. Array bypass are very simple, normal takes many different forms. Connections are also have recursive hierarchy and can be mixed in a lot of forms.

## Sender/Receiver vs Inport/Outport

In this page we going to use letters like `s` and `r` to specify sender and receiver respectfully. Senders and receivers are terms that exist outside of the inport/outport. E.g. when you send message to your sub-component (node) you use it inport as receiver, but when message is received inside that sub-node, if it has it's own network (components are also recursive), it will use its inport as sender, because there's need to _send_ message from inport to somewhere else.

> Fan implementation fact: It's even better to think about inports and outports inside the network as of separate `io` node with it's inports and outports. Then for the network itself inports become outports of such node, while outports become inports.

## Array Bypass Connection

Connects array-inport of the component to either self array-outport or array-inport of the sub-node. Array-bypass connection never has slot index specified because it connects all the existing slots together. E.g. here we connect all slots of `s` with all slots of `r`.

```neva
s => r
```

## Normal

Normal connections are all connections that are not array-bypass.

### Normal Pipe

Simple one to one connection.

```neva
s -> r
```

### Normal Fan-In

Explicit fan-in - multiple senders and one receiver.

```neva
[s1, s2] -> r
```

Implicit fan-in:

```neva
s1 -> r
s2 -> r
```

### Normal Fan-Out

Explicit:

```neva
s -> [r1, r2]
```

Implicit:

```neva
s -> r1
s -> r2
```

### Normal Fan-in + Fan-Out

Explicit:

```neva
[s1, s2] -> [r1, r2]
```

You can imagine implicit version youself.

### Deferred

Each connection has left and right side. On the left side we have sender (port address, constant reference, etc.), on a right side receiver.

What if we could have another connection as a right side? And that connection would be "deferred" until message from left side comes? That would be deferred connections:

```neva
s1 -> (s2 -> r)
```

Deferred connection could be _any_ connection, it can even contain other deferred connections - this way deferred connections are _nested_.

```neva
s1 -> (s2 -> (s3 -> r))
```

on of the `s1 -> (s2 -> r)` looks:

Because deferred connection is form of right side, we omit different forms of the left side. Left side could be anything normal connection allows.

### Chained

Chained connection is, just like deferred one, a form of a right side of the connection:

```neva
s1 -> s2 -> r1
```

`s2 -> r1` is _chained_ connection here. Unlike deferred connections we do not use `(...)` braces. Note that even though chained and deferred connections look almost the same, they have different meaning. Deferred connection inserts implicit lock node in the middle. Chained connection does not insert anything. Is just a way of writing two connections like one. Here's desugared version of the connection above:

```neva
s1 -> s2
s2 -> r1
```

Chained connection only possible if intermediate node:

1. Have 1 (in/out)port and/or
2. Inport and outport with the same name are used

Example 1:

```
Lock, Println
---
42 -> lock:data -> println
```

`Lock` has 2 inports (`data` and `sig`, we don't show `sig` usage here) and 1 outport, chained connection is possible because it have inport and outport named `data`

Example 2:

```neva
nodes { println Println }
...
42 -> println -> :stop
```

`Println` have 1 inport `data` and 1 outport `sig`. Even though port names are different chaining is possible if we omit them. We can do that because compiler doesn't have to guess which port to pick. -->
