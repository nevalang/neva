# Operators

Operators are set of implicitly imported _native components_ with one type parameter `t` and single array inport `v[] t`. Output interface is component specific.

> Implementation detail:
> Any operator have

## Logic

Logic operators receives messages of type `t` and sends `bool` messages. All these operators share some "truthy/falsy" logic that is the following:

- `bool` is truthy if it's `true`
- `u8, u16, u32, u64, i8, i16, i32, i64, f32, f64, c64, c128` are truthy if they're `> 0`
- `str`, `vec` and `map` are truthy if they're not empty
- `arr` and `struct` are truthy if all their fields/elements are "truthy"

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

These operators called "arithmetic" but they operates on all types except `bool` and `struct`. They receive messages of type `t` and send messages of type `t`. The logic is specific for operator and type but they all have the same interface:

> num = u8 | u16 | u32 | u64 | i8 | i16 | i32 | i64 | f32 | f64 | c64 | c128

```
+<t> (v[] t ~ num | str | vec | map | arr | {}) (v t)
```

### Sum `+`

#### Numeric `8 | u16 | u32 | u64 | i8 | i16 | i32 | i64 | f32 | f64 | c64 | c128`

Sum adds numeric values as soon as messages arrives and sends result when all the messages have arrived. Then new iteration starts.

#### String

Sum concatenates string values in the order from left to right. When the last string arrived it sends the result to outport. Then new iteration starts.

#### Vector

Sum concatenates vectors in the order from left to right. When the last vector arrived it sends the result to outport. Then new iteration starts.

#### Map

Sum merges maps in the order from left to right. If some key was already present then its value is overwritten with the value from righter map. When last map is merged it sends the result to outport. Then new iteration starts.

#### Struct

Sum receives structs in the order from left to right. The first struct becomes an accumulator. Every next struct is summed with the accumulator so the result becomes new accumulator. Summing two structs means summing all its fields. When the last accumulator is calculated sum sends result to outport and new iteration starts.

#### Array

Works exactly like for structs.

### Difference `-`

#### Numeric `8 | u16 | u32 | u64 | i8 | i16 | i32 | i64 | f32 | f64 | c64 | c128`

Difference extracts numeric values in the order from left to right. When all the messages have arrived it sends result to outport. Then new iteration starts.

#### String

...

#### Vector

...

#### Map

Removes all key:value pairs found in right map from left map any time new message arrives. Processes messages as fast as they arrive. After that starts new iteration.

#### Struct

Works like Summing but does Difference.

#### Array

Works exactly like for structs.

### Product `*`

#### Numeric `8 | u16 | u32 | u64 | i8 | i16 | i32 | i64 | f32 | f64 | c64 | c128`

Product multiplies numeric values as soon as messages arrive. When all the messages have arrived it sends result to outport. Then new iteration starts.

#### String

...

#### Vector

...

#### Map

Works like Summing but if does not replaces already presented pairs.

#### Struct

Works like Summing but does Product and does not respect order.

#### Array

Works exactly like for structs.

### Quotient `/`

#### Numeric `8 | u16 | u32 | u64 | i8 | i16 | i32 | i64 | f32 | f64 | c64 | c128`

Quotient devides numeric values in the order from left to right. When all the messages have been arrived it sends result to outport. Then new iteration starts.

#### String

...

#### Vector

...

#### Map

Finds pairs that are presented in all maps. Does not respect the order.

#### Struct

Works like Difference but does devision.

#### Array

Works exactly like for structs.

### Remainder `%`

...

<!-- ^ = XOR -->
