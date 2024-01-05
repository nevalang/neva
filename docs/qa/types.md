# Type-system

## Why structural subtyping?

1. It allowes write less code, especially mappings between records, vectors and maps of records
2. Nominal subtyping doesn't protect from mistake like passing wrong value to type-cast

## Why have `any`?

First of all it's more like Go's `any`, not like TS's `any`. It's similar to TS's `unknown`. It means you can't do anything with `any` except _receive_, _send_ or _store_ it. There are some [critical cases](https://github.com/nevalang/neva/issues/224) where you either make your type-system super complicated or simply introduce any. Keep in mind that unlike Go where generics were introduced almost after 10 years of language release, Neva has type parameters from the beggining. Which means in 90% of cases you can avoid using of `any` and panicking.

## Why there's no Option/Maybe data type?

## Short answer

We don't have that problem in FBP that `Option` types solves in conventional langauges. Because components can have multiple inports and outports for different cases and it's hard to mix different flows together.

## Long answer

In FBP data is data and code is code. This means we can't pass functions (or components) as parameters or store them inside objects. As a result we cannot have objects with "behavior" since they cannot have methods. Since there are no OOP-objects in the language, having `Option/Either/Maybe/Result/etc.` doesn't really brings any advantages. The good part is **there's actually no need for that**.

For example in conventional language (e.g. Go), where we _control execution flow_, it's possible to do this:

```go
type Foo struct { bar int } // define type Foo that is a structure with integer field bar
var foo *Foo := f() // call function f that returns pointer to value of type Foo
print(foo.bar) // dereference foo and access bar field
```

There's no guarantee that `f()` won't return `nil`. This code will crash with panic even though Go is statically typed. This is the problem that `Option` types solve. **And that problem doesn't exists in FBP.** The source of this problem is control over execution flow - _we use low-level primitives like variables and pass them around expecting them to have some specific state_. And we encounter the problem where a flow for non-nil value actually faces nil value.

But as soon as we look at this program from the _dataflow_ perspective, _where we control data flow_ instead, we'll see that we have 2 flows here - one for `nil` value and one for non-nil. In Neva such program looks like this:

```neva
Main() () {
    net {
        f.foo[bar] -> print.v
    }
}
```

If we want we can cover both flows:

```neva
Main() () {
    net {
        f.foo[bar] -> print.v
        f.err -> print.v
    }
}
```

The thing is - there's a separate flow for `err` and for `foo`. There's no way we unintentionally mix them and use `err` instead of `foo`. All we need to do is to make sure our `f` returns `Foo`. There also no need to introduce absence of value with pointers or nils. If there's no `foo` then simply nothing happens. No value is sent from `f.foo` outport until there's an actual value of type `Foo`.

Of course there's unions so nothing stops as from using `Foo | nil`. We need this to process external data e.g. json from the network. But for programming dataflows? There are inports and outports connected to each other. It's that simple.
