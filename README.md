![Big Header](./assets/header/big.svg "Big header with nevalang logo")

**<p align="center">Flow Based Programming Language</p>**

## Neva

A general-purpose, flow-based programming language with static typing and implicit parallelism, designed with visual programming in mind, that compiles to machine code and Go.

Website: https://nevalang.org

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

See [CONTRIBUTING.md](./CONTRIBUTING.md) and [ARCHITECTURE.md](./ARCHITECTURE.md)

---

> WARNING: This project is under heavy development and not production ready yet.
