# Networks

A component network is a directed, unweighted graph describing message passing. Vertices represent senders and receivers, while edges are connections between them. Messages flow through these connections, triggering data transformations and side effects. The graph may contain cycles.

Simplest connection with port-address `in:start` on sender-side and `out:stop` on receiver-side - each time message is send from `in:start` it goes to `out:stop`.

```neva
:start -> :stop
```

Computation occurs through as messages move from one place to another. Starting from `Main`, execution flows through sub-components recursively, with native components as leaf nodes. Unlike controlflow, execution happens concurrently across multiple connections and components simultaneously.

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
def Foo([data]) (sig) {
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
def Foo(sig) (sig) {
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

## Senders and Receivers

Each connection always has a sender and receiver side. There are 3 types of each, leading to 9 possible forms of a single connection. Connections can be infinitely nested, resulting in countless options. This is similar to how control flow languages combine expressions, operators, and statements. Don't be intimidated; we don't need to learn every possible combination.

### Senders

There are 5 sender-side forms:

1. Port address
2. Constant reference
3. Message literal
4. Struct selector
5. Range expression

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
$foo.bar -> // imported constant
```

When used as a standalone sender, it acts as an infinite loop, repeatedly sending the same message at the receiver's speed. Usually you need to synchronize constant sender with some event.

One way to do it is to have a node with multiple inports, where at least one is connected to a port, not a constant. This limits the constants' speed to that of the port. Here's a simple example:

We'll create a component that increments a number using an addition component with 2 inports: `:left` and `:right`. We'll use a constant `$one` for `:left`, while `:right` receives dynamic values:

```neva
const one int = 1

def Inc(data int) (res int) {
    Add
    ---
    $one -> add:left
    :data -> add:right
    add -> :res
}
```

In this example, `add:left` and `add:right` are synchronized. When `:data -> add:right` has a message, `add:left` can receive. If the parent of `Inc` sends `1, 2, 3`, `add` will receive `left=1 right=1; left=1 right=2; left=1 right=3` and produce `2, 3, 4` respectively.

Second way to synchronize constant senders is to use deferred connections.

```neva
:start -> { $msg -> println }
```

In this example `$msg` immidietely sends to `println` but locks by implicit `Lock` node, that waits for signal from `:start` inport. We'll learn how this works in details in the receiver-side forms section.

When used in a chained connection, the constant sender is triggered by the incoming message.

```neva
:start -> $msg -> println
```

In this example, when `:start` fires, it triggers sending the constant `$msg` to `println`. This is preferred (easier to read and with better performance) way. It's not always possible to use it though. We'll learn about chained connections in receiver side forms section.

**Internal Implementation**

> Implementation details of constant senders:

Const-ref and msg-literal senders are syntax sugar. In the desugared program, all senders and receivers are port-addresses. For constants, there are two components used depending on the context:

1. For standalone constant senders, a `New` component is used:

```neva
#extern(new)
pub def New<T>() (msg T)
```

It's one of the few components without inports/outports (this is only possible inside stdlib). `New` instances require the `#bind` directive to associate a constant with the node, allowing the runtime to use it throughout the program's lifecycle.

2. For constants in chained connections, a `NewV2` component is used:

```neva
#extern(new)
pub def NewV2<T>(sig any) (msg T)
```

`NewV2` also requires `#bind`.

> Message literal senders are implemented similarly, with the compiler inserting a virtual constant for the bind directive.

#### Message Literal Sender

Sometimes it's convenient to refer to message values directly in the network without creating a dedicated constant. This works the same as using constants, as both are syntax sugar for creating an emitter-node with a bound message.

```neva
def Inc(data int) (res int) {
    Add
    ---
    1 -> add:left
    :data -> add:right
    add -> :res
}
```

Only primitive data-types (`bool`, `int`, `float`, `string` and tagged union literals) can be used like this. `struct`, `list`, and `dict` literals are not allowed in the network.

- bool: `true -> ...` or `false -> ...`
- int: `42 -> ...`
- float: `42.0 -> ...`
- tagged union: `Day::Friday ->` or `Day::Friday(value) ->`

#### Struct Selector

In Nevalang, a struct selector sender allows you to access specific fields within a structured data type directly within the network. This feature is particularly useful when you need to extract and send specific data fields from a struct to other nodes in the network.

**Syntax**

The struct selector is used in conjunction with a chained connection. It uses the dot `.` notation to access fields within a struct. Here's the general form:

```neva
foo -> .bar.baz -> bax
```

In this example, `foo` is a sender that outputs a struct. The struct selector `.bar.baz` accesses nested fields within the struct, and the value of `baz` is sent to the receiver `bax`.

Connection that has struct selector as sender must be chained. In this example `.bar.baz -> bax` is chained to `foo ->`.

Compiler will ensure that `foo` is a struct that has `bar` field, which is itself also a struct with a `baz` field, and that type of `baz` field is compatible to what `bax` expects.

**Example**

```neva
import { fmt }

type Person struct { age int }

def Foo(person Person) (sig any) {
    fmt.Println
    ---
    :person -> .age -> println
    [println:res, println:err] -> :sig
}
```

**`Field` component**

Struct selectors are a syntactic sugar. Under the hood, they are translated into a series of `Field` components that access the specific fields within a struct. Each `Field` component extracts a single field from the struct and passes it along the network.

```neva
def Field<T>(data struct {}) (res T)
```

For example, the struct selector `.bar.baz` in the connection `foo -> .bar.baz -> bax` is desugared into:

```neva
#bind(path1)
field1 Field1<T1>
#bind(path2)
field2 Field2<T2>
#bind(path3)
field3 Field3<T3>
---
foo -> field1
field1 -> field2
field2 -> bax
```

As you can see `Field` is one of few components that are expected to be used with `#bind` directive so it's much better to just use `.` dot notation instead.

#### Range Expression

A range expression sender allows you to generate a `stream<int>` of messages within a specified range.

```
sig -> 0..100 -> receiver
```

In this example we generate stream of 100 integers from `0` up to `99` - that is, range is exclusive.

> Only message literals (integers) are supported at the moment, but in the future we'll allow different kinds (e.g. port-addresses) of senders for more flexible ranging

Negative ranging is also supported

```
sig -> 100..0 -> receiver
```

**How it works**

Range expressions is syntax sugar over explicit `Range`:

```neva
def Range(from int, to int, sig any) (res stream<int>)
```

`Range` component waits for all 3 inports to fire, then emits a stream of `N` messages. You are free to use range as a normal component, but you should prefer `..` syntax whenever possible.

### Receivers

There are 4 types of receivers:

1. Port-Address
2. Chained Connection
3. Deferred Connection
4. Switch

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
def Foo(a) (b)
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
def Foo(data) (sig) {
    Println, Panic
    ---
    :data -> println -> :sig
    println:err -> panic
}
```

`:data -> println -> :sig` combines port-address on the sender-side and chained connection on the receiver-side. `println -> :stop` is chained to `:data ->`. Here's a desugared version:

```neva
:data -> println
println -> :sig
```

Components don't need matching inport and outport names. Chained connections require one port per side. Both `def Foo(bar) (bar)` and `def Foo (bar) (baz)` are valid.

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

Turns out, this program is indeterministic and could give different outputs. The problem is that `42 ->` acts like an emitter sending messages in an infinite loop. Therefore, `42` might reach `print` twice if the program doesn't terminate quickly enough:

1. `42` received and printed by `println`
2. signal sent from `println` to `:stop`
3. new `42` sent and printed again
4. runtime processed `:stop` signal and terminated the program

To ensure `42` is printed once, synchronize it with `:start` using "defer". Here's the fix:

```neva
:start -> { 42 -> println -> :stop }
```

This syntax sugar inserts a `Lock` node between `:start` and `42`. Here's the desugared version:

```neva
def Main(start any) (stop any) {
    Lock, Println, Panic
    ---
    :start -> lock:sig
    42 -> lock:data
    lock:data -> println -> :stop
    println:err -> panic
}
```

**Deferred connections defer receiving, not sending**. In `foo -> { bar -> baz }`, `bar` sends immediately, but `baz` receives the message only after `foo` unlocks it.

Deferred connections can nest infinitely:

```neva
foo -> {bar -> {baz -> bax}}
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
a -> b -> {c -> d}
a -> {b -> c -> d}
a -> b -> {c -> d -> e}
a -> {b -> {c -> d -> e}}
```

#### Switch

Another type of sender was added to simplify the use of the `Switch` component. This is necessary when you need to trigger different branches of the network based on the value of an incoming connection. The syntax is as follows:

```neva
s -> switch {
    c1 -> r1
    c1 -> r2
    c1 -> r3
    _ -> r4
}
```

Here `s` means sender, which could be any sender. `c1, c2, c3` are "case senders" - they are also senders, and any senders will work as long as they are type-safe. Finally, `_` is the default sender. The default branch is required, making each switch expression exhaustive. The compiler ensures that the incoming `s ->` and all `c` and `_` senders are compatible with their corresponding receiver parts.

If one branch is triggered, other branches will not be (until the next message, if the corresponding pattern fires) - one way to think about this is that every branch has a "break" (and there's no way to "fallthrough").

> I don't like this explanation because it's control-flow centric, while Nevalang's switch is pure dataflow, but it makes sense as an analogy.

The switch is syntactic sugar for the `Switch` component:

```neva
def Switch<T>(data T, [case] T) ([case] T, else T)
```

> You are allowed to use `Switch` as a component, if you need to, but prefer statement syntax if possible

To better understand the `switch` statement, let's look at a few examples:

```neva
// simple
sender -> switch {
    true -> receiver1
    false -> receiver2
    _ -> receiver3
}

// multiple senders, multuple receivers
sender -> switch {
    [a, b] -> [receiver1, receiver2]
    c -> [receiver3, receiver4]
    _ -> receiver5
}

// nested
eq Eq
gt Gt
---
... -> 1 -> eq:left
... -> 1 -> eq:right
... -> 2 -> gt:left
... -> 3 -> gt:right
sender -> switch {
    true -> switch {
        eq:res -> receiver1
        gt:res -> receiver2
        _ -> receiver3
    }
    false -> receiver4
    _ -> receiver5
}

// as chained connection
sender -> .field -> switch {
    true -> receiver1
    false -> receiver2
    _ -> receiver3
}
```

**Multuple senders/receivers**

In this example

```neva
sender -> switch {
    [a, b] -> [receiver1, receiver2]
    c -> [receiver3, receiver4]
    _ -> receiver5
}
```

Case senders `a` and `b` are concurrent to each other, the one that will send faster, will be used by switch as a case value. This might be counter intuitive, because one might expect that this works like in controlflow languages where multple cases on a same line means "either".

Multiple receivers on the other hand work as expected. I.e. if `sender` message is equal to `c` in this example, then it will be sent to both `receiver3` and `receiver5`.

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
def FanInWrap([data]) (res)
```

Such components can't refer to their ports by index (e.g., `data[i]`). To operate on these ports, we use array-bypass.

```neva
def FanInWrap([data]) (res) {
    FanIn
    ---
    :data => fanIn
}
```

The `=>` operator indicates an array-bypass connection, where both sender and receiver are always port-addresses without indexes. This connects all array-port slots, not just two specific slots. Array-bypass effectively creates multiple connections, one for each used slot.

Let's examine a specific example to understand how it works:

```neva
def Main() () {
    wrap FanInWrap, Println
    ---
    1 -> wrap[0]
    2 -> wrap[1]
    3 -> wrap[2]
    wrap -> println
    [println:res, println:err] -> :stop
}
```

In this example, `:data => fanIn` in `FanInWrap` expands to:

```neva
:data[0] -> fanIn[0]
:data[1] -> fanIn[1]
:data[2] -> fanIn[2]
```

Array-bypass connections adapt to the number of slots used by the parent. Without this feature, we'd need to create numerous variations of `FanInWrap` for different slot counts, potentially up to 255 (the maximum for array-ports), and even more for components with multiple array-ports.

## All Together

Here's an intentionally complex example combining different features. This level of complexity should be avoided in practice; consider decomposing such logic into separate components. However, it demonstrates the potential intricacy of connections. Note the recommended formatting for nested connections.

```
a -> [
    b -> c -> d,
    {
        $e -> f -> 42 -> g -> [
            h,
            [i, j] -> k -> { l -> m },
        ]
    }
]
```
