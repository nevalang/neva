<div align="center">
  <img src="./assets/logo/light_gradient.svg" alt="Nevalang logo">
</div>

<div align="center" style="display:grid;place-items:center;">

<h1>Neva Programming Language</h1>

[**Documentation**](./docs/README.md)
| [**Examples**](./examples/)
| [**Community**](#-community)
| [**Releases**](https://github.com/nevalang/neva/releases)
| [**Contributing**](./CONTRIBUTING.md)
| [**Architecture**](./ARCHITECTURE.md)

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

</div>

<div align="center">
  <i>Next-generation programming language that solves programmer's problems.</i>
</div>

## 🤔 What Is Nevalang?

Nevalang is a new kind of programming language where instead of writing step-by-step instructions you create **networks** where data flows between **nodes** as **immutable messages** and everything runs **in parallel by default**. After **type-checking**, your program is compiled into **machine code** and can be distributed as a **single executable** with zero dependencies.

All this combined with native **stream processing** support and features such as **advanced error handling** makes Nevalang a perfect choice for **cloud applications** with **high concurrency**.

Future updates will include **visual programming** and **Go interoperability** to allow gradual adoption and leverage existing ecosystem.

## 🔥 Features

- 📨 **Dataflow Programming** - Write programs as message-passing graphs
- 🔀 **Implicit Parallelism** - Everything is parallel by default, no async-await/threads/goroutines/etc.
- 🛡️ **Strong Static Typing** - Robust type system with generics and pattern-matching
- 🚀 **Machine Code Compilation** - Compile for any Go-supported platform, including WASM
- ⚡️ **Stream Processing** - Handle real-time data with streams as first class citizens
- 🧯 **Advanced Error Handling** - Errors as values with `?` operator to avoid boilerplate
- 🪶 **Minimal Core** - Simple language with limited abstractions
- 📦 **Package Manager** - Publish packages by pushing a git-tag
- 🧩 **Functional Patterns** - Immutability and higher-order components
- 🔌 **Dependency Injection** - Modularity with interfaces and DI
- ♻️ **Garbage Collection** - Automatic memory management using Go's low-latency GC
- 🌈 **Visual Programming** (WIP): Edit programs as visual graphs
- 🔄 **Go Interoperability** (WIP): Call Go from Neva and Neva from Go
- 🕵 **NextGen Debugging** (WIP): Observe execution in realtime and intercept messages on the fly

> ⚠️ This project is under active development and not yet production-ready.

## 👋 Hello, World!

```neva
import { fmt }

def Main(start any) (stop any) {
	println fmt.Println<string>
	---
	:start -> 'Hello, World!' -> println -> :stop
}
```

What’s happening here:

- `import { fmt }` loads the `fmt` package for printing
- `def Main` defines the main component with input port `start` and output port `stop`
- `:start -> ‘Hello, World!’ -> println -> :stop` defines a connection that sends `Hello, World!` string to the `println` printer-node and then terminates the program

## 📢 Community

As you can see, this is quite an ambitious project. Typically, such projects are backed by companies, but Nevalang is maintained by a very small group of enthusiasts. Your support by joining us will show interest and motivate us to continue.

- [Discord](https://discord.gg/dmXbC79UuH)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)

Also, **please give us a star ⭐️** to increase our chances of getting into GitHub's trending repositories and tell your friends about the project. The more attention Nevalang gets, the higher our chances of actually making a difference!

## 💭 What's Next?

- [Documentation](./docs/README.md) - Install and learn the language basics
- [Examples](./examples/) - Learn the language by small programs

> Please keep in mind that these resources might not be ready or may be outdated due to the current state of the project. However, rest assured that we take development seriously. We simply don't have enough time to keep everything up to date all the time. Please don't feel intimidated and contact us on our social platforms if you have any questions. We welcome _any_ feedback, no matter what.

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).

Check out our [roadmap](https://github.com/nevalang/neva/milestones?direction=asc&sort=due_date&state=open) to see what we are planning next. You can follow the development process on our [Kanban board](https://github.com/orgs/nevalang/projects/2/views/3?filterQuery=).

Also please check our [CoC](./CODE_OF_CONDUCT.md).
