![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow Based Programming Language</p>**

# Neva (Proof of Concept)

A general-purpose flow-based programming language with static types and implicit parallelism. Compiles to machine code and Go.

Website: https://nevalang.org/

> ⚠️ WARNING: This project is under heavy development and **not production ready** yet.

## Hello World

```neva
component Main(start any) (stop any) {
	nodes { printer Printer<string> }
	net {
		:start -> ('Hello, World!' -> printer:data)
		printer:sig -> :stop
	}
}
```

Nevalang programs are executable dataflow graphs where components are connected to each other through input and output ports. No state, no control flow, just message passing in fully asynchronous environment.

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
- Garbage Collection

Please note that even though these features are technically implemented, **developer-experience could be very bad** due to current project state which is pre-MVP. **No backward-compatibility** guarantees at the time.

### Roadmap (🚧 WIP)

Nevalang is at extrimely early stage but with the help of community it can become feature-reach mature langauge.

- Building a Community
- Core Standard Library
- Feature-Rich LSP-compatible Language Server
- Go Interop (import go from neva, import neva from go)
- DAP-compatible Debugger
- No Runtime Exceptions (If it runs then it works)
- Visual Programming in VSCode (Nevalang becomes hybrid langauge)

Nevalang needs **your** help - it only have one maintainer currently.

## Community

Join community. **Together** we can change programming for the better:

- [Discord](https://discord.gg/8fhETxQR)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).

Please note that due to early stage of development documentation can be sometimes outdated. You can reach main maintainer to find help.
