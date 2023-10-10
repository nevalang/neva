# FAQ

Frequently asked questions

## General

### What is this?

This is Neva - visual general-purpose data-flow (flow-based) programming language with static structural typing and implicit parallelism that compiles to machine code.

### Why yet another language?

The goal is to create a system that is so powerful and easy to use at the same time, that even a person with zero programming skills can create effective and maintainable programs. Imagine what a professional could do with such a tool.

To achieve this we need 2 things: _visual programming_ and _implicit parallelism_. First will take maintainability to the next level by making diagrams first class source code and the second will effortlessly unlock the power of multi-core processors.

### Why FBP/Dataflow?

Conventional programming paradigms served us well by taking us so far but it's time to admit that they have failed at least at 2 things: visual programming and parallelism. Exactly 2 things we need to make a next-gen programming language.

#### Implicit parallelism

Any conventional program become more difficult when you add parallelism. As soon as you have more than one thread of execution bad things can happen such as deadlocks, race-conditions or memory-leaks. There are langauges that makes this simpler by introdusing concurrency primitives from dataflow world such as goroutines and channels in Go (CSP) or Erlang's processes (actor-model). However, it's still programmer's responsibility to manage those primitives. Concept of parallelism is simple, any adult understands it. But to make use of computer parallelism one must understand coroutines, channels, mutexes and atomics. We can do better.

#### Visual programming

...

First - text is also _visual_ representation (but using sounds or smells is by the way interesting idea). We recognize patterns by looking at code and parse the program's heirarchal structure with braces or offsets.

Second - argument that what we usually call visual programming is less maintanable is simply wrong. This is just different (also visual and more explicit one) form of representing a data of specific structure. Flow-based approach allowes to abstract things away exactly like text-based programming does.

Third - actually there's no dependency on visual programming. Neva designed with support for visual programming as a first-class citizen in mind but in fact it's possible to use text representation. Actually text is used as a source code for version control. There's also, by the way, no dependency on specific syntax.

Neva programs are, first of all, typed graphs that describes dataflows. The paradigm ([FBP]()) allowes to represent they in an infinite ways (including e.g. VR).

### Why Flow Based Programming and not OOP/FP/etc?

1. Higher level programming
2. Implicit concurrency
3. Easy to visualize

_Higher level programming_ means there's no: variables, functions, for loops, classes, methods, etc. All these things are low-level constructions that must be used under the hood as implementation details, but not as the API for the programmer. It's possible to have general purpose programming language with support for implicit concurrency and visual programming without using such things. Actually using of such low-level things is something that makes support for visual programming harder.

_Implicit concurrency_ means that programmer doesn't have to think about concurrency at all. At the same time any Neva program is concurrent by default. In fact there's no way to write non-concurrent programs. Explicit concurrency is like manual memory management - the great care must be put into. Concurrent programs in conventional langauges are always harder to maintain than regular ones. Not just all Neva programs are concurrent but programmer simply doesn't have a way to interact with concurrency. This is just how it works (thanks to FBP).

_Easy to visualize_ means that the nature of FBP programs is that we do not have [control flow](https://en.wikipedia.org/wiki/Control_flow), but instead we control [data flow](https://en.wikipedia.org/wiki/Dataflow_programming). This is how real electronic components works - there's electricity that flows through connections implementing specific logic. This is also how we document software with visual schemas - sort of boxes connected by arrowes where data flows from one component to another being transformed in someway. But those schemas are usually "dead" - they're not connected with the source code in anyway. FBP allowes to make diagrams source code itself.

## Design

### Why compiler has static ports and runtime has givers?

Because if compiler would have givers, they will be a special kind nodes which broke elegance of nodes being just component instances. Because giver is a regular component, it has a specific configuration - a message that it must distribute.

On the other hand, there's 2 types of effects at the runtime that are essentially different. Runtime anyway have concept of effects because if operators and giver is different than operator by the same reason.

### Why have static inports when we have const?

Static inports are actually syntactic sugar for `const`. If there wouldn't be static inports and only const then visual schemas would be complicated.

### Why have `fromRec`?

The reason is the same as with "static ports" vs "givers as special nodes". Otherwise there would be a special kind of nodes like "record builders" that are different from normal component nodes because they must have a specific configuration - record that they must build.

With `fromRec` feature (that is implemented outside of the typesystem, because type system doesn't know anything about ports) it's possible to say "hey compiler, I want a component with the same inports that this record has fields".

## Type-system

### Why structural subtyping?

1. It allowes write less code, especially mappings between records, vectors and maps of records
2. Nominal subtyping doesn't protect from mistake like passing wrong value to type-cast

### Why there's no Option/Maybe data type?

Because 1) there's `T | nil` (that's the reason shy `nil` is a type) 2) Having `Option` doesn't really introcice any advantage in FBP because we don't have objects and behaviour This means no methods for `Option` instances and thus no handy OOP mechanics for it.

## Implementation

### Why Go?

It's a perfect match. Go has builtin green threads, scheduler and garbage collector. Even more than that - it has goroutines and channels that are 1-1 mappings to FBP's ports and connections. Last but not least is that it's a pretty fast compiled language. Having Go as a compile target allowes to reuse its stdlib and increase performance by just updating compiler i.e. for free.

## FBP

### Is Neva "classical FBP"?

No. But it inherits so many ideas from it that it would be better to use word "FBP" than anything else. There's a [great article](https://jpaulm.github.io/fbp/fbp-inspired-vs-real-fbp.html) by mr. J. Paul Morrison (inventor of FBP) on this topic.

Now here's what makes Neva different from classical FBP:

- Neva has C-like syntax for its textual representation while FBP syntax is somewhat esoteric. It's important to node though that despite C-like syntax Neva programs are 100% declarative
- Neva doesn't let you program in "implementation-level" language like Go (similar to how Python doesn't let you program in assembly). FBP on the other hand forces you to program in langauges like Go or Java to implement elementary components.
- Neva introduces builtin observability via runtime interceptor and messages tracing, FBP has nothing like that
- Existing FBP implementations are essentially interpreters. Neva has both compiler and interpreter.
- Neva is statically typed while FBP isn't. FBP's idea is that you write code by hand in statically typed langauge like Go or Java and then use it in a non-typed FBP program, introducing runtime type-checks where needed
- Neva have _runtime functions_. In FBP there's just _elementary components_ that are written by programmer. Mr. Morrison did not like the idea of having "atomic" components like e.g. "numbers adder"
- Neva introduces hierarchical structure of program entities and package management similar to Go. Entities are packed into reusable packages and could be either public or private.
- Neva leverages existing Go's GC, FBP on the other hand introduces IP's life-cycle
- Neva's concurrency model runs on top of Go's scheduler which means it uses CSP as a lower-level fundament. FBP implementations on the other hand used to use shared state with mutex locks
- Neva has low-level program representation (LLR). FBP on the other hand doesn't describe anything like that

Also there's differences in naming:

- _Message_ instead of _IP (information package)_ not to be confused with "IP" as _internet protocol_
- _Node_ instead of _process_ 1) not to be confused with _OS processes_
- _Bound inports_ instead of _IIPs_ because of not using word _IP_
