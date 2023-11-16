# ARCHITECTURE

## Production

```mermaid
flowchart LR
    builder-->|raw-program|compiler
    subgraph builder
        git-client
    end
    subgraph compiler
        parser-->|program|analyzer
        subgraph analyzer
            typesystem
        end
        analyzer-->|analyzed-program|irgen
        irgen-->|ir|decoder
    end
    compiler-->|protobuf|VM
    subgraph VM
        loader-->|ir|runtime
        runtime
        subgraph runtime
            connector
            func-runner
        end
    end
```

## Development

```mermaid
flowchart LR
    subgraph interpreter
        builder-->compiler
        compiler-->runtime
    end
```

## VSCode Extension

See [web/ARCHITECTURE.md](./web/ARCHITECTURE.md)