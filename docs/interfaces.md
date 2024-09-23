# Interfaces

Interfaces define the input-output of a component and may have type-parameters. They consist of input and output port definitions. It's recommended to prefix interface names with `I`:

```neva
interface IMain(start) (stop)
```

This is the simplest interface: one input and one output port. Types default to `any`. It describes a component that receives any message at `start` and sends any message from `stop`. The `Main` component implements this interface.

Interfaces can have multiple input and output ports, but it's recommended to limit the number of ports, especially input ports. More ports and specific data-types make implementation harder. Here's an example of an interface with 3 input and 3 output ports, all of type `any`:

```neva
interface IXXX(a, b, c) (d, e, f)
```

## Component IO and Dependency Injection

Interfaces in Nevalang are used for:

1. Defining component signatures - specifying IO (port names and data-types). Semantics should be described with naming and comments if necessary.
2. Defining dependencies (interface nodes) that parent components must provide, enabling dependency injection. This ensures a maintainable and testable codebase, allowing for easy replacement of implementations, which is crucial for polymorphic code and testing mocks.

## Type Parameters and Port Types

`any` is too generic; specific types are often needed. Let's define an interface for an integer adder:

```neva
interface IAdd(acc int, el int) (res int)
```

What if we want to add not just integers but also support `float` and `string`? You could use `union` for that

```neva
type addable int | float | string
interface IAdd(acc addable, el addable) (res addable)
```

This solution is problematic because it allows mixing types, e.g., `int` for `:acc` and `float` for `:el`. The `:res` message will have a union type `int | float | string`, making the code complex. To maintain type-safety, use type-parameters in the interface definition:

```neva
interface IAdd<T>(acc T, el T) (res T)
```

This ensures `:acc` and `:el` receive compatible types, and `:res` matches their type. However, our current definition allows any type, including `IAdd<bool, bool>`, which we don't want. To fix this, we can explicitly constrain `T` using our union:

```neva
interface IAdd<T int | float | string>(acc T, el T) (res T)
```

Now only `IAdd<int, int>`, `IAdd<float, float>` or `IAdd<string, string>` (and their compatible variants) are possible. `IAdd:res` will always be `int`, `float`, or `string`.

Type-expressions in interface definitions follow type-system rules, so you can pass `T` to other type-expressions. Example:

```
interface IAppend<T>(lst list<T>, el T) (list<T>)
```

These expressions can be complex and nested, but we'll keep this section brief.

## Array-Ports

Let's return to the original `IAdd` without type-parameters for simplicity

```neva
interface IAdd(acc int, el int) (res int)
```

It's suitable for combining 2 sources, but what if we need to combine any number of sources? Chaining multiple `IAdd` instances is tedious and sometimes impossible. Let's look at the array-inports solution:

```neva
interface IAdd([el] int) (res int)
```

In a component using `IAdd` in its network, we can do this:

```neva
add IAdd
---
1 -> add[0]
2 -> add[1]
3 -> add[2]
```

Syntax `add[i]` is shorthand for `add:el[i]`. The compiler infers the port name since there's only one.

Another example of a component that benefits from array-ports is `Switch`. It's used for routing - imagine we have message `data` and need to route it to different destinations based on its value. For example, if it's `a` we send it to the first destination, if `b` to the second, and `c` to the third. Otherwise, we send it to a default destination to handle unknown values. An adhoc solution with a fixed number of ports wouldn't scale. We need a generic component with dynamic port support. Here's the Switch signature:

```neva
pub flow Switch<T>(data T, [case] T) ([case] T, else T)
```

This allows code like this

```neva
:data -> switch:data
1 -> switch:case[0] -> dst1
2 -> switch:case[1] -> dst2
3 -> switch:case[2] -> dst3
switch:else -> ('!?' -> println)
```

Array-ports combine data from different sources. They are static, requiring the number of ports to be known at compile time for channel generation. This can always be determined from the source code. Limitations of array-ports will be discussed on the network page.
