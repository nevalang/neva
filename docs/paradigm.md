Here we focus on it's paradigm - Dataflow.

## Dataflow

There are 2 high level programming paradigms: Dataflow and Controlflow. Paradigms like OOP or FP are specific variations of Controlflow while Actors or CSP are specific variations of Dataflow.

Dataflow programming is paradigm, that explains computation in terms of directed graphs. It takes many forms: OOP (by Alan Kay), CSP, Actor-model, FBP, etc. FBP is the one that influenced Nevalang the most. However, Nevalang is not original FBP in many ways. FBP, for instance, is not pure dataflow because you intended to write atomic components, which means using controlflow.

Common things for dataflow are nodes, connections and some kind of programmed flow for (probably async) message passing. There are popular general purpose controlflow languages that support some dataflow subset: Go's Goroutines and Channels and Erlang's Actors. Problem with these languages is that controlflow and dataflow are very different and combine them is harder than to combine different controlflow paradigms (e.g. procedural+functional). E.g. concept of variables is very foreign to dataflow.

On the other hand there are (more or less pure) dataflow languages that are not general purpose: Unreal Blueprint, Labview, etc. Some of them are pure, some of them are not. What about general purpose pure dataflow?

### Dataflow in Nevalang

Nevalang takes a next step in dataflow by making it the main paradigm. Controlflow is supported but is far from "first class citizen". All the native API is pure dataflow i.e. expressed in terms of pure dataflow abstractions.

There are low-level components that are could be implemented in controlflow under the hood, but the API exposed for Nevalang programmer is always pure dataflow.

Nevalang has such dataflow concepts as:

- Components and Nodes (instances)
- Ports (input and output) and (buffered) Connections between
- (Immutable) Messages

### Code is Code, Data is Data

In many controlflow langauges code and data exist in a single "space" of memory so we can pass functions around. However, there's no such thing in Nevalang. What is passed between components are messages, not other components.

### Static Dependency Injection

In Nevalang you can write polymorphic code by creating interfaces and using them as dependencies. Parent component will have to provide implementation. This is something very common for controlflow programmers who came from statically typed languages. But because _code is not data_ you can't do it dynamically. You can't "pass component" as a message. So all DI happens statically. Good news is that this is usually enough.

## Controlflow

Controlflow is the idea that there is some execution thread that follows some instructions and mutates some external state. It has many faces: from assembler to Haskell.

Why e.g. functional programming (pure/impure) is a controlflow? Core idea is the same, there'is a "thread" executes something. Either it's binary instructions or evaluation of a callstack - thread is "moving" through the program.

In Dataflow it's vice versa - threads are just "there". What is moving is data. E.g. Go's CSP is pure dataflow (without the rest of the language).

## Imperative vs Declarative

Not to be confused with imperative/declarative dilemma. Nevalang is declarative dataflow language, Haskell is declarative controlflow, C is controlflow imperative.

Imperative answers the question "how?", declarative "what?". E.g. Go's loop is expressed in terms of explicit controlflow instructions like `break`, `continue`, etc. Nevalang comes abstractions for message passing.

On the other hand Haskell wants you to define three of expressions that is somehow evaluated to a value. However, Haskell has a concept of "call/return", that are clearly controlflow instructions.

Imperative Dataflow is probably possible but dataflow is usually declarative. It's hard to say why, maybe the fact that computation is expressed in terms of graphs has something to do with this.

## Controlflow vs Dataflow

Now when we know what are these, let's learn what are their upsides and downsides.

Dataflow is easy to visualize. It's also good for parallelism. Obviously, good for message passing.

It's hard to enforce order though - things like "A can happen before B but it must happen only after". In controlflow this is not even a question: `a(); b()` - in most controlflow guaranteed "A then _always_ B".

Obviously controlflow has benefit from being dominated platform for thinking for a lot of time. Of course you will easily find lots and lots of materials, libraries, articles, videos, etc. Not so much for dataflow. However, all it means is only that we have to research and implement good dataflow.
