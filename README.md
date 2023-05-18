<!-- # Neva Programming Language -->

![assets/header.png](assets/header.png)

> On shore abandoned, kissed by wave, he stood, of mighty thoughts the slave.

Neva is a general purpose [dataflow](https://en.wikipedia.org/wiki/Dataflow_programming) ([flow-based](https://en.wikipedia.org/wiki/Flow-based_programming)) programming language with static [structural](https://en.wikipedia.org/wiki/Structural_type_system) typing and [implicit parallelism](https://en.wikipedia.org/wiki/Implicit_parallelism) that compiles to machine code.

Oh, and a [visual programming](https://en.wikipedia.org/wiki/Visual_programming_language) finally done right!

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

### Complete asynchrony

Data gets processed as soon as it arrives. Concider this pseudocode:

```
sumRoots(
    getRootA(),
    getRootB(),
)
```

Both `getRootA()` and `getRootB()` return numbers, `sumRoots` find square roots for those numbers and returns their sum. In conventional programming we have to wait for `getRootA`, then for `getRootB` and then for `sumRoots`. In FBP all three `sumRoots`, `getRootA` and `getRootB` runs concurrently. No matter which finishes first, as soon as first number arrive, `sumRoots` start executing. At the time second number arrive, the first square root could already be calculated. Imagine the whole program work this way.

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

Thanks to graph-like nature of FBP programs and static types compiler can catch most of the possible errors. And everything compiler can't catch is checked at program's startup. So there's just the actual program's logic that programmer have to think about. Many situations like "error was handled like a valid result" or "class was initialized with null dependencies" are simply impossible.

### Interpreter mode

Iterate fast by using interpreter instead of compiler for development and see changes in realtime. Nevalang architecture allows to use the exact same code analyzer and runtime without the need to generate executable. So you can develop like it's Node.js with TS out of the box. Then compile it to fast executable when it's ready for production.

### Multiplayer mode

Nevalang comes with development server that allowes several developers at the same time to connect and modify the same program. Everyone can see changes in realtime.

### Observability out of the box

Message interceptor is built into the runtime and every message has trace so it's extremely easy to keep track of how data is moving across the program. Same goes for errors - you can see where error arose and what way did it went. You can debug probram by setting breakpoints on the graph connections and intercept messages. You can even substitute them on the fly to see what will happen.

### Visual programming

Last but not least. FBP is a perfect paradigm for visual programming because FBP programs are literally computational graphs - processes connected to each other through input and output ports. Making visual representation of a FBP program is simply rendering it's structure.

There's a big problem with visual representations of a software written in conventional langauge - they're dead. There's no connection with the real code. But FBP schema _is_ the code. So the most productive way to work with Nevalang must be by working with visual schemas. Move blocks around and wire stuff up.

It doesn't mean it's the only way to program in Nevaland. Think of visual editor as a source code generator. Generated code must be completely human-readable because we still need to review it. There probably will be usecases for hand-written code (e.g. REPL). Or maybe you just an old fashioned hacker.
