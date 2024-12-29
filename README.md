![Big Header](./assets/header/big_1.svg "Big header with nevalang logo")

<div align="center" style="display:grid;place-items:center;">

<h1>Dataflow Programming Language</h1>

[Documentation](./docs/README.md)
| [Examples](./examples/)
| [Community](#-community)
| [Releases](https://github.com/nevalang/neva/releases)
| [Contributing](./CONTRIBUTING.md)
| [Architecture](./ARCHITECTURE.md)

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

</div>

A general-purpose dataflow programming language with static types and implicit parallelism. Compiles to machine code and Go.

## 🚀 Features

- **Dataflow Programming** - write programs as message-passing graphs
- **Implicit Parallelism** - no threads/mutexes/coroutines/channels/etc
- **Reliable Type System** - strong static typing with unions and pattern matching
- **Native Stream Processing** - streams as first-class citizens
- **Compilation** - generate code for any platform Go supports including machine code and WASM
- **Functional Idioms** - immutability and composition via higher-order components
- **Garbage Collection** - utilizes Go runtime's efficient GC
- ... and much more!

> ⚠️ This project is currently under heavy development and is not yet ready for production use.

## Hello World

```neva
import { fmt }

def Main(start any) (stop any) {
	println fmt.Println<string>
	---
	:start -> 'Hello, World!' -> println -> :stop
}
```

The `import { fmt }` statement imports the standard library's `fmt` package which provides common formatting and printing functionality. The `Main` component has `start` input and `stop` output ports of type `any`, with a `println` node (instance of `fmt.Println`). The network after `---` shows: on `:start` message, `"Hello, World!"` is sent to `println`, then program terminates via `:stop` signal.

## What's Next?

- [Documentation](./docs/README.md) - learn how to install and use Nevalang
- [Examples](./examples/) - explore sample programs
- [Community](#community) - join us!

## 🚧 Roadmap

Check out our [full roadmap](https://github.com/nevalang/neva/milestones?direction=asc&sort=due_date&state=open). Here are a few highlights:

- **Visual Programming** - observe and debug your dataflow visually
- **Go Interop** - leverage the Go ecosystem and enable gradual adoption in Go projects

Nevalang is in its early stages, and community support can help it grow into a mature, feature-rich language. _We are looking for contributors_ to join our small team.

## 📢 Community

Join our community and help shape the future of programming:

- [Discord](https://discord.gg/dmXbC79UuH)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram](https://t.me/+H1kRClL8ppI1MWJi)
- [Twitter](https://x.com/neva_language)

Please also check our [Code of Conduct](./CODE_OF_CONDUCT.md).

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).
