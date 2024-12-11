# Directives

Compiler directives are special instructions for the compiler, not intended for daily use but important for understanding language features.

## `#extern`

Tells compiler a component lacks source code implementation and requires a runtime function call. Example:

```neva
#extern(println)
pub def Println<T>(data T) (sig T)
```

### Overloading

Native components can be overloaded using `#extern(t1 f1, t2 f2, ...)`. These components must have one type parameter with a union constraint. The compiler selects the appropriate implementation based on the data type. For instance:

```neva
#extern(int int_add, float float_add, string string_add)
pub def Add<T int | float | string>(left T, right T) (res T)
```

## `#bind`

Instructs compiler to insert a given message into a runtime function call for nodes with `extern` components. Example (desugared hello world):

```neva
const greeting string = 'Hello, World!'

def Main(start any) (stop any) {
	#bind(greeting)
	greeting New<string>
	println Println<string>
	lock Lock<string>
	---
	:start -> lock:sig
	greeting:res -> lock:data
	lock:data -> println:data
	println:res -> :stop
}
```

## `#autoports`

Derives component inports from its type-argument structure fields, rather than defining them in source code. Example:

```neva
#autoports
#extern(struct_builder)
pub def Struct<T struct {}> () (msg T)
```
