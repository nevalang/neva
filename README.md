# Neva Programming Language

**Neva** is a general purpose [flow-based](https://en.wikipedia.org/wiki/Flow-based_programming) programming language with static [structural](https://en.wikipedia.org/wiki/Structural_type_system) typing and [implicit parallelism](https://en.wikipedia.org/wiki/Implicit_parallelism) that compiles to machine code. Designed with [visual programming](https://en.wikipedia.org/wiki/Visual_programming_language) in mind. For the era of effortless concurrency.

> WARNING: This project is under heavy development and not production ready yet.

```neva
const {
	greeting string 'Hello, World!'
}

components {
	Main(enter any) (exit any) {
		nodes { printer Printer<string> }
		net {
			in:enter -> ($greeting -> printer:msg)
			printer:msg -> out:exit
		}
	}
}
```

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md)

## Architecture

See [ARCHITECTURE.md](./ARCHITECTURE.md)

## FAQ

See [FAQ.md](./docs/faq.md)
