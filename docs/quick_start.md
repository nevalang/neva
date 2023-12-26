# Quick Start

## Install Nevalang

```bash
...
```

## Prepare example project

```shell
mkdir neva_example
cd neva_example
neva init
```

This will create the following structure

```
neva.yml
cmd/
    main/
        main.neva
```

```neva
components {
	Main(enter) (exit) {
		net { in.enter -> out.exit }
	}
}
```