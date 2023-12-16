# About The Language

> This document describes the final state of the language. Please note that at the time of writing the language is under heavy development.

Nevalang is a programming language with unique set of characteristics. Let's break down each one of them

## Main Characteristics

### Flow-Based

The main thing about Nevalang is that it's [flow-based](https://en.wikipedia.org/wiki/Flow-based_programming). Flow-based programming is a paradigm from the family of [dataflow](https://en.wikipedia.org/wiki/Dataflow_programming) programming paradigms. In dataflow programming, unlike the [control-flow](https://en.wikipedia.org/wiki/Control_flow) programming, we do not control execution flow. We control data flow instead. Basically this means that you don't have access to low-level instructions like "goto", "break" or even "return". Instead you can only connect inputs and outputs. And this changes everything we know about programming.

### Visual

Nevalang is language with strong visual programming support by design. Combination of static typing and dataflow nature of the language it's easy to create visual representations of the program. Nevalang isn't just visual language though. It's a _hybrid_ language. While being easy to visualize it has simple C-like syntax for text-based programming. There are sertain cases like version control, merge request or hacking in REPL, where text is king.

Right now visual editor is a supporting tool while text is the main way to program. It is assumed though that in the future when visual tools will be good enough things will change. Visual programming will be the norm and text programming will be supported alterrnative for edge-cases. Nevalang has everything to accomplish that.

### General Purpose

This means that Nevalang could be used to implement any kind of program. It doesn't mean that it's perfectly suitable for anything though. For instance it probably should't be used for low-level stuff like drivers because of the runtime overhead. But eventually user should be able to do so anyway.

### Statically Structurally Typed

Nevalang is statically typed which means that there is a compiler with static analyzer that ensures that the program satisfies the invariant. All types are known at compile time. Programmer can catch many bugs at compile-time.

Structural sub-typing means that e.g. if component wants object of form `{ age int }` as an input, then it's complitely file to send object of form `{ age int, name string }` to it. In this case `name` will simply be ignored. This allows to write much less "adater" code.

### Compiled

As been said there is a compiler. This means that source code is not interpreted by some kind of engine but instead there is a special entity that reads source code, analyzes and optimizes it and produces some lower-level code. Internally compiler has several intermidiate forms of the program but as the output (target) it can produce either byte-code (for VM), Go (for Go compiler) or executable binary machine code.

### Garbage-Collected

Used neither allocate and free memory manually not faces borrow-checker or other compile-time mechanism. Instead there is a garbage collector. This makes Nevalang [memory safe](https://en.wikipedia.org/wiki/Memory_safety) language.

## Extra Characteristics

We covered most important characteristics of the language but there is a couple more, a little bit less important but still unusual.

### VM-mode

As been said there is a byte-code which means there is a virtual machine that can execute it. That virtual machine is actually a very thing wrapper around the _runtime_.

### Interpreted-mode

Nevalang supports interpreter-mode that is suitable for development and debugging purposes. There is a special binary that can read and execute source code. There is still static types and compiler, but under the hood. No need to compile and then execute seperately.
