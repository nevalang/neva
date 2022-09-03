# Cancelation

Grasefull shutdown of all program (via context?)

# Many const nodes for every outport

Separated node for every const - it's simpler to look at

# Generics

```yaml
io:
  # easy to parse
  arg: [X, Y, Z] # args/param/params/generic/generics
  in:
    a: X
    b: Y
    c: Z
  out:
    d: Z
    e: Y
    f: X
```

## Alternative syntax

Harder to parse (how to unmarshal field with unknown name?)

```yaml
io<X Y Z>: # <X,Y,Z> / [X Y Z] / (X,Y,Z) / io(X Y Z)
  in:
    a: X
    b: Y
    c: Z
  out:
    d: Z
    e: Y
    f: X
```

# Unused output ports

Unused outports should be connected to special `/dev/null` analog (Eraser?)

# Checker

## Module validation rulesch

1. Module must have `in` and `out` with at least 1 port each
2. There must be no `Unknown` among port types
3. There must be `net` defining connections
4. There could be `workers` defing network's worker nodes
5. If there are `workers` then there must be also `deps` defining dependencies
6. Any `worker` must point to defined `dep`
7. Any `dep` must have valid `in` and `out` defined
8. There could be `const` defining constant values

## Checking network

9.  There must be no `deps` unused in `workers`
10. There must be no `worker` or `const` not used by `net`
