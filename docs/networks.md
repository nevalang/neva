# Networks

Component network is a directed graph connecting ports of component nodes for message passing. Vertices are senders and receivers, and edges are connections between them. Messages flow through connections. The graph may contain cycles.

Simplest connection with port-address `in:start` on sender-side and `out:stop` on receiver-side - each time message is send from `in:start` it goes to `out:stop`.

```neva
:start -> :stop
```

Program execution occurs through data transformations and side-effects as messages move from one node to another. Starting from Main, execution flows through sub-components recursively, with native components as leaf nodes. Unlike controlflow, execution happens concurrently across multiple connections and components simultaneously.

## Port Address

The basic form of sender and receiver sides. Port-address consists of node-name, port-name, and array-port-slot-index. At least node or port is always present. Node and port are separated by `:`, sometimes followed by slot index in `[]`.

### Only Port

```neva
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

#### Array-Ports Constraints

Components can't receive from their own array-inports by slot. Consider this scenario:

```neva
flow Foo([data]) (sig) {
    :data[0] -> ...
    :data[1] -> ...
    :data[3] -> ...
}
```

What if `Foo`'s parent only uses the first slot? Should it lead to deadlock or panic? `Foo` can't compute without all 3 inputs. To avoid this, we'd need to enforce that `Foo`'s parent always uses exactly 3 slots:

```neva
... -> foo[0]
... -> foo[1]
... -> foo[2]
```

This defeats the purpose of array-ports, which are needed for unknown situations and slot counts. Regular ports could handle fixed cases. To support this, we must restrict usage. Thus, components can't receive from their own array-inports by index.

However, components can send to and receive from other nodes' array-ports by index:

```neva
flow Foo(sig) (sig) {
    Bar
    ---
    :sig -> bar[0]
    2 -> bar[1]
    3 -> bar[2]
    bar -> :sig
}
```

But how do we operate with self array-ports? That's where "array-bypass" connections come in, which we'll cover later.

<!-- TODO explain outports -->

### Node and Index

The rule for omitting port names also applies to array-ports. If `Foo` has only one outport (even if it's an array-outport), we can omit its name. For simplicity, let's assume `baz` also has only one inport:

```neva
foo[0] -> baz
foo[1] -> baz
foo[3] -> baz
```

## Sender Side and Receiver Side

Each connection always has a sender and receiver side. There are 3 types of each, leading to 9 possible forms of a single connection. Connections can be infinitely nested, resulting in countless options. This is similar to how control flow languages combine expressions, operators, and statements. Don't be intimidated; we don't need to learn every possible combination.

### Sender Side

There are 3 sender-side forms:

1. Port address
2. Constant reference
3. Message literal

#### Port Address Sender

We've seen port-address senders:

```neva
:foo ->
foo ->
foo:bar ->
foo:bar[0] ->
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

**Internal Implementation**

> Implementation details of constant senders:

Const-ref and msg-literal senders are syntax sugar. In the desugared program, all senders and receivers are port-addresses. For constants and messages, a `New` component is used:

```neva
#extern(new)
pub flow New<T>() (msg T)
```

It's one of the few components without inports or outports, which are only allowed in stdlib. User-created components must have at least 1 inport and outport. New instances require the `#bind` directive to associate a constant with the node, allowing the runtime to use it throughout the program's lifecycle.

```
const p float = 3.14

flow Main(start) (stop) {
    #bind(p)
    New
    Println
    ---
    :start -> (new -> println -> :stop)
}
```

Message literal senders are implemented similarly, with the compiler inserting a virtual constant for the bind directive.

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

### Receivers Side

There are 3 types of receiver side:

1. Port-Address
2. Chained Connection
3. Deferred Connection

#### Port Address Receiver

We've seen port-address receivers:

```neva
-> :foo
-> foo
-> foo:bar
-> foo:bar[0]
-> foo[0]
```

#### Chained Connection

If a component has one inport and one outport

```neva
flow Foo(a) (b)
```

A "chained" connection is allowed

```neva
... -> foo -> ...
```

This is shorthand for two connections. The compiler infers port names when there's one port per side. Here's the expanded version:

```
... -> foo:a
foo:b -> ...
```

Here's an example using this feature:

```neva
flow Foo(data) (sig) {
    Println
    ---
    :data -> println -> :sig
}
```

`:data -> println -> :sig` combines port-address on the sender-side and chained connection on the receiver-side. `println -> :stop` is chained to `:data ->`. Here's a desugared version:

```neva
:data -> println
println -> :sig
```

Components don't need matching inport and outport names. Chained connections require one port per side. Both `flow Foo(bar) (bar)` and `flow Foo (bar) (baz)` are valid.

Chained connections can nest infinitely:

```neva
foo -> bar -> baz -> bax
```

Which translates to:

```neva
foo -> bar
bar -> baz
baz -> bax
```

#### Deferred Connection

In controlflow programming, instructions execute sequentially. In Nevalang's dataflow, everything happens concurrently. This might cause issues if we want to enforce a specific order of events.

Let's say we want to print 42 and then terminate.

