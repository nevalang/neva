# Array Ports

For tasks like "substract numbers from first to last" where:

1. Arguments count unknown
2. Order matters

We need a way to represent ordered variadic arguments.
Not lists because types of messages could be different.

## Array-Outport case

What is the case for array-outport?
The real array-outport, not the array-inport threated like outport by it's network.

> In J.Paul's version program can ask if port has enough space
> and the module's designer must make sure there is no 'holes'

### Routing

Warning: this maybe `mapping` sub-case.

### Mapping

Map `n` keys to `m` values

#### Tuples arrport solution

Generic `map<K,V>` component that has two inports: `cfg` and `k`, and one `v` outport`. `cfg`is an array-inport of type`Tuple<K,V>`and`k`is a regular port of type`K`. It takes a list of `tuples`and returns a value from its`v`outport of type`V`.

Problems:

- Allowes duplicate keys (is it possible to type-check?)
- Generic `map` component
  - Generics
- New `Tuple` data-type (could (should?) use plain `structs` instead (type-checking?))

#### Map-ports solution

Just like we have array-ports we could have array-outports.
We could specify key when attaching to the port.

`in[]: Tuple(key,value)`

## Array-Outport indexing

Imagine some component with an array-outport in a module's network.
It's dangerous to refer specific indexes of that outport because we don't know how many indexes will be used inside that component.

Another case is when some module has an array-inport.
It's dangerous to refer some specific indexes of that inport because we don't know how many indexes will be used by parent network.

But port reference assumes some index!

Actually these two cases are the same though we talking about outport in the first one and about inport in the second.
This happens because for module's network - inport is actually an outport, the port to read from!

### "ONLY Array-Bypass" solution

ArrOutPort indexing is forbidden.

The solution is to forbig array outport indexing and introduce second type of port reference - "array bypass".
It's a type of port reference that is used when we want to connect some array-outport
with some array-inport. Both part of the connection must be array-ports.

#### PortAddr vs Connection?

Is it `PortAddr` or `Connection` should have this information?

## Order (sorting)

As it turns out there's no order guarantee in current implementation.

### Solution

Use `map[string][]chan Msg` instead of `map[PortAddr]chan Msg` at the `core` level.

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

# Blockers

It should be impossible to compile a program that blocks.

# Struct fields

How to read struct fields?

## Component solution

Introduce `struct-reader` component that has `struct` and `field` inports and `value` outport

### Problems

- Generics/TypeChecking
- Lots of `const` just to read struct field