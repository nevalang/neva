![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow-Based Programming Language</p>**

![tests](https://github.com/nevalang/neva/actions/workflows/test.yml/badge.svg?branch=main) ![lint](https://github.com/nevalang/neva/actions/workflows/lint.yml/badge.svg?branch=main)

> ⚠️ Warning: This project is currently under heavy development and is not yet ready for production use.

# Neva

A general-purpose dataflow programming language with static types and implicit parallelism. Compiles to machine code and Go.

Website: https://nevalang.org/

## Features 🚀

- **Dataflow Programming**
- **Strong Static Typing**
- **Implicit Parallelism**
- **Compiles to Machine Code, Go and More**
- **Clean C-like Syntax**
- **Garbage Collection**
- **...And more!**

Please note that even though these features are technically implemented, _developer-experience may be bad_ due to the current state of the project. **No backward-compatibility guarantees** at the time.

## Quick Start

### Installation

For Mac OS and Linux:

```bash
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.sh | bash
```

If your device is connected to a chinese network:

```bash
curl -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/cnina/install.sh | bash
```

For Windows (please note there's an WIP [issue](https://github.com/nevalang/neva/issues/499) with Windows Defender, try manual download from releases if installed won't work):

```batch
curl -o installer.bat -sSL https://raw.githubusercontent.com/nevalang/neva/main/scripts/install.bat && installer.bat
```

### Creating a project

```bash
neva new test
```

### Running

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
flow Main(start) (stop) {
	nodes { Println }
	:start -> ('Hello, World!' -> println -> :stop)
}
```

Here we define a _flow_ `Main` with _inport_ `start` and _outport_ `stop`. It contains one _node_, `println`, an _instance_ of `Println`. The _network_ consist of one _connection_: upon receiving a message from `start`, "Hello, World!" is sent to `println`. After printing, the program terminates by signaling `stop`.

### Execute

Now run (make sure you are in the `test` directory with `neva.yml`):

```bash
neva run test/src # or neva run test/src/main.neva
```

You should see the following output:

```bash
Hello, World!
```

### What's Next?

- [Learn more about the language](https://nevalang.org/)
- [See more examples](./examples/) (`git clone` this repo and `neva run` them!)

## Roadmap (🚧 WIP)

Nevalang is at an extremely early stage but with the help of community it can become a feature-rich, mature language.

- Bigger Community and Better Documentation
- Batteries included stdlib (including Testing Framework)
- More language features including _Pattern-Matching_
- Good DX (Language Server, Debugger, Linter, etc)
- **Go Interop** (import Go from Neva and Neva from Go)
- **Visual Programming** in VSCode (Neva becomes hybrid langauge)

[See backlog for more details](https://github.com/orgs/nevalang/projects)

Nevalang needs your help - it currently has just a few maintainers.

## Community

Join community. _Together we can change programming_ for the better:

- [Discord](https://discord.gg/dmXbC79UuH)
- [Reddit](https://www.reddit.com/r/nevalang/)
- [Telegram group](https://t.me/+H1kRClL8ppI1MWJi)
- [Telegram channel](https://t.me/+H1kRClL8ppI1MWJi)

Also please check our [CoC](./CODE_OF_CONDUCT.md).

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md).

Neva is a relatively small and simple language. Don't be intimidated, feel free to dive in and hack around. Some directories have a `README.md`.

Note that, due to the early stage of development, the documentation can sometimes be outdated. Feel free to reach out to maintainers if you need _any_ help.