```neva
42 -> println -> :stop
```

Turns out, this program is indeterministic and could give different outputs. The problem is that `42 ->` acts like an emitter sending messages in an infinite loop. Therefore, `42` might reach `println` twice if the program doesn't terminate quickly enough:

1. `42` received and printed by `println`
2. signal sent from `println` to `:stop`
3. new `42` sent and printed again
4. runtime processed `:stop` signal and terminated the program

To ensure `42` is printed once, synchronize it with `:start` using "defer". Here's the fix:

```neva
:start -> (42 -> println -> :stop)
```

This syntax sugar inserts a `Lock` node between `:start` and `42`. Here's the desugared version:

```neva
flow Main(start) (stop) {
    Lock, Println
    ---
    :start -> lock:sig
    42 -> lock:data
    lock:data -> println -> :stop
}
```

**Deferred connections defer receiving, not sending**. In `foo -> (bar -> baz)`, `bar` sends immediately, but `baz` receives the message only after `foo` unlocks it.

Deferred connections can nest infinitely:

```neva
foo -> (bar -> (baz -> bax))
```

Which translates to:

```neva
foo -> lock1:sig
bar -> lock1:data

lock1:data -> lock2:sig
baz -> lock2:data

lock2:data -> bax
```

#### Deferred + Chained

Deferred and chained connections can be combined in various ways. Here are a few examples:

```neva
a -> b -> (c -> d)
a -> (b -> c -> d)
a -> b -> (c -> d -> e)
a -> (b -> (c -> d -> e))
```

## Fan-in and Fan-out

Connections can have multiple senders and receivers, not just one-to-one. We'll explore these many-to-one (fan-in) and one-to-many (fan-out) scenarios next.

### Fan-in

Fan-in occurs when multiple senders share a single receiver. Messages are merged and received in the order they were sent. If messages are sent simultaneously, their order is randomized.

Here's a simple fan-in example:

```neva
[foo, bar] -> baz
```

Baz receives messages from `foo` and `bar` in the order they were sent. For example, if the sending order is `f1, b1, b2, f2`, Baz will receive them in this exact sequence.

**Internal Implementation**

The `[...] ->` syntax is syntactic sugar. In the desugared version, fan-in is implemented using the `FanIn` component from stdlib:

```neva
foo -> fanIn[0]
bar -> fanIn[1]
fanIn -> baz
```

### Fan-out

Fan-out occurs when one sender has multiple receivers. Messages are copied and sent to all receivers simultaneously. The sender waits for all receivers to process the message before sending the next one. This synchronization means faster receivers are limited by slower ones. To allow different processing speeds without data loss, programmers can explicitly add buffer nodes where needed.

```neva
foo -> [bar, baz]
```

In this scenario, `foo` and `bar` are fast, while `baz` is slow. The speed of all components is limited by the slowest one, `baz`. When `foo` sends message `1`, both `bar` and `baz` receive it simultaneously. `bar` processes it quickly and waits for the next message, while `baz` processes slowly. Only when `baz` is ready does `foo` send the next message `2`. This cycle continues until `foo` stops sending messages. Later, you'll learn how to optimize such bottlenecks using buffers.

**Internal Implementation**

`-> [...]` syntax is syntactic sugar. In the desugared version, fan-out is implemented using the `FanOut` component from stdlib:

```neva
foo -> fanOut
fanOut[0] -> bar
fanOut[1] -> baz
```

#### Fan-out + Deferred Connections

Fan-out can also be used with deferred connections:

```neva
:sig -> [
    foo,
    (bar -> baz)
]
```

This means `sig` sends messages to both `foo` and the lock-node controlling the `bar -> baz` connection.

> Fan-out + chained connection is WIP

### FanIn + FanOut

Fan-in and fan-out can be combined:

```neva
[a, b] -> [c, d]
```

## Array Bypass

There's one more connection type to discuss: array-bypass for components with array-ports.

```neva
flow FanInWrap([data]) (res)
```

Such components can't refer to their ports by index (e.g., `data[i]`). To operate on these ports, we use array-bypass.

```neva
flow FanInWrap([data]) (res) {
    FanIn
    ---
    :data => fanIn
}
```

The `=>` operator indicates an array-bypass connection, where both sender and receiver are always port-addresses without indexes. This connects all array-port slots, not just two specific slots. Array-bypass effectively creates multiple connections, one for each used slot.

Let's examine a specific example to understand how it works:

```neva
flow Main() () {
    wrap FanInWrap, Println
    ---
    1 -> wrap[0]
    2 -> wrap[1]
    3 -> wrap[2]
    wrap -> println -> :stop
}
```

In this example, `:data => fanIn` in `FanInWrap` expands to:

```neva
:data[0] -> fanIn[0]
:data[1] -> fanIn[1]
:data[2] -> fanIn[2]
```

Array-bypass connections adapt to the number of slots used by the parent. Without this feature, we'd need to create numerous variations of `FanInWrap` for different slot counts, potentially up to 255 (the maximum for array-ports), and even more for components with multiple array-ports.
