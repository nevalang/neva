Только компонент у которого нет ввода-вывода, может быть использован
как рут. Такой компонент также нельзя использовать никаким другим образом, ибо у него нет ввода-вывода, а любая зависимость должна его иметь. Компилятор проверит, что пакет не содержит "мертвый код".


# Generate DevServer-SDK

```shell
protoc api/devserver.proto \
    --js_out=import_style=commonjs,binary:web/src/sdk \
    --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:web/src/sdk \
    --go_out=pkg/devserversdk \
    --go-grpc_out=pkg/devserversdk
```

# Generate Runtime-SDK

```shell
protoc api/runtime.proto --go_out=pkg/runtimesdk
```

```
BASE_PATH # path to directory where programs (packages) can be found
compiler will use recoursion to find all dirs with pkg.yml files.
```

# Respect

Flow Based Programming Environment

- https://spec.openapis.org/oas/latest.html
- https://openapi-generator.tech/docs/generators/go-server
- https://openapi-generator.tech/docs/generators/typescript

### Deps

- Docker
- Go
- NodeJS and NPM

# Rules

0. Runtime runs programs
1. Compiler creates programs for runtime
2. Source code is generated

# Module validation rules

1. Module must have `in` and `out` with at least 1 port each
2. There must be no `Unknown` among port types
3. There must be `net` defining connections
4. There could be `workers` defing network's worker nodes
5. If there are `workers` then there must be also `deps` defining dependencies
6. Any `worker` must point to defined `dep`
7. Any `dep` must have valid `in` and `out` defined
8. There could be `const` defining constant values

## Checking network

9.  There must be no `deps` unused in `workers`
10. There must be no `worker` or `const` not used by `net`

## Getting started

# Flow Based Programming

There's 2 types of programming:

1. Where you control control flow
2. Where you don't

Most of the languages we use are belong to first.
I now only one of second kind, it's called FBP.
Imagine
