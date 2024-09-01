# Motivation

There are thousands of programming languages, why need another one? This document describes the motivation behind Nevalang in detail.

Before we dig into the details, here's super high-level idea: it's possible to create a better language, where programs are easier to understand, debug and create. Unique set of features properly combined together could lead to such language. Nevalang is nothing but an attempt to create such language.

## General Purpose Dataflow

Almost all general purpose languages are controlflow. Dataflow languages on the other hand are usually domain specific. There are general purpose dataflow but very few of them and they are either not pure (you have to write controlflow) or not "good enough" (i.e. only works with JavaScript or has no static types).

## General Purpose Visual Programming

It's a shame that we don't have at least one visual programming language in out TOP-10. Imagine having something as good as Java or Python but visual.

Turns out dataflow and visual programming are connected. Controlflow is hard to visualize so we need good general purpose dataflow for strong visual programming.

## Text/Visual Hybrid

One of the reasons visual programming haven't succeeded is the fact that we've created a ecosystem that bound to text. E.g. version control systems like Git are text-based. Also, you need to be able to do code review.

Some visual languages store code in JSON, Yaml or XML. However, this is not easy to read and write. Syntax is important, and that's why Nevalang has it's own minimalistic C-like syntax. Visual representation fully reflect this syntax.

## Implicit Parallelism

There are not many implicitly parallel languages - languages where you don't have to (or even more than that, you are not allowed) to think in terms of threads, locks, channels, etc. You "just write your code" and "system" somehow automatically figures out how to insert parallelism where possible.

Dataflow is the perfect model for that. You just have nodes and they, because there's no shared state, can operate completely parallel. And this is exactly how Nevalang works.

Imagine a pipeline of 3 nodes. There's input data that goes to node A, then from A to B and finally from B to C. Then imagine infinite stream of data going to A. In Nevalang all A, B and C can work at the same time. This is an example of implicit parallelism without explicit parallel connections in network topology. And of course you can have literally (and visually) parallel connections to nodes that process data in parallel. You can use technics like round-robin to handle highload at the level of your executable binary.

## Better Language

It's hard to explain the core motivation without telling that I always felt a need for a better language. How great would it be to collect all the best working technics into one perfect language? Besides all the listed features, I had a very specific opinion on topics such as type-system, error handling, syntax, etc. In Nevalang I tried to do my best.

Core motivation for creating Nevalang is belief that it's possible to create a more optimized tool for our daily tasks. Do more with less effort. There were always questions like "what is the best way to do X?". All the things like dataflow basically consequence from that.

### Type-System

It's very hard to find a perfect balance for a type-system (and semantic analysis in general) so you don't "fight" it, but compiler helps you instead. In my personal opinion Go is very close to that, but it has too much data-types and it uses nominative sub-typing for them. On the other hand TypeScript uses structure sub-typing but it has `any`. Go also has but in Go any is a type you can't do anything with. In TS you _can_ do anything with any. Also in TS you usually interop with unsound JS ecosystem. Also TS is too complex, its type-system literally turing complete.

Nevalang tries to find balance between complexity and reliability. Its type-system looks like something between TypeScript and Go, probably closer to Go, but with just a little bit more strictness. It inherits structural sub-typing for data-types from TypeScript though. Sub-typing for interfaces is structural just like in Go. It also supports generics with constraints.

### Error Handling

Nevalang follows Go's "errors are values" mantra. In Nevalang error is just a data-type for a message. A structure (of a very simple shape), to be more specific. Just like Go's functions can return several values and one of them (last usually) can have type `error` (and the variable that store this value usually called `err`), Nevalang's components can have `err error` outport. In Go (and any controlflow) you check error (using `if` or `try-catch` sequence) and either handle or throw/return. And Nevalang you _send_ error somewhere. It could be panic component, logger, or maybe your own `error` outport. This is analog for throwing. Maybe you want to add some context before that, that's also easily possible.

Rust's approach to errors is even better than Go's. It has `Result` type that's either error or result. This is not very much different than having several return values but Rust also have special operator that allows you to "ignore" error-handling while still handling error. That's a syntax sugar that implicitly inserts error-handling into your code everywhere you use this special operator. Nevalang also have this, it's called `?` operator.

Finally, Go _does_ allow you to ignore errors. In Nevalang you always have to use `error` outport of your sub-nodes (and your own) if it's present.

### Tracing

Each message has its own path that is always updated whenever it moves from one place to another. This is something like a stacktrace that you get with exception in languages with exceptions like Python or JavaScript, but not just for exceptional situations, but for literally every message in the program.

Go doesn't have stacktraces for errors (only for `panic`) and that leads to abuse of `fmt.Errorf` with obscure directives like `%w` that you have to memorize. It's good that nowadays LLMs can more or less automate that but it's not perfect. In Nevalang you don't have to do anything to get a full trace of the message, from the place it was created up to the point we log the trace.

It makes error handling even simpler, most of the time you should be ok with nothing but just `?`. However, it's not only for error handling. You can track _any_ message this way. Any! Imagine what debugging possibilities in opens.

### Next-Gen Debugger

Dataflow + tracing opens the door for next-generation debugging tool where you can visually set breakpoints on a specific connections in your network graph and stop the program when message arrives there. You can then observe the message or even update its value.
