# Neva

Flow based programming environment.

```shell
$ git clone git@github.com:emil14/neva.git
$ cd neva/cmd/neva
$ go install
$ neva run ../../examples/arr.yml
```

## Messages are immutable

Message is interface where only getters are defined.
There is no "behaviour" in messages, they are only data.
This data is once created and there for can only be used for creating new data.

Components 