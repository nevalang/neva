---
Title: FAQ
---

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

## Is Neva "classical FBP"?

No. But it inherits so many ideas from it that it would be better to use word "FBP" than anything else. There's a [great article](https://jpaulm.github.io/fbp/fbp-inspired-vs-real-fbp.html) by mr. J. Paul Morrison (inventor of FBP) on this topic.

Now here's what makes Neva different from classical FBP:

- Neva has C-like syntax for its textual representation while FBP syntax is somewhat esoteric. It's important to node though that despite C-like syntax Neva programs are 100% declarative
- Neva doesn't let you program in "implementation-level" language like Go (similar to how Python doesn't let you program in assembly). FBP on the other hand forces you to program in langauges like Go or Java to implement elementary components.
- Neva introduces builtin observability via runtime interceptor and messages tracing, FBP has nothing like that
- Existing FBP implementations are essentially interpreters. Neva has both compiler and interpreter.
- Neva is statically typed while FBP isn't. FBP's idea is that you write code by hand in statically typed langauge like Go or Java and then use it in a non-typed FBP program, introducing runtime type-checks where needed
- Neva have _runtime functions_. In FBP there's just _elementary components_ that are written by programmer. Mr. Morrison did not like the idea of having "atomic" components like e.g. "numbers adder"
- Neva introduces hierarchical structure of program entities and package management similar to Go. Entities are packed into reusable packages and could be either public or private.
- Neva leverages existing Go's GC, FBP on the other hand introduces IP's life-cycle
- Neva's concurrency model runs on top of Go's scheduler which means it uses CSP as a lower-level fundament. FBP implementations on the other hand used to use shared state with mutex locks
- Neva has low-level program representation (LLR). FBP on the other hand doesn't describe anything like that

Also there's differences in naming:

- _Message_ instead of _IP (information package)_ not to be confused with "IP" as _internet protocol_
- _Node_ instead of _process_ 1) not to be confused with _OS processes_
- _Bound inports_ instead of _IIPs_ because of not using word _IP_

## Why need array-ports?

Every time we need to somehow combine/accumulate/reduce several sources of data into one e.g.

- create list of 3 elements based on outputs of 3 outports
- create structure with field-values from several outports
- substract values from left to right

