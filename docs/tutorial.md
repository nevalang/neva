# Tutorial

Welcome to a tour of the Neva programming language. This tutorial will introduce you to Nevalang through a series of guided examples.

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
   - [Switch](#switch)

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

It should emit something like `0.33.0`

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
	:start -> 'Hello, World!' -> println
	[println:res, println:err] -> :stop
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

> Execute `neva build --help` to learn more - how to compile to Go, WASM or how to do cross-compilation e.g. compile linux binaries in windows.

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
    :start -> println
    [println:res, println:err] -> :stop
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
    :start -> 'Hello, World!' -> println
    [println:res, println:err] -> :stop
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
    :start -> $greeting -> println
    [println:res, println:err] -> :stop
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
neva: 0.33.0
```

This defines the Nevalang version for our project. As your project grows, you can include dependencies on third-party modules here.

#### Packages

A package is a directory with `.neva` files. In our Hello World example, the `src` package is our _main_ package, used as the compilation entry point with `neva run my_awesome_project/src` or `neva build my_awesome_project/src`. The main package must include a `Main` component, which serves as the program's entry point. Here's our Hello World program:

```neva
import { fmt }

def Main(start any) (stop any) {
   println fmt.Println<string>
   ---
   :start -> 'Hello, World!' -> println
   [println:res, println:err] -> :stop
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
	@:src/utils // new import
}

def Main(start any) (stop any) {
	greet utils.Greet // new node
	println fmt.Println<string>
	---
	:start -> 'World' -> greet -> println
	[println:res, println:err] -> :stop // new connection
}
```

Notice how we can have multiple imports:

- `fmt` from the standard library for printing
- `@:src/utils` from our local module (`@` is module name, `:` separates module/package)

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
    @:src/utils
}

def Main(start any) (stop any) {
    greet utils.Greet
    exclaim AddExclamation  // same package, no import needed
    println fmt.Println<string>
    ---
    :start -> 'World' -> greet -> exclaim -> println
    [println:res, println:err] -> :stop
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

This allowed us to chain nodes together:

```neva
:start -> 'World' -> greet -> println
```

When chaining nodes, we actually reference their ports implicitly. The chain could be written more verbosely as:

```neva
:start -> 'World' -> greet:data
greet:res -> println:data
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
    @:src/utils
}

def Main(start any) (stop any) {
    concat utils.Concat
    println fmt.Println<string>
    ---
    :start -> 'Hello, ' -> concat:prefix
    'World' -> concat:suffix
    concat -> println
    [println:res, println:err] -> :stop
}
```

Notice that we can omit `concat:res ->` and write just `concat ->` since `Concat` has one outport.

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
    @:src/utils
}

def Main(start any) (stop any) {
    concat utils.Concat
    println fmt.Println<string>
    ---
    :start -> 'Hello, ' -> concat:prefix
    'World' -> concat:suffix
    concat:res -> println
    [println:res, println:err] -> :stop // concat:debug is not used
}
```

### Fan-In/Fan-Out

Sometimes we need to handle multiple senders or receivers. While we've primarily used one-to-one connections (pipelines), Nevalang also supports many-to-one (fan-in) and one-to-many (fan-out) connections.

#### Fan-In

Fan-in allows multiple senders to connect to a single receiver using square brackets on the sender side. The receiver processes messages in FIFO (first in, first out) order, based on when senders emit their messages.

Let's explore this using `strconv.Atoi` from the standard library, which converts strings to integer numbers:

```neva
// strconv package
pub def Atoi(data string) (res int, err error)
```

Note that it has an `err` outport of type `error`. While we can usually ignore node outports as long as we use at least one, the `err` port is special - we must always handle potential errors. You may have already seen the fan-in pattern with `Println` earlier, since it also may return an error. But let's take a look into another example.

Let's try converting `'42'` to `42` and print both the result and any potential errors:

```neva
import {
    fmt
    strconv
}

def Main(start any) (stop any) {
    parse strconv.Atoi
    println fmt.Println<any>
    ---
    :start -> '42' -> parse
    [parse:res, parse:err] -> println
    [println:res, println:err] -> :stop // fan-in
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
[parse:res, parse:err] -> println
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
    parse_left strconv.Atoi
    parse_right strconv.Atoi
    ---
    :left -> parse_left
    :right -> parse_right
    [parse_left:err, parse_right:err] -> :err // fan-in with error propagation
    (parse_left:res + parse_right:res) -> :res
}
```

Key points:

1. We create two instances of `strconv.Atoi` to process both connections - `parse_left` for `:left` and `parse_right` for `:right`
2. We use fan-in to connect both error outputs to our `:err` port for error propagation

Now let's update `src/main.neva` to add `'21'` to itself using our new `utils.AddIntStrings`. We'll need to send `'21'` to both input ports simultaneously using fan-out:

