# Nevalang

> On shore abandoned, kissed by wave, he stood, of mighty thoughts the slave.

Nevalang is general purpose [dataflow](https://en.wikipedia.org/wiki/Dataflow_programming) ([flow-based](https://en.wikipedia.org/wiki/Flow-based_programming)) programming language with static [structural](https://en.wikipedia.org/wiki/Structural_type_system) type-system and [implicit parallelism](https://en.wikipedia.org/wiki/Implicit_parallelism) that compiles to machine code.

And also [visual programming](https://en.wikipedia.org/wiki/Visual_programming_language) done right!

## Safety

**If it runs - it works**.

Thanks to dataflow paradigm, where programmer doesn't interact with low-level promitives like variables, pointers or coroutines, and static type-system. If compiler accepted the program and the program started successfully, then there's a guarantee that there will be no exceptional situations including:

- [Race conditions](https://en.wikipedia.org/wiki/Race_condition)
- [Deadlocks](https://en.wikipedia.org/wiki/Deadlock)
- [Type errors](https://en.wikipedia.org/wiki/Type_system#Type_errors)
- [Null pointer](https://en.wikipedia.org/wiki/Null_pointer) aka [The Billion Dollar Mistake](https://www.infoq.com/presentations/Null-References-The-Billion-Dollar-Mistake-Tony-Hoare/)
- [Uninitialized variables](https://en.wikipedia.org/wiki/Uninitialized_variable)
- [Off-by-one](https://en.wikipedia.org/wiki/Off-by-one_error) and [indexing errors](https://en.wikipedia.org/wiki/Bounds_checking#Index_checking) in general
- [Stack overflow](https://en.wikipedia.org/wiki/Stack_overflow)
- [And](https://en.wikipedia.org/wiki/Dangling_pointer) [many](https://en.wikipedia.org/wiki/Buffer_overflow) [many](https://en.wikipedia.org/wiki/Segmentation_fault) [more](https://en.wikipedia.org/wiki/Stale_pointer_bug)

Any runtime error is threated as a compiler bug.

## Performance

### Implicit parallelism

FBP processes don't share memory and thus can run in parallel. Everything that can happen at the same time - will. Programmer doesn't intereact with threads, mutexes, coroutines or channels. Not only it eliminate concurrency-related bugs but also forces the program to utilize all CPU cores.

### As fast as Go

[Go](https://go.dev) is used as a low-level IR due to several reasons.

1. Perfect match. FBP's processes and ports maps 1-1 to Go's goroutines and channels (because [CSP](https://en.wikipedia.org/wiki/Communicating_sequential_processes), the formal model that Go based on, is a form of Dataflow programming). This makes Go perfect choice for low-level IR for FBP language.
2. State of the art [standard library](https://pkg.go.dev/std), coroutine scheduler, garbage collector and crossplatform machine code generation backed by huge community. Nevalang will become faster and safer by simply updating the underlying Go compiler.
3. One of the fastest compilers in the world. Compilation speed is [design](https://www.youtube.com/watch?v=rKnDgT73v8s#t=8m53) goal.

From some point Nevalang could be viewed as a Go code generator.

## Productivity

### Implicit parallelism

Implicit parallelism makes concurrent programming as simple as a regular one. Programmers used to think that concurrent programs are harder to reason about and thus test and maintain but that's not true in FBP.

### Static analysis

Thanks to graph-like nature of FBP programs and static types compiler can catch most of the possible errors. And everything compiler can't catch is checked at program's startup. So there's just the actual program's logic that programmer have to think about.

### Visual programming

Last but not least. FBP is a perfect paradigm for visual programming because FBP programs are literally computational graphs - processes connected to each other through input and output ports. Making visual representation of a FBP program is simply rendering it's structure.

There's a big problem with visual representations of a software written in conventional langauge - they're dead. There's no connection with the real code. But FBP schema _is_ the code. So the most productive way to work with Nevalang must be by working with visual schemas. Move blocks around and wire stuff up.

It doesn't mean it's the only way to program in Nevaland. Think of visual editor as a source code generator. Generated code must be completely human-readable because we still need to review it. There probably will be usecases for hand-written code (e.g. REPL). Or maybe you just an old fashioned hacker.