Ok but can't we substract values and do other stuff like that by simply passing lists around? Well, we have to create that list right somehow? It's fine if you already have it (let's say from JSON file you got from server) but what if you need to build it?

## Why component can't read from it's own array-inports by index?

Imagine you do stuff like:

```neva
in.foo[0] -> ...
in.foo[1] -> ...
```

Now what will happen if parent node will only use `0` slot of your `foo` array-inport? Should it block forever? Or maybe should the program crash? Sounds not too good.

The other way we could handle this is by making analyzer ensure that parent of your component uses your `foo` array-inport with exactly `0` and `1` slots. The problem is that makes array-ports useless. Why even have them then? The whole point of array-ports is that you don't know how many slots are going to be used. And that makes your component flexible. It allows you to create components that can do stuff like "sum all numbers from all inports no matter how many of them are present".

Besides, you can already say "use my component with exactly two values" already and you don't need any array-ports for that at all! All you need in that case is simply create two inports.

Having that said we must admit that it's impossible to allow component read form it's own array-inports by index and still having type-safety.

Also think about _variadic arguments_ in Go. It's not safe to refer to `...args` by index (even though it's possible because Go compiler won't protect you).

## Why component can read from sub-node's array-outports by index?

Isn't it unsafe to read from array-outports by index? We've restricted that for component itself by banning the ability to read form self outports by index. Why allow read from sub-node outports by index then?

Well, it turns out there are critical cases where we need that. One of them is "routing" - where you have some data on the input and you need to figure out, based on some predicate, where to send it further. Like if you have a web-server and you received a request, you need to route it to specific handler based on the path that this request contains.

It's possible to do that with sequence of if-else though but that would be really tedious and error-prone. That also would make your network more nested and less straightforward.

### Can't we implement syntax sugar for that?

It's possible to introduce some sort of syntax sugar where user interacts with array ports and under the hood it's just a bunch of if-elses. But that actually makes no sense. If we have array-outports as a part of the language interface, we have them anyway. We also have use-cases for array-inports which means there are other reasons why have array ports. And finally it would be better for performance to have one low-level control-flow component implemented in implementation langauge and not Nevalang instead of compiling one high-level component to another big high-level component. One might ask - but we did that for Lock, what's the difference? The thing is with lock we are not replacing one component usage with the another, like we would in case of replacing some kind of "router" with bunch if if-elses. We simply insert implicit code, that is assumed by the higher level constructs like only exist at the level of the source code and not the real machinery.

## Why outports usage is optional and inport usage is required?

Indeed when component `A` uses `B` as it's sub-component (when it instantiates a _node_ with it) in it's _network_ it's _enforced_ to use _all_ the inports of `B` and it's _at least one_ outport. It doesn't have to use all the outports though.

This is because inports are requirements - they are needed to receive the data that component _needs_ to produce result. Outports on the other hands are options. They are results that parent network might need to a sertain degree. For instance if `B` have outports `foo` and `bar`, it's completely possible that `A` only needs `foo` and have nothing to do with `bar`.

This leads us to the need of the `Void` (builtin) component. This is the only component that doesn't have outports. It is used for discarding the unwanted data. If there would be no syntactic sugar for that, then we would have to explicitly create `void` nodes and use it in places like this:

```neva
nodes {
    b B
    void Void
}
net {
    // ...
    b.bar -> void.v // discard all messages from `bar` outport
}
```

It's not the problem that it's tedious (even though it is, imagine having 10 unwanted outports in your network which is completely possible). The real problem is that by discarding some outports user is in danger of programming the dataflow in the wrong way.

Imagine that `B` has outports `v` (for valid results) and `err` (for error messages). It fires either `v` or `err` and never both at the same time. And we want out program to terminate if there's nothing to do left. Consider this code:

```neva
Main(enter) (exit) {
    nodes {
        b B
        void Void
        print Print
    }
    net {
        in.enter -> b.sig
        b.err -> void.v // ignore the `err` outport, only handle happy path
        b.v -> print.v
        print.v -> out.exit
    }
}
```

We print the success result and then terminate. If there is no success result and only error we well... do nothing. And that's bad. What we should do instead is this:

```neva
// ...
net {
    in.enter -> b.sig

    // print both result and error
    b.err -> print.v
    b.v -> print.v

    // and then exit
    print.v -> out.exit
}
```

As you can see it's easy to get in trouble by ignoring some outports (especially the error ones). If user wouldn't have the ability to do so he would have to do _something_ with `err` message. Anyway there would still be two problems...

1. Even then user still _can_ send the data in the wrong way. E.g. send the `err` message back to `b.sig` or `print` it but then send the `print.v` back to the `print` forming an endless loop. This kind of _logical_ mistakes are hard to catch. Making the language _that_ safe would also make it much more complicated (think of Haskell or Rust (where we still have such kinds of problems!)).
2. Sometimes we have _nothing to do_ with unwanted data. We don't wanna print it or even send downstream (because that would simply delay the question what to do with unwanted data). This is the reason why `Void` doesn't have outports. Otherwise a paradox arises.

This leads us to a conclusions:

- There must be a way to omit unwanted data, whether it's explicit (`Void`) or implicit sugar
- It's impossible to make langauge 100% safe without sacrificing the simplicity of use

As we saw explicit Void doesn't solve these problems so why not introduce sugar? Let's allow user to simply omit unwanted data and let the compiler implicitly insert `Void` under the hood. The logical mistakes? Well... They are "unsolvable" anyway.

## Why there's no int32, float32, etc?

Because that's a simple language. Lack of ability to configure bit-size of the number but still being able to choose between integers and floats is the compromise that seems to be reasonable. Probably Python is a good example of that too.

## Why have integers and floats and not just numbers?

1. Overflow issues: if you only have `number`, probably represented as a `float64` in memory, your maximum safe number is bigest float64. Integer can store "bigger" values because it doesn't have to store (floating) precision. This is especially important when you work with big numbers.

2. Performance Overhead: Floating-point operations are generally slower than integer operations. In a system where all numbers are floating-point, operations that could have been more efficient with integers suffer a performance penalty.

3. Predictability in Comparisons: Floating-point arithmetic can lead to non-intuitive results due to precision errors, making comparisons and certain calculations (like summing a large list of numbers) less predictable.

4. Lack of Type Safety: The absence of distinct types can lead to bugs that are hard to detect, as the language won't provide errors or warnings when performing potentially erroneous operations between different kinds of numeric values.

## Why no literal senders in component networks?

In conventional languages like e.g. Python we can simply spell

```python
print(42)
```

To do same thing in Nevalang you must create `const`:

```neva
const {
    msg int 42
}
components {
    Main(enter) (exit) {
        nodes {
            print Print
        }
        net {
            msg -> print.v
        }
    }
}
```

Wouldn't it be great to allow user to simply say?

```neva
42 -> print.v
```

Turns out there's a problem with that approach. Under the hood (just like with `const` sender) we need to create `Const` node. But in the first case we use name of the constant `msg` as the node name so it desugares down to

```neva
nodes {
    #runtime_func_msg(msg)
    msg Const<int>
}
```

In case of `42` there's no identifier that we can use and thus we have to generate it. That's not a problem until we debug our program but as soon as we will we have to face some autogenerated node name that we have no idea where came from.

This will probably happen quite often because when you don't have to create constant you probably won't. On the other hand with current approach you have to do that all the time. As a good thing - you won't have to deal with ambiguity - need a static value? Create const!

## What is the motivation behind putting exactly these entities under builtin package?

1. Frequently used
2. Used internally by compiler (desugarer)

## Why constants are implemented this way? (Lock-semantic instead of trigger semantic)

Const nodes are implemented like infinite loops that constantly sends messags across their receivers. This covers all the usecases but also requires locks because we usually want control when we send constant messages.

Alternative to that design would be "trigger" semantics where we have some kind of `sig` inport for const nodes. Instead of starting at the program startup such trigger component would wait for external signal and then do the usual const stuff (infinite loop with sending messages).

**The problem #1 with this approach - we still needs locks**. It's not enough to trigger infinite loop. E.g. in "hello world" example nothing would stop `msg` const node to send more than 1 `hello world` message to `print`.

**Possible solution for that would be to change semantics and remove infinite loop logic**. Make const node send signal only after we trigger it via sig port. The problem with this approach is that there is many usecases where we want infinite loop behavior. Think of initial inports - e.g. `requestSender` component with `data` and `url` inports where `data` is dynamic and `url` is static. It's not enough to send static url value just once (`requestSender` must remember it, we don't go that way because that leads to bad design where components know where they get data from - this is huge violation of transport vs logic separation).

