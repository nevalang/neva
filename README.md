<div align="center">
  <img src="./assets/logo/light_gradient.svg" alt="Nevalang logo">
</div>

<div align="center" style="display:grid;place-items:center;">

<h1>Neva Programming Language</h1>

<p align="center">
  <img src="./assets/animations/dataflow.gif" alt="Dataflow">
</p>

<div align="center">
  Next-generation programming language that solves programmer's problems
</div>
<br/>

[**Documentation**](./docs/README.md)
| [**Examples**](./examples/)
| [**Community**](#-community)
| [**Releases**](https://github.com/nevalang/neva/releases)
| [**Contributing**](./CONTRIBUTING.md)
| [**Architecture**](./ARCHITECTURE.md)

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

</div>

## 🤔 What Is Nevalang?

Nevalang is a new kind of programming language where instead of writing step-by-step instructions you create **networks** where data flows between **nodes** as **immutable messages** and everything runs **in parallel by default**. After **type-checking**, your program is compiled into **machine code** and can be distributed as a **single executable** with zero dependencies.

Combined with built-in **stream processing** support and features like **advanced error handling**, Nevalang is the perfect choice for **cloud-native applications** requiring **high concurrency**.

Future updates will include **visual programming** and **Go interoperability** to allow gradual adoption and leverage existing ecosystem.

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
- `:start -> 'Hello, World!' -> println -> :stop` creates a connection that sends the string to `println` when program starts and terminates after printing

> ℹ️ Runtime sends a message to `Main:start` at startup and then waits for `Main:stop` to terminate.

## 🔥 Features

- 📨 **Dataflow Programming** - Write programs as message-passing graphs
- 🔀 **Implicit Parallelism** - Everything is parallel by default, no async-await/threads/goroutines/etc.
- 🛡️ **Strong Static Typing** - Robust type system with generics and pattern-matching
- 🚀 **Machine Code Compilation** - Compile for any Go-supported platform, including WASM
- ⚡️ **Stream Processing** - Handle real-time data with streams as first class citizens
- 🧯 **Advanced Error Handling** - Errors as values with `?` operator to avoid boilerplate
- 🧩 **Functional Patterns** - Immutability and higher-order components
- 🔌 **Dependency Injection** - Modularity with interfaces and DI
- 🪶 **Minimal Core** - Simple language with limited abstractions
- 📦 **Package Manager** - Publish packages by pushing a git-tag
- ♻️ **Garbage Collection** - Automatic memory management using Go's low-latency GC
- 🌈 **Visual Programming** (WIP): Edit programs as visual graphs
- 🔄 **Go Interoperability** (WIP): Call Go from Neva and Neva from Go
- 🕵 **NextGen Debugging** (WIP): Observe execution in realtime and intercept messages on the fly

## 🧐 Why Nevalang?

Let's compare Nevalang with Go. We could compare it to any language but Go is a simple reference since Nevalang is written in Go.

| **Feature**              | **Neva**                                                           | **Go**                                                                            |
| ------------------------ | ------------------------------------------------------------------ | --------------------------------------------------------------------------------- |
| **Paradigm**             | Dataflow - nodes send and receive messages through connections     | Control flow - execution moves through instructions step by step                  |
| **Concurrency**          | Implicit - everything is concurrent by default                     | Explicit - goroutines, channels, and mutexes                                      |
| **Error Handling**       | Errors as values with `?` operator to avoid boilerplate            | Errors as values with `if err != nil {}` boilerplate                              |
| **Mutability**           | Immutable - no variables and pointers; data races are not possible | Mutable - variables and pointers; programmer must avoid data races                |
| **Null Safety**          | Yes - nil pointer dereference is impossible                        | No - nil pointer dereference is possible                                          |
| **Zero Values**          | No zero values - everything must be explicitly initialized         | Zero values by default - everything can be initialized implicitly                 |
| **Subtyping**            | Structural - types are equal by their shape                        | Nominal - types are equal by their name                                           |
| **Traceback**            | Automatic - every message traces its path                          | Manual - programmer must explicitly wrap every error to add context               |
| **Dependency Injection** | Built-in - any component with dependency expects injection         | Manual - programmer must create constructor function that takes dependencies      |
| **Stream Processing**    | Native support with components like `Map/Filter/Reduce`            | Programmer must manually implement dataflow patterns with goroutines and channels |

## 🏭 Architecture

> This is a high-level overview. For a more detailed overview of the architecture, please see [ARCHITECTURE.md](./ARCHITECTURE.md)

```mermaid
flowchart LR
  source_code-->compiler-->|go_code| go_compiler

  subgraph compiler
    parser-->analyzer-->backend
  end

  go_compiler-->machine_code
  go_compiler-->wasm
```

Nevalang compiles to dependency-free, human-readable Go code that uses goroutines and channels for message-passing with parallelism. The Go compiler then produces optimized platform-specific code for any supported platform. This approach gives our programs access to Go's production-grade runtime with an advanced scheduler, garbage collector, and battle-tested standard library. We stand on the shoulders of giants.

## 💭 What's Next?

> ℹ️ Please note that these resources may be outdated due to the current state of the project. We take development seriously but have limited time to keep everything current. Don't hesitate to reach out on our social platforms with any questions - we welcome all feedback!

- [Documentation](./docs/README.md) - Install and learn the language basics
- [Examples](./examples/) - Learn the language by small programs

### 📢 Community

As you can see, this is quite an ambitious project. Typically, such projects are backed by companies, but Nevalang is maintained by a very small group of enthusiasts. **Your support by joining us will show interest and motivate us to continue.**

- [**Discord**](https://discord.gg/dmXbC79UuH)
- [**Reddit**](https://www.reddit.com/r/nevalang/)
- [**Telegram**](https://t.me/+H1kRClL8ppI1MWJi)
- [**Twitter**](https://x.com/neva_language)

### 🤝 Contributing

1. See [contributing](./CONTRIBUTING.md) and [architecture](./ARCHITECTURE.md)
2. Check out [roadmap](https://github.com/nevalang/neva/milestones?direction=asc&sort=due_date&state=open) and [kanban-board](https://github.com/orgs/nevalang/projects/2/views/3?filterQuery=)
3. Also please read our [CoC](./CODE_OF_CONDUCT.md)
4. Join [discord server](https://discord.gg/dmXbC79UuH)

### 🙏 Support

Please **give us a star ⭐️** to increase our chances of getting into GitHub trends - the more attention Nevalang gets, the higher our chances of actually making a difference.

<p align="center">
  <img src="./assets/animations/github_star.gif" alt="GitHub Star">
</p>
