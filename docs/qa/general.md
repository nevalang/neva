# General

## What is this?

This is **Neva** - first visual general-purpose flow-based programming language with static structural typing and implicit parallelism.

## Why yet another language?

The goal is to create a system that is so powerful and easy to use at the same time, that even a person with zero programming skills can create effective and maintainable software with it. Imagine what a professional programmer could do with this tool.

To achieve this we need many things. 2 of them can highlighted among them all:

1. Visual programming - natural way of thinking
2. Implicit parallelism - elimination of manual concurrency management

## Why not yet another language?

Conventional programming paradigms served us well by taking us so far but it's time to admit that they have failed at visual programming and that parallelism is usually hard to implement right with them. Also dataflow is just what the things are in real world. This is the natural way of thinking about computation.

## No manual concurrency management

Any conventional program become more difficult when you add parallelism. As soon as you have more than one thread, bad things can happen - deadlocks, race-conditions, you name it. There are languages that makes this simpler by introdusing concurrency primitives from dataflow world such as goroutines and channels in Go (CSP) or Erlang's processes (actor-model). However, it's still programmer's responsibility to manage those primitives. Concept of parallelism is simple, any adult understands it. But to make use of computer parallelism one must understand coroutines, channels, mutexes and atomics.

## Visual programming

The argument that visual programming is less maintanable is wrong. This is just different form of representing a data. Flow-based approach allowes to abstract things away exactly like we used to with text-based programming.

Actually there's no dependency on visual programming. Neva designed with support for visual programming in mind but in fact it's possible to use text representation.

## Why FBP and not OOP/FP/etc?

1. Higher level programming
2. Implicit concurrency
3. Easy to visualize

_Higher level programming_ means there's no: variables, functions, for loops, classes, methods, etc. All these things are low-level constructions that must be used under the hood as implementation details, but not as the API for the programmer. It's possible to have general purpose programming language with support for implicit concurrency and visual programming without using such things. Actually using of such low-level things is something that makes support for visual programming harder.

_Implicit concurrency_ means that programmer doesn't have to think about concurrency at all. At the same time any Neva program is concurrent by default. In fact there's no way to write non-concurrent programs. Explicit concurrency is like manual memory management - the great care must be put into. Concurrent programs in conventional langauges are always harder to maintain than regular ones. Not just all Neva programs are concurrent but programmer simply doesn't have a way to interact with concurrency. This is just how it works (thanks to FBP).

_Easy to visualize_ means that the nature of FBP programs is that we do not have [control flow](https://en.wikipedia.org/wiki/Control_flow), but instead we control [data flow](https://en.wikipedia.org/wiki/Dataflow_programming). This is how real electronic components works - there's electricity that flows through connections implementing specific logic. This is also how we document software with visual schemas - sort of boxes connected by arrowes where data flows from one component to another being transformed in someway. But those schemas are usually "dead" - they're not connected with the source code in anyway. FBP allowes to make diagrams source code itself.