This problem by itself is fixable with using external sources like signals. When we have some static inport we usually have some kind of dynamic data that must be used in combination with it. Even though it would lead to making networks more complicated (locks do this too though), it's possible solution. But we have another problem:

**Problem #2** - `$` syntax sugar.

Another problem with previous solution (const nodes have sig inport and they send one message per one signal) is how use `$` syntax sugar.

Currently it's possible to _refer constants_ in network like this:`$msg -> ...`

This won't be the thing because we have to have not just entity reference but regular ports like `$msg.sig` and `$msg.v`. This is not a disaster but rather strange API.

Where this `$msg` node come from? Is it clear to me as a user that there are implicit nodes $ prefix for every constant that I can refer? Why these `sig` and `v` inports? Because this is how `std/builtin.Const` works? Why do I have to know this? Why do I have to know how syntax sugar is implemented under the hood in compiler?

Finally another possible solution to that could be `Trigger` components in combination with `Const` components. The difference would be that const behaves like infinite loops that requires locks and triggers behaves like single sending triggers (no lock required).

Problems with this solution are:

1. Now we have 2 ways to do the same thing. Do I need to use const? Or trigger? When to use what?
2. Trigger either must be used in combination with `#runtime_func_msg` directive (this violates principle that user must be able to program without directvies) or there must be sugar for triggers.

It's possible in theory to create sugar for triggers but the language could be too complicated with so many syntax features. This conclusion by itself is questionable but in combination with the first problem - having 2 ways to use static values. Looks like it's better not to have triggers.

All this leads to a conclusion that the only good semantic for constants is the current ones - infinite loops that requires locks. The need to have locks is not a fancy thing by itself, but "then connections" sugar made this pretty simple.

## Why have sub-streams?

In programming we range over collections all the time. We do that either via _loops_ or higher order functions (recursion). There's a problem with both of these in FBP:

We don't have _code as data_ (we can't pass components around like we pass functions in conventional languages)

To implement loop we need:

