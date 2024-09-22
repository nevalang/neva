# Contributing

Neva is a relatively small and simple language, don't be intimidated and feel free to hack around and reach out to maintainers if you need help.

Start from reading [ARCHITECTURE.md](./ARCHITECTURE.md) and [Makefile](./Makefile).

## Requirements

- Go: https://go.dev/doc/install
- Make: https://www.gnu.org/software/make/#download
- NodeJS and NPM: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm/
- Antlr: `pip install antlr4-tools`
<!-- - Tygo: `go install github.com/gzuidhof/tygo@latest` -->

### VSCode

Not required but recommended:

- [nevalang](https://marketplace.visualstudio.com/items?itemName=nevalang.vscode-nevalang)
- [antlr4](https://marketplace.visualstudio.com/items?itemName=mike-lischke.vscode-antlr4)
- [tmlanguage](https://marketplace.visualstudio.com/items?itemName=pedro-w.tmlanguage)
- [markdown-mermaid](https://marketplace.visualstudio.com/items?itemName=bierner.markdown-mermaid)

## Development

## Syntax (ANTLR, Parser)

1. Make changes to `neva.g4` and corresponding `*.neva` files in the repo
2. If something doesn't work, run `/parser/smoke_test`
3. To debug deeper, make sure `neva.g4` is opened in the editor and launch VSCode's `ANTLR` debug task

<!-- ## VSCode Extension

VSCode extension depends on types defined in the `sourcecode` and `typesystem` packages so it's dangerous to rename those. If you going to do so, make sure you did't brake TS types generation.

Check out [tygo.yaml](./tygo.yaml). and `CONTRIBUTING.md` in "vscode-neva" repo. -->

## Learning Resources

### Dataflow

- [Nevalang's Documentation](./docs/README.md)
- [Flow-Based Programming: A New Approach to Application Development](https://jpaulmorrison.com/fbp/1stedchaps.html)
- [Dataflow and Reactive Programming Systems: A Practical Guide](https://www.amazon.com/Dataflow-Reactive-Programming-Systems-Practical/dp/1497422442)

### Golang

Advanced understanding of concurrency is very helpful:

- [Concurrency is not parallelism](https://go.dev/blog/waza-talk)
- [Share Memory By Communicating](https://go.dev/blog/codelab-share)
- [Go Concurrency Patterns: Timing out, moving on](https://go.dev/blog/concurrency-timeouts)
- [Go Concurrency Patterns: Context](https://go.dev/blog/context)
- [Go Concurrency Patterns: Pipelines and cancellation](https://go.dev/blog/pipelines)

## Design Principles

Nevalang adheres to the following key principles:

1. **Fail-fast**: Programs should fail at compile-time or startup, not during runtime.
2. **Unsafe but efficient runtime**: Runtime assumes program correctness for speed and flexibility.
3. **Compiler directives**: Powerful but unsafe tools for advanced users and language developers.
4. **Visual programming ready**: Language design considers future visual programming tools.

## Internal Implementation Q&A

Here you'll find explanations for specific implementation choices.

### Why structures are not represented as Go structures?

It would take generating Go types dynamically which is either makes use of reflection or codegeneration (which makes interpreter mode impossible). Maps have their overhead but they are easy to work with.

### Why nested structures are not represented as flat maps?

Indeed it's possible to represent `{ foo {bar int } }` like `{ "foo/bar": 42 }`. The problem arise when when we access the whole field. Let's take this example:

```
types {
    User {
        pet {
            name str
        }
    }
}

...

$u.pet -> foo.bar
```

What will `foo.bar` actually receive? This design makes impossible to actually send structures around and allows to operate on non-structured data only.

### Why Go?

It's a perfect match. Go has builtin green threads, scheduler and garbage collector. Even more than that - it has goroutines and channels that are 1-1 mappings to FBP's ports and connections. Last but not least is that it's a pretty fast compiled language. Having Go as a compile target allows to reuse its state of the art standart library and increase performance for free by just updating the underlaying compiler.

### Why compiler operates on multi-module graph (build) and not just turns everything into one big module?

Imagine you have `foo.bar` in your code. How does compiler figures out what that actually is? In order to do that it needs to _resolve_ that _reference_. And this is how _reference resolution_ works:

First, find out what `foo` is. Look at the `import` section in the current file. Let's say we see something like:

```neva
import {
    github.com/nevalang/x/foo
}
```

This is how we now that `foo` is actually `github.com/nevalang/x/foo` imported package. Cool, but when version of the `github.com/nevalang/x` we should use? Well, to figure that out we need to look out current _module_'s _manifest_ file. There we can find something like:

```yaml
deps:
  - github.com/nevalang/x 0.0.1
```

Cool, now we now what _exactly_ `foo` is. It's a `foo` package inside of `0.0.1` version of the `github.com/nevalang/x` module. So what's the point of operating on a nested multi-module graph instead of having one giant module?

Now let's consider another example. Instead of depending on `github.com/nevalang/x` your code depends on `submodule` and that sub-module itself depends on `github.com/nevalang/x`

You still have that `foo.bar` in your code and your module still depends on `github.com/nevalang/x` module. But now you also depends on another `submod` sub-module that also depends on `github.com/nevalang/x`. But your module depends on `github.com/nevalang/x` of the `0.0.1` version and `submod` depends on `1.0.0`.

Now we have a problem. When compiler sees `foo.bar` in some file it does import lookup and sees `github.com/nevalang/x` and... does not know what to do. To solve this issue we need to lookup current module manifest and check what version `github.com/nevalang/x` _this current module_ uses. To do that we need to preserve the multi-module structure of the program.

One might ask can't we simply import things like:

```neva
import {
    github.com/nevalang/x@0.0.1
}
```

That actually could solve the issue. The problem is that now we have to update the source code _each time we update our dependency_. That's a bad solution. We simply made probramming harder to avoid working on a compiler. We can do better.

### Why `#bind` does not accept literals?

Indeed it would be handy to be able to do stuff like this:

```neva
nodes {
    #bind(str "hello world!")
    const Const<str>
}
```

This would make desugarer much simpler (no need to create all this virtual constants), and not just for const senders but for struct selectors too.

However, to implement this we need to be able to parse literals inside `irgen`. Right now we already introduce dependency for parsing entity references, but for arbitrary expressions we need the whole parser.

Of course, it's possible to hide actual parser implementation behind some kind of interface defined by irgen but that would make code more complicated. Besides, the very idea of having parser inside code-generator sounds bad. Parsing references is the acceptable compromise on the other hand.

### Why Analyzer knows about stdlib? Isn't it bad design?

At first there was a try to implement analyzer in a way that it only knows about the core of the language.

But turns out that some flows in stdlib (especially `builtin` package, especially the ones that uses `#extern` and `#bind` directives) are actually part of the core of the language.

E.g. when user uses struct selectors like `foo.bar/baz -> ...` and then desugarer replaces this with `foo.bar -> structSelectorNode("baz") -> ...` (this is pseudocode) we must ensure that type of the `bar` is 1) a `struct` 2) has field `baz` and 3) `baz` is compatible with whatever `...` is. _This is static semantic analysis_ and that's is work for analyzer.

Actually every time we use compiler directive we depend on implicit contract that cannot be expressed in the terms of the language itself (except we introduce abstractions for that, which will make language more complicated). That's why we have to analyze such things by injecting knowledge about stdlib.

Designing the language in a way where analyzer has zero knowledge about stdlib is possible in theory but would make the language more complicated and would take much more time.

### Why desugarer comes after analyzer in compiler's pipeline?

Two reasons:

1. Analyzer should operate on original "sugared" program so it can found errors in user's source code. Otherwise found errors can relate to desugar implementation (compiler internals) which is not the compilation error but debug info for compiler developers. Finally it's much easier to make end-user errors readable and user-friendly this way.
2. Desugarer that comes before analysis must duplicate some validation because it's unsafe to desugar some constructs before ensuring they are valid. E.g. desugar struct selectors without knowing fir sure that outport's type is a valid structure. Also many desugaring transformations are only possible on analyzed program with all type expressions resolved.

Actually it's impossible to have desugarer before analysis. It's possible to have two desugarers - one before and one after. But that would make compiler much more complicated without visible benefits.

### Why union types are allowed for constants at syntax level?

You indeed can declare `const foo int | string = 42` and that won't make much sense. The problem it's not enough to restrict that at root level, you also have to recursively check every complex type like `struct`, `list` or `map`. And that is impossible to make at syntax level and require work in analyzer. This is could be done in the future when we cover more important cases.

### Why we have special syntax for union?

We don't have sugar for `maybe<T>` and `list<T>` so why would we have this for unions? The reason is union is special for the type system. It's handled differently at the level of compatibility checking and resolving.

However it's not `struct` where we _technically_ have to have some "literal" syntax. It's possible in theory to have just `union<T1, T2, ... Tn>` like e.g. in Python but would require _type-system_ known about `union` name and handle this reference expressions very differently. In fact this will only make design more complicated because we _pretend_ like it's regular type instantiation consisting of reference and arguments but in fact it's not.

Lastly it's just common to have `|` syntax for unions.

### Why type system supports arrays?

Because type-system is public package that can be used by others to implement languages (or something else constraint-based).

Since there's no arrays at the syntax and internal representation levels then there's no performance overhead. Also having arrays in type system is not the most complicated thing so removing them won't save us much.

### Why isn't Nevalang self-hosted?

- Runtime will never be written in Nevalang itself because of the overhead of FBP runtime on to of Go's runtime. Go provides exactly that level of control we needed to implement FBP runtime for Nevalang.
- Compiler will be someday rewritten in Nevalang itself but we need several years of active usage of the language before that

There's 2 reasons why we don't rewrite compiler in Nevalang right now:

1. Language is incredibly unstable. Stdlib and even the core is massively changing these days. Compiler will be even more unstable and hard to maintain if we do that, until Nevalang is more or less stable.
2. Languages that are mostly used for writing compilers are eventually better suited for that purpose. While it's good to be able to write compiler in Nevalang without much effort, it's not the goal to create a language for compilers. Writing compilers is a good thing but it's not very popular task for programmers. Actually it's incredibly rare to write compilers at work. We want Nevalang to be good language for many programmers.
