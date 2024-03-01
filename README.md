![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow-Based Programming Language</p>**

# Neva

A general-purpose flow-based programming language with static types and implicit parallelism. Compiles to machine code and Go.

Website: https://nevalang.org/

> ⚠️ Warning: This project is currently under heavy development and is not yet ready for production use.

## Hello World

```neva
component Main(start any) (stop any) {
	nodes { Printer<string> }
	net {
		:start -> ('Hello, World!' -> printer:data)
		printer:sig -> :stop
	}
}
```

Nevalang programs are executable dataflow graphs with components connected through input and output ports. No state, no control flow, just message passing in a fully asynchronous environment.

[Read more about the language](https://nevalang.org/docs/about)

## Features 🚀

- Flow-Based Programming
- Implicit Parallelism
- Strong Static Typing
- Multi-Target Compilation
- Clean C-like Syntax
- Interpreter Mode
- Builtin Dependency Injection
- Builtin Observability
- Package Management
- Garbage Collection

Please note that even though these features are technically implemented, **developer-experience may be bad** due to current project state. **No backward-compatibility** guarantees at the time.

### Roadmap (🚧 WIP)

Nevalang is at an extremely early stage but with the help of community it can become a feature-rich, mature language.

- Building a Community
- Core Standard Library
- Feature-Rich LSP-compatible Language Server
- Go Interop (import go from neva and neva from go)
- DAP-compatible Debugger
- Testing Framework
- No Runtime Exceptions (If it runs then it works)
- Visual Programming in VSCode (Nevalang becomes hybrid langauge)

Nevalang needs your help - it currently has only one maintainer.

## Community

Join community. **Together** we can change programming for the better:

- [Discord](https://discord.gg/8fhETxQR)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).

Please note that, due to the early stage of development, the documentation can sometimes be outdated. You can reach main maintainer to find help.
