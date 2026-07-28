# Directives

Compiler directives are special instructions for the compiler, not intended for daily use but important for understanding language features.

## `#extern`

Tells compiler a component lacks source code implementation and requires a runtime function call. Example:

```neva
#extern(println)
pub def Println<T>(data T) (res T, err error)
```

### Overloading

Overloading uses multiple declarations with the same component name. Each
declaration has one ordinary `#extern` directive; the compiler selects the
compatible declaration.

```neva
#extern(int_add)
pub def Add(left int, right int) (res int)
#extern(float_add)
pub def Add(left float, right float) (res float)
#extern(string_add)
pub def Add(left string, right string) (res string)
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
	[println:res, println:err] -> :stop
}
```

## `#autoports`

Derives component inports from its type-argument structure fields, rather than defining them in source code. Example:

```neva
#autoports
#extern(struct_builder)
pub def Struct<T struct {}> () (msg T)
```
