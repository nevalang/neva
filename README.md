# Nevalang

> On shore abandoned, kissed by wave
>
> He stood, of mighty thoughts the slave

Nevalang is a general purpose [flow-based](https://en.wikipedia.org/wiki/Flow-based_programming) programming language with static [structural](https://en.wikipedia.org/wiki/Structural_type_system) type-system and [implicit parallelism](https://en.wikipedia.org/wiki/Implicit_parallelism) that compiles to machine code. And also [visual programming](https://en.wikipedia.org/wiki/Visual_programming_language) done right!

## Safety

**If it runs - it works**.

This is possible because of flow-based paradigm (programmer doesn't interact with low-level promitives like variables, pointers or coroutines) and static type-system. If compiler accepted the program and the program started successfully, then there's a guarantee that there will be no exceptional situations including:

- Data-races (and mutations in general)
- Deadlocks
- Null pointer dereferences

Any runtime error must be threated as a bug in compiler.

## Performance

### Implicit parallelism

In FBP processes don't share memory and thus can run in parallel. Everything that can happen at the same time - will. Programmer doesn't have to worry about things like threads, mutexes, coroutines and channels. It not only eliminate concurrency-related bugs but also forces the program to utilize all CPU cores.

### As fast as Go

Nevalang uses [Go](https://go.dev) as a compile-target due to several reasons. Firstly because it's a perfect match - FBP's processes and ports maps 1-1 to Go's goroutines and channels. Secondly - it's always good to stand on a giant's shoulders. It would be impossible to implement and maintain from scratch something like Go compiler backed go Google and its huge community, with its state of the art standard library. Ability to reuse this ecosystem is a great gift - Nevalang will become faster by simply updating the underlying Go compiler.

## Productivity

### Implicit parallelism

Implicit parallelism makes concurrent programming as simple as a regular one. Programmers used to think that concurrent programs are harder to reason about and thus test and maintain but that's not true in FBP.

### Static analysis

Thanks to graph-like nature of FBP programs and static types compiler can catch most of the possible errors. And everything compiler can't catch is checked at program's startup. So there's just the actual program's logic that programmer have to think about.

### Visual programming

Last but not least. FBP is a perfect paradigm for visual programming because FBP programs are literally graphs of processes connected to each other through input and output ports. So making visual representation of a FBP program is literally rendering it's structure.

There's a big problem with visual representation of a software written in conventional langauge - those representation are dead, there's no connection with the real code. But FBP schema _is_ the code. So the most productive way to work with Nevalang must be by interacting with visual schemas. Move blocks around and wire stuff up.

It doesn't mean it's the only way to program in Nevaland. Think of visual editor as a code generator. Generated code must be completely human-readable because we still need to review it. There's probably will be a cases for hand-written code like in REPL. Or maybe you just old fashioned.