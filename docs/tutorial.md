# Tutorial

Welcome to a tour of the Nevalang programming language. This tutorial will introduce you to Nevalang through a series of guided examples.

1. [Welcome](#welcome)
   - [What kind of language is this?](#what-kind-of-language-is-this)
   - [Installation](#installation)
   - [Hello, World!](#hello-world)
   - [Compiling programs](#compiling-programs)
2. [Basic Concepts](#basic-concepts)
   - [Components, Ports, Nodes and Connections](#components-nodes-ports-and-connections)
   - [Messages and Basic Types](#messages-and-basic-types)
   - [Constants](#constants)
   - [Modules and Packages](#modules-and-packages)
   - [Imports and Visibility](#imports-and-visibility)
   <!-- 3. [Dataflow](#dataflow) -->

````shell
# will name output file "my_awesome_binary"
neva build my_awesome_project/src --output=my_awesome_binary

# will generate go code instead of machine code
neva build my_awesome_project/src --target=go

# will output my_awesome_wasm.wasm
neva build my_awesome_project/src --target=wasm --output=my_awesome_wasm
``` -->
<!-- Components part 2 (Native vs Normal) -->

## Welcome

### What Kind of Language is This?

Nevalang is a general-purpose programming language that uses dataflow instead of control flow, lacking variables and functions, and expressing programs through pure message passing with nodes, connections, and ports. It is implicitly parallel, meaning all nodes operate in parallel by default, eliminating the need for threads, coroutines, or async-await, which simplifies some tasks while complicating others.

Influenced by functional programming, Nevalang embraces immutability and higher-order components, disallowing data mutation and shared state. It is a compiled, strongly statically-typed, garbage-collected language, sharing Go's abstraction level but aligning more with Rust's strictness. Using Go as a backend, it supports all Go targets, including WASM and cross-compiled machine code.

### Installation

#### Requirements

Make sure you have [Go compiler](https://go.dev/dl/) installed.

#### Via Shell Script

For Mac OS and Linux:

```shell
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash
````

If your device is connected to a chinese network:

```shell
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/cnina/install.sh | bash
```

For Windows (see [issue](https://github.com/nevalang/neva/issues/499) with Windows Defender, try manual download from [releases](https://github.com/nevalang/neva/releases) if installation won't work):

```batch
curl -o installer.bat -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.bat && installer.bat
```

#### From Source

Here's how you can build Nevalang for all supported platforms

```
git clone github.com/nevalang/neva
cd neva
make build
```

After building is finished, pick the one for your architecture and put it in your `PATH`. The rest of the binaries can be removed.

If you don't have Make or you only want to build for your platform, open the `Makefile` in the root of the repository and check what the `build` command does. For example, to build for Mac OS, the instruction is `GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o neva-darwin-amd64 ./cmd/neva`

#### Testing

After installation is finished, you should be able to run the `neva` CLI from your terminal

```shell
neva version
```

It should emit something like `0.26.0`

### Hello, World!

Once you've installed the neva-cli, you are able to use the `new` command to scaffold new Nevalang projects

```shell
neva new my_awesome_project
```

Each new project contains a Hello World program, so we can just run it

```shell
neva run my_awesome_project/src
```

You should see the following output:

```shell
Hello, World!
```

If you open `my_awesome_project/src/main.neva` with your favorite IDE, you'll see this:

```neva
import { fmt }

def Main(start any) (stop any) {
	println fmt.Println<any>
	---
	:start -> { 'Hello, World!' -> println -> :stop }
}
```

Congratulations, you have just compiled and executed your first Nevalang program!

### Compiling Programs

As mentioned, `neva run` builds and runs the executable, then cleans up by removing the temporary binary. This is useful for development, but for production, we usually prefer separate compilation and execution. You can achieve this with the `neva build` command.

```neva
neva build my_awesome_project/src
```

This will produce an `output` file in the directory where neva-cli was executed, typically the project's root. Let's run our executable.

```shell
./output
```

Once again you should see `Hello, World!`.

## Basic Concepts

### Components, Nodes, Ports and Connections

Components are the basic building blocks in Nevalang. Let's look at the simplest possible Nevalang program:

```neva
def Main(start any) (stop any) {
    :start -> :stop
}
```

This program defines a `Main` component with:

- An input port `start` that accepts any type
- An output port `stop` that outputs any type
- A connection `->` that passes messages from `start` to `stop`

When this program runs:

1. Runtime sends a message to `start`
2. Message flows through the connection to `stop`
3. Program terminates

Most components do more interesting work by using nodes to process data:

```neva
import { fmt }

def Main(start any) (stop any) {
    println fmt.Println<any>
    ---
    :start -> println -> :stop
}
```

The `---` separator divides the component into two sections:

- Above: Node declarations (components used)
- Below: Network connections (data flow)

This program:

1. Creates a node `println` using `fmt.Println`
2. Sends the message from `start` to `println` to print it
3. Prints `{}` (start message) and terminates

### Messages and Basic Types

Back to hello world:

```neva
import { fmt }

def Main(start any) (stop any) {
    println fmt.Println<any>
    ---
    :start -> { 'Hello, World!' -> println -> :stop }
}
```

We sent `'Hello, World!'` to the `println` node. This is a string message literal, one of Nevalang's 4 basic types. Ignore the `-> { ... }` part for now; we'll cover it later.

```neva
// 1. `bool` - Boolean values: true or false
true -> println   // prints: true
false -> println  // prints: false

// 2. `int` - 64-bit signed integer numbers
42 -> println    // prints: 42
-100 -> println  // prints: -100

// 3. `float` - 64-bit floating-point numbers
3.14 -> println    // prints: 3.14
-0.5 -> println    // prints: -0.5

// 4. `string` - UTF-8 encoded text
'Hello!' -> println              // prints: Hello!
'Numbers: 123' -> println        // prints: Numbers: 123
'Special chars: @#$' -> println  // prints: Special chars: @#$
```

These primitive types are the basis for sending messages between nodes. We'll cover complex types later.

### Constants

Nevalang has no variables, only constants. Constants allow you to reuse values across your program. They must have explicit types and be known at compile-time. A constant's value cannot change during execution. Define constants using the `const` keyword:

```neva
const is_active bool = true
const age int = 25
const pi float = 3.14
const greeting string = 'Hello!'
```

Use `$` to prefix a constant in a network:

```neva
def Main(start any) (stop any) {
    println fmt.Println<any>
    ---
    :start -> { $greeting -> println -> :stop }
}
```

### Modules and Packages

Here's the structure of our Hello World project:

```
my_awesome_project/
├── src/
│   └── main.neva
└── neva.yaml
```

This structure introduces two fundamental concepts in Nevalang: modules and packages.

#### Modules

A module is a set of packages with a manifest file (`neva.yaml`). When we created our project with `neva new`, it generated a basic module with the following manifest file:

```yaml
neva: 0.26.0
```

This defines the Nevalang version for our project. As your project grows, you can include dependencies on third-party modules here.

#### Packages

A package is a directory with `.neva` files. In our Hello World example, the `src` package is our _main_ package, used as the compilation entry point with `neva run my_awesome_project/src` or `neva build my_awesome_project/src`. The main package must include a `Main` component, which serves as the program's entry point. Here's our Hello World program:

```neva
import { fmt }

def Main(start any) (stop any) {
   println fmt.Println<any>
   ---
   :start -> { 'Hello, World!' -> println -> :stop }
}
```

- Importing the `fmt` package from the standard library
- Defining the `Main` component in the entry package
- Using `fmt.Println` from the imported package

Let's add a `utils` package with helper components:

```
my_awesome_project/
├── src/
│   ├── main.neva
│   └── utils/
│       └── strings.neva
└── neva.yaml
```

Here's an example of a string utility and its usage in the main program:

```neva
// src/utils/strings.neva
pub def Greet(data string) (res string) { // <- new component
	('Hello, ' + :data) -> :res
}

// src/main.neva
import {
	fmt
	@/utils // <- new import
}

def Main(start any) (stop any) {
	greet utils.Greet // <- new node
	println fmt.Println<any>
	---
	:start -> { 'World' -> greet -> println -> :stop } // <- new connection
}
```

Notice how we can have multiple imports:

- `fmt` from the standard library for printing
- `@/utils` from our local module using the `@` symbol

This modular structure keeps your code organized and reusable as your projects grow.

### Imports and Visibility

In `utils`, we used `pub` keyword:

```neva
// src/utils/strings.neva
pub def Greet(data string) (res string) {
    ('Hello, ' + :data) -> :res
}
```

The `pub` keyword makes `Greet` public for imports. Without `pub`, `Greet` is private, causing compilation failure. This system encapsulates package details while defining the public API.

Let's show how components in the same package are used. Updated project structure:

```
my_awesome_project/
├── src/
│   ├── main.neva
│   ├── exclaim.neva
│   └── utils/
│       └── strings.neva
└── neva.yaml
```

Add a component in `exclaim.neva` to add exclamation marks to strings:

```neva
def AddExclamation(data string) (res string) {
    (:data + '!!!') -> :res
}
```

We can use `Greet` (import needed) and `AddExclamation` (no import needed) in our `src/main.neva`:

```neva
import {
    fmt
    @/utils
}

def Main(start any) (stop any) {
    greet utils.Greet
    exclaim AddExclamation  // No import needed - same package
    println fmt.Println<any>
    ---
    :start -> {
        'World' -> greet -> exclaim -> println -> :stop
    }
}
```

Output:

```
Hello, World!!!
```

<!-- ## Dataflow

Let's explore some common patterns for connecting components in Nevalang.

### Chained Connections

Remember our simple string utility in `src/utils/strings.neva`:

```neva
pub def Greet(data string) (res string) {
    ('Hello, ' + :data) -> :res
}
```

Since `Greet` has exactly one input port and one output port, we can use a "chained" connection syntax:

```neva
def Main(start any) (stop any) {
    greet utils.Greet
    println fmt.Println<any>
    ---
    :start -> { 'World' -> greet -> println -> :stop }
}
```

The `'World' -> greet -> println -> :stop` is shorthand for:

```neva
'World' -> greet:data
greet:res -> println:data
println:res -> :stop
```

However, if we modify `Greet` to have multiple ports, this won't work:

```neva
pub def Greet(data string, prefix string) (res string) {
    (:prefix + :data) -> :res
}
```

Now we must explicitly connect each port:

```neva
:start -> { 'Hi' -> greet:prefix }
'World' -> greet:data
greet -> println -> :stop
```

Since `greet` has one output port, we are allowed not to specify it, and `println` still has one input and one output port, we can still use the chained connection syntax for that part of the network.

### Port Usage Rules

There are two important rules about port usage in Nevalang:

1. Components must use ALL their own input and output ports in their network
2. When using nodes (instances of other components), all input ports must be connected, but at least one output port is sufficient

Let's look at examples of both rules:

#### Rule 1: Using Component's Own Ports

This won't compile because the component doesn't use its `:data` input port in the network:

```neva
pub def Greet(data string) (res string) {
    'Hello!' -> :res  // Error: input port 'data' is not used
}
```

This also won't compile because the component doesn't use its `:res` output port:

```neva
pub def Greet(data string) (res string) {
    // Error: output port 'res' is not used
    :data -> fmt.Println
}
```

The correct version uses both ports:

```neva
pub def Greet(data string) (res string) {
    ('Hello, ' + :data) -> :res
}
```

#### Rule 2: Using Node Ports

When using nodes (other components) in your network, the rules are different:

```neva
pub def Greet(data string) (res string, log string) {
    msg ('Hello, ' + :data)
    logger fmt.Println<string>
    ---
    msg -> [
        :res,
        { ('Greeted: ' + msg) -> logger }  // Only using logger's input port is fine
    ]
}
```

Here:
- We must connect to all of `logger`'s input ports (it has one)
- We don't need to use `logger`'s output port
- But we must use all of `Greet`'s own ports (`:data`, `:res`, `:log`)

This distinction between component's own ports and node ports is important for building modular programs.

### Port Name Omission

When a component has only one port on either side, you can omit the port name:

```neva
// Instead of foo:data -> bar:input
foo -> bar
```

This works even when chained syntax isn't available (like when a component has multiple ports on the other side).

### Fan-out and Fan-in

Sometimes you need to send the same message to multiple receivers (fan-out) or combine multiple senders into one receiver (fan-in).

Fan-out broadcasts a message to all receivers:

```neva
def Main(start any) (stop any) {
    p1 fmt.Println<any>
    p2 fmt.Println<any>
    ---
    :start -> { 'Hello!' -> [p1, p2] -> :stop }
}
```

This prints "Hello!" twice. The message is copied to both `p1` and `p2`.

Fan-in merges multiple senders:

```neva
def Main(start any) (stop any) {
    println fmt.Println<any>
    ---
    :start -> {
        ['Hello', 'World'] -> println -> :stop
    }
}
```

Messages are received in the order they were sent. If multiple messages arrive simultaneously, their order is non-deterministic.

You can't reuse senders or receivers in fan-in/fan-out. If you need to send the same message multiple times, you must use fan-out explicitly.

### Deferred Connections

Let's revisit our Hello World program:

```neva
def Main(start any) (stop any) {
    println fmt.Println<any>
    ---
    :start -> { 'Hello, World!' -> println -> :stop }
}
```

The `{ ... }` syntax creates a "deferred" connection. Since all input ports must be used, we need to wait for the `:start` signal before sending our message. The deferred connection ensures that 'Hello, World!' is only sent after receiving the start signal.

Without deferral, the program would be non-deterministic - the string might be sent multiple times before the program terminates. Deferred connections defer receiving rather than sending, ensuring proper synchronization.

This section introduces key dataflow patterns while maintaining the tutorial's focus on practical examples. It builds on previous concepts and prepares readers for more advanced topics in the book. -->
