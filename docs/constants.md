# Constants

Constants in Nevalang are static entities with compile-time known values. They must have explicit types and can be nested or reference each other. Constants provide static values for network computations and must be compatible with the inport type they are sent to.

## Constant References as Network Senders

To create a component that increments a number, we can use an addition component with 2 inports: `:acc` and `:el`. We'll use a constant `$one` for the `:acc` input, while `:el` will receive dynamic values:

```neva
const one int = 1

flow Inc(data int) (res int) {
    Add
    ---
    $one -> add:acc
    :data -> add:el
    add -> :res
}
```

Only primitive data-types are allowed to be used like this: `bool`, `int`, `float` and `string`. You can't use `struct`, `list` or `dict` literals in the network.

## Message Literals as Network Senders

You can omit explicit constants; the compiler will create and refer to them implicitly.

```neva
flow Inc(data int) (res int) {
    Add
    ---
    1 -> add:acc
    :data -> add:el
    add -> :res
}
```

## Semantics

Constants are sent in an infinite loop. Imagine a `SendOne` component that constantly sends ones. Its speed is limited by the receiver: `sendOne -> add:acc` sends a message each time `add:acc` can receive. In this example, `add:acc` and `add:el` are synchronized. When `:data -> add:el` has a message, `add:acc` can receive. If the parent of `Inc` sends `1, 2, 3`, `add` will receive `acc=1 el=1; acc=1 el=2; acc=1 el=3` and produce `2, 3, 4` respectively.

### Internal Implementation (`New` and `#bind`)

> You can skip this section. The above explanation is enough for writing programs. Here we discuss the implementation details of constant sending.

Both forms are syntax sugar. Here's the desugared form:

```neva
const one int = 1

flow Inc(data int) (res int) {
    #bind(one)
    New
    Add
    ---
    new -> add:acc
    :data -> add:el
    add -> :res
}
```

`New` with `#bind(one)` binds the `one` constant to the component instance, making it an emitter node that infinitely sends `1` to the receiver. This is similar to `SendOne` from the previous section. Constants cannot be rebound, and just a few components need `#bind`. The compiler usually handles this automatically.

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
