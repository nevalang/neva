# Constants

Constants in Nevalang are static entities with compile-time known values. They must have explicit types and can be nested or reference each other. Constants provide static values for network computations and must be compatible with the inport type they are sent to. Constant definition starts with the `const` keyword, following by id and the message literal expression.

```neva
// primitive types
const a bool = true
const b int = 42
const c float = 42.0
const d string = 'Hello, world!'

// complex types
const e list<int> = [1, 2, 3]
const f dict<float> = { one: 1.0, two: 2.0 }
const g struct { b int, c float } = { a: 42, b: 42.0 }
```

## As Network Senders

This section briefly outlines how constants are used in networks. For detailed semantics, see the [network page](./networks.md).

### Constant References

It's possible to use constant as a network-sender by referring to it with `$` prefix.

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

### Message Literals

You can omit explicit constants; the compiler will create and refer to them implicitly.

```neva
def Inc(data int) (res int) {
    Add
    ---
    1 -> add:left
    :data -> add:right
    add -> :res
}
```

## Nesting and Referensing

Non-primitive constants can to other constants and implement infinite nesting. Examples:

```neva
type NestedStruct struct {
    a int
    b int
}

const constA int = 5
const constB int = 10

const nested NestedStruct = {
    a: constA,
    b: constB
}

const outerList list<NestedStruct> = [nested, nested]
```
