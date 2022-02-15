# Array Ports

For tasks like "substract numbers from first to last" where:

1. Arguments count unknown
2. Order matters

We need a way to represent ordered variadic arguments.
Not lists because types of messages could be different.

<!-- ### Batch alternative

Introduce "accumulator components" that allows to collect substreams into batches:

```
d6, d5, d4, d3, d2, d1 ->
3 ->
  f(data, count)
    -> x
```

These components always take at least 2 args where one is `data` and second is `count`.

#### Pros

- no need for array-ports
  - no array-outport problem
    - no array-bypass problem
  - no array ordering problem

#### Cons

- `generics` must be implemented
- generic `list` must be implemented
- `lists` implies synchronization
- general purpose `batch` must be implemented
- `constants` must be used as not data
  - they need to be updated everytime algorithm changes
- sometimes multiple `batch` instances because messages could be of different types
- (bonus: P.Morrison decided to implement array-ports) -->

### SubStream alternative

How to merge different sources and keep order?

## Array-Outport case

Routing (map paths to handlers)

```
const.path_a ->
const.path_b ->
const.path_c ->
  router.inport

router.outport[0] -> handler_a
router.outport[1] -> handler_b
router.outport[2] -> handler_c
```

What is the case for array-outport?
The truth array-outport, not the array-inport threated like outport by it's network.

## Array-Outport usage

Imagine some component with an array-outport in a module's network.
It's dangerous to refer specific indexes of that outport because we don't know how many indexes will be used inside that component.

Another case is when some module has an array-inport.
It's dangerous to refer some specific indexes of that inport because we don't know how many indexes will be used by parent network.

But port reference assumes some index!

Actually these two cases are the same though we talking about outport in the first one and about inport in the second.
This happens because for module's network - inport is actually an outport, the port to read from!

### Array-Bypass solution

The solution is to introduce a second type of port reference that called "array bypass".
It's a type of port reference that is used when we want to connect some array-outport
with some array-inport. Both part of the connection must be array-ports.

#### PortAddr vs Connection?

Is it `PortAddr` or `Connection` should have this information?

<!-- ### Compiler solution

Maybe it's possible to allow indexes references in module's networks
Because compiler can check such connections?
TODO THINK

### Problem

Reading from array inport:

If module required specific count of elements for its array inport,
then it's not really array-inport because arguments no longer variadic.
Such module should use normal ports instad.

Actually this breaks the idea of variadic arguments because parent network
forced to use specific amount of inports. -->

## Outport case

1. Something for routing?
2. It has something to do with array bypass?

## Order (sorting)

As it turns out there is no order guarantees in current implementation

### Solution

Use `map[string][]chan Msg` instead of `map[PortAddr]chan Msg` at the `core` level.

#### Problem

More work to collector (not a big deal).

# SubStreams

Allows to arrange messages into peaces

-> . . . . ->

## Structures instead of substreams?

## Sig for SubStream?

Is it possible to use signals for infinite nesting of streams?

```
-> m m m m m ->
-> m m m m m ->
```
