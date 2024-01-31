---
Title: About the language
---

Nevalang is a programming language with unique set of characteristics. Let's break down each one of them

## Main Characteristics

### Flow-Based (Dataflow)

In most _conventional_ languages programmer [controls execution flow](https://en.wikipedia.org/wiki/Control_flow). In [dataflow programming](https://en.wikipedia.org/wiki/Dataflow_programming) on the other hand programmer _controls data flow_. Basically it means that you don't have access to "low-level" instructions like "goto", "break" or even "return". Instead you can only connect inputs and outputs. And this changes everything we know about programming.

Dataflow programming is very different from the conventional. Some things are easier with conventional programming and some with dataflow. One thing for sure - dataflow programming is much more suitable for at least 2 things: parallel computations and visual programming.

The main thing about Nevalang is that it's [flow-based](https://en.wikipedia.org/wiki/Flow-based_programming). Flow-based programming is a subset of dataflow programming. It's not a classical FBP though.

### Visual

Nevalang is language with strong visual programming support by design. Combination of static typing and dataflow nature of the language it's easy to create visual representations of the program. Nevalang isn't just visual language though. It's a _hybrid_ language. While being easy to visualize it has simple C-like syntax for text-based programming. There are sertain cases like version control, merge request or hacking in REPL, where text is king.

Right now visual editor is a supporting tool while text is the main way to program. It is assumed though that in the future when visual tools will be good enough things will change. Visual programming will be the norm and text programming will be supported alterrnative for edge-cases. Nevalang has everything to accomplish that.

### General Purpose

Most visual programming languages are not general purpose. They are designed for specific purposes. E.g. "Scratch" is for educating, "Unreal blueprints" are for games, "Labview" for science, etc. There are a few _general purpose visual programming languages_ but none of them are popular as well as none of them are flow-based.

Nevalang on the other hand is general purpose visual flow-based language. It could be used for anything. It doesn't mean that it's perfectly suitable for anything though. For instance it probably should't be used for low-level stuff like drivers (because of the runtime overhead). But eventually user should be _able_ to do so.

### Statically Structurally Typed

Nevalang is statically typed which means that there is a compiler with static analyzer that ensures that the program satisfies the invariant. All types are known at compile time. Programmer can catch many bugs at compile-time.

Structural sub-typing means that e.g. if component wants object of form `{ age int }` as an input, then it's complitely file to send object of form `{ age int, name string }` to it. In this case `name` will simply be ignored. This allows to write much less "adater" code.

### Compiled

As been said there is a compiler. This means that source code is not interpreted by some kind of engine but instead there is a special entity that reads source code, analyzes and optimizes it and produces some lower-level code. Internally compiler has several intermidiate forms of the program but as the output (target) it can produce either byte-code (for VM), Go (for Go compiler) or executable binary machine code.

### Garbage-Collected

Used neither allocate and free memory manually not faces borrow-checker or other compile-time mechanism. Instead there is a garbage collector. This makes Nevalang [memory safe](https://en.wikipedia.org/wiki/Memory_safety) language.

## Extra Characteristics

We covered most important characteristics of the language but there is a couple more, a little bit less important but still unusual.

### Interpreted-mode

Nevalang supports interpreter-mode that is suitable for development and debugging purposes. There is a special binary that can read and execute source code. There is still static types and compiler, but under the hood. No need to compile and then execute seperately.
