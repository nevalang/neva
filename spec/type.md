# Types

```
complex = c64 | c128
float = float32 | float64
uint = u8 | u16 | u32 | u64
int = i8 | i16 | i32 | i64
num = complex | float | uint | int

error = {
    parent ?error
    msg str
}
```