# Nevalang

Welcome to Nevalang - general purpose dataflow programming language.

In this language we don't have control-flow. It means we can't call, return, break, continue or goto. There are no functions, no for loops, no variables. Instead we have nodes connected through their input and output ports for message passing. This is what's called dataflow programming.

There are no coroutines, channels or mutexes. You write concurrent program just by having parallel connections in your network. If there is machine capacity, code will be executed in parallel. It's called implicit parallelism.

Compiler performs strong static type-checking. Language has interfaces and generics for polymorphic code. You can emit ready to deploy machine code or Go, to integrate with existing codebase.
