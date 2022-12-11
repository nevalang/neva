# Builtin components

Special subset of components in [standard library](std.md) that provide basic functionality. Can be represented as [native](native.md) or [custom](custom.md) components depending on implementation.

## Logic

Logic operators receive messages of type `t` and send `bool` messages. These operators share "truthy/falsy" logic that is the following:

- `bool` is truthy if it's `true`
- `u8, u16, u32, u64, i8, i16, i32, i64, f32, f64, c64, c128` are truthy if they're `> 0`
- `str`, `vec` and `map` are truthy if they're not empty
- `arr` and `struct` are truthy if all their fields/elements are "truthy"

### Conditional AND `&&`

```
&&<t> (v[] t) (v bool)
```

Checks messages, as soon as they arrive, whether they're "truthy". If message is "falsy" it sends `false` to outport and waits for the rest inports to receive their messages, but those messages do not have any effect. Then new iteration starts.

### Conditional OR `||`

```
&&<t> (v[] t) (v bool)
```

Works like `&&` but vice versa - it checks messages, as soon as they arrive, whether they're "falsy" and as soon as first "truthy" message arrives, it sends `true` to outport and waits for the rest inports to receive their messages but those messages do not have any effect. It only sends `false` in if all inports received "falsy" messages. Then new iteration starts.

### NOT `!`

```
!<t> (v[] t) (v[] bool)
```

Checks messages, as soon as they arrive, whether they're "truthy" or "falsy" and sends a an inverted `bool` to corresponding slot. As soon as all the inports have received their messages, new iteration starts.

## Arithmetic

Despite their name they operate not only on numeric types. They receive messages of type `t` and send messages of the same type `t` but logic and interface is specific for every operator and type.

### Sum `+`

```
+<t num | str | vec | map | [] | {}> (v[] t) (v t)
```

#### Numeric `num`

Sum adds numeric values as soon as messages arrives and sends result when all the messages have arrived. Then new iteration starts.

#### String `str`

Sum concatenates string values in the order from left to right. When the last string arrived it sends the result to outport. Then new iteration starts.

#### Vector `vec`

Sum concatenates vectors in the order from left to right. When the last vector arrived it sends the result to outport. Then new iteration starts.

#### Map `map`

Sum merges maps in the order from left to right. If some key was already present its value is replaced with the value from arrived map. When last map is merged it the result to outport. Then new iteration starts. Being applied to set (map of empty structs) it acts like a union from set theory.

#### Struct `struct`

Sum receives structs in the order from left to right. The first struct becomes an accumulator. Every next struct is summed with the accumulator so the result becomes new accumulator. Summing two structs means summing all its fields. When the last accumulator is calculated sum sends result to outport and new iteration starts.

#### Array `arr`

Works exactly the same as with structures.

### Difference `-`

```
-<t num | map | arr | struct> (v[] t) (v t)
```

#### Numeric

Difference extracts numeric values in the order from left to right. When all the messages have arrived it sends result to outport. Then new iteration starts.

#### Map

Removes from left map all keys found in right map any time new message arrives. Processes messages as soon as they arrive. After that starts new iteration. Being applied to set acts like a complement from set theory.

#### Struct

Works like Summing but looks for difference.

#### Array

Works exactly the same as with structures.

### Product `*`

```
*<t num | map | [] | {}> (v[] t) (v t)
```

#### Numeric

Product multiplies numeric values as soon as messages arrive. When all the messages have arrived it sends result to outport. Then new iteration starts.

#### Map

Works like Summing but doesn't replace already presented pairs.

#### Struct

Works like Summing but does Product and does not respect order.

#### Array

Works exactly the same as with structures.

### Quotient `/`

#### Numeric `num`

Quotient devides numeric values in the order from left to right. When all the messages have been arrived it sends result to outport. Then new iteration starts.

#### Map `map`

Finds pairs that are presented in all maps. Does not respect the order.

#### Struct `struct`

Works like Difference but does devision.

#### Array `arr`

Works exactly the same as with structures.

### Remainder `%`

```
<t num> (v[] t) (v t)
```

#### Numeric

Does [euclidean division](https://en.wikipedia.org/wiki/Euclidean_division) from left to right. When last number arrived it sends result to outport and starts new iteration.

### Bitwise AND `&`

define for int

### Bitwise OR `|`

define for int

### Bit clear (AND NOT) `&^`

define for int

### Bitwise XOR `^`

The only bitwise operator that also works for `bool`. TODO define for int and bool

## Conversion

### Complex to float `c2f`

```
(v c128) (r, i f64)
```

### To string `2str`

```
<t>(v t) (v str)
```

### Float to complex `f2c`

```
(r, i f64) (v c128)
```

## Rest

### Len `len`

```
<t str | vec | map> (v t) (v int)
```

### Del `del`

```
<t, y> (v map<t, y>, k t) (v y?)
```

### Append `append`

```
append<t> (vv vec<t>, v t) (v t)
```
