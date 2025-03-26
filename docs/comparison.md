# Comparison (Nevalang vs X)

## Neva vs Go

Neva is built on top of Go and it tries to borrow as much good things from there as possible, no surprise they have a lot in common: both statically typed, compiles to machine code, have garbage collector, builtin concurrency, both tries to be small and simple, and both aim for great developer experience with dependency management and other tooling.

The difference is that Go's paradigm is mixed - control-flow plus dataflow subset. Nevalang on the other hand is purely dataflow. It means that Go mostly operates in abstractions such as call/return, callstack, expression evaluation. It's actually expects you to control execution flow at the level of imperative instructions like `break`, `continue` and even `goto`. The dataflow subset is what's known as [CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes) implemented as goroutines and channels. CSP is indeed dataflow but it's usage is limited in Go and programs are mostly control-flow. Go also allows control-flow concurrency with mutexes and shared state. Concurrent Go code is usually considered as harder to reason about and more error prone, despite Go having state of the art concurrency support among all popular languages.

| **Feature**              | **Neva**                                                            | **Go**                                                                            |
| ------------------------ | ------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| **Paradigm**             | Pure Dataflow - nodes send and receive messages through connections | Mixed - control-flow (imperative) + dataflow subset (CSP)                         |
| **Concurrency**          | Defaults to concurrency. Requires explicit synchronicity            | Defaults to synchronicity. Requires explicit concurrency.                         |
| **Error Handling**       | Errors as values with `?` operator to avoid boilerplate             | Errors as values with `if err != nil {}` boilerplate                              |
| **Mutability**           | Immutable - no variables and pointers; data races are not possible  | Mutable - variables and pointers; programmer must avoid data races                |
| **Null Safety**          | Yes - nil pointer dereference is impossible                         | No - nil pointer dereference is possible                                          |
| **Zero Values**          | No zero values - everything must be explicitly initialized          | Zero values by default - everything can be initialized implicitly                 |
| **Subtyping**            | Structural - types are equal by their shape                         | Nominal - types are equal by their name                                           |
| **Traceback**            | Automatic - every message traces its path                           | Manual - programmer must explicitly wrap every error to add context               |
| **Dependency Injection** | Built-in - any component with dependency expects injection          | Manual - programmer must create constructor function that takes dependencies      |
| **Stream Processing**    | Native support with components like `Map/Filter/Reduce`             | Programmer must manually implement dataflow patterns with goroutines and channels |
| **Visual Programming**   | Aims for hybrid programming with Visual Editor (WIP)                | Textual language, very little visual tooling support                              |

## Neva vs Erlang/Elixir

People often compare Neva to BEAM langauges because of message-passing, immutability, concurrency and stream-processing:

| **Feature**            | **Neva**                                                            | **Erlang/Elixir**                                             |
| ---------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------- |
| **Paradigm**           | Pure Dataflow - no functions, no call-stack; Just message-passing   | Mixed - control-flow (FP) with dataflow subset (actors)       |
| **Message Passing**    | Static connections defined at compile time                          | Dynamic message passing to PIDs                               |
| **Type-System**        | Strongly Typed - types are always required                          | Dynamic (Erlang) / Gradually-typed (Elixir)                   |
| **Execution Model**    | Compiles to machine code and can be deployed as a single executable | Needs Virtual-Machine (BEAM) to be installed on the server    |
| **Syntax**             | C-like syntax with curly-braces                                     | Esoteric (Erlang) / Ruby-like (Elixir)                        |
| **Error Tolerance**    | Everything must be type-safe; Errors must be explicitly handled     | "Let it crash"                                                |
| **Visual Programming** | Aims for hybrid programming with Visual Editor (WIP)                | Textual language, very little visual tooling support          |
| **Interop**            | Aims for interopability with Golang (WIP)                           | Interopable with BEAM-compatible family (Erlang/Elixir/Gleam) |

## Neva vs Gleam

BEAM family includes language that is even closer to Neva because of static types. However, there are differences:

| **Feature**            | **Neva**                                                            | **Gleam**                                                     |
| ---------------------- | ------------------------------------------------------------------- | ------------------------------------------------------------- |
| **Paradigm**           | Pure Dataflow - no functions, no call-stack; Just message-passing   | Mixed - control-flow (FP) with dataflow subset (actors)       |
| **Subtyping**          | Structural - types are equal by their shape                         | Nominal - types are equal by their name                       |
| **Execution Model**    | Compiles to machine code and can be deployed as a single executable | Needs Virtual-Machine (BEAM) to be installed on the server    |
| **Interop**            | Aims for interopability with Golang (WIP)                           | Interopable with BEAM-compatible family (Erlang/Elixir/Gleam) |
| **Visual Programming** | Aims for hybrid programming with Visual Editor (WIP)                | Textual language, very little visual tooling support          |
