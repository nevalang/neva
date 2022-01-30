# Cons nodes rewriting

There should be separate node for every const.
It's simpler to look at.

# Generics

```yaml
io:
  arg: [X]
  in:
    s: str
    msg: X
  out:
    msg: X
    err: str
```

# Unused output ports

Unused outports should be connected to special `/dev/null` analog (Eraser?)

