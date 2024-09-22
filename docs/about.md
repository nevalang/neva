# About

Nevalang is a general-purpose programming language with the following key characteristics:

- **Pure Dataflow**: Programs are expressed using dataflow abstractions like nodes, ports, connections, and messages. Traditional programming constructs (variables, functions, loops, etc.) are absent.
- **Implicit Parallelism**: Parallelism is built-in, allowing concurrent execution without explicit primitives like mutexes or threads.
- **Compilation**: Nevalang compiles to Go as an intermediate representation, then to machine code or WASM, offering performance benefits and easier deployment.
- **Strong Static Typing**: The language employs a strong, static, and structural type system, ensuring type safety at compile-time.
- **Hybrid Programming (WIP)**: Nevalang supports both textual and visual programming as first-class citizens, allowing flexibility in development approaches.
- **Go Interoperability (WIP)**: Nevalang can interact with existing Go codebases, allowing for gradual adoption and integration with Go projects.

These features combine to create a unique language that emphasizes dataflow, type safety, and declarative programming while leveraging Go's ecosystem and compilation targets. Nevalang aims to provide a fresh approach to programming, particularly suited for concurrent and parallel applications.

## Motivation

Nevalang aims to create a better programming language where programs are easier to understand, debug, and create. Here are the key motivations behind its development:

### General Purpose Dataflow

Most general-purpose languages are controlflow-based, while dataflow languages are often domain-specific. Nevalang seeks to bridge this gap by offering a pure, general-purpose dataflow language with static typing.

### Implicit Parallelism

Nevalang introduces _implicit parallelism_, automating concurrency control like garbage collection automated memory management. Key features:

- Components exchange messages via buffered queues
- Blocking limited to specific parts, allowing others to continue
- First-class stream processing for efficient data handling
- No low-level concurrency primitives, enhancing safety

This approach makes concurrent programming the default, simplifying development.

### Text/Visual Hybrid (WIP)

Despite humans' visual nature, text-based coding has dominated for decades. A general-purpose visual programming language is needed, as existing visual languages lack the popularity of text-based ones.

Visual programming's necessity is evident in software design diagrams and the popularity of visual tools. People naturally think in interconnected processes, yet static diagrams quickly become outdated.

The claim that visual programming is less maintainable is unfounded. Nevalang offer abstraction similar to text-based programming through modules, packages and components.

Nevalang aims to combine textual and visual programming in the future, addressing the limitations of purely visual languages. The plan is to implement a minimalistic C-like syntax that can be fully represented visually, making it compatible with text-based tools like version control systems.

### Optimized Type System

Nevalang aims for a balanced type system that aids developers without being overly restrictive. It combines elements from Go and TypeScript, featuring structural subtyping for data types and interfaces, along with generics with constraints.

#### Structural Subtyping

Unlike Go’s nominative subtyping, Nevalang checks type compatibility based on structure, not name, avoiding unnecessary casting. In Go, a `readBook` function taking a `Book` struct won’t accept a `Magazine` struct, even though `Magazine` has the necessary fields.

```go
type Book struct { title string, author string }
type Magazine struct { title string, author string, number int }
func main() {
    readBook(Magazine{title: "FBP time", author: "Emil Valeev"}) // compile error!
}
func readBook(book Book) { fmt.Println(book) }
```

Go requires explicit casting:

```go
readBook(Book{
    title: magazine.title,
    author: magazine.author,
})
```

In Nevalang, structural typing eliminates this problem. For example, in web apps using API boundaries like gRPC or GraphQL, where developers often write mapping code between generated types and domain models, structural typing automatically handles compatibility, avoiding hundreds of lines of manual conversions.

### Better Debugging

#### Improved Error Handling

Following Go's "errors are values" approach, Nevalang treats errors as data types. It incorporates Rust-like error handling with a `?` operator, while ensuring that errors are always handled when present.

#### Advanced Tracing

Every message in Nevalang has a path that updates as it moves through the program. This provides comprehensive tracing capabilities, similar to stack traces in exception-based languages, but for all messages, not just errors.

#### Next-Generation Debugging (WIP)

The combination of dataflow architecture and advanced tracing enables powerful debugging tools. Developers can set visual breakpoints on specific connections in the network graph, observe messages, and even update their values during runtime.

By combining these features, Nevalang strives to offer a more efficient and intuitive programming experience, pushing the boundaries of what's possible in language design.

## Paradigm

Nevalang adopts the dataflow paradigm, which explains computation in terms of directed graphs. While influenced by Flow-Based Programming (FBP), Nevalang diverges from it in several ways.

### Dataflow vs Controlflow

Two high-level programming paradigms exist: Dataflow and Controlflow. Dataflow includes variations like Actors and CSP, while Controlflow encompasses OOP and FP.

Dataflow programming typically involves nodes, connections, and asynchronous message passing. Some controlflow languages (e.g., Go, Erlang) support dataflow subsets, but combining the two paradigms can be challenging.

### Nevalang's Approach

Nevalang makes dataflow its primary paradigm, with controlflow support as a secondary feature. Key dataflow concepts in Nevalang include:

- Components and Nodes
- Ports (input/output) and Connections
- Immutable Messages

Unlike many controlflow languages, Nevalang strictly separates code and data. Messages, not components, are passed between components.

### Dependency Injection

Nevalang supports polymorphism through interfaces and static dependency injection. While this approach is familiar to controlflow programmers from statically typed languages, Nevalang's implementation is purely static due to the separation of code and data.

### Declarative Nature

Nevalang is a declarative dataflow language, focusing on "what" rather than "how". This approach contrasts with imperative languages and even declarative controlflow languages like Haskell, which use controlflow concepts like call/return.

### Advantages and Challenges

Dataflow excels in visualization and parallelism but can struggle with enforcing strict order of operations. While controlflow benefits from its long-standing dominance and extensive resources, dataflow presents opportunities for research and innovation in programming paradigms.

### Flow-Based Programming

> Feel free to skip this section if you are not familiar with the original concept of fbp.

Nevalang's dataflow is influenced by FBP but differs in several key aspects:

#### Granularity and Controlflow

Nevalang is designed for general-purpose programming, expecting entire programs to be written in dataflow, unlike FBP which is used for high-level orchestration. Nevalang provides low-level components (written in Go) for operations like math, eliminating the need for users to write controlflow code except when contributing to the stdlib or integrating with Go.

#### Garbage Collection

Nevalang is garbage-collected with immutable data, avoiding ownership concepts. This prevents data races but may impact performance. Mutations are possible via unsafe packages (WIP) but are discouraged. FBP, in contrast, uses ownership and allows mutations.

#### Node Behavior

Nevalang's nodes are always running, automatically starting, suspending, and restarting as needed. FBP processes have explicit states (start, suspend, restart, shutdown) that can be manipulated.

#### Static Typing

Nevalang features a static type system with generics and structural sub-typing, improving IDE support and reducing runtime validations. FBP is dynamically typed in its dataflow part.

#### Similarities

Both paradigms support dataflow and implicit parallelism, sharing much terminology.

#### Terminology Comparison

| FBP                              | Nevalang                  |
| -------------------------------- | ------------------------- |
| Component (Atomic/Complementary) | Component (Native/Normal) |
| Process                          | Node                      |
| Connection                       | Connection                |
| Ports                            | Ports                     |
| IP                               | Message                   |
| IIP                              | Constant                  |
| IP Tree                          | Structure or Dictionary   |
| Sub-Stream                       | Stream                    |

---

- [Back: Table of Contents](./README.md#table-of-contents)
- [Next: Program Structure](./program_structure.md)
