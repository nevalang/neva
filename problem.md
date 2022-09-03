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



