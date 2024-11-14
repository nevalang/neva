![Big Header](./assets/header/big_1.svg "Big header with nevalang logo")

<div align="center" style="display:grid;place-items:center;">

<h1>Dataflow Programming Language</h1>

[Documentation](./docs/)
| [Examples](./examples/)
| [Community](#community)
| [Releases](https://github.com/nevalang/neva/releases)
| [Contributing](./CONTRIBUTING.md)
| [Architecture](./ARCHITECTURE.md)

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

</div>

A general-purpose dataflow programming language with static types and implicit parallelism. Compiles to machine code and Go.

## 🚀 Features

- Dataflow programming
- Implicit parallelism
- Compiles to machine code and Go
- Garbage collection
- Strong static typing
- Clean C-like syntax
- ...And more!

> ⚠️ This project is currently under heavy development and is not yet ready for production use.

## 🔧 Quick Start

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
import { fmt }

def Main(start) (stop) {
	fmt.Println
	---
	:start -> { 'Hello, World!' -> println -> :stop }
}
```

The `import { fmt }` statement imports the standard library's `fmt` package which provides common formatting and printing functionality. The `Main` component has `start` inport and `stop` outport, with a `println` node (instance of stdlib's `fmt.Println`). The network after `---` shows: on `start` message, `"Hello, World!"` is sent to `println`, then program terminates via `stop` signal.

### What's Next?

- [Documentation](./docs/README.md)
- [Examples](./examples/)
- [Community](#community)

## 🚧 Roadmap

Nevalang is in its early stages, but community support can help it grow into a mature, feature-rich language.

- Grow community and improve docs
- Expand stdlib (including test-framework)
- Better syntax and more features
- Enhance developer experience (lsp, debugger, etc)
- Implement **Go interoperability** (call Go from Neva and vice versa)
- Enable **visual programming** in VSCode

We seek contributors to join our small team.

## 📢 Community

Join our community and help shape the future of programming:

- [Discord](https://discord.gg/dmXbC79UuH)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

Also please check our [CoC](./CODE_OF_CONDUCT.md).

## 🤝 Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).
