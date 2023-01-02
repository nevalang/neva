# Types

## Fundamental

### Scalars

- boolean: bool
- complex numbers: c64, c128
- floating point numbers: f32, f64
- unsigned integers: u8, u16, u32, u64
- signed integers: i8, i16, i32, i64
- string: str
- enum: { V1 V2 ... Vn } (u8)

### Composite

- array: [N]T 
- vector: []T
- record: { k1 T1, k2 T2, ... kn Tn }
- map: map[T1]T2
- union: T1 | T2 ... | Tn
- optional: T?

## Non-fundamental  

- complex: c64 | c128
- float: float32 | float64
- uint: u8 | u16 | u32 | u64
- int: i8 | i16 | i32 | i64
- num: complex | float | uint | int
- error: { parent error?, trace TODO }
- substreamItem: <T>{isOpenBracket: bool, value: T?}