## Port Runtime Representation

At runtime there is no `array-ports` but any port has `idx` instead.

It turns out that we need `idx` for only one thing - to allow operator get arr-port channels. Can we omit them?

### New ConnectionPoint

Introduce `ConnectionPoint` type and use it instead of new `PortAddr`:

```
Connection {
  from ConnectionPoint
  to []ConnectionPoint
}

ConnectionPoint {
  PortAddr PortAddr
  Meta ConnectionPointMeta
}

ConnectionPointMeta {
  Type ConnectionPointType
  StructField []string
}

ConnectionPointType = NormFieldReceiver | StructFieldReceiver
```

Problems:

- only `to` should be allowed to have `StructFieldReceiver` type

# Maps

Map `kk` keys to `vv` values and read `v` value by `k` key.
Duplicates not allowed

## Tuples arrport solution

Generic `map<K,V>` component that has two inports: `cfg` and `k`, and one `v` outport`. `cfg`is an array-inport of type`Tuple<K,V>`and`k`is a regular port of type`K`. It takes a list of `tuples`and returns a value from its`v`outport of type`V`.

Problems:

- Allowes duplicate keys
- Generic `map` component
  - Generics
- New `Tuple` data-type (could (should?) use `struct` instead (type-checking?))

## Map-ports solution

Just like we have array-ports we could have map-ports.
We could specify key when attaching to the port.

TODO: poorly thought of

## New ConnectionPoint type

# Runtime Port Addr Idx

## Getting arrport for operators

# Debugging

...

# SubStreams

## Structures instead of substreams?

## Signals instead of brackets?

Is it possible to use `sig` for infinite nesting of streams?

### Type Checking

In the classical FBP port can receive messages of different types and in `neva` it's forbidden.
