![Big Header](./assets/header/big_1.svg "Big header with nevalang logo")

**<p align="center">Dataflow Programming Language</p>**

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

> ⚠️ **Warning**: This project is currently under heavy development and is not yet ready for production use.

# Nevalang

A general-purpose dataflow programming language with static types and implicit parallelism. Compiles to machine code and Go.

## Features 🚀

- **Dataflow Programming**
- **Implicit Parallelism**
- **Compiles to Machine Code, Go and More**
- **Garbage Collection**
- **Strong Static Typing**
- **Clean C-like Syntax**
- **...And more!**

_Note: Features are implemented but may have poor developer experience. No backward-compatibility guaranteed at the moment._

## Quick Start

### Installation

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

### Hello World

First, use Neva CLI to create a project template

```bash
neva new my_awesome_project
```

Then run it

```bash
neva run my_awesome_project/src
```

You should see the following output

```bash
Hello, World!
```

If you open `my_awesome_project/src/main.neva` with your favorite IDE you'll see this

```neva
flow Main(start) (stop) {
	Println
	---
	:start -> ('Hello, World!' -> println -> :stop)
}
```

The `Main` component has `start` inport and `stop` outport, with a `println` node (instance of stdlib's `Println`). The network after `---` shows: on `start` message, `"Hello, World!"` is sent to `println`, then program terminates via `stop` signal.

### What's Next?

- [Documentation](./docs/README.md)
- [Examples](./examples/)
- [Community](#community)

## Roadmap (🚧 WIP)

Nevalang is in its early stages, but community support can help it grow into a mature, feature-rich language.

- Grow community and improve docs
- Expand stdlib (including test-framework)
- New language features (e.g. true patter-matching)
- Enhance developer experience (lsp, debugger, etc)
- Implement **Go interoperability** (call Go from Neva and Neva from Go)
- Enable **visual programming** in VSCode

Nevalang seeks contributors to join its small team of maintainers.

## Community

Join our community and help shape the future of programming:

- [Discord](https://discord.gg/dmXbC79UuH)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

Also please check our [CoC](./CODE_OF_CONDUCT.md).

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).