```neva
import {
    fmt
    @:src/utils
}

def Main(start any) (stop any) {
    add utils.AddIntStrings
    println fmt.Println<any>
    ---
    :start -> '21' -> [add:left, add:right] // chain + fan-out
    [add:res, add:err] -> println //fan-in
    [println:res, println:err] -> :stop // fan-in
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

The error from `strconv.Atoi` propagates through `utils.AddIntStrings` up to `Main`, demonstrating proper error handling.

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
pub def TriangleArea(b int, h int) (res int) {
    ((:b * :h) / 2) -> :res
}
```

Here, `(:b * :h)` is a binary expression used as the left operand of another binary expression. The calculation proceeds once both `:b` and `:h` are ready. We can test it by changing content in `src/main.neva`

```neva
import {
    fmt
    @:src/utils
}

def Main(start any) (stop any) {
    area utils.TriangleArea
    println fmt.Println<any>
    ---
    :start -> [
        10 -> area:b,
        20 -> area:h
    ]
    area -> println
    [println:res, println:err] -> :stop
}
```

Outputs:

```neva
100
```

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
    @:src/utils
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
    ((format == 'true') ? 'Big' : 'Small') -> println
    [println:res, println:err] -> :stop
}
```

Outputs:

```
Big
```

This example calculates a triangle's area (base=20, height=10), checks if it's larger than 50, and prints either "Big" or "Small" accordingly. While contrived, it demonstrates how the ternary operator can be used in more complex scenarios.

### Switch

So far we've learned how to _select_ sources based on conditions, but the message's _route_ was always the same. For example, in `utils.FormatBool` we selected either `'true'` or `'false'` but the destination was always `:res`:

```neva
(:data ? 'true' : 'false') -> :res
```

To write real programs we need to be able to configure both sources and destinations. In other words, we need "routers" in addition to "selectors", and `switch` is one of them. It has the following syntax:

```neva
condition_sender -> switch {
    case_sender_1 -> case_receiver_1
    case_sender_2 -> case_receiver_2
    ...
    _ -> default_receiver
}
```

Switch consists of a "condition" sender and "case" sender/receiver pairs, including a required default case with `_`. It compares the condition message with case messages for equality and executes the first matching branch. Once triggered, other branches won't fire until the next iteration.

Let's see switch in action. We're going write a program that reads name from standard output and if name is 'Alice' makes it upper case and prints, if its 'Bob' it makes it lowercase and prints (sorry, Bob), otherwise it panics because it only knows these two names:

```neva
// src/main.neva

import {
    fmt
    strings
}

def Main(start any) (stop any) {
    print fmt.Print
    scanln fmt.Scanln
    upper strings.ToUpper
    lower strings.ToLower
    println fmt.Println
    panic1, panic2 Panic
    ---
    :start -> 'Enter the name: ' -> print -> scanln -> switch {
        'Alice' -> upper
        'Bob' -> lower
        _ -> panic1
    }
    [upper, lower] -> println
    println:res -> :stop
    println:err -> panic2
}
```

We used several new things here. First, the `strings` package from the standard library contains components for string manipulation. In this example we use `strings.ToUpper` and `strings.ToLower` to convert text case.

The `fmt` package is used again - `fmt.Print` works like `Println` but without adding `\n` at the end, and `fmt.Scanln` waits for keyboard input followed by Enter.

Finally, there's the builtin `Panic` component. It immediately terminates the program with a non-zero status code when its node receives a message. We use it to 'panic', when println fails (e. g. couldn't print to stdout).

The program prompts for a name, converts it to uppercase for "Alice" or lowercase for "Bob" (panicking for any other input), then prints the result.

#### Multiple Destinations

Let's modify our program. For "Alice", we'll uppercase and lowercase simultaneously, concatenate the results and print. Any other name terminates with an error (sorry Bob!).

```neva
import {
    fmt
    strings
}

def Main(start any) (stop any) {
    print fmt.Print
    scanln fmt.Scanln
    upper strings.ToUpper
    lower strings.ToLower
    println fmt.Println
    panic1, panic2 Panic
    ---
    :start -> 'Enter the name: ' -> print -> scanln -> switch {
        'Alice' -> [upper, lower]
        _ -> panic1
    }
    (upper + lower) -> println
    println:res -> :stop
    println:err -> panic2
}
```

Things to notice:

- Fan-out to both `upper` and `lower` nodes for the 'Alice' branch
- Binary expression `(upper + lower)` connects to `println -> :stop` - inside switch we refer to their inports, inside binary expression to outports
- **Implicit parallelism utilized** - `upper` and `lower` will work in parallel, not sequentially

#### If/Else

While switch can route messages by comparing values to multiple cases, it also serves as Nevalang's if-else when working with boolean conditions. Rather than having a separate if-else construct, we use switch with a boolean condition and two branches - one for true and one for false. Let's add one more component to `src/utils/utils.neva`:

```neva
pub def ClassifyInt(data int) (neg any, pos any) {
    (:data >= 0) -> switch {
        true -> :pos
        _ -> :neg
    }
}
```

Things to notice:

1. Both outports accept `bool` messages since they're typed as `any`
2. We use `_` as default case since negative is the only other option
3. The `_` case naturally handles `false` values

Let's update `src/main.neva` to see how it can be used:

```neva
import {
    fmt
    @:src/utils
}

