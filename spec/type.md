# Types

All the types described in this secion are _builtins_ in a sence that they're available in a global scope from any package. No need (and it's impossible) to import them.

## Base

These are types that don't have definition in source code, compiler is awared of them. 

### Scalars

- boolean: `bool`
- complex numbers: `c64, c128`
- floating point numbers: `f32, f64`
- unsigned integers: `u8, u16, u32, u64`
- signed integers: `i8, i16, i32, i64`
- string: `str`
- enum: `{ V1 V2 ... Vn }`

### Composite

- array: `[N]T `
- vector: `[]T`
- record: `{ k1 T1, k2 T2, ... kn Tn }`
- map: `map[T1]T2`
- union: `T1 | T2 ... | Tn`
- optional: `T?`

## Non-base

These are types that user could express himself but for convinience they're represented in a standard library.

- complex: `c64 | c128`
- float: `float32 | float64`
- uint: `u8 | u16 | u32 | u64`
- int: `i8 | i16 | i32 | i64`
- num: `complex | float | uint | int`
- error: `{ parent error?, trace TODO }`
- subStreamItem: `<T>{isOpenBracket: bool, value: T?}`