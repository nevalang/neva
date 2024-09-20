# About

Nevalang is a general purpose programming language: documentation, compiler, runtime, cli, editor plugins, etc. What kind of language it is?

- Pure Dataflow (FBP-like)
- Compiled (Machine code, Golang, WASM)
- Statically (Strongly) Typed
- Concurrent/Parallel by default
- Hybrid (textual/visual)
- Pure Declarative
- Go Interop (call Neva from Go and vice versa)

We will discover what this means what piece by piece.

## Pure Dataflow

Pure Dataflow means that your program is expressed in terms of dataflow abstractions such as nodes, ports, connections and messages.

There are no variables, no functions, no classes, no objects, no methods, no loops, no ifs, no switches, no exceptions, no mutable state.

Nevalang's dataflow mostly influenced by Flow-Based Programming paradigm but Nevalang is not pure FBP. You can find more info about this in specialized section.

## Compilation

Nevalang uses Go as low-level IR (Intermediate Representation). This means your Nevalang program is firstly compiled to Go and then Go compiler is used to compile it to any target that Go compiler supports, such as machine code or WASM.

This firstly has performance benifits and secondly it's easier to deploy if all you need to ship is a binary, unlike Python or Java where you need interpreter/VM of a specific (correct) version.

There's a temptation to say that Nevalang is as fast as Go but that statement is incorrect because Nevalang is higher abstraction. It has it's own (very-very thin, but still) runtime on top of Go's.

## Strong Static Typing

Nevalang has _strong static structural_ type-system. What does it mean?

- _Strong_ - means there are no implicit type-casts. E.g. `float` will never be implicitly converted to `int` and you will always have to convert (cast) them explicitly, by using special components
- _Static_ - means that types are known at compile time. E.g. - you can (and you must) explicitly define what data type your port is sending or receiving, or what data-type some constant have
- _Structural_ - means _structural sub-typing_ is used. In this case it used for data-types and interfaces/components implementation relation.

If you are confused by something, don't worry. Type-system will be explained in detail on related section.

## Implicit Parallelism

Nevalang supports implicit parallelism. It means that you don't have to (and you can't) access parallelism/concurrency primitives such as mutexes, threads, coroutines, channels, etc. You just create dataflow and everything that can happen in parallel (better word is _concurrently_) - will happen this way. It's good because of 2 things: performance and maintainance. You basically get parallelism for free.

If you want to know how exactly Nevalang executes your program then don't worry, there will be a special section about that. This is just a high level overview of a specific feature.

## Hybrid (Textual and Visual Programming)

Textual and visual programming in Nevalang are combined into single workflow. You can use only text if you e.g. work in environment that doesn't support visuals or you just don't like visual programming. However, Nevalang is and always will be hybrid langauge, which means visual programming is a first class citizen. If some feature cannot be supported in visual programming it might not be added or can be removed. Even language abstractions are done in a way so it's possible to visually represent them.

## Pure Declarative

Nevalang doesn't have a single imperative bit except go-interop (including both writing extensions and maintaining stdlib).

Declarative Programming means you write your program by describing what do you want to get as a result, instead of explaining how to actually do it under the hood. E.g. explicit controlflow instructions like `for`, `break` or `goto` are imperative while things like `[1,2,3].map(double)` are declarative. You are not explaining _how_ to actually `map` stuff, you just say that this is _what_ result should be.

In Nevalang declarative programming is expressed in a way you write your program - you define entities. Computation is expressed in terms of graph that is also declarative. Network show what should happen (e.g. message from X should be moved to Y) and not how (send message to channel, block coroutine, etc).

Of course there are different problems and imperative programming is better for some of them. However, if we have imperative languages, why not have declarative? Besides, there are tons of declarative languages, they are just (mostly) not dataflow. Let's take functional programming paradigm for instance. If we are talking about pure FP then we are automatically talking about pure declarative. Same goes for dataflow - pure dataflow is always pure declarative.

## Go Interop

As being said Nevalang uses Go as low-level IR which opens the door for two-way interop with existing Go codebases. That is, Nevalang has potential to call Go code from neva code and vice versa - call Nevalang from Go.

This is very important thing about Nevalang - it's intended for gradual adaptation among Go developers. Otherwise it would be close to impossible to adapt the language to such a highly competitive market.

However, if you are not Go developer, don't worry. Nevalang won't have interop with your language so you will have to use network for communication. However, it's completely fine to use Nevalang without knowing anything about Go. Nothing requires you Go knowledge to write Nevalang programs. Maybe couple of things might not seem obvious without Go background or as a non-Go developer you couldn't appreciate some decisions. That doesn't matter. Nevalang is separate language. Nevalang to Go is more what C to assembler is, rather than TypeScript to JavaScript.
