![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow Based Programming Language</p>**

# Neva (Proof of Concept)

A general-purpose flow-based programming language with static types and implicit parallelism. Compiles to machine code and Go.

Website: https://nevalang.org/

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

When `Main` receives a `start` signal, it sends a `'Hello, World!'` message to `data` inport of the node `printer`, which is an instance of a `Printer` component. Then `printer` sends a signal to it's `sig` outport, we use it as a signal to our `stop` outport.

## How it works?

Program starts with the single `start` message and ends after first `stop` message. This is our only way to control execution flow. For everything else we control flow of the data.

We create networks of components that are then executed inside a special asynchronous message passing runtime where everything happens in parallel. The compiler analyzes these networks for various semantic errors, including type-safety.

The command line interface can execute such source code directly using an interpreter, which is handy for development and debugging. For production use, it emits either machine code, Golang, or WASM.

[Read more about the language](https://nevalang.org/docs/about)

## Features 🚀

- Flow-Based Programming
- Effortless Concurrency
- Strong Static Typing
- Multi-Target Compilation (PoC)
- Clean C-like Syntax
- Interpreter Mode
- First-Class Dependency Injection
- Builtin Observability
- Garbage Collection

### Roadmap (🚧 WIP)

> ⚠️ WARNING: This project is under heavy development and **not production ready** yet.

- Go Interop
- No Runtime Exceptions
- Visual Programming

## Contributing

Nevalang needs your help, it has only one maintainer. Join community. Together we can change programming for the best.

- [Discord](https://discord.gg/8fhETxQR)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md)

---
