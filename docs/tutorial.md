# Tutorial

Welcome to a tour of the Nevalang programming language. This tutorial will introduce you to Nevalang through a series of guided examples.

1. [Welcome](#welcome)
   - [What kind of language is this?](#what-kind-of-language-is-this)
   - [Installation](#installation)
   - [Hello, World!](#hello-world)
   - [Compiling programs](#compiling-programs)
2. [Core Concepts](#core-concepts)
   - [Components](#components)
   - [Messages and Basic Types](#messages-and-basic-types)
   - [Constants](#constants)
   - [Modules and Packages](#modules-and-packages)
   - [Imports and Visibility](#imports-and-visibility)
3. [Dataflow Basics](#dataflow)
   - [Chained Connections](#chained-connections)
   - [Multiple Ports](#multiple-ports)
   - [Fan-In/Fan-Out](#fan-infan-out)
   - [Binary Operators](#binary-operators)
   - [Ternary Operator](#ternary-operator)
   - [Switch Expression](#switch-expression)
   - [Deferred Connections](#deferred-connections)

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
```

If your device is connected to a chinese network:

```shell
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/china/install.sh | bash
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

It should emit something like `0.27.1`

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
	println fmt.Println<string>
	---
	:start -> 'Hello, World!' -> println -> :stop
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

## Core Concepts

### Components

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
    println fmt.Println<string>
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
    println fmt.Println<string>
    ---
    :start -> 'Hello, World!' -> println -> :stop
}
```

We sent `'Hello, World!'` to the `println` node. This is a string message literal, one of Nevalang's 4 basic types.

```neva
// `bool` - Boolean values: true or false
true -> println   // prints: true
false -> println  // prints: false

// `int` - 64-bit signed integer numbers
42 -> println    // prints: 42
-100 -> println  // prints: -100

// `float` - 64-bit floating-point numbers
3.14 -> println    // prints: 3.14
-0.5 -> println    // prints: -0.5

// `string` - UTF-8 encoded text
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
const greeting string = 'Hello!'

def Main(start any) (stop any) {
    println fmt.Println<string>
    ---
    :start -> $greeting -> println -> :stop
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
neva: 0.27.1
```

This defines the Nevalang version for our project. As your project grows, you can include dependencies on third-party modules here.

#### Packages

A package is a directory with `.neva` files. In our Hello World example, the `src` package is our _main_ package, used as the compilation entry point with `neva run my_awesome_project/src` or `neva build my_awesome_project/src`. The main package must include a `Main` component, which serves as the program's entry point. Here's our Hello World program:

```neva
import { fmt }

def Main(start any) (stop any) {
   println fmt.Println<string>
   ---
   :start -> 'Hello, World!' -> println -> :stop
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
│       └── utils.neva
└── neva.yaml
```

Let's add a string utility `Greet` that receives a `data` string, prefixes it with `"Hello, "` and sends to `res` output port. We'll use binary expression `('Hello, ' + :data)` to concatenate strings:

```neva
// src/utils/utils.neva
pub def Greet(data string) (res string) { // new component
	('Hello, ' + :data) -> :res
}

// src/main.neva
import {
	fmt
	@:utils // new import
}

def Main(start any) (stop any) {
	greet utils.Greet // new node
	println fmt.Println<string>
	---
	:start -> 'World' -> greet -> println -> :stop // new connection
}
```

Notice how we can have multiple imports:

- `fmt` from the standard library for printing
- `@:utils` from our local module (`@` is module name, `:` separates module/package)

This modular structure keeps your code organized and reusable as your projects grow.

### Imports and Visibility

In `utils`, we used `pub` keyword:

```neva
// src/utils/utils.neva
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
│       └── utils.neva
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
    @:utils
}

def Main(start any) (stop any) {
    greet utils.Greet
    exclaim AddExclamation  // same package, no import needed
    println fmt.Println<string>
    ---
    :start -> 'World' -> greet -> exclaim -> println -> :stop
}
```

Output:

```
Hello, World!!!
```

## Dataflow

### Chained Connections

Nodes send and receive messages through ports. Each port is referenced with a `:` prefix followed by its name:

```neva
def Main(start any) (stop any) {
    :start -> :stop
}
```

We refer to input ports as "inports" and output ports as "outports". In this example, we connect the `start` inport with the `stop` outport. This single inport/outport pattern is also seen in `utils.Greet`:

```neva
// src/utils/utils.neva

pub def Greet(data string) (res string) {
    ('Hello, ' + :data) -> :res
}
```

Same true for `fmt.Println`:

```neva
// fmt package
pub def Println<T>(data T) (sig struct{})
```

This allowed us to chain nodes together:

```neva
:start -> 'World' -> greet -> println -> :stop
```

When chaining nodes, we actually reference their ports implicitly. The chain could be written more verbosely as:

```neva
:start -> 'World' -> greet:data
greet:res -> println:data
println:sig -> :stop
```

Both versions are equivalent, but the chained syntax is preferred for readability.

### Multiple Ports

Let's look at components with multiple ports. A component can have any number of inports and outports. Here's another component we can add to `src/utils/utils.neva`:

```neva
pub def Concat(prefix string, suffix string) (res string) {
    (:prefix + :suffix) -> :res
}
```

Components must use all their ports within their network. For example, if we remove `:suffix`, the program won't compile:

```neva
def Concat(prefix string, suffix string) (res string) {
    :prefix -> :res // ERROR: suffix inport is not used
}
```

When using nodes with multiple inports, we can't use chain syntax because compiler won't know which port to use. Instead, we must specify ports explicitly:

```neva
import {
    fmt
    @:utils
}

def Main(start any) (stop any) {
    concat utils.Concat
    println fmt.Println<string>
    ---
    :start -> 'Hello, ' -> concat:prefix
    'World' -> concat:suffix
    concat -> println -> :stop
}
```

Notice that:

1. We can omit `concat:res ->` and write just `concat ->` since `Concat` has one outport
2. We can chain `-> println ->` since it still has a single port

Let's add a `debug` outport to `Concat`:

```neva
pub def Concat(prefix string, suffix string) (res string, debug string) {
    (:prefix + :suffix) -> :res
    'Debug: concatenating strings' -> :debug
}
```

Unlike self outports, we can ignore node outports we don't need (like `concat:debug`), but we must now specify `concat:res` explicitly since `concat` has multiple outports:

```neva
import {
    fmt
    @:utils
}

def Main(start any) (stop any) {
    concat utils.Concat
    println fmt.Println<string>
    ---
    :start -> 'Hello, ' -> concat:prefix
    'World' -> concat:suffix
    concat:res -> println -> :stop // concat:debug is not used
}
```

### Fan-In/Fan-Out

Sometimes we need to handle multiple senders or receivers. While we've primarily used one-to-one connections (pipelines), Nevalang also supports many-to-one (fan-in) and one-to-many (fan-out) connections.

#### Fan-In

Fan-in allows multiple senders to connect to a single receiver using square brackets on the sender side. The receiver processes messages in FIFO (first in, first out) order, based on when senders emit their messages.

Let's explore this using `strconv.ParseNum` from the standard library, which converts strings to numbers:

```neva
// strconv package
pub def ParseNum<T int | float>(data string) (res T, err error)
```

Note that it has an `err` outport of type `error`. While we can usually ignore node outports as long as we use at least one, the `err` port is special - we must always handle potential errors.

Let's try converting `'42'` to `42` and print both the result and any potential errors:

```neva
import {
    fmt
    strconv
}

def Main(start any) (stop any) {
    parse strconv.ParseNum<int>
    println fmt.Println<any>
    ---
    :start -> '42' -> parse
    [parse:res, parse:err] -> println -> :stop // fan-in
}
```

The fan-in connection `[parse:res, parse:err] -> println` connects both outports to the `println` receiver. Running this produces:

```
42
```

The output is consistent because `'42'` is a valid integer string. In this case, only `parse:res` sends a message while `parse:err` remains silent. With error ports, only one will ever fire - either the success result or the error.

If we try an invalid number:

```neva
:start -> 'forty two' -> parse
[parse:res, parse:err] -> println -> :stop
```

We'll see:

```
parsing "forty two": invalid syntax
```

Now only the `parse:err` port fires, demonstrating the exclusive nature of success and error outputs.

#### Fan-Out

Fan-out allows a single sender to connect with multiple receivers using square brackets on the receiver side. Each receiver gets an identical message. The sender blocks until all receivers process the message, meaning the connection speed is limited by the slowest receiver.

Let's add a component to `src/utils/utils.neva` that receives two strings, parses them as integers, and returns their sum as a result if successful, or an error otherwise:

```neva
import { strconv }

// ...existing code...

pub def AddIntStrings(left string, right string) (res int, err error) {
    parse_left strconv.ParseNum<int>
    parse_right strconv.ParseNum<int>
    ---
    :left -> parse_left
    :right -> parse_right

    [parse_left:err, parse_right:err] -> :err // fan-in with error propagation

    parse:res -> :res
}
```

Key points:

1. We create two instances of `strconv.ParseNum` to process both connections - `parse_left` for `:left` and `parse_right` for `:right`
2. We use fan-in to connect both error outputs to our `:err` port for error propagation

Now let's update `src/main.neva` to add `'21'` to itself using our new `utils.AddIntStrings`. We'll need to send `'21'` to both input ports simultaneously using fan-out:

```neva
import {
    fmt
    @:utils
}

def Main(start any) (stop any) {
    add utils.AddIntStrings
    ---
    :start -> '21' -> [add:left, add:right] // fan-out
    [add:res, add:err] -> println -> :stop
}
```

Note that we can use both fan-out (`'21' -> [...]`) and fan-in (`[add:res, add:err] -> ...`) in the same network. Running this program outputs:

```
42
```

This works because `'21'` is a valid integer string and `21 + 21` equals `42`. If we try an invalid input:

```neva
:start -> 'twenty one' -> [add:left, add:right]
```

We get:

```
parsing "twenty one": invalid syntax
```

The error from `strconv.ParseNum` propagates through `utils.AddIntStrings` up to `Main`, demonstrating proper error handling.

### Binary Operators

We've already seen binary expressions e.g. in `utils.Greet`:

```neva
// src/utils/utils.neva
pub def Greet(data string) (res string) {
	('Hello, ' + :data) -> :res
}
```

A binary expression consists of a left operand, operator, and right operand. In this case, `'Hello, '` is the left operand and `:data` is the right operand. When both operands are ready, the operator transforms the data and sends the result forward.

Both operands must share the same type, and the operator must support that type. For example, string concatenation works (`'Hello, ' + 'World'`), but adding an integer to a string (`21 + '21'`) will fail compilation. Nevalang is strongly typed with no implicit conversions.

Operands can be any senders e.g. ports, constants, even other binary expressions. Let's add one more utility component to our `src/utils/utils.neva`:

```neva
pub def TriangleArea(b int, h int) (res float) {
    ((:b * :h) / 2) -> :res
}
```

Here, `(:b * :h)` is a binary expression used as the left operand of another binary expression. The calculation proceeds once both `:b` and `:h` are ready.

Nevalang supports these binary operators:

```neva
// arithmetic
(5 + 3) -> println // addition: 8
(5 - 3) -> println // subtraction: 2
(5 * 3) -> println // multiplication: 15
(6 / 2) -> println // division: 3
(7 % 3) -> println // modulo: 1
(2 ** 3) -> println // power: 8

// comparison
(5 == 5) -> println // equal: true
(5 != 3) -> println // not equal: true
(5 > 3) -> println // greater than: true
(5 < 8) -> println // less than: true
(5 >= 5) -> println // greater or equal: true
(5 <= 8) -> println // less or equal: true

// logic
(true && true) -> println // AND: true
(true || false) -> println // OR: true

// bitwise
(5 & 3) -> println // AND: 1
(5 | 3) -> println // OR: 7
(5 ^ 3) -> println // XOR: 6
```

### Ternary Operator

Ternary operator allows to select between two sources based on a condition, using the syntax `(if ? then : else) -> receiver`. Operator waits for all operands, selects the message and sends it downstream. Let's add one more component to `src/utils/utils.neva`:

```neva
pub def FormatBool(data bool) (res string) {
    (:data ? 'true' : 'false') -> :res
}
```

All three operands can be any valid senders (ports, literals, constants, expressions, etc.), as long as "if" sends a `bool` and "then/else" are compatible with the receiver. Let's update our `src/main.neva` and use our new utility components to see how a more complex ternary operator looks:

```neva
import {
    fmt
    @:utils
}

def Main(start any) (stop any) {
    println fmt.Println
    area utils.TriangleArea
    format utils.FormatBool
    ---
    :start -> [
        20 -> area:b,
        10 -> area:h
    ]
    (area > 50) -> format
    ((format == 'true') ? 'Big' : 'Small') -> println -> :stop
}
```

This example calculates a triangle's area (base=20, height=10), checks if it's larger than 50, and prints either "Big" or "Small" accordingly. While contrived, it demonstrates how the ternary operator can be used in more complex scenarios.

### Switch Expression

So far we've learned how to _select_ sources based on conditions, but the message's _route_ was always the same. For example, in `utils.FormatBool` we selected either `'true'` or `'false'` but the destination was always `:res` - `(:data ? 'true' : 'false') -> :res`. To write real programs we need to be able to select both sources and destinations. In other words, we need "routers" in addition to "selectors", and `switch` is one of them. It has the following syntax:

```neva
condition_sender -> switch {
    case_sender_1 -> case_receiver_1
    case_sender_2 -> case_receiver_2
    ...
    _ -> default_receiver
}
```

Switch consists of a condition sender and pairs of case senders/receivers, including a required default case with `_`. It waits for all senders, compares the condition message with case messages for equality, and executes the first matching branch. Once triggered, other branches won't fire until the next message arrives. Let's add another component to `src/utils/utils.neva` to see `switch` in action:

```neva
pub def HttpStatus(code int) (
    continue any,
    ok any,
    multiple_choices any,
    bad_request any,
    internal_err any,
    unknown any
) {
    :code -> switch {
        100 -> :continue
        200 -> :ok
        300 -> :multiple_choices
        400 -> :bad_request
        500 -> :internal_err
        _ -> :unknown
    }
}
```

This component receives an integer `code` and based on its value sends to one of its 6 outports.

<!-- TODO make http.get and compare status -->
<!-- TODO: fan-in, fan-out inside switch and get rid of utility -->
<!-- TODO: explain how switch works -->

#### If/Else

So far switch branches were chosen by comparing messages, but sometimes we just need to branch on a boolean value. Since Nevalang has no if-else construct, this pattern is implemented using `switch`:

```neva
pub def ClassifyInt(data int) (neg any, pos any) {
    :start -> (data >= 0) -> switch {
        true -> :pos
        _ -> :neg
    }
}
```

Things to notice:

1. We send `bool` messages to both outports `neg` and `pos`, which works since they're typed as `any`
2. We use `_` as default case rather than `(:data < 0)` since switch requires a default and negative is the only other option
3. This works naturally with switch's behavior - we compare the `bool` condition with `true`, and `_` handles `false`

#### Switch True

...

<!-- TODO -->
