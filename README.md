![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow Based Programming Language</p>**

# Neva (Proof of Concept)

A general-purpose flow-based programming language with static types and implicit parallelism. Compiles to machine code and Go.

Website: https://nevalang.org/

> ⚠️ WARNING: This project is under heavy development and **not production ready** yet.

## Hello World

Nevalang programs are executable dataflow graphs. Execution starts with the single `start` message and ends after first `stop` message.

```neva
component Main(start any) (stop any) {
	nodes { printer Printer<string> }
	net {
		:start -> ('Hello, World!' -> printer:data)
		printer:sig -> :stop
	}
}
```

When `Main` receives a `:start` signal, it sends a `'Hello, World!'` message to `data` inport of the node `printer`, which is an instance of a `Printer` component. Then `printer` sends a signal to it's `sig` outport and we use that as a signal to our `:stop` outport.

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

Please note that even though these features are technically implemented, **developer-experience could be very bad** due to current project state which is pre-MVP.

### Roadmap (🚧 WIP)

Nevalang is at extrimely early stage but with the help of community it can become feature-reach mature langauge.

- Building a Community
- Core Standard Library
- Feature-Rich LSP-compatible Language Server
- Go Interop (import go from neva, import neva from go)
- DAP-compatible Debugger
- No Runtime Exceptions (If it runs then it works)
- Visual Programming in VSCode (Nevalang becomes hybrid langauge)

Nevalang needs **your** help - it only have one maintainer.

## Community

Join community. **Together** we can change programming for the better:

- [Discord](https://discord.gg/8fhETxQR)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).

Please note that due to early stage of development documentation can be sometimes outdated. You can reach main maintainer to find help.
