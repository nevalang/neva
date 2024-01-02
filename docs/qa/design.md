# Design

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

<!-- ## Why have `fromRec`?

The reason is the same as with "static ports" vs "givers as special nodes". Otherwise there would be a special kind of nodes like "record builders" that are different from normal component nodes because they must have a specific configuration - record that they must build.

With `fromRec` feature (that is implemented outside of the typesystem, because type system doesn't know anything about ports) it's possible to say "hey compiler, I want a component with the same inports that this record has fields". -->
