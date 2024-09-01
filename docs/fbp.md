# Flow Based Programming

Nevalang's dataflow is not pure FBP but FBP influenced Nevalang the most. This document describes differences between Nevalang's and FBP's dataflow.

## Granularity

This might be the most important difference between FBP and Nevalang. FBP is intended for high-level orchestration of lower-level components, written in controlflow languages, while Nevalang is truly general purpose and expects you to write your whole program in Dataflow. For instance, in FBP you are not expected to write math logic, but instead to do it in controlflow language and then use FBP as a glue. In Nevalang, on the other hand, you have math components and you are expected to use them instead.

## Controlflow

Granularity difference leads to differences in how we threat controlflow in our dataflow environment.

FBP has an idea that there are atomic and complementary components. Atomic components are implemented in the "host" platform (e.g. Java) while complementary ones are dataflow.

Nevalang also has atomic/complementary (it's called low-level/native and high-level/custom/normal). However, the difference is that Nevalang provides you with ready to use low level components (written in Go) and you never supposed to write controlflow by yourself. This is exactly the way operator and builtin functions works in controlflow languages. E.g. in Go some functions are implemented in assembler.

The only time when you are intended to write controlflow in Nevalang is when you work on stdlib (which means you are a language contributor) or you write some glue code to integrate Go code into Nevalang or vice versa. Anyway, "normal" way of doing things in Neva is to stay in pure dataflow land as long as possible.

## Garbage Collection, Ownership and Immutability

FBP is not garbage collected and has concept of ownership. Each message has one "process" owner. Message should be either copied or ownership must be passed. Thanks to this data mutations are possible. However, this has some downsides like not being able to implement fan-out without using some special components, because FBP tries to avoid mutation and performance (copying) problems.

Nevalang is garbage collected, there is no concept of ownership/borrowing and all data is always immutable. This has performance downside but on the other hand you don't have to think about ownership or mutation problems. E.g. data race is not possible in Nevalang (don't be confused with race condition which is possible in both Nevalang and pure FBP).

Mutations are possible with using special unsafe package from stdlib but should be considered as unwanted optimization tricks when it's clear that nothing else can't help. However, even this way component can only mutate messages inside its own network, not child or parent, so there's some visibility scope protection. Another way to say that Nevalang allows data mutation is to say that it's implemented in controlflow langauge. That's true from implementation point, but not end-user experience. Exact same argument applies for the controlflow interop.

## Nodes Shutdown and Restart

FBP and Nevalang both has concept of nodes, but in FBP it's called "process". The similarity is that in both cases we are talking about components instances. However, FBP's process has concept of "state" which means process could be started, suspended, restarted or shutted down. There are some technics that allow manipulate this behaviour by enabling/disabling some specific parts of the network.

On the other hand, Nevalang's nodes always run. After the program started, all the nodes there are started. Node is not doing any work until there's some condition met e.g. message received to some port, but this happens automatically without user's intention. In other words, Nevalang's node is always started, suspended, restarted and disabled automatically when needed.

## Static Typing

FBP doesn't have concept of static typing except coding atomic components in statically typed controlflow. Dataflow part of the FBP is dynamically typed like JavaScript.

Obviously, you have to write more runtime validations and tests. Also IDE support will be worse. Because of that Nevalang comes with static type system.

Nevalang's type system is relatively powerful, e.g. it supports generics and structural sub-typing. However, it's also relatively simple, which means you shouldn't spend a lot of time working with types, but also some things cannot be expressed - it may be "unsound" for something who came from something like static typed FP.

## Similarities

Both dataflow, both support implicit parallelism. Nevalang's dataflow looks more like FBP than something else (given the differences we discussed). Also a lot of terminology is shared.

## Terminology

**FBP / Nevalang**

- Component / Component
  - Atomic / Native
  - Complementary / Normal
- Process / Node
- Connection / Connection
- Ports / Ports
- IP / Message
- IIP / Constant
- IP Tree / Structure or Dictionary
