# Components

Components are blueprints for nodes, which are computational units in Nevalang. Components can have multiple instances, each potentially differing in name, type-arguments, and dependencies. Unlike other languages with various computational constructs, Nevalang uses only components for all computations, including data transformations and side-effects. There are 2 types of components: native and normal.

## Native Components

Native components have only a signature (interface) and no implementation in Nevalang. They use the `#extern` directive to indicate implementation in the runtime's language. There's no way to tell if a component is native by its usage, only by its implementation. For example, the `And` boolean component uses a runtime function `and`. Native components exist only in the standard library and are not user-definable.

```neva
#extern(and)
flow And(a bool, b bool) (res bool)
```

> Native components may enable future Go-interop, allowing Go functions to be called from Neva using #extern.

### Overloading

Native components can be overloaded, allowing multiple implementations with the same signature for different data types. The compiler chooses the appropriate implementation based on the given data type. Overloading is limited to native components in the standard library and not available for user-defined components.

Overloaded native components use a modified extern directive: `#extern(t1 f1, t2 f2, ...)`. These components must have exactly one type parameter with a union constraint. Example:

```neva
#extern(int int_add, float float_add, string string_add)
pub flow Add<T int | float | string>(acc T, el T) (res T)
```

## Normal Components

Normal components are implemented in Nevalang source code. As a Nevalang programmer, you'll primarily work with these components, which are also found in the standard library alongside native ones. Normal components don't use the `#extern` directive and include an implementation consisting of a required network and optional nodes section. The network must use all of the component's inports and outports, enabling at least basic routing.

Minimal nevalang program is one normal component `Main` without any nodes and with network of a single connection:

```neva
flow Main(start) (stop) {
   :start -> :stop
}
```

### Nodes

Components however typically perform data transformations or side-effects using nodes, which are instances of other components. This means normal components usually depend on other components, directly or indirectly (through interfaces).

Normal component `Main` with a `println` node (instance of `Println`):

```neva
flow Main(start) (stop) {
   Println
   ---
   :start -> (42 -> println -> :stop)
}
```

As you can see, we refer to the instance of `Println` as `println`. The compiler implicitly assigns a lowercase version of the component name to its instance. However, with multiple instances of the same component, this leads to name collisions. To avoid this, the compiler requires explicit naming for nodes in such cases. Example:

```neva
flow Main(start) (stop) {
   p1 Println
   p2 Println
   ---
   :start -> (42 -> p1 -> p2 -> :stop)
}
```

#### IO (Implicit) Nodes

Normal components actually have implicit `in` and `out` nodes. Even in our `Main` example with a single connection `:start -> :stop`, the compiler interprets this as `in:start -> out:stop`. Interestingly, `in` only has outports (you can only send _from_ it) while `out` only has inports (you can only send _to_ it). The compiler automatically generates these in/out nodes with necessary ports based on the component's interface.

#### Interface Nodes (Dependency Injection)

Consider an application that performs some business logic with logging:

```neva
flow App(data) (sig) {
   Logic, Log
   ---
   :data -> logic -> log -> :sig
}
```

Now imagine we want to replace `Logger` with another component based on a condition. Let's say we want to use real logger in production and a mock in testing. Without dependency injection (DI), we'd have to extend the interface with a `flag bool` inport and check it.

```neva
flow App(data any, prod bool) (sig any) {
   Cond, Logic, ProdLogger, MockLogger
   ---
   :data -> businessLogic -> cond:data
   :prod -> cond:if
   cond:then -> prodLogger
   cond:else -> mockLogger
   [prodLogger, mockLogger] -> :sig
}
```

This not only makes the code more complex but also means we have to initialize both implementations: `ProdLogger` in the test environment and `MockLogger` in the production environment, even though they are not needed in those respective contexts. What if you need to read environment variables to initialize a component? For example, your logger might need to send requests to a third-party service to collect errors. And finally, imagine if it were not a boolean flag but an enum with several possible states. The complexity would increase dramatically.

> As you can see it's possible to write nodes in a single line, separated by comma: `Cond, Logic, Println, Mock`. Don't abuse this style - Nevalang is not about clever one-liners.

Let's implement this using dependency injection. First, define an interface:

```neva
type ILog(data) (sig)
```

Next, define the dependency (interface node):

```neva
flow App(data) (sig) {
   Logic, ILog
   ---
   :data -> logic -> iLog -> :sig
}
```

Our network is again a simple pipeline. Finally, provide the dependency.

Before dependency injection:

```neva
// cmd/app
flow Main(start) (stop)
   App
   ...
}
```

After dependency injection:

```neva
flow Main(start) (stop)
   App{ProdLogger}
   ...
}

flow Test(start) (stop)
   App{MockLogger}
   ...
}
```

`App{ProdLogger}` syntax sugar for `App{iLog: MockLogger}`, same for `App{MockLogger}`. Compiler is able to infer name of the dependency we provide if there's only one dependency. Syntax for providing several dependencies looks like a structure or dictionary initialization: `Component{dep1: nodeExpr1, dep2: nodeExpr2, ..., depN: nodeExprN}`.

**Component and Interface Compatibility**

Component `C1` implements interface `I1` if:

- Type parameters are compatible: same amount and each parameter of `C1` is compatible with corresponding `I1` parameter
- Inports are compatible: full match by name plus each inport of `C1` is compatible with corresponding `I1` inport
- Outports are compatible
  - `C1` has full match by name or superset of `I1` outports. Example: `C1(a) (b, c)` is compatible with `I1(a) (b)`, but not with `I1(a) (b, d)`. Types must be compatible too.

### Main Component

Main component is the entry point of a nevalang program. A package containing this component is called a "main package" and serves as the compilation entry point. Each main package must have exactly one non-public `Main` component with no interface nodes, implementing the `(start any) (stop any)` interface. Nevalang's runtime sends a message to `Main:start` at startup and waits for a message from `Main:stop`. Upon receiving the stop signal, it terminates the program.

## Type Parameters

Components contain interfaces, which may have type parameters. Type arguments must be provided during initialization. For interfaces, this occurs when initializing an interface node, while for components, it happens when initializing concrete nodes.

Component uses type-parameters for its IO:

```neva
flow Foo<T>(data T) (sig any)
```

Components can pass type parameters from their interface to node expressions:

```neva
flow Bar<T>(data T) (sig any) {
   Println<T>
   ---
   :data -> println -> :sig
}
```

This means when we initialize `Bar` with a type argument, it replaces `T` in `Println<T>`. For example, if a parent component initializes `Bar<int>`, then inside Bar, `Println<T>` becomes `Println<int>`.
