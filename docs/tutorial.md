# Tutorial

Welcome to a tour of the Nevalang programming language. This tutorial will introduce you to Nevalang through a series of guided examples.

1. [Welcome](#welcome)
   - [What kind of language is this?](#what-kind-of-language-is-this)
   - [Installation](#installation)
   - [Hello, World!](#hello-world)
   - [Compiling programs](#compiling-programs)

<!-- Advanced topics: tracing -->

## Welcome

### What Kind of Language is This?

Nevalang is a general-purpose programming language that uses dataflow instead of control flow, lacking variables and functions, and expressing programs through pure message passing with nodes, connections, and ports. It is implicitly parallel, meaning all nodes operate in parallel by default, eliminating the need for threads, coroutines, or async-await, which simplifies some tasks while complicating others.

Influenced by functional programming, Nevalang embraces immutability and higher-order components, disallowing data mutation and shared state. It is a compiled, strongly statically-typed, garbage-collected language, sharing Go's abstraction level but aligning more with Rust's strictness. Using Go as a backend, it supports all Go targets, including WASM and cross-compiled machine code.

### Installation

#### Requirements

Make sure you have [Go compiler](https://go.dev/dl/) installed.

#### Via Shell Script

For Mac OS and Linux:

```bash
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash
```

If your device is connected to a chinese network:

```bash
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

```bash
neva new my_awesome_project
```

Each new project contains a Hello World program, so we can just run it

```bash
neva run my_awesome_project/src
```

You should see the following output:

```bash
Hello, World!
```

If you open `my_awesome_project/src/main.neva` with your favorite IDE, you'll see this:

```neva
import { fmt }

def Main(start any) (stop any) {
	println fmt.Println
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

This will produce an `output` file in the directory where neva-cli was executed, typically the project's root.

Run `neva build --help` to see various options. Here are some examples:

```shell
# will name output file "my_awesome_binary"
neva build my_awesome_project/src --output=my_awesome_binary

# will generate go code instead of machine code
neva build my_awesome_project/src --target=go

# will output my_awesome_wasm.wasm
neva build my_awesome_project/src --target=wasm --output=my_awesome_wasm
```