1. Mutable state. (We don't have one! We can simulate one using network loops tho). After we simulate mutable
2. Condition to check whether current cursor is still less than the length of the list
3. Length of the list

It would be very verbose to do such things all the time so we can imagine some kind of generic `for` component that takes `list<T>` and sends single values `T`.

The biggest problem is - how do we know that list ended? How do we know that the previous element was the last one and the current element is the first element of the new list that just arrived? Without knowing this we loos information about the list boundary. And that is huge problem. In conventional programming we always know that. Without this simplest iteration patterns like `map` are impossible to implement.

Possible solution to this (without introducing sub-streams) would be adding some kind of _signal_ that "the list just ended". One might think that `For` component could simply have two outports `v` and `sig`. No, it can't. In this case `sig` cannot be separate port because it needs to be in the exact same _steam_ as the elements themselves. Otherwise it's unclear how to be sure that we synchronized both streams (flows) together. Streams are concurrent and the order of messages across different streams is often unpredictable.

That leads us to conclusion - such `For` component must have one outport (or at least it must not have separate `sig` outport). It instead must send not just `T` values (single elements of the colletion), but instead it must send some kind of structures. The shape must be something like this

```neva
types {
    Element<T> {
        v T
        isLast bool
    }
}
```

Congratulations! You just discovered sub-streams.

## Why sub-streams are not like in classical FBP?

John Paul Morrison, creator of Flow-Based Programming created _Sub streams_ as a part of the FBP. The problem that he was solving wasn't just iteration over collections. It was work with structured data. Sub-streams are the way we transfer structured objects in his paradigm.

In Nevalang we have `struct`, `map` and `list` for that. We don't need to create "flat nested" sub-stream like this `( (1 2 3) (4 5 6) ) ( (7 8 9) )` to move two lists, we can simply move them like regular messages across the stream `-> l1 l2 ->`. The downstream component that receives `l1` and `l2` can then unpack them into sub-streams and process their individual elements.

At the end it must pack them pack tho. This is _maybe_ where classical FBP outperforms Nevalang. We have to spend time on destructuring and structuring back. However, the **data in the outside world is structured**. We usually work with some kind of relational data, JSON, Protobuf, etc.

## Isn't it the problem that component that works with type `T` cannot operate on `SubStreamItem<T>`

This is just how type-system works. We don't want to have a lot of special cases here and there. It's not a big deal also.

If you have a component `C1` that takes `T` and you want to operate on `SubStreamItem<T>` all you need is to create a wrapper. That wrapper will receive `SubStreamItem<T>` and use `C1` inside of it with `.v` struct selection.

If you need to continue sub-stream you simply send `SubStreamItem<T>` from you wrapper component downstream. Or `SubStreamItem<WhateverYouWant>` (probably preserving `isLast` value).

It's either you continue sub-stream or you do not. Depending on what your're doing (maybe you're counting sub-stream items so you just sends `int` eachtime sub-stream ends).


## Why `out:exit` of the `Main` is't `int`?

This is the question about why `out:exit` isn't interpreted as exit code.

The things is - you don't always have `int` as your exit condition. That's why it's `any`.

Ok, but why then we don't check if that `any` is actually `int` under the hood and then interpret it as an exit code?

Well, we can do that. But that would lead to situations where you accidentally have `int` like your exit condition but don't actually want it to be your exit code. Such cases are non obvious and will require you to somehow check that you send exit code you want.

This problem gets bigger when you have `any` or _union_ `... | int` outport that is directed to `out:exit` - you'll have to check whether value isn't an `int`. Otherwise you're at risk of terminating with wrong code.

**Exit codes are important**. Shell scripts and CI/CD depends on that. Most of the time you want your exit code to be `zero`. Non-zero exit code is not happypath, it's more rare. Having corner case like a base design decision is not what we want.

## Why structural subtyping?

1. It allowes write less code, especially mappings between records, vectors and maps of records
2. Nominal subtyping doesn't protect from mistake like passing wrong value to type-cast

## Why have `any`?

First of all it's more like Go's `any`, not like TS's `any`. It's similar to TS's `unknown`. It means you can't do anything with `any` except _receive_, _send_ or _store_ it. There are some [critical cases](https://github.com/nevalang/neva/issues/224) where you either make your type-system super complicated or simply introduce any. Keep in mind that unlike Go where generics were introduced almost after 10 years of language release, Neva has type parameters from the beggining. Which means in 90% of cases you can avoid using of `any` and panicking.