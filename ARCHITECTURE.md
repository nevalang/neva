# Architecture

## Compiler

```mermaid
flowchart LR
    package-manager -->|raw-build| compiler

    subgraph package-manager
        git-client
        file-system
    end

    subgraph compiler
        subgraph backend
            go-code-generator
        end

        parser -->|parsed-build| analyzer

        subgraph parser
            antlr
        end

        analyzer -->|analyzed-build| desugarer

        subgraph analyzer
            typesystem
        end

        desugarer -->|desugared-build| irgen

        irgen -->|ir| backend
    end

    compiler -->|go-code| go-compiler
```

## Runtime

```mermaid
flowchart LR
    program-->runtime-->side-effects

    subgraph runtime
        func-runner
    end

    subgraph func-runner
        func-registry[(func-registry)]
    end
```

## VSCode Extension

```mermaid
flowchart LR
    language-server-->|jsonrpc|vscode
    vscode-->|jsonrpc|language-server

    subgraph language-server
        indexer
    end

    subgraph indexer
        compiler-frontend
    end
```
