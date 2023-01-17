# FAQ

## General

### Why yet another language?

First - because there's no general-purpose [flow-based](https://en.wikipedia.org/wiki/Flow-based_programming) statically typed visual programming language with meta-programming,runtime-introspection and code generation in mind.

Second - it's not just a language. This project is more like a platform for a flow-based programming. It's a compiler, runtime, type-system, syntax, IDE plugins. All these components are more or less independent from eachother and could be used in separate.

The goal is to create a system that is so easy to use that even a not programmer could create efficient concurrent program that is easy to maintain. So the real answer is - it's not really _yet another_ programming language.

### Why visual language? Aren't visual programs less unmaintainable?

Because [visual cortex exists](https://youtu.be/8Ab3ArE8W3s?t=1220).

First - text is also _visual_ representation (but using sounds or smells is by the way interesting idea). We recognize patterns by looking at code and parse the program's heirarchal structure with braces or offsets.

Second - argument that what we usually call visual programming is less maintanable is simply wrong. This is just different (also visual and more explicit one) form of representing a data of specific structure. Flow-based approach allowes to abstract things away exactly like text-based programming does.

Third - actually there's no dependency on visual programming. Neva designed with support for visual programming as a first-class citizen in mind but in fact it's possible to use text representation. Actually text is used as a source code for version control. There's also, by the way, no dependency on specific syntax.

Neva programs are, first of all, typed graphs that describes dataflows. The paradigm ([FBP]()) allowes to represent they in an infinite ways (including e.g. VR).

### Why Flow Based Programming and not OOP/FP/etc?

1. Higher level programming
2. Implicit concurrency
3. Easy to visualize

_Higher level programming_ means there's no: variables, functions, for loops, classes, methods, etc. All these things are low-level constructions that must be used under the hood as implementation details, but not as the API for the programmer. It's possible to have general purpose programming language with support for implicit concurrency and visual programming without using such things. Actually using of such low-level things is something that makes support for visual programming harder.

_Implicit concurrency_ means that programmer doesn't have to think about concurrency at all. At the same time any Neva program is concurrent by default. In fact there's no way to write non-concurrent programs. Explicit concurrency is like manual memory management - the great care must be put into. Concurrent programs in conventional langauges are always harder to maintain than regular ones. Not just all Neva programs are concurrent but programmer simply doesn't have a way to interact with concurrency. This is just how it works (thanks to FBP).

_Easy to visualize_ means that the nature of FBP programs is that we do not have [control flow](https://en.wikipedia.org/wiki/Control_flow), but instead we control [data flow](https://en.wikipedia.org/wiki/Dataflow_programming). This is how real electronic components works - there's electricity that flows through connections implementing specific logic. This is also how we document software with visual schemas - sort of boxes connected by arrowes where data flows from one component to another being transformed in someway. But those schemas are usually "dead" - they're not connected with the source code in anyway. FBP allowes to make diagrams source code itself.

## Design

### Why compiler has static ports and runtime has givers?

Because if compiler would have givers, they will be a special kind nodes which broke elegance of nodes being just component instances. Because giver is a regular component, it has a specific configuration - a message that it must distribute.

On the other hand, there's 2 types of effects at the runtime that are essentially different. Runtime anyway have concept of effects because if operators and giver is different than operator by the same reason.

### Why have `fromRec`?

The reason is the same as with "static ports" vs "givers as special nodes". Otherwise there would be a special kind of nodes like "record builders" that are different from normal component nodes because they must have a specific configuration - record that they must build.

With `from rec` feature (that is implemented outside of the typesystem, because type system doesn't know anything about ports, it only knows about types) it's possible to say "hey compiler, I want a component with the same inports that this record has fields".

## Type-system

### Why structural subtyping?

Because it's allowes to write much less code.

## Implementation

### Why Go?

It's a perfect match. Go has builtin green threads, scheduler and garbage collector. Even more than that - it has goroutines and channels that are 1-1 mappings to FBP's ports and connections. Last but not least is that it's a pretty fast compiled language.

## FBP

### Are there any differences between Neva and classical FBP

First of all, there's a [great article](https://jpaulm.github.io/fbp/fbp-inspired-vs-real-fbp.html) by J Paul Morrison (inventor of FBP).

There are many differences, among them:

- Unlike in classical FBP, programmer should not write code by hand. Code should be generated. (But it's possible of course to do so)
- FBP itself doesn't have static typing. The idea is that you write code by hand in statically typed langauge like Java or Go and then reuse it as a component
- FBP doesn't have _operators_. There's just _elementary components_ that are written by hand and do specific work. JPaulm did not like the idea of e.g. summing numbers using FBP components
- Neva has lower-level runtime program structure that is more primitive representation of a program. FBP on the other hand doesn't describe anything like that. So there's no `effects` in classical FBP.
- There's no _start port_ and _sig_ interface in FBP
- FBP describes life-cycle of IPs (messages), Neva just uses Go's GC to avoid overhead of using extra memory management on top
- FBP allowes some sort of process lifecycle control. E.g. process can _terminate_ and it's possible to _revive_ it. On the other hand Neva allowes to design component in a way that it can pause and resume its work at signal if that's needed. Otherwise Neva's runtime takes care of when node must be shutdown.

### Why some things have different naming?

- _Message_ instead of information package not to be confused with _IP_ as internet protocol
- _Node_ instead of _process_ 1) not to be confused with _OS processes_ and 2) there are _io nodes_ that are not component instances but part of the _network_
- _Static ports_ instead of _IIPs_ because of not using word _IP_
- Word _worker_ used to highlight nodes that are component instances


### Why optional is base type

Because if it would be a regular record then it would be possible to read it's internal field with value of type `T` that could not be there. 