![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow-Based Programming Language</p>**

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

> ⚠️ Warning: This project is currently under heavy development and is not yet ready for production use.

# Neva

A general-purpose flow-based programming language with static types and implicit parallelism. Compiles to machine code and Go.

Website: https://nevalang.org/

## Features 🚀

- **Flow-Based Programming**
- **Strong Static Typing**
- **Implicit Parallelism**
- **Compiles to Machine Code, Go and More**
- **Clean C-like Syntax**
- **Interpreter Mode**
- **First-Class Dependency Injection**
- **Builtin Observability**
- **Package Management**
- **Garbage Collection**

Please note that even though these features are technically implemented, _developer-experience may be bad_ due to current project state. _No backward-compatibility_ guarantees at the time.

## Quick Start

### Download Neva CLI

For Mac OS and Linux:

```bash
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash
```

if you device connected to China network:

```bash
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash
```

For Windows (please note there's an WIP [issue](https://github.com/nevalang/neva/issues/499) with Windows Defender, try manual download from releases if installed won't work):

```batch
curl -o installer.bat -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.bat && installer.bat
```

### Create test project

```bash
neva new test
```

Replace the code in `cat test/src/main.neva` with the following:

### Execute

```bash
neva run test/src
```

You should see the following output:

```bash
Hello, World!
```

### What's inside?

If you open `test/src/main.neva` with your favorite IDE you'll see this

```neva
component Main(start) (stop) {
	nodes { Println }
	net {
		:start -> ('Hello, World!' -> println -> :stop)
	}
}
```

Here we define component `Main` with inport `start` and outport `stop`. It has 1 node `println` that is an instance of the `Println` component. Then we define `net` - set of connections that describe dataflow. When message from `start` received, a string literal `Hello, World!` is sent to node `println`. When that message is printed, program is terminated by sending to `stop`.

### What's Next?

- [Learn more about the language](https://nevalang.org/)
- [See more examples](./examples/) (`git clone` this repo and `neva run` them!)

## Roadmap (🚧 WIP)

Nevalang is at an extremely early stage but with the help of community it can become a feature-rich, mature language.

- Building a Community
- Core Standard Library
- Feature-Rich LSP-compatible Language Server
- **Go Interop** (import go from neva and neva from go)
- DAP-compatible Debugger
- Testing Framework
- No Runtime Exceptions (If it runs then it works)
- **Visual Programming** in VSCode (Nevalang becomes hybrid langauge)

[See backlog for more details](https://github.com/orgs/nevalang/projects)

Nevalang needs _your_ help - it currently just a few maintainers.

## Community

Join community. Together we can change programming for the better:

- [Discord](https://discord.gg/dmXbC79UuH)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

Also please check our [CoC](./CODE_OF_CONDUCT.md).

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).

Neva is a relatively small and simple language. Don't be intimidated, feel free to dive in and hack around. Some directories have a `README.md`.

Note that, due to the early stage of development, the documentation can sometimes be outdated. Feel free to reach maintainers if you need _any_ help.