def Main(start any) (stop any) {
    classify utils.ClassifyInt
    println1 fmt.Println
    println2 fmt.Println
    panic Panic
    ---
    :start -> -42 -> classify
    classify:pos -> 'positive :)' -> println1
    classify:neg -> 'negative :(' -> println2
    [println1:res, println2:res] -> :stop
    [println1:err, println2:err] -> panic
}
```

Outputs:

```
negative :(
```

#### Switch True

So far we've explored message routing through comparison with set of values and boolean branching with if-else pattern. However, sometimes we need to chain multiple conditional branches where each condition is independent and not just comparing against an input value. This pattern, known as "switch-true", allows us to check multiple conditions in sequence and route messages accordingly.

Let's add one more component to `src/utils/utils.neva` and call it `CommentOnUser`. If user's name "Bob" it will comment on that, because that's the most important thing, otherwise if user's age is under 18, it will comment about that. Otherwise, if there's nothing to comment, it will just panic.

```neva
import {
    fmt
    strconv
}

// ...existing code...

pub def CommentOnUser(name string, age int) (sig any, err error) {
    println1 fmt.Println
    println2 fmt.Println
    panic Panic
    ---
    true -> switch {
        (:name == 'Bob') -> 'Beautiful name!' -> println1
        (:age < 18) -> 'Young fellow!' -> println2
        _ -> panic1
    }
    [println1:res, println2:res] -> :sig
    [println1:err, println2:err] -> :err
}
```

Here's how it can be used in `src/main.neva`

```neva
import {
    fmt
    @:src/utils
}

def Main(start any) (stop any) {
    comment utils.CommentOnUser
    panic Panic
    ---
    :start -> [
        'Bob' -> comment:name,
        17 -> comment:age
    ]
    comment:sig -> :stop
    comment:err -> panic
}
```

Output:

```
Beautiful name!
```

Note that `utils.CommentOnUser` ignored age of the user, even though it was 17. This is because how switch works - it doesn't trigger several branches in a single iteration, and once it selects branch to execute, it will ignore other branches, until next iteration will start. We can test it by replacing `Bob` with e.g. `Alice` - our switch isn't interested in Alice, but age is still 17 and it will comment on that instead.

```neva
:start -> [
    'Alice' -> comment:name,
    17 -> comment:age
]
```

Output:

```
Young fellow!
```

By the way, there's another way to solve this problem. We can use if-else pattern and nest switches one inside another like this:

```neva
:age < 18 -> switch {
    true -> 'Young fellow!' -> println1
    _ -> (:name == 'Bob') -> switch {
        true -> 'Beauteful name!' -> println2
        _ -> panic
    }
}
```

You should never do that if it's possible to follow "switch-true" pattern, because it's much easier to read and doesn't envolve two switch nodes.

#### Multiple Sources

> WARNING: This is important section. **Don't skip** it!

One might ask, why didn't we cover multiple case senders if we covered multiple receivers? When using switch with multiple case receivers, it works differently than in control flow languages. For example:

```neva
switch {
    ['Alice', 'Bob'] -> upper
    _ -> lower
}
```

Is **not** "if either Alice or Bob then do uppercase". It's a fan-in, meaning `Alice` and `Bob` are concurrent. Switch will select the first value sent as a case, which is random since both are message literals.

> These semantics might change in the future. There's an [issue](https://github.com/nevalang/neva/issues/788) about that.

**How To Cover this Case?**

```neva
pass1 Pass{Upper}
pass2 Pass{Upper}
---
... -> switch {
    'Alice' -> pass1
    'Bob' -> pass2
    _ -> lower
}
[pass1, pass2] -> upper
```

With a little bit of boilerplate

**WARNING** (Possible Concurrency Issues)

`[pass1, pass2] -> upper` guarantees that `upper` receives messages in the exact same order as they are sent by `pass1` and `pass2` nodes. And `switch` guarantees that it will never reiceve next message from it's sender `... ->` until previously selected receiver has received.

BUT switch doesn't know about `[pass1, pass2] -> upper`. It will wait for `pass1` and `pass2` to receive yes, but it will not wait `pass1` and `pass2` to _send_!

Example: let's say we send `'bob', 'alice'` sequence to our switch. This sequence of events is totally possible:

```
pass1 received alice and blocked
pass2 received bob and blocked
pass2 send BOB and unlocked
pass1 send ALICE and unlocked
```

For `[pass1, pass2] -> upper` connection it will mean that `upper` must receive `BOB` _before_ `Alice` which is not correct!

_Safe Solution_

This situation will never happen if you guarantee that your component never receive next input until it previous sent result is received. But this can only be guaranteed by the parent of your compnent.

Idiomatic way would be to use higher-order component, that can wrap your concurrent component and use it in a safe way, so it will process messages fully sequentionally - synchronized.
