# Connections

Network is a graph that connects ports together for message passing.

## Array Bypass

Array bypass connection is a way to connect two array-ports together without explicitly specifing their slots.

```neva
s => r
```

Unlike _normal_ connection, members of array-bypass connection are always port addresses (without slot indexes).

## Normal

Normal connection connects one or more senders to one or more receivers or another connection that is deferred or chained.

### Pipe

Simple one to one connection.

```neva
s -> r
```

### Fan-In

Multiple senders and one receiver.

```neva
[s1, s2] -> r
```

Fan-In messages are _serialized_ - language guarantees that receiver will receive incoming messages in the exact same order they were sent by senders.

### Fan-Out

One sender, multiple receivers.

```neva
s -> [r1, r2]
```

This is a syntactic sugar over explicit `FanOut`. Desugared version of this connection looks like this:

```neva
s -> fanOut
fanOut[0] -> r1
fanOut[1] -> r2
```

### Fan-in + Fan-Out

Multiple senders, multiple receivers.

```neva
[s1, s2] -> [r1, r2]
```

Desugared version:

```neva
[s1, s2] -> fanOut
fanOut[0] -> r1
fanOut[1] -> r2
```

> Note fan-in pattern on line `[s1, s2] -> fanOut`. It means `fanOut` receives and sends messages in the same order they were sent. Therefore `r1` and `r2` will receive them in the same order.

### Deferred

Instead of specifing receiver(s) in the right side, we put another connection decorated by `(...)` braces. That connection is called _deferred_.

```neva
s1 -> (s2 -> r)
```

Deferred connection could be any connection, it can even contain other deferred connections - this way deferred connections are _nested_.

```neva
s1 -> (s2 -> (s3 -> r))
```

Deferred connection is a syntactic sugar over explicit `Lock`. It's called deferred because it right side sending until left-side sending happened by inserting lock node in the middle. This is how desugared version of the `s1 -> (s2 -> r)` looks:

```neva
s1 -> lock:sig
s2 -> lock:data
lock -> r
```

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
