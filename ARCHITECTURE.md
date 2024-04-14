# Architecture

This document only keeps visual schemas of the core language components. For details see [CONTRIBUTING.md](./CONTRIBUTING.md).

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
    program-->runtime

    subgraph runtime
        connector-->|msg|func-runner
        func-runner-->|msg|connector
    end

    subgraph connector
        event-listener
    end

    subgraph func-runner
        func-registry[(func-registry)]
    end
```

### Connector Algorithm

> WARNING: Algorithm has changed, update is needed.

```mermaid
flowchart TB
    cond{are there still connections?}
    cond -->|yes| broadcast[spawn broadcast goroutine]
    cond -->|no| exit[wait for all broadcast goroutines to finish]
    broadcast --> cond
```

#### Broadcast

```mermaid
flowchart TB
    msg[await new message from sender] --> inc[semaphore increment]
    inc --> distribute[spawn distribute goroutine]
    distribute --> |first receiver processed| msg
    distribute --> |all receivers processed| dec[semaphore decrement]
```

#### Distribute

```mermaid
flowchart TB
    q{is receivers queue empty?}
    q --> |yes| exit
    q --> |no| pick[pick receiver]
    pick --> try[try to send message to current receiver]
    try --> busy{is current receiver busy?}
    busy --> |yes| next[go to the next one]
    busy --> |no| remove[remove this receiver from queue]
    remove --> next
    next --> q
```

## Interpreter

```mermaid
flowchart LR
    source-code-->interpreter

    subgraph interpreter
        compiler-->|ir|adapter-->|program|runtime
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

    subgraph vscode
        webview-->extension
        extension-->webview
    end

    subgraph indexer
        compiler-frontend
    end
```
