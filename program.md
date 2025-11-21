# Program: hello_world

**Compiler:** 0.32.0

## 1. Visual Flow

```mermaid
flowchart TD
    %% Style definitions
    classDef port fill:#fff,stroke:#000,stroke-width:1px,min-width:0,padding:2px,rx:5,ry:5;
    classDef node fill:#f5f5f5,stroke:#333,stroke-width:1px,color:#333;
    classDef invisible display:none;
    subgraph __newv2__13 ["__newv2__13"]
        direction TB
        __newv2__13__sig(sig):::port
        __newv2__13__res(res):::port
    end
    class __newv2__13 node

    subgraph in ["in"]
        direction TB
        in__start(start):::port
    end
    class in node

    subgraph out ["out"]
        direction TB
        out__stop(stop):::port
    end
    class out node

    subgraph panic ["panic"]
        direction TB
        panic__data(data):::port
    end
    class panic node

    subgraph println_renamed ["println_renamed"]
        direction TB
        println_renamed__data(data):::port
        println_renamed__err(err):::port
        println_renamed__res(res):::port
    end
    class println_renamed node

    __newv2__13__res --> println_renamed__data
    in__start --> __newv2__13__sig
    println_renamed__err --> panic__data
    println_renamed__res --> out__stop
```
## 2. Components

| Node | Ref | Config | Ports |
| :--- | :--- | :--- | :--- |
| `__newv2__13` | `new_v2` | `"Hello, World!"` | `in:sig, out:res` |
| `in` | - | - | `out:start` |
| `out` | - | - | `in:stop` |
| `panic` | `panic` | - | `in:data` |
| `println_renamed` | `println` | - | `in:data, out:err, out:res` |

## 3. Metrics
* **Nodes:** 5
* **Connections:** 4
