# Operators

Operators are set of implicitly imported _native components_ with one type parameter `t` and single array inport `v[] t`. Output interface is component specific.

> Implementation detail:
> Any operator have

## Logic

Logic operators receives messages of type `t` and sends `bool` messages. All these operators share some "truthy/falsy" logic that is the following:

- `bool` is truthy if it's `true`
- `uint`, `int` and `float` are truthy if they're `> 0`
- `string`, `vector` and `dict` are truthy if they're not empty
- `struct` is truthy if all its fields are "truthy"

### Conditional AND `&&`

```
&&<t> (v[] t) (v bool)
```

Conditional AND checks messages, as soon as they arrive, whether they're "truthy". If some message is "falsy" it sends `false` to `v` outport and waits for the rest inports to receive their messages, but those messages do not have any effect. Then new iteration starts.

### Conditional OR `||`

```
||<t> (v[] t) (v bool)
```

Conditional OR works like `&&` but vice versa - it checks messages, as soon as they arrive, whether they're "falsy" and as soon as first "truthy" message arrives, it sends `true` to `v` outport and waits for the rest inports to receive their messages. Those messages do not have any effect. Then new iteration starts. It only sends `false` in one case - if all the inports received "falsy" messages.

### NOT `!`

```
!<t> (v[] t) ([]v bool)
```

NOT checks messages, as soon as they arrive, whether they're "truthy" or "falsy" and sends a an inverted boolean value to corresponding outport. As soon as all the inports have received their messages, new iteration starts.

## Arithmetic operators

These operators called "arithmetic" but they operates on all types. They receive messages of type `t` and send messages of type `t`. The logic is specific for operator and type.

### Sum `+`

```
+<t> (v[] t) (v t)
```

#### Bool

...

#### Numeric `uint, int, float`

...

#### String

...

#### Vector

...

#### Dict

...

#### Struct

...

### Difference `-`

...

### Product `*`

...

### Quotient `/`

...

### Remainder `%`

...

<!-- ## Relation

- ==
- !=
- <
- >
- <=
- > =

## Addition

- -
- -
- |
- ^

## Multiplication

- -
- /
- %
- <<
- > >
- &
- &^ -->
