![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow Based Programming Language</p>**

# Neva

A general-purpose, flow-based programming language with static typing and implicit parallelism, designed with visual programming in mind, that compiles to machine code and Go.

Website: https://nevalang.org

```neva
component Main(start any) (stop any) {
	nodes { printer Printer<string> }
	net {
		in:start -> ('Hello, World!' -> printer:msg)
		printer:msg -> out:stop
	}
}
```

## 🚀 Features

- 🌊 Flow-Based Programming
- 🔀 Effortless Concurrency
- 🛡️ Static Type System
- 🎯 Multi-Target Compilation
- ✨ Simple and Clean C-like Syntax
- 🏃‍♂️ Interpreter Mode
- 💉 First-Class Dependency Injection
- 🕵️‍♂️ Builtin Observability
- ♻️ Garbage Collection
- 🌈 Visual Programming (WIP)
- 🦫 Go Interop (WIP)
- 🦺 No Runtime Exceptions (WIP)

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md)

---

> ⚠️ WARNING: This project is under **heavy** development and not production ready yet.
