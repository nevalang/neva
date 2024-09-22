# Networks

Connections forms component's network. There are array-bypass connections and normal connections. Array bypass are very simple, normal takes many different forms. Connections are also have recursive hierarchy and can be mixed in a lot of forms.

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
nodes { lock Lock }
...
42 -> lock:data -> println
```

`Lock` has 2 inports (`data` and `sig`, we don't show `sig` usage here) and 1 outport, chained connection is possible because it have inport and outport named `data`

Example 2:

```neva
nodes { println Println }
...
42 -> println -> :stop
```

`Println` have 1 inport `data` and 1 outport `sig`. Even though port names are different chaining is possible if we omit them. We can do that because compiler doesn't have to guess which port to pick.
